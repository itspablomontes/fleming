package timeline

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
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
}

type GraphData struct {
	Events []TimelineEvent `json:"events"`
	Edges  []EventEdge     `json:"edges"`
}

type service struct {
	repo         Repository
	auditService audit.Service
}

func NewService(repo Repository, auditService audit.Service) Service {
	return &service{
		repo:         repo,
		auditService: auditService,
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
