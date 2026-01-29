package timeline

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/apps/backend/internal/storage"
	protocol "github.com/itspablomontes/fleming/pkg/protocol/audit"
	"github.com/itspablomontes/fleming/pkg/protocol/timeline"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type Service interface {
	// Protocol-compliant methods (preferred)
	CreateEvent(ctx context.Context, event *timeline.Event) error
	GetEventByID(ctx context.Context, id types.ID) (*timeline.Event, error)
	GetTimelineForPatient(ctx context.Context, patientID types.WalletAddress) ([]timeline.Event, error)
	UpdateEventProtocol(ctx context.Context, event *timeline.Event) error
	DeleteEventByID(ctx context.Context, id types.ID) error
	LinkEventsProtocol(ctx context.Context, fromID, toID types.ID, relType timeline.RelationshipType) (*timeline.Edge, error)
	UnlinkEventsByID(ctx context.Context, edgeID types.ID) error

	// Legacy methods returning backend types (for backward compatibility with handlers)
	GetTimeline(ctx context.Context, patientID string) ([]TimelineEvent, error)
	GetEvent(ctx context.Context, id string) (*TimelineEvent, error)
	AddEvent(ctx context.Context, event *TimelineEvent) error
	UpdateEvent(ctx context.Context, event *TimelineEvent) error
	DeleteEvent(ctx context.Context, id string) error
	LinkEvents(ctx context.Context, fromID, toID string, relType timeline.RelationshipType) (*EventEdge, error)
	UnlinkEvents(ctx context.Context, edgeID string) error
	GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error)
	GetGraphData(ctx context.Context, patientID string) (*GraphData, error)

	UploadFile(ctx context.Context, eventID string, fileName string, contentType string, reader io.Reader, size int64, wrappedDEK []byte, metadata common.JSONMap) (*EventFile, error)
	GetFile(ctx context.Context, fileID string, actor string) (*EventFile, io.ReadCloser, error)

	StartMultipartUpload(ctx context.Context, eventID string, fileName string, contentType string) (string, string, error)
	UploadMultipartPart(ctx context.Context, objectName string, uploadID string, partNumber int, reader io.Reader, size int64) (string, error)
	CompleteMultipartUpload(ctx context.Context, eventID string, objectName string, uploadID string, parts []storage.Part, fileName string, contentType string, size int64, wrappedDEK []byte, metadata common.JSONMap) (*EventFile, error)

	GetFileKey(ctx context.Context, fileID string, actor string, patientID string) ([]byte, error)
	SaveFileAccess(ctx context.Context, fileID string, grantee string, wrappedDEK []byte) error
}

type service struct {
	repo         Repository
	auditService audit.Service
	storage      storage.Storage
}

func NewService(repo Repository, auditService audit.Service, storage storage.Storage) Service {
	return &service{
		repo:         repo,
		auditService: auditService,
		storage:      storage,
	}
}

// CreateEvent implements protocol-compliant event creation.
func (s *service) CreateEvent(ctx context.Context, event *timeline.Event) error {
	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	// Record action
	_ = s.auditService.Record(ctx, event.PatientID.String(), protocol.ActionCreate, protocol.ResourceEvent, event.ID.String(), nil)

	return nil
}

// GetEventByID implements protocol-compliant event retrieval.
func (s *service) GetEventByID(ctx context.Context, id types.ID) (*timeline.Event, error) {
	event, err := s.repo.GetEvent(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get event %s: %w", id, err)
	}
	return event, nil
}

// GetTimelineForPatient implements protocol-compliant timeline retrieval.
func (s *service) GetTimelineForPatient(ctx context.Context, patientID types.WalletAddress) ([]timeline.Event, error) {
	allEvents, err := s.repo.GetTimeline(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("get timeline for patient %s: %w", patientID, err)
	}

	// Filter replaced events and tombstones
	replacedIDs := make(map[types.ID]bool)
	for _, evt := range allEvents {
		// Check if this event is replaced by querying related events
		_, edges, err := s.repo.GetRelated(ctx, evt.ID, 1)
		if err == nil {
			for _, edge := range edges {
				if edge.Type == timeline.RelReplaces && edge.ToID == evt.ID {
					replacedIDs[evt.ID] = true
					break
				}
			}
		}
	}

	activeEvents := make([]timeline.Event, 0, len(allEvents))
	for _, evt := range allEvents {
		if !replacedIDs[evt.ID] && evt.Type != timeline.EventTombstone {
			activeEvents = append(activeEvents, evt)
		}
	}

	return activeEvents, nil
}

// Legacy methods for backward compatibility

// GetTimeline returns active events for a patient, filtering superseded ones.
func (s *service) GetTimeline(ctx context.Context, patientID string) ([]TimelineEvent, error) {
	addr, err := types.NewWalletAddress(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID: %w", err)
	}

	events, err := s.GetTimelineForPatient(ctx, addr)
	if err != nil {
		return nil, err
	}

	entities := make([]TimelineEvent, len(events))
	for i, e := range events {
		entities[i] = *ToTimelineEvent(&e)
	}
	return entities, nil
}

// GetEvent retrieves a specific event.
func (s *service) GetEvent(ctx context.Context, id string) (*TimelineEvent, error) {
	eventID, err := types.NewID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	event, err := s.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return ToTimelineEvent(event), nil
}

// AddEvent persists a new event (legacy method).
func (s *service) AddEvent(ctx context.Context, event *TimelineEvent) error {
	protocolEvent, err := ToProtocolEvent(event)
	if err != nil {
		return fmt.Errorf("convert event: %w", err)
	}

	return s.CreateEvent(ctx, protocolEvent)
}

// UpdateEventProtocol implements append-only correction using protocol types.
func (s *service) UpdateEventProtocol(ctx context.Context, event *timeline.Event) error {
	if event.ID.IsEmpty() {
		return fmt.Errorf("update event: id required")
	}

	originalID := event.ID

	err := s.repo.Transaction(ctx, func(repo Repository) error {
		if _, err := repo.GetEvent(ctx, originalID); err != nil {
			return fmt.Errorf("find original: %w", err)
		}

		// Create new event with cleared ID
		correction := *event
		correction.ID = types.ID("")
		correction.CreatedAt = time.Time{}
		correction.UpdatedAt = time.Time{}

		if err := repo.CreateEvent(ctx, &correction); err != nil {
			return fmt.Errorf("create correction: %w", err)
		}

		if correction.ID.IsEmpty() {
			return fmt.Errorf("empty id after creation")
		}

		// Create replacement edge
		edge, err := timeline.NewEdgeBuilder().
			WithFromID(correction.ID).
			WithToID(originalID).
			WithType(timeline.RelReplaces).
			Build()
		if err != nil {
			return fmt.Errorf("build edge: %w", err)
		}

		if err := repo.CreateEdge(ctx, edge); err != nil {
			return fmt.Errorf("link correction: %w", err)
		}

		// Update event ID
		event.ID = correction.ID
		return nil
	})
	if err != nil {
		return fmt.Errorf("correct event %s: %w", originalID, err)
	}

	// Record action
	_ = s.auditService.Record(ctx, event.PatientID.String(), protocol.ActionUpdate, protocol.ResourceEvent, event.ID.String(), nil)

	slog.InfoContext(ctx, "timeline event corrected", "original", originalID, "replacement", event.ID)
	return nil
}

// DeleteEventByID implements append-only deletion using protocol types.
func (s *service) DeleteEventByID(ctx context.Context, id types.ID) error {
	err := s.repo.Transaction(ctx, func(repo Repository) error {
		original, err := repo.GetEvent(ctx, id)
		if err != nil {
			return fmt.Errorf("find original: %w", err)
		}

		// Create tombstone event
		tombstone, err := timeline.NewEventBuilder().
			WithPatientID(original.PatientID).
			WithType(timeline.EventTombstone).
			WithTitle("Deleted Event").
			WithTimestamp(time.Now()).
			Build()
		if err != nil {
			return fmt.Errorf("build tombstone: %w", err)
		}

		if err := repo.CreateEvent(ctx, tombstone); err != nil {
			return fmt.Errorf("create tombstone: %w", err)
		}

		// Create replacement edge
		edge, err := timeline.NewEdgeBuilder().
			WithFromID(tombstone.ID).
			WithToID(id).
			WithType(timeline.RelReplaces).
			Build()
		if err != nil {
			return fmt.Errorf("build edge: %w", err)
		}

		if err := repo.CreateEdge(ctx, edge); err != nil {
			return fmt.Errorf("link tombstone: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("delete event %s: %w", id, err)
	}

	// Record action
	_ = s.auditService.Record(ctx, id.String(), protocol.ActionDelete, protocol.ResourceEvent, id.String(), nil)

	return nil
}

// LinkEventsProtocol implements protocol-compliant edge creation.
func (s *service) LinkEventsProtocol(ctx context.Context, fromID, toID types.ID, relType timeline.RelationshipType) (*timeline.Edge, error) {
	if fromID == toID {
		return nil, fmt.Errorf("link events: self-loop not allowed")
	}

	edge, err := timeline.NewEdgeBuilder().
		WithFromID(fromID).
		WithToID(toID).
		WithType(relType).
		Build()
	if err != nil {
		return nil, fmt.Errorf("build edge: %w", err)
	}

	if err := s.repo.CreateEdge(ctx, edge); err != nil {
		return nil, fmt.Errorf("link events: %w", err)
	}

	return edge, nil
}

// UnlinkEventsByID implements protocol-compliant edge deletion.
func (s *service) UnlinkEventsByID(ctx context.Context, edgeID types.ID) error {
	if err := s.repo.DeleteEdge(ctx, edgeID); err != nil {
		return fmt.Errorf("unlink events %s: %w", edgeID, err)
	}
	return nil
}

// Legacy methods for backward compatibility

// UpdateEvent implements append-only correction by creating a new version (legacy).
func (s *service) UpdateEvent(ctx context.Context, event *TimelineEvent) error {
	protocolEvent, err := ToProtocolEvent(event)
	if err != nil {
		return fmt.Errorf("convert event: %w", err)
	}

	if err := s.UpdateEventProtocol(ctx, protocolEvent); err != nil {
		return err
	}

	// Update entity with new ID
	event.ID = protocolEvent.ID.String()
	return nil
}

// DeleteEvent implements append-only deletion (legacy).
func (s *service) DeleteEvent(ctx context.Context, id string) error {
	eventID, err := types.NewID(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	return s.DeleteEventByID(ctx, eventID)
}

// LinkEvents connects two existing events (legacy).
func (s *service) LinkEvents(ctx context.Context, fromID, toID string, relType timeline.RelationshipType) (*EventEdge, error) {
	from, err := types.NewID(fromID)
	if err != nil {
		return nil, fmt.Errorf("invalid from ID: %w", err)
	}

	to, err := types.NewID(toID)
	if err != nil {
		return nil, fmt.Errorf("invalid to ID: %w", err)
	}

	edge, err := s.LinkEventsProtocol(ctx, from, to, relType)
	if err != nil {
		return nil, err
	}

	return ToEventEdge(edge), nil
}

// UnlinkEvents removes a connection edge (legacy).
func (s *service) UnlinkEvents(ctx context.Context, edgeID string) error {
	id, err := types.NewID(edgeID)
	if err != nil {
		return fmt.Errorf("invalid edge ID: %w", err)
	}

	return s.UnlinkEventsByID(ctx, id)
}

// GetRelatedEvents finds connected events by traversing the graph (legacy).
func (s *service) GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error) {
	if maxDepth < 1 {
		maxDepth = 2
	}
	if maxDepth > 5 {
		maxDepth = 5
	}

	id, err := types.NewID(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}

	events, _, err := s.repo.GetRelated(ctx, id, maxDepth)
	if err != nil {
		return nil, fmt.Errorf("get related for %s: %w", eventID, err)
	}

	entities := make([]TimelineEvent, len(events))
	for i, e := range events {
		entities[i] = *ToTimelineEvent(&e)
	}

	return entities, nil
}

// GetGraphData returns the adjacency list of nodes and edges.
func (s *service) GetGraphData(ctx context.Context, patientID string) (*GraphData, error) {
	events, edges, err := s.repo.GetGraphData(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("get graph data for %s: %w", patientID, err)
	}

	return &GraphData{
		Events: events,
		Edges:  edges,
	}, nil
}

func (s *service) UploadFile(ctx context.Context, eventID string, fileName string, contentType string, reader io.Reader, size int64, wrappedDEK []byte, metadata common.JSONMap) (*EventFile, error) {
	blobRef, err := s.storage.Put(ctx, "fleming-blobs", fileName, reader, size, contentType)
	if err != nil {
		return nil, fmt.Errorf("storage put: %w", err)
	}

	file := &EventFile{
		EventID:    eventID,
		BlobRef:    blobRef,
		FileName:   fileName,
		MimeType:   contentType,
		FileSize:   size,
		WrappedDEK: wrappedDEK,
		Metadata:   metadata,
	}

	if err := s.repo.CreateFile(ctx, file); err != nil {
		return nil, fmt.Errorf("repo create file: %w", err)
	}

	eventIDTyped, _ := types.NewID(eventID)
	if event, err := s.repo.GetEvent(ctx, eventIDTyped); err == nil && event != nil {
		auditMetadata := common.JSONMap{
			"eventId":   eventID,
			"fileName":  fileName,
			"fileSize":  size,
			"mimeType":  contentType,
			"isMultipart": false,
		}
		_ = s.auditService.Record(ctx, event.PatientID.String(), protocol.ActionUpload, protocol.ResourceFile, file.ID, auditMetadata)
	}

	return file, nil
}

func (s *service) GetFile(ctx context.Context, fileID string, actor string) (*EventFile, io.ReadCloser, error) {
	file, err := s.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("repo get file %s: %w", fileID, err)
	}

	reader, err := s.storage.Get(ctx, "fleming-blobs", file.BlobRef)
	if err != nil {
		return nil, nil, fmt.Errorf("storage get %s: %w", file.BlobRef, err)
	}

	if actor != "" {
		auditMetadata := common.JSONMap{
			"eventId":  file.EventID,
			"fileName": file.FileName,
			"fileSize": file.FileSize,
			"mimeType": file.MimeType,
		}
		_ = s.auditService.Record(ctx, actor, protocol.ActionDownload, protocol.ResourceFile, file.ID, auditMetadata)
	}

	return file, reader, nil
}

func (s *service) StartMultipartUpload(ctx context.Context, eventID string, fileName string, contentType string) (string, string, error) {
	objectName := fmt.Sprintf("%s/%s", eventID, fileName)
	uploadID, err := s.storage.CreateMultipartUpload(ctx, "fleming-blobs", objectName, contentType)
	if err != nil {
		return "", "", err
	}
	return uploadID, objectName, nil
}

func (s *service) UploadMultipartPart(ctx context.Context, objectName string, uploadID string, partNumber int, reader io.Reader, size int64) (string, error) {
	return s.storage.UploadPart(ctx, "fleming-blobs", objectName, uploadID, partNumber, reader, size)
}

func (s *service) CompleteMultipartUpload(
	ctx context.Context,
	eventID string,
	objectName string,
	uploadID string,
	parts []storage.Part,
	fileName string,
	contentType string,
	size int64,
	wrappedDEK []byte,
	metadata common.JSONMap,
) (*EventFile, error) {
	blobRef, err := s.storage.CompleteMultipartUpload(ctx, "fleming-blobs", objectName, uploadID, parts)
	if err != nil {
		return nil, err
	}

	file := &EventFile{
		EventID:    eventID,
		BlobRef:    blobRef,
		FileName:   fileName,
		MimeType:   contentType,
		FileSize:   size,
		WrappedDEK: wrappedDEK,
		Metadata:   metadata,
	}

	if err := s.repo.CreateFile(ctx, file); err != nil {
		return nil, fmt.Errorf("repo create file: %w", err)
	}

	eventIDTyped, _ := types.NewID(eventID)
	if event, err := s.repo.GetEvent(ctx, eventIDTyped); err == nil && event != nil {
		auditMetadata := common.JSONMap{
			"eventId":     eventID,
			"fileName":    fileName,
			"fileSize":    size,
			"mimeType":    contentType,
			"isMultipart": true,
		}
		_ = s.auditService.Record(ctx, event.PatientID.String(), protocol.ActionUpload, protocol.ResourceFile, file.ID, auditMetadata)
	}

	return file, nil
}

func (s *service) GetFileKey(ctx context.Context, fileID string, actor string, patientID string) ([]byte, error) {
	file, err := s.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	if actor == patientID {
		return file.WrappedDEK, nil
	}

	access, err := s.repo.GetFileAccess(ctx, fileID, actor)
	if err != nil {
		return nil, err
	}
	return access.WrappedDEK, nil
}

func (s *service) SaveFileAccess(ctx context.Context, fileID string, grantee string, wrappedDEK []byte) error {
	access := &EventFileAccess{
		FileID:    fileID,
		Grantee:   grantee,
		WrappedDEK: wrappedDEK,
	}
	if err := s.repo.UpsertFileAccess(ctx, access); err != nil {
		return err
	}

	file, err := s.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return err
	}
	eventIDTyped, _ := types.NewID(file.EventID)
	if event, err := s.repo.GetEvent(ctx, eventIDTyped); err == nil && event != nil {
		auditMetadata := common.JSONMap{
			"eventId":  file.EventID,
			"fileName": file.FileName,
			"grantee":  grantee,
		}
		_ = s.auditService.Record(ctx, event.PatientID.String(), protocol.ActionShare, protocol.ResourceFile, fileID, auditMetadata)
	}

	return nil
}
