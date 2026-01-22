package auth

import (
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

	c.JSON(http.StatusOK, LoginResponse{Success: true})
}

func (h *Handler) HandleMe(c *gin.Context) {
	address, exists := c.Get("user_address")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"status":  "authenticated",
	})
}
