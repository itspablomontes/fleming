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
)

type Service interface {
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
	GetFile(ctx context.Context, fileID string) (*EventFile, io.ReadCloser, error)

	StartMultipartUpload(ctx context.Context, eventID string, fileName string, contentType string) (string, string, error)
	UploadMultipartPart(ctx context.Context, objectName string, uploadID string, partNumber int, reader io.Reader, size int64) (string, error)
	CompleteMultipartUpload(ctx context.Context, eventID string, objectName string, uploadID string, parts []storage.Part, fileName string, contentType string, size int64, wrappedDEK []byte, metadata common.JSONMap) (*EventFile, error)

	GetFileKey(ctx context.Context, fileID string, actor string, patientID string) ([]byte, error)
	SaveFileAccess(ctx context.Context, fileID string, grantee string, wrappedDEK []byte) error
}

type GraphData struct {
	Events []TimelineEvent `json:"events"`
	Edges  []EventEdge     `json:"edges"`
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

// GetTimeline returns active events for a patient, filtering superseded ones.
func (s *service) GetTimeline(ctx context.Context, patientID string) ([]TimelineEvent, error) {
	allEvents, err := s.repo.GetByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("get timeline for patient %s: %w", patientID, err)
	}

	replacedIDs := make(map[string]bool)
	for _, evt := range allEvents {
		for _, edge := range evt.IncomingEdges {
			if edge.RelationshipType == timeline.RelReplaces {
				replacedIDs[evt.ID] = true
			}
		}
	}

	activeEvents := make([]TimelineEvent, 0, len(allEvents))
	for _, evt := range allEvents {
		// Filter replaced events and tombstones.
		if !replacedIDs[evt.ID] && evt.Type != timeline.EventTombstone {
			activeEvents = append(activeEvents, evt)
		}
	}

	return activeEvents, nil
}

// GetEvent retrieves a specific event.
func (s *service) GetEvent(ctx context.Context, id string) (*TimelineEvent, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get event %s: %w", id, err)
	}
	return event, nil
}

// AddEvent persists a new event.
func (s *service) AddEvent(ctx context.Context, event *TimelineEvent) error {
	if event.PatientID == "" {
		return fmt.Errorf("add event: patient id required")
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return fmt.Errorf("add event: %w", err)
	}

	// Record action
	_ = s.auditService.Record(ctx, event.PatientID, protocol.ActionCreate, protocol.ResourceEvent, event.ID, nil)

	return nil
}

// UpdateEvent implements append-only correction by creating a new version.
func (s *service) UpdateEvent(ctx context.Context, event *TimelineEvent) error {
	if event.ID == "" {
		return fmt.Errorf("update event: id required")
	}

	originalID := event.ID

	err := s.repo.Transaction(ctx, func(repo Repository) error {
		if _, err := repo.GetByID(ctx, originalID); err != nil {
			return fmt.Errorf("find original: %w", err)
		}

		event.ID = ""
		event.CreatedAt = time.Time{}
		event.UpdatedAt = time.Time{}

		if err := repo.Create(ctx, event); err != nil {
			return fmt.Errorf("create correction: %w", err)
		}

		if event.ID == "" {
			return fmt.Errorf("empty id after creation")
		}

		edge := &EventEdge{
			FromEventID:      event.ID,
			ToEventID:        originalID,
			RelationshipType: timeline.RelReplaces,
		}

		if err := repo.CreateEdge(ctx, edge); err != nil {
			return fmt.Errorf("link correction: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("correct event %s: %w", originalID, err)
	}

	// Record action
	_ = s.auditService.Record(ctx, event.PatientID, protocol.ActionUpdate, protocol.ResourceEvent, event.ID, nil)

	slog.InfoContext(ctx, "timeline event corrected", "original", originalID, "replacement", event.ID)
	return nil
}

// DeleteEvent implements append-only deletion by replacing the target with a Tombstone.
func (s *service) DeleteEvent(ctx context.Context, id string) error {
	err := s.repo.Transaction(ctx, func(repo Repository) error {
		original, err := repo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("find original: %w", err)
		}

		tombstone := &TimelineEvent{
			PatientID: original.PatientID,
			Type:      timeline.EventTombstone,
			Title:     "Deleted Event",
			Timestamp: time.Now(),
		}

		if err := repo.Create(ctx, tombstone); err != nil {
			return fmt.Errorf("create tombstone: %w", err)
		}

		if tombstone.ID == "" {
			return fmt.Errorf("empty id after tombstone creation")
		}

		edge := &EventEdge{
			FromEventID:      tombstone.ID,
			ToEventID:        id,
			RelationshipType: timeline.RelReplaces,
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
	_ = s.auditService.Record(ctx, id, protocol.ActionDelete, protocol.ResourceEvent, id, nil)

	return nil
}

// LinkEvents connects two existing events.
func (s *service) LinkEvents(ctx context.Context, fromID, toID string, relType timeline.RelationshipType) (*EventEdge, error) {
	if fromID == toID {
		return nil, fmt.Errorf("link events: self-loop not allowed")
	}

	edge := &EventEdge{
		FromEventID:      fromID,
		ToEventID:        toID,
		RelationshipType: relType,
	}

	if err := s.repo.CreateEdge(ctx, edge); err != nil {
		return nil, fmt.Errorf("link events: %w", err)
	}

	return edge, nil
}

// UnlinkEvents removes a connection edge.
func (s *service) UnlinkEvents(ctx context.Context, edgeID string) error {
	if err := s.repo.DeleteEdge(ctx, edgeID); err != nil {
		return fmt.Errorf("unlink events %s: %w", edgeID, err)
	}
	return nil
}

// GetRelatedEvents finds connected events by traversing the graph.
func (s *service) GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error) {
	if maxDepth < 1 {
		maxDepth = 2
	}
	if maxDepth > 5 {
		maxDepth = 5
	}

	events, err := s.repo.GetRelatedEvents(ctx, eventID, maxDepth)
	if err != nil {
		return nil, fmt.Errorf("get related for %s: %w", eventID, err)
	}

	return events, nil
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

	return file, nil
}

func (s *service) GetFile(ctx context.Context, fileID string) (*EventFile, io.ReadCloser, error) {
	file, err := s.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("repo get file %s: %w", fileID, err)
	}

	reader, err := s.storage.Get(ctx, "fleming-blobs", file.BlobRef)
	if err != nil {
		return nil, nil, fmt.Errorf("storage get %s: %w", file.BlobRef, err)
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
	return s.repo.UpsertFileAccess(ctx, access)
}
