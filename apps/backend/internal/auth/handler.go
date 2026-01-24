package auth

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type ChallengeRequestDTO struct {
	Address string `json:"address" binding:"required"`
	Domain  string `json:"domain" binding:"required"`
	URI     string `json:"uri" binding:"required"`
	ChainID int    `json:"chainId" binding:"required"`
}

type ChallengeResponse struct {
	Message string `json:"message"`
}

func (h *Handler) HandleChallenge(c *gin.Context) {
	var req ChallengeRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	message, err := h.service.GenerateChallenge(c.Request.Context(), ChallengeRequest{
		Address: req.Address,
		Domain:  req.Domain,
		URI:     req.URI,
		ChainID: req.ChainID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate challenge"})
		return
	}

	c.JSON(http.StatusOK, ChallengeResponse{Message: message})
}

type LoginRequest struct {
	Address   string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}
type LoginResponse struct {
	Success bool `json:"success"`
}

func (h *Handler) HandleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := h.service.ValidateChallenge(c.Request.Context(), req.Address, req.Signature)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	secure := os.Getenv("ENV") == "production"
	c.SetCookie("auth_token", token, 3600*24, "/", "", secure, true)
	c.SetCookie("fleming_has_session", "true", 3600*24, "/", "", secure, false)

	c.JSON(http.StatusOK, LoginResponse{Success: true})
}

func (h *Handler) HandleLogout(c *gin.Context) {
	secure := os.Getenv("ENV") == "production"
	c.SetCookie("auth_token", "", -1, "/", "", secure, true)
	c.SetCookie("fleming_has_session", "", -1, "/", "", secure, false)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) HandleMe(c *gin.Context) {
	env := os.Getenv("ENV")
	overrideAddress := os.Getenv("DEV_OVERRIDE_WALLET_ADDRESS")

	var address string
	var exists bool

	if env == "dev" && overrideAddress != "" {
		address = overrideAddress
		exists = true
		slog.Debug("auth: HandleMe using dev override", "address", address)
	} else {
		val, ok := c.Get("user_address")
		if ok {
			address = val.(string)
			exists = true
		}
	}

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"status":  "authenticated",
	})
}
