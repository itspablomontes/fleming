package consent

import (
	"fmt"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

const SchemaVersionConsent = protocol.SchemaVersionConsent

// GrantBuilder provides a fluent API for building Grant instances with validation.
type GrantBuilder struct {
	grant *Grant
	errs  types.ValidationErrors
}

// NewGrantBuilder creates a new GrantBuilder instance.
func NewGrantBuilder() *GrantBuilder {
	return &GrantBuilder{
		grant: &Grant{
			State:         StateRequested,
			SchemaVersion: SchemaVersionConsent,
		},
	}
}

// WithID sets the grant ID.
func (b *GrantBuilder) WithID(id types.ID) *GrantBuilder {
	b.grant.ID = id
	return b
}

// WithGrantor sets the grantor wallet address.
func (b *GrantBuilder) WithGrantor(grantor types.WalletAddress) *GrantBuilder {
	b.grant.Grantor = grantor
	return b
}

// WithGrantee sets the grantee wallet address.
func (b *GrantBuilder) WithGrantee(grantee types.WalletAddress) *GrantBuilder {
	b.grant.Grantee = grantee
	return b
}

// WithScope sets the scope (list of event IDs).
func (b *GrantBuilder) WithScope(scope []types.ID) *GrantBuilder {
	b.grant.Scope = scope
	return b
}

// AddToScope adds an event ID to the scope.
func (b *GrantBuilder) AddToScope(eventID types.ID) *GrantBuilder {
	b.grant.Scope = append(b.grant.Scope, eventID)
	return b
}

// WithPermissions sets the permissions.
func (b *GrantBuilder) WithPermissions(permissions Permissions) *GrantBuilder {
	b.grant.Permissions = permissions
	return b
}

// AddPermission adds a permission.
func (b *GrantBuilder) AddPermission(permission Permission) *GrantBuilder {
	if !permission.IsValid() {
		b.errs.Add("permissions", fmt.Sprintf("invalid permission: %s", permission))
		return b
	}
	b.grant.Permissions = append(b.grant.Permissions, permission)
	return b
}

// WithState sets the grant state.
func (b *GrantBuilder) WithState(state State) *GrantBuilder {
	b.grant.State = state
	return b
}

// WithExpiresAt sets the expiration time.
func (b *GrantBuilder) WithExpiresAt(expiresAt time.Time) *GrantBuilder {
	b.grant.ExpiresAt = expiresAt
	return b
}

// WithReason sets the reason for the grant.
func (b *GrantBuilder) WithReason(reason string) *GrantBuilder {
	b.grant.Reason = reason
	return b
}

// WithCreatedAt sets the creation timestamp.
func (b *GrantBuilder) WithCreatedAt(createdAt time.Time) *GrantBuilder {
	b.grant.CreatedAt = createdAt
	return b
}

// WithUpdatedAt sets the update timestamp.
func (b *GrantBuilder) WithUpdatedAt(updatedAt time.Time) *GrantBuilder {
	b.grant.UpdatedAt = updatedAt
	return b
}

// Build validates and returns the Grant, or returns an error if validation fails.
func (b *GrantBuilder) Build() (*Grant, error) {
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	if err := b.grant.Validate(); err != nil {
		return nil, fmt.Errorf("grant validation failed: %w", err)
	}

	// Set timestamps if not provided
	if b.grant.CreatedAt.IsZero() {
		b.grant.CreatedAt = time.Now().UTC()
	}
	if b.grant.UpdatedAt.IsZero() {
		b.grant.UpdatedAt = b.grant.CreatedAt
	}

	return b.grant, nil
}
