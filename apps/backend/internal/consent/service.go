package consent

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	protocol "github.com/itspablomontes/fleming/pkg/protocol/audit"
	"github.com/itspablomontes/fleming/pkg/protocol/consent"
)

// Service defines the business logic for patient consent.
type Service interface {
	RequestConsent(ctx context.Context, grantor, grantee, reason string, permissions []string, expiresAt time.Time) (*ConsentGrant, error)
	ApproveConsent(ctx context.Context, grantID string) error
	DenyConsent(ctx context.Context, grantID string) error
	RevokeConsent(ctx context.Context, grantID string) error
	GetActiveGrants(ctx context.Context, grantee string) ([]ConsentGrant, error)
	CheckPermission(ctx context.Context, grantor, grantee string, permission string) (bool, error)
}

type service struct {
	repo         Repository
	auditService audit.Service
}

// NewService creates a new consent service.
func NewService(repo Repository, auditService audit.Service) Service {
	return &service{
		repo:         repo,
		auditService: auditService,
	}
}

func (s *service) RequestConsent(ctx context.Context, grantor, grantee, reason string, permissions []string, expiresAt time.Time) (*ConsentGrant, error) {
	grant := &ConsentGrant{
		Grantor:     grantor,
		Grantee:     grantee,
		Reason:      reason,
		Permissions: permissions,
		State:       consent.StateRequested,
		ExpiresAt:   expiresAt,
	}

	if err := s.repo.Create(ctx, grant); err != nil {
		return nil, err
	}

	metadata := common.JSONMap{
		"grantee":     grant.Grantee,
		"permissions": grant.Permissions,
		"expiresAt":   grant.ExpiresAt,
	}
	_ = s.auditService.Record(ctx, grantor, protocol.ActionConsentRequest, protocol.ResourceConsent, grant.ID, metadata)
	return grant, nil
}

func (s *service) ApproveConsent(ctx context.Context, grantID string) error {
	grant, err := s.repo.GetByID(ctx, grantID)
	if err != nil {
		return err
	}

	if err := consent.TryTransition(grant.State, consent.StateApproved); err != nil {
		return fmt.Errorf("invalid transition: %w", err)
	}

	grant.State = consent.StateApproved
	if err := s.repo.Update(ctx, grant); err != nil {
		return err
	}

	_ = s.auditService.Record(ctx, grant.Grantor, protocol.ActionConsentApprove, protocol.ResourceConsent, grant.ID, nil)
	return nil
}

func (s *service) DenyConsent(ctx context.Context, grantID string) error {
	grant, err := s.repo.GetByID(ctx, grantID)
	if err != nil {
		return err
	}

	if err := consent.TryTransition(grant.State, consent.StateDenied); err != nil {
		return fmt.Errorf("invalid transition: %w", err)
	}

	grant.State = consent.StateDenied
	if err := s.repo.Update(ctx, grant); err != nil {
		return err
	}

	_ = s.auditService.Record(ctx, grant.Grantor, protocol.ActionConsentDeny, protocol.ResourceConsent, grant.ID, nil)
	return nil
}

func (s *service) RevokeConsent(ctx context.Context, grantID string) error {
	grant, err := s.repo.GetByID(ctx, grantID)
	if err != nil {
		return err
	}

	if err := consent.TryTransition(grant.State, consent.StateRevoked); err != nil {
		return fmt.Errorf("invalid transition: %w", err)
	}

	grant.State = consent.StateRevoked
	if err := s.repo.Update(ctx, grant); err != nil {
		return err
	}

	_ = s.auditService.Record(ctx, grant.Grantor, protocol.ActionConsentRevoke, protocol.ResourceConsent, grant.ID, nil)
	return nil
}

func (s *service) GetActiveGrants(ctx context.Context, grantee string) ([]ConsentGrant, error) {
	all, err := s.repo.GetByGrantee(ctx, grantee)
	if err != nil {
		return nil, err
	}

	active := make([]ConsentGrant, 0)
	now := time.Now()
	for _, g := range all {
		if g.State == consent.StateApproved {
			if g.ExpiresAt.IsZero() || g.ExpiresAt.After(now) {
				active = append(active, g)
			}
		}
	}
	return active, nil
}

func (s *service) CheckPermission(ctx context.Context, grantor, grantee string, permission string) (bool, error) {
	if grantor == grantee {
		return true, nil
	}

	latest, err := s.repo.FindLatest(ctx, grantor, grantee)
	if err != nil {
		return false, err
	}

	if latest == nil || latest.State != consent.StateApproved {
		return false, nil
	}

	if !latest.ExpiresAt.IsZero() && latest.ExpiresAt.Before(time.Now()) {
		latest.State = consent.StateExpired
		_ = s.repo.Update(ctx, latest)
		_ = s.auditService.Record(ctx, latest.Grantor, protocol.ActionConsentExpire, protocol.ResourceConsent, latest.ID, nil)
		return false, nil
	}

	if slices.Contains(latest.Permissions, permission) {
		return true, nil
	}

	return false, nil
}
