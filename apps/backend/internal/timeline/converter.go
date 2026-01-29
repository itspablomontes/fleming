package timeline

import (
	"fmt"

	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/pkg/protocol/timeline"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// ToTimelineEvent converts a protocol Event to a GORM TimelineEvent entity.
func ToTimelineEvent(protocolEvent *timeline.Event) *TimelineEvent {
	if protocolEvent == nil {
		return nil
	}

	// Convert Codes
	codes := make(common.JSONCodes, len(protocolEvent.Codes))
	copy(codes, protocolEvent.Codes)

	// Convert Metadata
	metadata := make(common.JSONMap)
	for k, v := range protocolEvent.Metadata {
		metadata[k] = v
	}

	entity := &TimelineEvent{
		ID:          protocolEvent.ID.String(),
		PatientID:   protocolEvent.PatientID.String(),
		Type:        protocolEvent.Type,
		Title:       protocolEvent.Title,
		Description: protocolEvent.Description,
		Provider:    protocolEvent.Provider,
		Codes:       codes,
		Timestamp:   protocolEvent.Timestamp,
		Metadata:    metadata,
		CreatedAt:   protocolEvent.CreatedAt,
		UpdatedAt:   protocolEvent.UpdatedAt,
	}

	return entity
}

// ToProtocolEvent converts a GORM TimelineEvent entity to a protocol Event.
func ToProtocolEvent(entity *TimelineEvent) (*timeline.Event, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity is nil")
	}

	id, err := types.NewID(entity.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	patientID, err := types.NewWalletAddress(entity.PatientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID: %w", err)
	}

	// Convert Codes
	codes := make(types.Codes, len(entity.Codes))
	copy(codes, entity.Codes)

	// Convert Metadata
	metadata := types.NewMetadata()
	for k, v := range entity.Metadata {
		metadata = metadata.Set(k, v)
	}

	protocolEvent := &timeline.Event{
		ID:          id,
		PatientID:   patientID,
		Type:        entity.Type,
		Title:       entity.Title,
		Description: entity.Description,
		Provider:    entity.Provider,
		Codes:       codes,
		Timestamp:   entity.Timestamp,
		Metadata:    metadata,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}

	return protocolEvent, nil
}

// ToEventEdge converts a protocol Edge to a GORM EventEdge entity.
func ToEventEdge(protocolEdge *timeline.Edge) *EventEdge {
	if protocolEdge == nil {
		return nil
	}

	// Convert Metadata
	metadata := make(common.JSONMap)
	for k, v := range protocolEdge.Metadata {
		metadata[k] = v
	}

	entity := &EventEdge{
		ID:               protocolEdge.ID.String(),
		FromEventID:      protocolEdge.FromID.String(),
		ToEventID:        protocolEdge.ToID.String(),
		RelationshipType: protocolEdge.Type,
		Metadata:         metadata,
	}

	return entity
}

// ToProtocolEdge converts a GORM EventEdge entity to a protocol Edge.
func ToProtocolEdge(entity *EventEdge) (*timeline.Edge, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity is nil")
	}

	id, err := types.NewID(entity.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	fromID, err := types.NewID(entity.FromEventID)
	if err != nil {
		return nil, fmt.Errorf("invalid from ID: %w", err)
	}

	toID, err := types.NewID(entity.ToEventID)
	if err != nil {
		return nil, fmt.Errorf("invalid to ID: %w", err)
	}

	// Convert Metadata
	metadata := types.NewMetadata()
	for k, v := range entity.Metadata {
		metadata = metadata.Set(k, v)
	}

	protocolEdge := &timeline.Edge{
		ID:       id,
		FromID:   fromID,
		ToID:     toID,
		Type:     entity.RelationshipType,
		Metadata: metadata,
	}

	return protocolEdge, nil
}

// ToProtocolEvents converts a slice of TimelineEvent entities to protocol Events.
func ToProtocolEvents(entities []TimelineEvent) ([]*timeline.Event, error) {
	events := make([]*timeline.Event, 0, len(entities))
	for i := range entities {
		event, err := ToProtocolEvent(&entities[i])
		if err != nil {
			return nil, fmt.Errorf("convert event at index %d: %w", i, err)
		}
		events = append(events, event)
	}
	return events, nil
}

// ToProtocolEdges converts a slice of EventEdge entities to protocol Edges.
func ToProtocolEdges(entities []EventEdge) ([]*timeline.Edge, error) {
	edges := make([]*timeline.Edge, 0, len(entities))
	for i := range entities {
		edge, err := ToProtocolEdge(&entities[i])
		if err != nil {
			return nil, fmt.Errorf("convert edge at index %d: %w", i, err)
		}
		edges = append(edges, edge)
	}
	return edges, nil
}
