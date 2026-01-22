package auth

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/challenge", h.HandleChallenge)
	rg.POST("/login", h.HandleLogin)
	rg.GET("/me", h.HandleMe)
}
