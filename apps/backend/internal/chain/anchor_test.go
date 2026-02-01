package chain

import (
	"testing"
)

func TestHexToBytes32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid 64 char hex",
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr: false,
		},
		{
			name:    "valid with 0x prefix",
			input:   "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr: false,
		},
		{
			name:    "too short",
			input:   "0123456789abcdef",
			wantErr: true,
		},
		{
			name:    "too long",
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef00",
			wantErr: true,
		},
		{
			name:    "invalid hex chars",
			input:   "0123456789GHIJKL0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := hexToBytes32(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("hexToBytes32() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHexToBytes32_ReturnsCorrectBytes(t *testing.T) {
	// SHA-256 of "test" is e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
	// But let's use a simpler pattern for verification
	input := "0000000000000000000000000000000000000000000000000000000000000001"

	result, err := hexToBytes32(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Last byte should be 0x01
	if result[31] != 0x01 {
		t.Errorf("expected last byte to be 0x01, got 0x%02x", result[31])
	}

	// All other bytes should be 0x00
	for i := 0; i < 31; i++ {
		if result[i] != 0x00 {
			t.Errorf("expected byte %d to be 0x00, got 0x%02x", i, result[i])
		}
	}
}
