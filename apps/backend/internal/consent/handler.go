package consent

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	}
}

type ConsentRequestDTO struct {
	Grantor     string   `json:"grantor" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
	Reason      string   `json:"reason"`
	Duration    int      `json:"durationDays"` // Optional: how long access should last
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
	address, _ := c.Get("user_address")
	grantee := address.(string)

	grants, err := h.service.GetActiveGrants(c.Request.Context(), grantee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch active grants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"grants": grants})
}
