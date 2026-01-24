package consent

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Repository defines the interface for consent grant persistence.
type Repository interface {
	Create(ctx context.Context, grant *ConsentGrant) error
	GetByID(ctx context.Context, id string) (*ConsentGrant, error)
	GetByGrantee(ctx context.Context, grantee string) ([]ConsentGrant, error)
	GetByGrantor(ctx context.Context, grantor string) ([]ConsentGrant, error)
	Update(ctx context.Context, grant *ConsentGrant) error
	FindLatest(ctx context.Context, grantor, grantee string) (*ConsentGrant, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository creates a new GORM repository for consent.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, grant *ConsentGrant) error {
	if err := r.db.WithContext(ctx).Create(grant).Error; err != nil {
		return fmt.Errorf("create consent grant: %w", err)
	}
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id string) (*ConsentGrant, error) {
	var grant ConsentGrant
	if err := r.db.WithContext(ctx).First(&grant, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get consent grant %s: %w", id, err)
	}
	return &grant, nil
}

func (r *gormRepository) GetByGrantee(ctx context.Context, grantee string) ([]ConsentGrant, error) {
	var grants []ConsentGrant
	if err := r.db.WithContext(ctx).Where("grantee = ?", grantee).Find(&grants).Error; err != nil {
		return nil, fmt.Errorf("list grants for grantee %s: %w", grantee, err)
	}
	return grants, nil
}

func (r *gormRepository) GetByGrantor(ctx context.Context, grantor string) ([]ConsentGrant, error) {
	var grants []ConsentGrant
	if err := r.db.WithContext(ctx).Where("grantor = ?", grantor).Find(&grants).Error; err != nil {
		return nil, fmt.Errorf("list grants from grantor %s: %w", grantor, err)
	}
	return grants, nil
}

func (r *gormRepository) Update(ctx context.Context, grant *ConsentGrant) error {
	if err := r.db.WithContext(ctx).Save(grant).Error; err != nil {
		return fmt.Errorf("update consent grant: %w", err)
	}
	return nil
}

func (r *gormRepository) FindLatest(ctx context.Context, grantor, grantee string) (*ConsentGrant, error) {
	var grant ConsentGrant
	err := r.db.WithContext(ctx).
		Where("grantor = ? AND grantee = ?", grantor, grantee).
		Order("created_at DESC").
		Limit(1).
		First(&grant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &grant, nil
}
