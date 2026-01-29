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
		{StateSuspended, true},
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
		{StateSuspended, false}, // Suspended is NOT terminal
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

func TestState_IsSuspended(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{StateSuspended, true},
		{StateRequested, false},
		{StateApproved, false},
		{StateDenied, false},
		{StateRevoked, false},
		{StateExpired, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.IsSuspended(); got != tt.want {
				t.Errorf("IsSuspended() = %v, want %v", got, tt.want)
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
		{"approved to suspended", StateApproved, StateSuspended, true},
		{"suspended to approved", StateSuspended, StateApproved, true},
		{"suspended to revoked", StateSuspended, StateRevoked, true},
		{"requested to revoked", StateRequested, StateRevoked, false},
		{"approved to denied", StateApproved, StateDenied, false},
		{"denied to approved", StateDenied, StateApproved, false},
		{"revoked to approved", StateRevoked, StateApproved, false},
		{"suspended to denied", StateSuspended, StateDenied, false},
		{"suspended to expired", StateSuspended, StateExpired, false},
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
	tests := []struct {
		name     string
		from     State
		to       State
		want     string
		wantOk   bool
	}{
		{"requested to approved", StateRequested, StateApproved, "approve", true},
		{"requested to denied", StateRequested, StateDenied, "deny", true},
		{"approved to revoked", StateApproved, StateRevoked, "revoke", true},
		{"approved to expired", StateApproved, StateExpired, "expire", true},
		{"approved to suspended", StateApproved, StateSuspended, "suspend", true},
		{"suspended to approved", StateSuspended, StateApproved, "resume", true},
		{"suspended to revoked", StateSuspended, StateRevoked, "revoke", true},
		{"denied to approved", StateDenied, StateApproved, "", false},
		{"invalid transition", StateRequested, StateRevoked, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, ok := GetAction(tt.from, tt.to)
			if ok != tt.wantOk {
				t.Errorf("GetAction() ok = %v, want %v", ok, tt.wantOk)
				return
			}
			if ok && action != tt.want {
				t.Errorf("GetAction() action = %v, want %v", action, tt.want)
			}
		})
	}
}

func TestTryTransition(t *testing.T) {
	tests := []struct {
		name    string
		from    State
		to      State
		wantErr bool
	}{
		{"requested to approved", StateRequested, StateApproved, false},
		{"approved to suspended", StateApproved, StateSuspended, false},
		{"suspended to approved", StateSuspended, StateApproved, false},
		{"suspended to revoked", StateSuspended, StateRevoked, false},
		{"denied to approved", StateDenied, StateApproved, true},
		{"requested to revoked", StateRequested, StateRevoked, true},
		{"suspended to denied", StateSuspended, StateDenied, true},
		{"invalid state", "invalid", StateApproved, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TryTransition(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
