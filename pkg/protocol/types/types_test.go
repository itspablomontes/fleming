package types

import "testing"

func TestNewID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid ID", "abc123", false},
		{"UUID format", "550e8400-e29b-41d4-a716-446655440000", false},
		{"empty ID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id.String() != tt.input {
				t.Errorf("NewID() = %v, want %v", id.String(), tt.input)
			}
		})
	}
}

func TestNewWalletAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid lowercase", "0x1234567890abcdef1234567890abcdef12345678", false},
		{"valid uppercase", "0x1234567890ABCDEF1234567890ABCDEF12345678", false},
		{"valid mixed case", "0x1234567890AbCdEf1234567890AbCdEf12345678", false},
		{"missing 0x prefix", "1234567890abcdef1234567890abcdef12345678", true},
		{"too short", "0x1234567890abcdef", true},
		{"too long", "0x1234567890abcdef1234567890abcdef1234567890", true},
		{"invalid chars", "0x1234567890ghijkl1234567890ghijkl12345678", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := NewWalletAddress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWalletAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && addr.IsEmpty() {
				t.Error("NewWalletAddress() returned empty for valid input")
			}
		})
	}
}

func TestWalletAddress_Equals(t *testing.T) {
	addr1, _ := NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")
	addr2, _ := NewWalletAddress("0x1234567890ABCDEF1234567890ABCDEF12345678")
	addr3, _ := NewWalletAddress("0xabcdef1234567890abcdef1234567890abcdef12")

	if !addr1.Equals(addr2) {
		t.Error("Expected case-insensitive equality")
	}
	if addr1.Equals(addr3) {
		t.Error("Expected different addresses to not be equal")
	}
}

func TestMetadata(t *testing.T) {
	m := NewMetadata()
	m.Set("name", "test")
	m.Set("count", 42)

	if m.GetString("name") != "test" {
		t.Error("GetString failed")
	}
	if m.GetInt("count") != 42 {
		t.Error("GetInt failed")
	}
	if m.GetString("missing") != "" {
		t.Error("GetString should return empty for missing key")
	}
	if m.GetInt("missing") != 0 {
		t.Error("GetInt should return 0 for missing key")
	}
}

func TestTimestamp(t *testing.T) {
	ts := Now()
	if ts.IsZero() {
		t.Error("Now() should not return zero timestamp")
	}

	parsed, err := ParseTimestamp("2026-01-23T12:00:00Z")
	if err != nil {
		t.Errorf("ParseTimestamp failed: %v", err)
	}
	if parsed.IsZero() {
		t.Error("Parsed timestamp should not be zero")
	}

	_, err = ParseTimestamp("invalid")
	if err == nil {
		t.Error("ParseTimestamp should fail on invalid input")
	}
}
