package consent

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestGrantBuilder_WithGrantor(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	builder := NewGrantBuilder()

	builder.WithGrantor(validAddr)
	if builder.grant.Grantor != validAddr {
		t.Error("WithGrantor() did not set grantor")
	}
}

func TestGrantBuilder_WithGrantee(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")
	builder := NewGrantBuilder()

	builder.WithGrantee(validAddr)
	if builder.grant.Grantee != validAddr {
		t.Error("WithGrantee() did not set grantee")
	}
}

func TestGrantBuilder_AddPermission(t *testing.T) {
	builder := NewGrantBuilder()

	builder.AddPermission(PermRead)
	if len(builder.grant.Permissions) != 1 {
		t.Errorf("AddPermission() expected 1 permission, got %d", len(builder.grant.Permissions))
	}

	// Invalid permission should add error
	builder2 := NewGrantBuilder()
	builder2.AddPermission("invalid")
	if !builder2.errs.HasErrors() {
		t.Error("AddPermission() with invalid permission should add error")
	}
}

func TestGrantBuilder_AddToScope(t *testing.T) {
	eventID1, _ := types.NewID("event-1")
	eventID2, _ := types.NewID("event-2")
	builder := NewGrantBuilder()

	builder.AddToScope(eventID1)
	builder.AddToScope(eventID2)

	if len(builder.grant.Scope) != 2 {
		t.Errorf("AddToScope() expected 2 events, got %d", len(builder.grant.Scope))
	}
}

func TestGrantBuilder_Build(t *testing.T) {
	grantor, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	grantee, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")

	tests := []struct {
		name    string
		builder func() *GrantBuilder
		wantErr bool
	}{
		{
			name: "valid grant",
			builder: func() *GrantBuilder {
				return NewGrantBuilder().
					WithGrantor(grantor).
					WithGrantee(grantee).
					AddPermission(PermRead)
			},
			wantErr: false,
		},
		{
			name: "missing grantor",
			builder: func() *GrantBuilder {
				return NewGrantBuilder().
					WithGrantee(grantee).
					AddPermission(PermRead)
			},
			wantErr: true,
		},
		{
			name: "missing grantee",
			builder: func() *GrantBuilder {
				return NewGrantBuilder().
					WithGrantor(grantor).
					AddPermission(PermRead)
			},
			wantErr: true,
		},
		{
			name: "no permissions",
			builder: func() *GrantBuilder {
				return NewGrantBuilder().
					WithGrantor(grantor).
					WithGrantee(grantee)
			},
			wantErr: true,
		},
		{
			name: "self-grant",
			builder: func() *GrantBuilder {
				return NewGrantBuilder().
					WithGrantor(grantor).
					WithGrantee(grantor).
					AddPermission(PermRead)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.builder()
			grant, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if grant == nil {
					t.Error("Build() returned nil for valid grant")
				}
				if grant.CreatedAt.IsZero() {
					t.Error("Build() should set CreatedAt if not provided")
				}
			}
		})
	}
}

func TestGrantBuilder_WithExpiresAt(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	builder := NewGrantBuilder()

	builder.WithExpiresAt(future)
	if builder.grant.ExpiresAt != future {
		t.Error("WithExpiresAt() did not set expiration")
	}
}
