package audit

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/apps/backend/internal/config"
	"github.com/itspablomontes/fleming/pkg/protocol/audit"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// Handler handles HTTP requests for audit logs.
type Handler struct {
	service              Service
	chainClient          ChainAnchorer
	chainEndpointEnabled bool
}

// NewHandler creates a new audit handler.
func NewHandler(service Service, chainClient ChainAnchorer) *Handler {
	env := config.NormalizeEnv(os.Getenv("ENV"))

	enabled := env == "dev"
	if !enabled {
		enabled = os.Getenv("ENABLE_CHAIN_ANCHOR_ENDPOINT") == "true"
	}

	return &Handler{
		service:              service,
		chainClient:          chainClient,
		chainEndpointEnabled: enabled,
	}
}

// RegisterRoutes registers audit endpoints.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	audit := rg.Group("/audit")
	{
		audit.GET("", h.HandleGetLogs)
		audit.GET("/entries/:id", h.HandleGetEntry)
		audit.GET("/resource/:resourceId", h.HandleGetByResource)
		audit.GET("/query", h.HandleQuery)
		audit.GET("/verify", h.HandleVerify)
		audit.GET("/verify/:root", h.HandleVerifyRoot)
		audit.POST("/merkle/build", h.HandleBuildMerkle)
		audit.GET("/merkle/batches", h.HandleListMerkleBatches)
		audit.GET("/merkle/:batchId", h.HandleGetMerkleBatch)
		audit.POST("/merkle/verify", h.HandleVerifyMerkle)
		if h.chainEndpointEnabled {
			audit.POST("/merkle/:batchId/anchor", h.HandleAnchorMerkleBatch)
		}
	}
}

// HandleGetLogs returns the latest audit entries for the current user.
func (h *Handler) HandleGetLogs(c *gin.Context) {
	address, exists := c.Get("user_address")
	actor, ok := address.(string)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	entries, err := h.service.GetLatestEntries(c.Request.Context(), actor, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

// HandleVerify performs a check of the entire chain integrity.
func (h *Handler) HandleVerify(c *gin.Context) {
	valid, err := h.service.VerifyIntegrity(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "integrity check failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   valid,
		"message": "Audit chain integrity verified",
	})
}

func (h *Handler) HandleGetEntry(c *gin.Context) {
	entryID := c.Param("id")
	if entryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entry ID is required"})
		return
	}

	entry, err := h.service.GetEntryByID(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch entry"})
		return
	}
	if entry == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"entry": entry})
}

func (h *Handler) HandleGetByResource(c *gin.Context) {
	resourceID := c.Param("resourceId")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource ID is required"})
		return
	}

	entries, err := h.service.GetEntriesByResource(c.Request.Context(), resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch entries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

func (h *Handler) HandleQuery(c *gin.Context) {
	filter := audit.NewQueryFilter()

	if actor := c.Query("actor"); actor != "" {
		address, err := types.NewWalletAddress(actor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actor address"})
			return
		}
		filter.Actor = address
	}

	if resourceID := c.Query("resourceId"); resourceID != "" {
		filter.ResourceID = types.ID(resourceID)
	}

	if resourceType := c.Query("resourceType"); resourceType != "" {
		filter.ResourceType = audit.ResourceType(resourceType)
	}

	if action := c.Query("action"); action != "" {
		act := audit.Action(action)
		if !act.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
			return
		}
		filter.Action = act
	}

	if start := c.Query("startTime"); start != "" {
		ts, err := types.ParseTimestamp(start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startTime"})
			return
		}
		filter.StartTime = &ts
	}

	if end := c.Query("endTime"); end != "" {
		ts, err := types.ParseTimestamp(end)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endTime"})
			return
		}
		filter.EndTime = &ts
	}

	if limit := c.Query("limit"); limit != "" {
		value, err := strconv.Atoi(limit)
		if err != nil || value < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		filter.Limit = value
	}

	if offset := c.Query("offset"); offset != "" {
		value, err := strconv.Atoi(offset)
		if err != nil || value < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		filter.Offset = value
	}

	entries, err := h.service.QueryEntries(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query entries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

type merkleBuildRequest struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func (h *Handler) HandleBuildMerkle(c *gin.Context) {
	address, exists := c.Get("user_address")
	actor, ok := address.(string)
	if !exists || !ok || actor == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req merkleBuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var startTime time.Time
	var endTime time.Time
	if req.StartTime != "" {
		ts, err := types.ParseTimestamp(req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startTime"})
			return
		}
		startTime = ts.Time
	}
	if req.EndTime != "" {
		ts, err := types.ParseTimestamp(req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endTime"})
			return
		}
		endTime = ts.Time
	}

	batch, tree, err := h.service.BuildMerkleTree(c.Request.Context(), actor, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build merkle tree"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"batch": batch,
		"root":  tree.Root,
	})
}

func (h *Handler) HandleListMerkleBatches(c *gin.Context) {
	address, exists := c.Get("user_address")
	actor, ok := address.(string)
	if !exists || !ok || actor == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 25
	if limitStr != "" {
		v, err := strconv.Atoi(limitStr)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		limit = v
	}

	offset := 0
	if offsetStr != "" {
		v, err := strconv.Atoi(offsetStr)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		offset = v
	}

	batches, err := h.service.ListBatches(c.Request.Context(), actor, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list merkle batches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"batches": batches})
}

func (h *Handler) HandleGetMerkleBatch(c *gin.Context) {
	address, exists := c.Get("user_address")
	actor, ok := address.(string)
	if !exists || !ok || actor == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID := c.Param("batchId")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch ID is required"})
		return
	}

	batch, err := h.service.GetBatch(c.Request.Context(), actor, batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch merkle batch"})
		return
	}
	if batch == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"batch": batch})
}

type merkleVerifyRequest struct {
	Root      string      `json:"root" binding:"required"`
	EntryHash string      `json:"entryHash" binding:"required"`
	Proof     audit.Proof `json:"proof" binding:"required"`
}

func (h *Handler) HandleVerifyMerkle(c *gin.Context) {
	var req merkleVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	valid := h.service.VerifyMerkleProof(req.Root, req.EntryHash, &req.Proof)
	c.JSON(http.StatusOK, gin.H{"valid": valid})
}

func (h *Handler) HandleAnchorMerkleBatch(c *gin.Context) {
	if h.chainClient == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "chain anchoring is not configured (set ANCHOR_RPC_URL, ANCHOR_CONTRACT_ADDRESS, ANCHOR_PRIVATE_KEY)",
		})
		return
	}

	address, exists := c.Get("user_address")
	actor, ok := address.(string)
	if !exists || !ok || actor == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID := c.Param("batchId")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch ID is required"})
		return
	}

	batch, err := h.service.AnchorBatch(c.Request.Context(), actor, batchID, h.chainClient)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "anchor batch failed", "batchId", batchID, "error", err)
		// If we have a batch, return it for visibility (e.g., failed status + error).
		if batch != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to anchor batch", "batch": batch})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to anchor batch"})
		return
	}
	if batch == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"batch": batch})
}
