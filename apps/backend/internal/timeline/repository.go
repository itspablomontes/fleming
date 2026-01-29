package timeline

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/itspablomontes/fleming/pkg/protocol/timeline"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// Repository implements protocol GraphReader and GraphWriter interfaces.
// It converts between protocol types and GORM entities at the boundary.
type Repository interface {
	// Protocol interfaces
	timeline.GraphReader
	timeline.GraphWriter

	// File operations (backend-specific, not in protocol)
	CreateFile(ctx context.Context, file *EventFile) error
	GetFileByID(ctx context.Context, id string) (*EventFile, error)
	GetFilesByEventID(ctx context.Context, eventID string) ([]EventFile, error)
	UpsertFileAccess(ctx context.Context, access *EventFileAccess) error
	GetFileAccess(ctx context.Context, fileID string, grantee string) (*EventFileAccess, error)

	// Graph data for visualization (backend-specific)
	GetGraphData(ctx context.Context, patientID string) ([]TimelineEvent, []EventEdge, error)

	// Transaction support
	Transaction(ctx context.Context, fn func(repo Repository) error) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// GetEvent implements timeline.GraphReader.
func (r *GormRepository) GetEvent(ctx context.Context, id types.ID) (*timeline.Event, error) {
	var entity TimelineEvent
	err := r.db.WithContext(ctx).
		Preload("Files").
		First(&entity, "id = ?", id.String()).Error
	if err != nil {
		return nil, fmt.Errorf("get timeline event %s: %w", id, err)
	}
	return ToProtocolEvent(&entity)
}

// GetTimeline implements timeline.GraphReader.
func (r *GormRepository) GetTimeline(ctx context.Context, patientID types.WalletAddress) ([]timeline.Event, error) {
	var entities []TimelineEvent
	err := r.db.WithContext(ctx).
		Where("patient_id = ?", patientID.String()).
		Preload("OutgoingEdges").
		Preload("IncomingEdges").
		Order("timestamp DESC").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get timeline events by patient %s: %w", patientID, err)
	}

	events, err := ToProtocolEvents(entities)
	if err != nil {
		return nil, fmt.Errorf("convert events: %w", err)
	}

	// Convert []*timeline.Event to []timeline.Event
	result := make([]timeline.Event, len(events))
	for i, e := range events {
		result[i] = *e
	}
	return result, nil
}

// GetRelated implements timeline.GraphReader.
func (r *GormRepository) GetRelated(ctx context.Context, eventID types.ID, depth int) ([]timeline.Event, []timeline.Edge, error) {
	var entities []TimelineEvent
	query := `
		WITH RECURSIVE related_events AS (
			SELECT e.id, e.patient_id, e.type, e.title, e.description, e.provider, e.codes,
			       e.timestamp, e.blob_ref, e.is_encrypted, e.metadata, e.created_at, e.updated_at,
			       0 as depth, ARRAY[e.id] as path
			FROM timeline_events e
			WHERE e.id = ?

			UNION ALL

			SELECT e2.id, e2.patient_id, e2.type, e2.title, e2.description, e2.provider, e2.codes,
			       e2.timestamp, e2.blob_ref, e2.is_encrypted, e2.metadata, e2.created_at, e2.updated_at,
			       re.depth + 1, re.path || e2.id
			FROM related_events re
			JOIN event_edges ee ON (ee.from_event_id = re.id OR ee.to_event_id = re.id)
			JOIN timeline_events e2 ON (
				e2.id = CASE 
					WHEN ee.from_event_id = re.id THEN ee.to_event_id 
					ELSE ee.from_event_id 
				END
			)
			WHERE re.depth < ?
			  AND NOT e2.id = ANY(re.path)
		)
		SELECT DISTINCT id, patient_id, type, title, description, provider, codes,
		       timestamp, blob_ref, is_encrypted, metadata, created_at, updated_at
		FROM related_events
		ORDER BY timestamp DESC
	`

	if err := r.db.WithContext(ctx).Raw(query, eventID.String(), depth).Scan(&entities).Error; err != nil {
		return nil, nil, fmt.Errorf("query related events for %s: %w", eventID, err)
	}

	events, err := ToProtocolEvents(entities)
	if err != nil {
		return nil, nil, fmt.Errorf("convert events: %w", err)
	}

	// Get edges between related events
	if len(entities) == 0 {
		return []timeline.Event{}, []timeline.Edge{}, nil
	}

	eventIDs := make([]string, len(entities))
	for i, e := range entities {
		eventIDs[i] = e.ID
	}

	var edgeEntities []EventEdge
	err = r.db.WithContext(ctx).
		Where("from_event_id IN ? AND to_event_id IN ?", eventIDs, eventIDs).
		Find(&edgeEntities).Error
	if err != nil {
		return nil, nil, fmt.Errorf("query edges: %w", err)
	}

	edges, err := ToProtocolEdges(edgeEntities)
	if err != nil {
		return nil, nil, fmt.Errorf("convert edges: %w", err)
	}

	// Convert []*timeline.Event to []timeline.Event
	resultEvents := make([]timeline.Event, len(events))
	for i, e := range events {
		resultEvents[i] = *e
	}

	// Convert []*timeline.Edge to []timeline.Edge
	resultEdges := make([]timeline.Edge, len(edges))
	for i, e := range edges {
		resultEdges[i] = *e
	}

	return resultEvents, resultEdges, nil
}

// CreateEvent implements timeline.GraphWriter.
func (r *GormRepository) CreateEvent(ctx context.Context, event *timeline.Event) error {
	entity := ToTimelineEvent(event)
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("create timeline event: %w", err)
	}
	// Update event ID from generated entity ID
	event.ID, _ = types.NewID(entity.ID)
	return nil
}

// UpdateEvent implements timeline.GraphWriter.
func (r *GormRepository) UpdateEvent(ctx context.Context, event *timeline.Event) error {
	entity := ToTimelineEvent(event)
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return fmt.Errorf("update timeline event %s: %w", event.ID, err)
	}
	return nil
}

// DeleteEvent implements timeline.GraphWriter.
func (r *GormRepository) DeleteEvent(ctx context.Context, id types.ID) error {
	if err := r.db.WithContext(ctx).Delete(&TimelineEvent{}, "id = ?", id.String()).Error; err != nil {
		return fmt.Errorf("delete timeline event %s: %w", id, err)
	}
	return nil
}

// CreateEdge implements timeline.GraphWriter.
func (r *GormRepository) CreateEdge(ctx context.Context, edge *timeline.Edge) error {
	entity := ToEventEdge(edge)
	if entity.FromEventID == entity.ToEventID {
		return fmt.Errorf("create edge: self-loops not allowed")
	}
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("create event edge: %w", err)
	}
	// Update edge ID from generated entity ID
	edge.ID, _ = types.NewID(entity.ID)
	return nil
}

// DeleteEdge implements timeline.GraphWriter.
func (r *GormRepository) DeleteEdge(ctx context.Context, id types.ID) error {
	if err := r.db.WithContext(ctx).Delete(&EventEdge{}, "id = ?", id.String()).Error; err != nil {
		return fmt.Errorf("delete event edge %s: %w", id, err)
	}
	return nil
}

// GetByID is a convenience method that returns backend entity.
// Use GetEvent() for protocol-compliant access.
func (r *GormRepository) GetByID(ctx context.Context, id string) (*TimelineEvent, error) {
	idTyped, err := types.NewID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	event, err := r.GetEvent(ctx, idTyped)
	if err != nil {
		return nil, err
	}

	return ToTimelineEvent(event), nil
}

// GetByPatientID is a convenience method that returns backend entities.
// Use GetTimeline() for protocol-compliant access.
func (r *GormRepository) GetByPatientID(ctx context.Context, patientID string) ([]TimelineEvent, error) {
	addr, err := types.NewWalletAddress(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID: %w", err)
	}

	events, err := r.GetTimeline(ctx, addr)
	if err != nil {
		return nil, err
	}

	entities := make([]TimelineEvent, len(events))
	for i, e := range events {
		entities[i] = *ToTimelineEvent(&e)
	}
	return entities, nil
}

// Create is a convenience method that accepts backend entity.
// Use CreateEvent() for protocol-compliant access.
func (r *GormRepository) Create(ctx context.Context, event *TimelineEvent) error {
	protocolEvent, err := ToProtocolEvent(event)
	if err != nil {
		return fmt.Errorf("convert event: %w", err)
	}

	if err := r.CreateEvent(ctx, protocolEvent); err != nil {
		return err
	}

	// Update entity ID from protocol event
	event.ID = protocolEvent.ID.String()
	return nil
}

// Update is a convenience method that accepts backend entity.
// Use UpdateEvent() for protocol-compliant access.
func (r *GormRepository) Update(ctx context.Context, event *TimelineEvent) error {
	protocolEvent, err := ToProtocolEvent(event)
	if err != nil {
		return fmt.Errorf("convert event: %w", err)
	}

	return r.UpdateEvent(ctx, protocolEvent)
}

// Delete is a convenience method that accepts string ID.
// Use DeleteEvent() for protocol-compliant access.
func (r *GormRepository) Delete(ctx context.Context, id string) error {
	idTyped, err := types.NewID(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	return r.DeleteEvent(ctx, idTyped)
}

// CreateEdge is a convenience method that accepts backend entity.
// Use CreateEdge() from protocol interface for protocol-compliant access.
func (r *GormRepository) CreateEdgeLegacy(ctx context.Context, edge *EventEdge) error {
	protocolEdge, err := ToProtocolEdge(edge)
	if err != nil {
		return fmt.Errorf("convert edge: %w", err)
	}

	if err := r.CreateEdge(ctx, protocolEdge); err != nil {
		return err
	}

	// Update entity ID from protocol edge
	edge.ID = protocolEdge.ID.String()
	return nil
}

// DeleteEdge is a convenience method that accepts string ID.
// Use DeleteEdge() from protocol interface for protocol-compliant access.
func (r *GormRepository) DeleteEdgeLegacy(ctx context.Context, id string) error {
	idTyped, err := types.NewID(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	return r.DeleteEdge(ctx, idTyped)
}

// GetRelatedEvents is a convenience method that returns backend entities.
// Use GetRelated() for protocol-compliant access.
func (r *GormRepository) GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error) {
	var events []TimelineEvent
	query := `
		WITH RECURSIVE related_events AS (
			SELECT e.id, e.patient_id, e.type, e.title, e.description, e.provider, e.codes,
			       e.timestamp, e.blob_ref, e.is_encrypted, e.metadata, e.created_at, e.updated_at,
			       0 as depth, ARRAY[e.id] as path
			FROM timeline_events e
			WHERE e.id = ?

			UNION ALL

			SELECT e2.id, e2.patient_id, e2.type, e2.title, e2.description, e2.provider, e2.codes,
			       e2.timestamp, e2.blob_ref, e2.is_encrypted, e2.metadata, e2.created_at, e2.updated_at,
			       re.depth + 1, re.path || e2.id
			FROM related_events re
			JOIN event_edges ee ON (ee.from_event_id = re.id OR ee.to_event_id = re.id)
			JOIN timeline_events e2 ON (
				e2.id = CASE 
					WHEN ee.from_event_id = re.id THEN ee.to_event_id 
					ELSE ee.from_event_id 
				END
			)
			WHERE re.depth < ?
			  AND NOT e2.id = ANY(re.path)
		)
		SELECT DISTINCT id, patient_id, type, title, description, provider, codes,
		       timestamp, blob_ref, is_encrypted, metadata, created_at, updated_at
		FROM related_events
		ORDER BY timestamp DESC
	`

	if err := r.db.WithContext(ctx).Raw(query, eventID, maxDepth).Scan(&events).Error; err != nil {
		return nil, fmt.Errorf("query related events for %s: %w", eventID, err)
	}

	return events, nil
}

func (r *GormRepository) GetGraphData(ctx context.Context, patientID string) ([]TimelineEvent, []EventEdge, error) {
	var events []TimelineEvent
	err := r.db.WithContext(ctx).
		Where("patient_id = ?", patientID).
		Preload("Files").
		Order("timestamp DESC").
		Find(&events).Error
	if err != nil {
		return nil, nil, fmt.Errorf("query events for graph: %w", err)
	}

	if len(events) == 0 {
		return events, []EventEdge{}, nil
	}

	eventIDs := make([]string, len(events))
	for i, e := range events {
		eventIDs[i] = e.ID
	}

	var edges []EventEdge
	err = r.db.WithContext(ctx).
		Where("from_event_id IN ? AND to_event_id IN ?", eventIDs, eventIDs).
		Find(&edges).Error
	if err != nil {
		return nil, nil, fmt.Errorf("query edges for graph: %w", err)
	}

	return events, edges, nil
}

func (r *GormRepository) Transaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&GormRepository{db: tx})
	})
}

func (r *GormRepository) CreateFile(ctx context.Context, file *EventFile) error {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return fmt.Errorf("create event file: %w", err)
	}
	return nil
}

func (r *GormRepository) GetFileByID(ctx context.Context, id string) (*EventFile, error) {
	var file EventFile
	if err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get event file %s: %w", id, err)
	}
	return &file, nil
}

func (r *GormRepository) GetFilesByEventID(ctx context.Context, eventID string) ([]EventFile, error) {
	var files []EventFile
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&files).Error; err != nil {
		return nil, fmt.Errorf("get files for event %s: %w", eventID, err)
	}
	return files, nil
}

func (r *GormRepository) UpsertFileAccess(ctx context.Context, access *EventFileAccess) error {
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "file_id"}, {Name: "grantee"}},
		DoUpdates: clause.AssignmentColumns([]string{"wrapped_dek", "updated_at"}),
	}).Create(access).Error; err != nil {
		return fmt.Errorf("upsert file access: %w", err)
	}
	return nil
}

func (r *GormRepository) GetFileAccess(ctx context.Context, fileID string, grantee string) (*EventFileAccess, error) {
	var access EventFileAccess
	if err := r.db.WithContext(ctx).
		Where("file_id = ? AND grantee = ?", fileID, grantee).
		First(&access).Error; err != nil {
		return nil, fmt.Errorf("get file access for %s: %w", fileID, err)
	}
	return &access, nil
}
