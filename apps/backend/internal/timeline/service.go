package timeline

import (
	"context"
	"fmt"
)

type Service interface {
	GetTimeline(ctx context.Context, patientID string) ([]TimelineEvent, error)
	GetEvent(ctx context.Context, id string) (*TimelineEvent, error)
	AddEvent(ctx context.Context, event *TimelineEvent) error
	UpdateEvent(ctx context.Context, event *TimelineEvent) error
	DeleteEvent(ctx context.Context, id string) error

	LinkEvents(ctx context.Context, fromID, toID string, relType RelationshipType) (*EventEdge, error)
	UnlinkEvents(ctx context.Context, edgeID string) error
	GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error)
	GetGraphData(ctx context.Context, patientID string) (*GraphData, error)
}

type GraphData struct {
	Events []TimelineEvent `json:"events"`
	Edges  []EventEdge     `json:"edges"`
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTimeline(ctx context.Context, patientID string) ([]TimelineEvent, error) {
	events, err := s.repo.GetByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get timeline: %w", err)
	}
	return events, nil
}

func (s *service) GetEvent(ctx context.Context, id string) (*TimelineEvent, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get event: %w", err)
	}
	return event, nil
}

func (s *service) AddEvent(ctx context.Context, event *TimelineEvent) error {
	if event.PatientID == "" {
		return fmt.Errorf("service: patient ID is required")
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return fmt.Errorf("service: failed to add event: %w", err)
	}

	return nil
}

func (s *service) UpdateEvent(ctx context.Context, event *TimelineEvent) error {
	if event.ID == "" {
		return fmt.Errorf("service: event ID is required")
	}

	if err := s.repo.Update(ctx, event); err != nil {
		return fmt.Errorf("service: failed to update event: %w", err)
	}

	return nil
}

func (s *service) DeleteEvent(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service: failed to delete event: %w", err)
	}
	return nil
}

func (s *service) LinkEvents(ctx context.Context, fromID, toID string, relType RelationshipType) (*EventEdge, error) {
	if fromID == toID {
		return nil, fmt.Errorf("service: cannot link event to itself")
	}

	edge := &EventEdge{
		FromEventID:      fromID,
		ToEventID:        toID,
		RelationshipType: relType,
	}

	if err := s.repo.CreateEdge(ctx, edge); err != nil {
		return nil, fmt.Errorf("service: failed to create edge: %w", err)
	}

	return edge, nil
}

func (s *service) UnlinkEvents(ctx context.Context, edgeID string) error {
	if err := s.repo.DeleteEdge(ctx, edgeID); err != nil {
		return fmt.Errorf("service: failed to delete edge: %w", err)
	}
	return nil
}

func (s *service) GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error) {
	if maxDepth < 1 {
		maxDepth = 2 // Default to 2 hops as per approved plan
	}
	if maxDepth > 5 {
		maxDepth = 5 // Safety limit
	}

	events, err := s.repo.GetRelatedEvents(ctx, eventID, maxDepth)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get related events: %w", err)
	}

	return events, nil
}

func (s *service) GetGraphData(ctx context.Context, patientID string) (*GraphData, error) {
	events, edges, err := s.repo.GetGraphData(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get graph data: %w", err)
	}

	return &GraphData{
		Events: events,
		Edges:  edges,
	}, nil
}
