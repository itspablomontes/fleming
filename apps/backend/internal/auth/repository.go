package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Repository interface {
	SaveChallenge(ctx context.Context, challenge *Challenge) error
	FindChallenge(ctx context.Context, address string) (*Challenge, error)
	DeleteChallenge(ctx context.Context, address string) error
	DeleteExpiredChallenges(ctx context.Context) (int64, error)

	SaveUser(ctx context.Context, user *User) error
	FindUser(ctx context.Context, address string) (*User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) SaveChallenge(ctx context.Context, challenge *Challenge) error {
	return r.db.WithContext(ctx).Save(challenge).Error
}

func (r *GormRepository) FindChallenge(ctx context.Context, address string) (*Challenge, error) {
	var challenge Challenge
	if err := r.db.WithContext(ctx).Where("address = ?", address).First(&challenge).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find challenge: %w", err)
	}
	return &challenge, nil
}

func (r *GormRepository) DeleteChallenge(ctx context.Context, address string) error {
	return r.db.WithContext(ctx).Delete(&Challenge{}, "address = ?", address).Error
}

func (r *GormRepository) DeleteExpiredChallenges(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&Challenge{})
	return result.RowsAffected, result.Error
}

func (r *GormRepository) SaveUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *GormRepository) FindUser(ctx context.Context, address string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("address = ?", address).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}
