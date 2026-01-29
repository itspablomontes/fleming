package consent

import (
	"slices"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type Permission string

const (
	PermRead  Permission = "read"
	PermWrite Permission = "write"
	PermShare Permission = "share"
)

func ValidPermissions() []Permission {
	return []Permission{PermRead, PermWrite, PermShare}
}

func (p Permission) IsValid() bool {
	return slices.Contains(ValidPermissions(), p)
}

type Permissions []Permission

func (pp Permissions) Has(p Permission) bool {
	return slices.Contains(pp, p)
}

type Grant struct {
	ID          types.ID            `json:"id"`
	Grantor     types.WalletAddress `json:"grantor"`
	Grantee     types.WalletAddress `json:"grantee"`
	Scope       []types.ID          `json:"scope,omitempty"`
	Permissions Permissions         `json:"permissions"`
	State       State               `json:"state"`
	ExpiresAt   time.Time           `json:"expiresAt,omitempty"`
	Reason      string              `json:"reason,omitempty"`
	SchemaVersion string            `json:"schemaVersion,omitempty"` // Protocol schema version (e.g., "consent.v1")
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

func (g *Grant) Validate() error {
	var errs types.ValidationErrors

	if g.Grantor.IsEmpty() {
		errs.Add("grantor", "grantor address is required")
	}

	if g.Grantee.IsEmpty() {
		errs.Add("grantee", "grantee address is required")
	}

	if g.Grantor.Equals(g.Grantee) {
		errs.Add("grantee", "cannot grant consent to self")
	}

	if len(g.Permissions) == 0 {
		errs.Add("permissions", "at least one permission is required")
	}

	for _, p := range g.Permissions {
		if !p.IsValid() {
			errs.Add("permissions", "invalid permission: "+string(p))
		}
	}

	if !g.State.IsValid() {
		errs.Add("state", "invalid state")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func (g *Grant) IsExpired() bool {
	if g.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(g.ExpiresAt)
}

func (g *Grant) IsActive() bool {
	return g.State.IsActive() && !g.IsExpired()
}

func (g *Grant) HasPermission(p Permission) bool {
	if !g.IsActive() {
		return false
	}
	return g.Permissions.Has(p)
}

func (g *Grant) CanAccess(eventID types.ID) bool {
	if !g.IsActive() {
		return false
	}

	if len(g.Scope) == 0 {
		return true
	}

	return slices.Contains(g.Scope, eventID)
}

func (g *Grant) AddToScope(eventID types.ID) {
	if slices.Contains(g.Scope, eventID) {
		return
	}
	g.Scope = append(g.Scope, eventID)
}

func (g *Grant) RemoveFromScope(eventID types.ID) {
	for i, id := range g.Scope {
		if id == eventID {
			g.Scope = append(g.Scope[:i], g.Scope[i+1:]...)
			return
		}
	}
}

func (g *Grant) Transition(newState State) error {
	if err := TryTransition(g.State, newState); err != nil {
		return err
	}
	g.State = newState
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Grant) Approve() error {
	return g.Transition(StateApproved)
}
func (g *Grant) Deny() error {
	return g.Transition(StateDenied)
}

func (g *Grant) Revoke() error {
	return g.Transition(StateRevoked)
}

func (g *Grant) Expire() error {
	return g.Transition(StateExpired)
}
