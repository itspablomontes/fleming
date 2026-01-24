package consent

import "testing"

func TestState_IsValid(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{StateRequested, true},
		{StateApproved, true},
		{StateDenied, true},
		{StateRevoked, true},
		{StateExpired, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_IsTerminal(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{StateRequested, false},
		{StateApproved, false},
		{StateDenied, true},
		{StateRevoked, true},
		{StateExpired, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.IsTerminal(); got != tt.want {
				t.Errorf("IsTerminal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name string
		from State
		to   State
		want bool
	}{
		{"requested to approved", StateRequested, StateApproved, true},
		{"requested to denied", StateRequested, StateDenied, true},
		{"approved to revoked", StateApproved, StateRevoked, true},
		{"approved to expired", StateApproved, StateExpired, true},
		{"requested to revoked", StateRequested, StateRevoked, false},
		{"approved to denied", StateApproved, StateDenied, false},
		{"denied to approved", StateDenied, StateApproved, false},
		{"revoked to approved", StateRevoked, StateApproved, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanTransition(tt.from, tt.to); got != tt.want {
				t.Errorf("CanTransition(%v, %v) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestGetAction(t *testing.T) {
	action, ok := GetAction(StateRequested, StateApproved)
	if !ok {
		t.Error("Expected to find action for valid transition")
	}
	if action != "approve" {
		t.Errorf("Expected action 'approve', got '%s'", action)
	}

	_, ok = GetAction(StateDenied, StateApproved)
	if ok {
		t.Error("Expected no action for invalid transition")
	}
}

func TestTryTransition(t *testing.T) {
	err := TryTransition(StateRequested, StateApproved)
	if err != nil {
		t.Errorf("Expected valid transition, got error: %v", err)
	}

	err = TryTransition(StateDenied, StateApproved)
	if err == nil {
		t.Error("Expected error from terminal state transition")
	}

	err = TryTransition(StateRequested, StateRevoked)
	if err == nil {
		t.Error("Expected error for invalid transition")
	}
}
