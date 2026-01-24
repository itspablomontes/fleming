package consent

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func newValidGrant() *Grant {
	grantor, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	grantee, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")

	return &Grant{
		ID:          "grant-1",
		Grantor:     grantor,
		Grantee:     grantee,
		Permissions: Permissions{PermRead},
		State:       StateRequested,
		CreatedAt:   time.Now(),
	}
}

func TestGrant_Validate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Grant)
		wantErr bool
	}{
		{
			name:    "valid grant",
			modify:  func(g *Grant) {},
			wantErr: false,
		},
		{
			name: "missing grantor",
			modify: func(g *Grant) {
				g.Grantor = ""
			},
			wantErr: true,
		},
		{
			name: "missing grantee",
			modify: func(g *Grant) {
				g.Grantee = ""
			},
			wantErr: true,
		},
		{
			name: "self-grant",
			modify: func(g *Grant) {
				g.Grantee = g.Grantor
			},
			wantErr: true,
		},
		{
			name: "no permissions",
			modify: func(g *Grant) {
				g.Permissions = nil
			},
			wantErr: true,
		},
		{
			name: "invalid permission",
			modify: func(g *Grant) {
				g.Permissions = Permissions{"invalid"}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newValidGrant()
			tt.modify(g)
			err := g.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGrant_IsExpired(t *testing.T) {
	g := newValidGrant()

	if g.IsExpired() {
		t.Error("Grant without expiration should not be expired")
	}

	g.ExpiresAt = time.Now().Add(time.Hour)
	if g.IsExpired() {
		t.Error("Grant with future expiration should not be expired")
	}
	g.ExpiresAt = time.Now().Add(-time.Hour)
	if !g.IsExpired() {
		t.Error("Grant with past expiration should be expired")
	}
}

func TestGrant_IsActive(t *testing.T) {
	g := newValidGrant()

	if g.IsActive() {
		t.Error("Requested grant should not be active")
	}
	g.State = StateApproved
	if !g.IsActive() {
		t.Error("Approved grant should be active")
	}

	g.ExpiresAt = time.Now().Add(-time.Hour)
	if g.IsActive() {
		t.Error("Expired grant should not be active")
	}
}

func TestGrant_CanAccess(t *testing.T) {
	g := newValidGrant()
	g.State = StateApproved

	if !g.CanAccess("any-event") {
		t.Error("Empty scope should allow access to any event")
	}
	g.Scope = []types.ID{"event-1", "event-2"}

	if !g.CanAccess("event-1") {
		t.Error("Should allow access to scoped event")
	}

	if g.CanAccess("event-3") {
		t.Error("Should deny access to non-scoped event")
	}
}

func TestGrant_Transitions(t *testing.T) {
	g := newValidGrant()

	if err := g.Approve(); err != nil {
		t.Errorf("Approve() error = %v", err)
	}
	if g.State != StateApproved {
		t.Errorf("Expected state approved, got %s", g.State)
	}

	if err := g.Revoke(); err != nil {
		t.Errorf("Revoke() error = %v", err)
	}
	if g.State != StateRevoked {
		t.Errorf("Expected state revoked, got %s", g.State)
	}

	if err := g.Approve(); err == nil {
		t.Error("Expected error when transitioning from terminal state")
	}
}

func TestGrant_ScopeManagement(t *testing.T) {
	g := newValidGrant()

	g.AddToScope("event-1")
	g.AddToScope("event-2")
	g.AddToScope("event-1")

	if len(g.Scope) != 2 {
		t.Errorf("Expected 2 events in scope, got %d", len(g.Scope))
	}

	g.RemoveFromScope("event-1")
	if len(g.Scope) != 1 {
		t.Errorf("Expected 1 event in scope after removal, got %d", len(g.Scope))
	}
}
