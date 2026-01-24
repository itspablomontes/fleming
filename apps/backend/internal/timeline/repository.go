package timeline

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetByPatientID(ctx context.Context, patientID string) ([]TimelineEvent, error)
	GetByID(ctx context.Context, id string) (*TimelineEvent, error)
	Create(ctx context.Context, event *TimelineEvent) error
	Update(ctx context.Context, event *TimelineEvent) error
	Delete(ctx context.Context, id string) error

	CreateEdge(ctx context.Context, edge *EventEdge) error
	DeleteEdge(ctx context.Context, id string) error
	GetRelatedEvents(ctx context.Context, eventID string, maxDepth int) ([]TimelineEvent, error)

	GetGraphData(ctx context.Context, patientID string) ([]TimelineEvent, []EventEdge, error)

	Transaction(ctx context.Context, fn func(repo Repository) error) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetByPatientID(ctx context.Context, patientID string) ([]TimelineEvent, error) {
	var events []TimelineEvent
	err := r.db.WithContext(ctx).
		Where("patient_id = ?", patientID).
		Preload("OutgoingEdges").
		Preload("IncomingEdges").
		Order("timestamp DESC").
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("get timeline events by patient %s: %w", patientID, err)
	}
	return events, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*TimelineEvent, error) {
	var event TimelineEvent
	err := r.db.WithContext(ctx).
		Preload("Files").
		First(&event, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("get timeline event %s: %w", id, err)
	}
	return &event, nil
}

func (r *GormRepository) Create(ctx context.Context, event *TimelineEvent) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return fmt.Errorf("create timeline event: %w", err)
	}
	return nil
}

func (r *GormRepository) Update(ctx context.Context, event *TimelineEvent) error {
	if err := r.db.WithContext(ctx).Save(event).Error; err != nil {
		return fmt.Errorf("update timeline event %s: %w", event.ID, err)
	}
	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&TimelineEvent{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete timeline event %s: %w", id, err)
	}
	return nil
}

func (r *GormRepository) CreateEdge(ctx context.Context, edge *EventEdge) error {
	if edge.FromEventID == edge.ToEventID {
		return fmt.Errorf("create edge: self-loops not allowed")
	}
	if err := r.db.WithContext(ctx).Create(edge).Error; err != nil {
		return fmt.Errorf("create event edge: %w", err)
	}
	return nil
}

func (r *GormRepository) DeleteEdge(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&EventEdge{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete event edge %s: %w", id, err)
	}
	return nil
}

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
