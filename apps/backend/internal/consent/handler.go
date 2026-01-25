package consent

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/itspablomontes/fleming/pkg/protocol/consent"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	consent := rg.Group("/consent")
	{
		consent.POST("/request", h.HandleRequest)
		consent.POST("/:id/approve", h.HandleApprove)
		consent.POST("/:id/deny", h.HandleDeny)
		consent.POST("/:id/revoke", h.HandleRevoke)
		consent.GET("/active", h.HandleGetActive)
		consent.GET("/grants", h.HandleGetMyGrants)
		consent.GET("/:id", h.HandleGetByID)
	}
}

type ConsentRequestDTO struct {
	Grantor     string   `json:"grantor" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
	Reason      string   `json:"reason"`
	Duration    int      `json:"durationDays"` // Optional: how long access should last
}

func getUserAddress(c *gin.Context) (string, bool) {
	address, ok := c.Get("user_address")
	if !ok {
		return "", false
	}
	value, ok := address.(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}

func filterGrantsByState(
	grants []ConsentGrant,
	states map[consent.State]struct{},
) []ConsentGrant {
	if len(states) == 0 {
		return grants
	}
	filtered := make([]ConsentGrant, 0, len(grants))
	for _, grant := range grants {
		if _, ok := states[grant.State]; ok {
			filtered = append(filtered, grant)
		}
	}
	return filtered
}

func (h *Handler) HandleRequest(c *gin.Context) {
	address, _ := c.Get("user_address")
	grantee := address.(string)

	var req ConsentRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var expiresAt time.Time
	if req.Duration > 0 {
		expiresAt = time.Now().AddDate(0, 0, req.Duration)
	}

	grant, err := h.service.RequestConsent(c.Request.Context(), req.Grantor, grantee, req.Reason, req.Permissions, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to request consent"})
		return
	}

	c.JSON(http.StatusCreated, grant)
}

func (h *Handler) HandleApprove(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.ApproveConsent(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) HandleDeny(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DenyConsent(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) HandleRevoke(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.RevokeConsent(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) HandleGetActive(c *gin.Context) {
	grantee, ok := getUserAddress(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	grants, err := h.service.GetActiveGrants(c.Request.Context(), grantee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch active grants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"grants": grants})
}

func (h *Handler) HandleGetMyGrants(c *gin.Context) {
	grantor, ok := getUserAddress(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	grants, err := h.service.GetGrantsByGrantor(c.Request.Context(), grantor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch consent grants"})
		return
	}

	stateFilters := c.QueryArray("state")
	if len(stateFilters) > 0 {
		states := make(map[consent.State]struct{}, len(stateFilters))
		for _, value := range stateFilters {
			state := consent.State(value)
			if !state.IsValid() {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state filter"})
				return
			}
			states[state] = struct{}{}
		}
		grants = filterGrantsByState(grants, states)
	}

	c.JSON(http.StatusOK, gin.H{"grants": grants})
}

func (h *Handler) HandleGetByID(c *gin.Context) {
	grantID := c.Param("id")
	if grantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing consent id"})
		return
	}

	address, ok := getUserAddress(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	grant, err := h.service.GetGrantByID(c.Request.Context(), grantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "consent not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch consent"})
		return
	}

	if grant.Grantor != address && grant.Grantee != address {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.JSON(http.StatusOK, grant)
}
