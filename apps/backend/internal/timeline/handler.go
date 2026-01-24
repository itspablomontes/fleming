package timeline

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/pkg/protocol/timeline"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// HandleGetTimeline returns the patient's history, excluding superseded events.
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

// HandleGetEvent returns a single event by ID.
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

// AddEventRequest defines the payload for creating a new event.
type AddEventRequest struct {
	EventType   string         `json:"eventType" binding:"required"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	Provider    string         `json:"provider"`
	Date        string         `json:"date"`
	Codes       []types.Code   `json:"codes"`
	BlobRef     string         `json:"blobRef"`
	IsEncrypted bool           `json:"isEncrypted"`
	Metadata    common.JSONMap `json:"metadata"`
}

// HandleAddEvent creates a new timeline event from JSON payload.
func (h *Handler) HandleAddEvent(c *gin.Context) {
	patientID, exists := c.Get("user_address")
	address, ok := patientID.(string)
	if !exists || !ok || address == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AddEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	timestamp, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		timestamp = time.Now()
	}

	if req.Title == "" {
		req.Title = req.EventType + " Record"
	}

	event := &TimelineEvent{
		PatientID:   address,
		Type:        timeline.EventType(req.EventType),
		Title:       req.Title,
		Description: req.Description,
		Provider:    req.Provider,
		Codes:       common.JSONCodes(req.Codes),
		Timestamp:   timestamp,
		BlobRef:     req.BlobRef,
		IsEncrypted: req.IsEncrypted,
		Metadata:    req.Metadata,
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

// HandleCorrectEvent implements the "Edit" logic using the Append-Only flow.
func (h *Handler) HandleCorrectEvent(c *gin.Context) {
	patientID, exists := c.Get("user_address")
	address, ok := patientID.(string)
	if !exists || !ok || address == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	var req AddEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	timestamp, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		timestamp = time.Now()
	}

	event := &TimelineEvent{
		ID:          eventID,
		PatientID:   address,
		Type:        timeline.EventType(req.EventType),
		Title:       req.Title,
		Description: req.Description,
		Provider:    req.Provider,
		Codes:       common.JSONCodes(req.Codes),
		Timestamp:   timestamp,
		BlobRef:     req.BlobRef,
		IsEncrypted: req.IsEncrypted,
		Metadata:    req.Metadata,
	}

	if err := h.service.UpdateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to correct event: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"event":   event,
	})
}

// HandleDeleteEvent removes an event.
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

// LinkRequest defines the payload for linking two events.
type LinkRequest struct {
	ToEventID        string `json:"toEventId" binding:"required"`
	RelationshipType string `json:"relationshipType" binding:"required"`
}

// HandleLinkEvents creates a semantic edge between two events.
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
		timeline.RelationshipType(req.RelationshipType),
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

// HandleUnlinkEvents removes a semantic edge.
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

// HandleGetRelatedEvents returns events connected to the given ID up to a depth.
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

// HandleGetGraphData returns the raw node/edge list for visualizers.
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
