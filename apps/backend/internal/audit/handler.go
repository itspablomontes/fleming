package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for audit logs.
type Handler struct {
	service Service
}

// NewHandler creates a new audit handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers audit endpoints.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	audit := rg.Group("/audit")
	{
		audit.GET("", h.HandleGetLogs)
		audit.GET("/verify", h.HandleVerify)
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
