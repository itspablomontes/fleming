package timeline

import (
	"context"
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/timeline"
)

type MockRepo struct {
	events []TimelineEvent
	edges  []EventEdge
}

func (m *MockRepo) GetByPatientID(ctx context.Context, patientID string) ([]TimelineEvent, error) {
	var result []TimelineEvent
	for _, e := range m.events {
		if e.PatientID == patientID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *MockRepo) GetByID(ctx context.Context, id string) (*TimelineEvent, error) {
	for _, e := range m.events {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, nil
}

func (m *MockRepo) Create(ctx context.Context, event *TimelineEvent) error {
	m.events = append(m.events, *event)
	return nil
}

func (m *MockRepo) Update(ctx context.Context, event *TimelineEvent) error {
	for i, e := range m.events {
		if e.ID == event.ID {
			m.events[i] = *event
			return nil
		}
	}
	return nil
}

func (m *MockRepo) Delete(ctx context.Context, id string) error {
	for i, e := range m.events {
		if e.ID == id {
			m.events = append(m.events[:i], m.events[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockRepo) CreateEdge(ctx context.Context, edge *EventEdge) error {
	m.edges = append(m.edges, *edge)
	return nil
}

func (m *MockRepo) DeleteEdge(ctx context.Context, id string) error {
	for i, e := range m.edges {
		if e.ID == id {
			m.edges = append(m.edges[:i], m.edges[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockRepo) GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error) {
	// Simplified mock: just return all events for testing
	return m.events, nil
}

func (m *MockRepo) GetGraphData(ctx context.Context, patientID string) ([]TimelineEvent, []EventEdge, error) {
	var events []TimelineEvent
	for _, e := range m.events {
		if e.PatientID == patientID {
			events = append(events, e)
		}
	}
	return events, m.edges, nil
}

func TestService_AddEvent(t *testing.T) {
	repo := &MockRepo{}
	svc := NewService(repo)

	tests := []struct {
		name    string
		event   *TimelineEvent
		wantErr bool
	}{
		{
			name: "valid event",
			event: &TimelineEvent{
				PatientID: "0x123",
				Type:      timeline.EventLabResult,
				Title:     "Blood Test",
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing patient ID",
			event: &TimelineEvent{
				Type:      timeline.EventLabResult,
				Title:     "Blood Test",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.AddEvent(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetTimeline(t *testing.T) {
	repo := &MockRepo{
		events: []TimelineEvent{
			{PatientID: "0x123", Title: "Event 1"},
			{PatientID: "0x456", Title: "Event 2"},
		},
	}
	svc := NewService(repo)

	tests := []struct {
		name      string
		patientID string
		wantCount int
		wantErr   bool
	}{
		{"existing patient", "0x123", 1, false},
		{"non-existing patient", "0x999", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.GetTimeline(context.Background(), tt.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("GetTimeline() count = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}
