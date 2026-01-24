package consent

import (
	"slices"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type State string

const (
	StateRequested State = "requested"

	StateApproved State = "approved"

	StateDenied State = "denied"

	StateRevoked State = "revoked"

	StateExpired State = "expired"
)

func ValidStates() []State {
	return []State{StateRequested, StateApproved, StateDenied, StateRevoked, StateExpired}
}

func (s State) IsValid() bool {
	return slices.Contains(ValidStates(), s)
}

func (s State) IsTerminal() bool {
	switch s {
	case StateDenied, StateRevoked, StateExpired:
		return true
	}
	return false
}

func (s State) IsActive() bool {
	return s == StateApproved
}

type Transition struct {
	From   State
	To     State
	Action string
}

var validTransitions = []Transition{
	{StateRequested, StateApproved, "approve"},
	{StateRequested, StateDenied, "deny"},
	{StateApproved, StateRevoked, "revoke"},
	{StateApproved, StateExpired, "expire"},
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
