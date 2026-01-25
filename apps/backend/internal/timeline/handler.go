package timeline

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"io"
	"log/slog"

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

// HandleAddEvent creates a new timeline event from Multipart Form Data.
func (h *Handler) HandleAddEvent(c *gin.Context) {
	patientID, exists := c.Get("user_address")
	address, ok := patientID.(string)
	if !exists || !ok || address == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
		return
	}

	form := c.Request.PostForm

	dateStr := form.Get("date")
	timestamp, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		timestamp = time.Now()
	}

	eventType := form.Get("eventType")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "eventType is required"})
		return
	}

	title := form.Get("title")
	if title == "" {
		title = eventType + " Record"
	}

	isEncrypted := form.Get("isEncrypted") == "true"

	event := &TimelineEvent{
		PatientID:   address,
		Type:        timeline.EventType(eventType),
		Title:       title,
		Description: form.Get("description"),
		Provider:    form.Get("provider"),
		Timestamp:   timestamp,
		BlobRef:     form.Get("blobRef"),
		IsEncrypted: isEncrypted,
	}

	metadataStr := form.Get("metadata")
	if metadataStr != "" {
		var meta common.JSONMap
		if err := json.Unmarshal([]byte(metadataStr), &meta); err == nil {
			event.Metadata = meta
		}
	}

	if err := h.service.AddEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save event"})
		return
	}

	// Handle File Upload if present
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		defer file.Close()
		wrappedKeyStr := form.Get("wrappedKey")
		wrappedKey, _ := common.HexToBytes(wrappedKeyStr) // Assume helper exists or add it

		_, err = h.service.UploadFile(
			c.Request.Context(),
			event.ID,
			header.Filename,
			header.Header.Get("Content-Type"),
			file,
			header.Size,
			wrappedKey,
		)
		if err != nil {
			slog.Warn("Failed to upload attached file", "error", err, "eventId", event.ID)
			// We don't fail the whole request because the event itself was created
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"event":   event,
	})
}

// HandleDownloadFile serves a file's ciphertext blob.
func (h *Handler) HandleDownloadFile(c *gin.Context) {
	fileID := c.Param("fileId")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file ID is required"})
		return
	}

	file, reader, err := h.service.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", "attachment; filename="+file.FileName)
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.FileSize, 10))

	if _, err := io.Copy(c.Writer, reader); err != nil {
		slog.Error("failed to pipe file content", "error", err)
	}
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

	const maxMultipartMemory = 32 << 20
	if err := c.Request.ParseMultipartForm(maxMultipartMemory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
		return
	}

	form := c.Request.PostForm

	dateStr := form.Get("date")
	timestamp, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		timestamp = time.Now()
	}

	eventType := form.Get("eventType")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "eventType is required"})
		return
	}

	title := form.Get("title")
	if title == "" {
		title = eventType + " Record"
	}

	isEncrypted := form.Get("isEncrypted") == "true"

	event := &TimelineEvent{
		ID:          eventID,
		PatientID:   address,
		Type:        timeline.EventType(eventType),
		Title:       title,
		Description: form.Get("description"),
		Provider:    form.Get("provider"),
		Timestamp:   timestamp,
		BlobRef:     form.Get("blobRef"),
		IsEncrypted: isEncrypted,
	}

	metadataStr := form.Get("metadata")
	if metadataStr != "" {
		var meta common.JSONMap
		if err := json.Unmarshal([]byte(metadataStr), &meta); err == nil {
			event.Metadata = meta
		}
	}

	// Parse codes if present
	// Note: In AddEvent codes were not parsed, but here we might want them.
	// For now keeping consistent with simple form fields.
	// If codes are needed, they would come as a JSON string too.

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
