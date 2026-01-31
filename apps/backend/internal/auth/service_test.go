package auth

import (
	"context"
	"testing"
	"time"

	internalAudit "github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/pkg/protocol/audit"
)

type MockAuditService struct{}

func (m *MockAuditService) Record(ctx context.Context, actor string, action audit.Action, resourceType audit.ResourceType, resourceID string, metadata common.JSONMap) error {
	return nil
}
func (m *MockAuditService) GetLatestEntries(ctx context.Context, actor string, limit int) ([]internalAudit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) VerifyIntegrity(ctx context.Context) (bool, error) {
	return true, nil
}
func (m *MockAuditService) BuildMerkleTree(ctx context.Context, actor string, startTime time.Time, endTime time.Time) (*internalAudit.AuditBatch, *audit.MerkleTree, error) {
	return nil, nil, nil
}
func (m *MockAuditService) GetBatch(ctx context.Context, actor string, batchID string) (*internalAudit.AuditBatch, error) {
	return nil, nil
}
func (m *MockAuditService) ListBatches(ctx context.Context, actor string, limit int, offset int) ([]internalAudit.AuditBatch, error) {
	return nil, nil
}
func (m *MockAuditService) AnchorBatch(ctx context.Context, actor string, batchID string, chainClient internalAudit.ChainAnchorer) (*internalAudit.AuditBatch, error) {
	return nil, nil
}
func (m *MockAuditService) GetBatchByRoot(ctx context.Context, actor string, rootHash string) (*internalAudit.AuditBatch, error) {
	return nil, nil
}
func (m *MockAuditService) VerifyMerkleProof(root string, entryHash string, proof *audit.Proof) bool {
	return true
}
func (m *MockAuditService) GetEntriesForMerkle(ctx context.Context, actor string, startTime time.Time, endTime time.Time) ([]internalAudit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) GetEntryByID(ctx context.Context, id string) (*internalAudit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) GetEntriesByResource(ctx context.Context, resourceID string) ([]internalAudit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) QueryEntries(ctx context.Context, filter audit.QueryFilter) ([]internalAudit.AuditEntry, error) {
	return nil, nil
}

type MockRepo struct {
	challenges map[string]*Challenge
	users      map[string]*User
}

func (m *MockRepo) SaveChallenge(ctx context.Context, c *Challenge) error {
	if m.challenges == nil {
		m.challenges = make(map[string]*Challenge)
	}
	m.challenges[c.Address] = c
	return nil
}

func (m *MockRepo) FindChallenge(ctx context.Context, address string) (*Challenge, error) {
	if c, ok := m.challenges[address]; ok {
		return c, nil
	}
	return nil, ErrNotFound
}

func (m *MockRepo) DeleteChallenge(ctx context.Context, address string) error {
	delete(m.challenges, address)
	return nil
}

func (m *MockRepo) DeleteExpiredChallenges(ctx context.Context) (int64, error) {
	var count int64
	for k, v := range m.challenges {
		if time.Now().After(v.ExpiresAt) {
			delete(m.challenges, k)
			count++
		}
	}
	return count, nil
}

func (m *MockRepo) SaveUser(ctx context.Context, u *User) error {
	if m.users == nil {
		m.users = make(map[string]*User)
	}
	m.users[u.Address] = u
	return nil
}

func (m *MockRepo) FindUser(ctx context.Context, address string) (*User, error) {
	if u, ok := m.users[address]; ok {
		return u, nil
	}
	return nil, ErrNotFound
}

func TestService_GenerateChallenge(t *testing.T) {
	repo := &MockRepo{}
	auditSvc := &MockAuditService{}
	svc := NewService(repo, "secret", auditSvc)

	tests := []struct {
		name    string
		req     ChallengeRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: ChallengeRequest{
				Address: "0x1234567890abcdef1234567890abcdef12345678",
				Domain:  "example.com",
				URI:     "https://example.com",
				ChainID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := svc.GenerateChallenge(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateChallenge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if msg == "" {
				t.Errorf("GenerateChallenge() returned empty message")
			}
			// Verify stored in repo
			if _, err := repo.FindChallenge(context.Background(), tt.req.Address); err != nil {
				t.Errorf("Challenge not stored in repo")
			}
		})
	}
}
