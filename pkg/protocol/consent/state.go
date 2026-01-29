package consent

import (
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type State string

const (
	StateRequested State = "requested" // Initial state - consent request pending
	StateApproved  State = "approved"  // Consent granted and active
	StateDenied    State = "denied"    // Consent request rejected (terminal)
	StateRevoked   State = "revoked"   // Consent revoked by grantor (terminal)
	StateExpired   State = "expired"   // Consent expired due to TTL (terminal)
	StateSuspended State = "suspended" // Consent temporarily suspended (can be resumed)
)

func (s State) IsValid() bool {
	return GetStateRegistry().IsValid(s)
}

// IsTerminal returns true if the state is final and cannot transition further.
// Suspended is NOT terminal - it can be resumed.
func (s State) IsTerminal() bool {
	switch s {
	case StateDenied, StateRevoked, StateExpired:
		return true
	}
	return false
}

// IsActive returns true if the consent is currently active (approved).
// Suspended grants are NOT active.
func (s State) IsActive() bool {
	return s == StateApproved
}

// IsSuspended returns true if the consent is temporarily suspended.
func (s State) IsSuspended() bool {
	return s == StateSuspended
}

type Transition struct {
	From   State
	To     State
	Action string
}

var validTransitions = []Transition{
	// From Requested
	{StateRequested, StateApproved, "approve"},
	{StateRequested, StateDenied, "deny"},

	// From Approved
	{StateApproved, StateRevoked, "revoke"},
	{StateApproved, StateExpired, "expire"},
	{StateApproved, StateSuspended, "suspend"}, // NEW: Temporarily suspend

	// From Suspended (can resume or permanently revoke)
	{StateSuspended, StateApproved, "resume"}, // NEW: Resume suspended consent
	{StateSuspended, StateRevoked, "revoke"},  // NEW: Permanently revoke from suspended
}

func ValidTransitions() []Transition {
	return validTransitions
}

func CanTransition(from, to State) bool {
	for _, t := range validTransitions {
		if t.From == from && t.To == to {
			return true
		}
	}
	return false
}

func GetAction(from, to State) (string, bool) {
	for _, t := range validTransitions {
		if t.From == from && t.To == to {
			return t.Action, true
		}
	}
	return "", false
}

type TransitionError struct {
	From State
	To   State
}

func (e TransitionError) Error() string {
	return "invalid transition from " + string(e.From) + " to " + string(e.To)
}

func TryTransition(from, to State) error {
	if !from.IsValid() {
		return types.NewValidationError("from", "invalid state")
	}
	if !to.IsValid() {
		return types.NewValidationError("to", "invalid state")
	}
	if from.IsTerminal() {
		return TransitionError{From: from, To: to}
	}
	if !CanTransition(from, to) {
		return TransitionError{From: from, To: to}
	}
	return nil
}
