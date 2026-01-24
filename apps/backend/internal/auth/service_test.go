package auth

import (
	"context"
	"testing"
	"time"
)

type MockRepo struct {
	challenges map[string]*Challenge
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

func TestService_GenerateChallenge(t *testing.T) {
	repo := &MockRepo{}
	svc := NewService(repo, "secret")

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
