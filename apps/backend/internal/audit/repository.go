package audit

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Repository defines the interface for audit log persistence.
type Repository interface {
	Create(ctx context.Context, entry *AuditEntry) error
	GetLatest(ctx context.Context) (*AuditEntry, error)
	List(ctx context.Context, actor string, limit int) ([]AuditEntry, error)
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
	query := r.db.WithContext(ctx).Order("timestamp DESC")
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
