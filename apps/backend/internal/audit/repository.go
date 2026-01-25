package audit

import (
	"context"
	"fmt"

	protocol "github.com/itspablomontes/fleming/pkg/protocol/audit"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
	"gorm.io/gorm"
)

// Repository defines the interface for audit log persistence.
type Repository interface {
	Create(ctx context.Context, entry *AuditEntry) error
	GetLatest(ctx context.Context) (*AuditEntry, error)
	List(ctx context.Context, actor string, limit int) ([]AuditEntry, error)
	GetByResource(ctx context.Context, resourceID types.ID) ([]AuditEntry, error)
	GetByActor(ctx context.Context, actor types.WalletAddress) ([]AuditEntry, error)
	GetByID(ctx context.Context, id types.ID) (*AuditEntry, error)
	Query(ctx context.Context, filter protocol.QueryFilter) ([]AuditEntry, error)
	CreateBatch(ctx context.Context, batch *AuditBatch) error
	GetBatchByID(ctx context.Context, id string) (*AuditBatch, error)
	GetBatchByRoot(ctx context.Context, rootHash string) (*AuditBatch, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository creates a new GORM-based repository for the audit protocol.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, entry *AuditEntry) error {
	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return fmt.Errorf("create audit entry: %w", err)
	}
	return nil
}

func (r *gormRepository) GetLatest(ctx context.Context) (*AuditEntry, error) {
	var entry AuditEntry
	err := r.db.WithContext(ctx).Order("timestamp DESC, id DESC").Limit(1).First(&entry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get latest audit entry: %w", err)
	}
	return &entry, nil
}

func (r *gormRepository) List(ctx context.Context, actor string, limit int) ([]AuditEntry, error) {
	var entries []AuditEntry
	query := r.db.WithContext(ctx).Order("timestamp DESC, id DESC")
	if actor != "" {
		query = query.Where("actor = ?", actor)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("list audit entries: %w", err)
	}
	return entries, nil
}

func (r *gormRepository) GetByResource(ctx context.Context, resourceID types.ID) ([]AuditEntry, error) {
	var entries []AuditEntry
	if err := r.db.WithContext(ctx).
		Where("resource_id = ?", resourceID.String()).
		Order("timestamp DESC, id DESC").
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("get audit entries by resource: %w", err)
	}
	return entries, nil
}

func (r *gormRepository) GetByActor(ctx context.Context, actor types.WalletAddress) ([]AuditEntry, error) {
	var entries []AuditEntry
	if err := r.db.WithContext(ctx).
		Where("actor = ?", actor.String()).
		Order("timestamp DESC, id DESC").
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("get audit entries by actor: %w", err)
	}
	return entries, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id types.ID) (*AuditEntry, error) {
	var entry AuditEntry
	if err := r.db.WithContext(ctx).First(&entry, "id = ?", id.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get audit entry by id: %w", err)
	}
	return &entry, nil
}

func (r *gormRepository) Query(ctx context.Context, filter protocol.QueryFilter) ([]AuditEntry, error) {
	var entries []AuditEntry
	query := r.db.WithContext(ctx).Order("timestamp DESC, id DESC")

	if !filter.Actor.IsEmpty() {
		query = query.Where("actor = ?", filter.Actor.String())
	}
	if !filter.ResourceID.IsEmpty() {
		query = query.Where("resource_id = ?", filter.ResourceID.String())
	}
	if filter.ResourceType != "" {
		query = query.Where("resource_type = ?", filter.ResourceType)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.StartTime != nil && !filter.StartTime.IsZero() {
		query = query.Where("timestamp >= ?", filter.StartTime.Time)
	}
	if filter.EndTime != nil && !filter.EndTime.IsZero() {
		query = query.Where("timestamp <= ?", filter.EndTime.Time)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("query audit entries: %w", err)
	}
	return entries, nil
}

func (r *gormRepository) CreateBatch(ctx context.Context, batch *AuditBatch) error {
	if err := r.db.WithContext(ctx).Create(batch).Error; err != nil {
		return fmt.Errorf("create audit batch: %w", err)
	}
	return nil
}

func (r *gormRepository) GetBatchByID(ctx context.Context, id string) (*AuditBatch, error) {
	var batch AuditBatch
	if err := r.db.WithContext(ctx).First(&batch, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get audit batch by id: %w", err)
	}
	return &batch, nil
}

func (r *gormRepository) GetBatchByRoot(ctx context.Context, rootHash string) (*AuditBatch, error) {
	var batch AuditBatch
	if err := r.db.WithContext(ctx).First(&batch, "root_hash = ?", rootHash).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get audit batch by root: %w", err)
	}
	return &batch, nil
}
