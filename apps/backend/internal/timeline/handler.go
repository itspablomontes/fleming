package timeline

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleGetTimeline(c *gin.Context) {
	patientID, exists := c.Get("user_address")
	address, ok := patientID.(string)
	if !exists || !ok || address == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing or invalid user address"})
		return
	}

	events, err := h.service.GetTimeline(c.Request.Context(), address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch timeline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

func (h *Handler) HandleGetEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	event, err := h.service.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *Handler) HandleAddEvent(c *gin.Context) {
	patientID, exists := c.Get("user_address")
	address, ok := patientID.(string)
	if !exists || !ok || address == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventType := c.PostForm("eventType")
	title := c.PostForm("title")
	description := c.PostForm("description")
	provider := c.PostForm("provider")
	dateStr := c.PostForm("date")
	code := c.PostForm("code")
	codingSystem := c.PostForm("codingSystem")

	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "eventType is required"})
		return
	}

	if title == "" {
		title = eventType + " Uploaded"
	}

	timestamp, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		timestamp = time.Now()
	}

	event := &TimelineEvent{
		PatientID:    address,
		Type:         TimelineEventType(eventType),
		Title:        title,
		Description:  description,
		Provider:     provider,
		Code:         code,
		CodingSystem: codingSystem,
		Timestamp:    timestamp,
		IsEncrypted:  true,
	}

	if err := h.service.AddEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"event":   event,
	})
}

func (h *Handler) HandleDeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	if err := h.service.DeleteEvent(c.Request.Context(), eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

type LinkRequest struct {
	ToEventID        string `json:"toEventId" binding:"required"`
	RelationshipType string `json:"relationshipType" binding:"required"`
}

func (h *Handler) HandleLinkEvents(c *gin.Context) {
	fromEventID := c.Param("id")
	if fromEventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	var req LinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	edge, err := h.service.LinkEvents(
		c.Request.Context(),
		fromEventID,
		req.ToEventID,
		RelationshipType(req.RelationshipType),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"edge":    edge,
	})
}

func (h *Handler) HandleUnlinkEvents(c *gin.Context) {
	edgeID := c.Param("edgeId")
	if edgeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "edge ID is required"})
		return
	}

	if err := h.service.UnlinkEvents(c.Request.Context(), edgeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete edge"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) HandleGetRelatedEvents(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	depthStr := c.DefaultQuery("depth", "2")
	depth, err := strconv.Atoi(depthStr)
	if err != nil {
		depth = 2
	}

	events, err := h.service.GetRelatedEvents(c.Request.Context(), eventID, depth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get related events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *Handler) HandleGetGraphData(c *gin.Context) {
	patientID, exists := c.Get("user_address")
	address, ok := patientID.(string)
	if !exists || !ok || address == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	graphData, err := h.service.GetGraphData(c.Request.Context(), address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get graph data"})
		return
	}

	c.JSON(http.StatusOK, graphData)
}
