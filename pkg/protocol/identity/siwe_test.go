package identity

import (
	"strings"
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestBuildSIWEMessage(t *testing.T) {
	addr, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")
	issuedAt := time.Date(2026, 1, 23, 12, 0, 0, 0, time.UTC)

	opts := SIWEOptions{
		Address:  addr,
		Domain:   "fleming.local",
		URI:      "https://fleming.local/auth",
		Nonce:    "abc123",
		ChainID:  1,
		IssuedAt: issuedAt,
	}

	msg := BuildSIWEMessage(opts)

	checks := []string{
		"fleming.local wants you to sign in with your Ethereum account:",
		"0x1234567890abcdef1234567890abcdef12345678",
		"Sign in to Fleming",
		"URI: https://fleming.local/auth",
		"Version: 1",
		"Chain ID: 1",
		"Nonce: abc123",
		"Issued At: 2026-01-23T12:00:00Z",
	}

	for _, check := range checks {
		if !strings.Contains(msg, check) {
			t.Errorf("Message missing expected part: %q\nGot:\n%s", check, msg)
		}
	}
}

func TestBuildSIWEMessage_WithExpiration(t *testing.T) {
	addr, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")
	expTime := time.Date(2026, 1, 23, 13, 0, 0, 0, time.UTC)

	opts := SIWEOptions{
		Address:        addr,
		Domain:         "fleming.local",
		URI:            "https://fleming.local/auth",
		Nonce:          "abc123",
		ChainID:        1,
		ExpirationTime: &expTime,
	}

	msg := BuildSIWEMessage(opts)

	if !strings.Contains(msg, "Expiration Time: 2026-01-23T13:00:00Z") {
		t.Errorf("Message missing expiration time\nGot:\n%s", msg)
	}
}

func TestBuildSIWEMessage_CustomStatement(t *testing.T) {
	addr, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")

	opts := SIWEOptions{
		Address:   addr,
		Domain:    "fleming.local",
		URI:       "https://fleming.local/auth",
		Nonce:     "abc123",
		ChainID:   1,
		Statement: "Custom message for testing",
	}

	msg := BuildSIWEMessage(opts)

	if !strings.Contains(msg, "Custom message for testing") {
		t.Errorf("Message missing custom statement\nGot:\n%s", msg)
	}
}

func TestSIWEOptions_Validate(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")

	tests := []struct {
		name    string
		opts    SIWEOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: SIWEOptions{
				Address: validAddr,
				Domain:  "fleming.local",
				URI:     "https://fleming.local/auth",
				Nonce:   "abc123",
				ChainID: 1,
			},
			wantErr: false,
		},
		{
			name: "missing address",
			opts: SIWEOptions{
				Domain:  "fleming.local",
				URI:     "https://fleming.local/auth",
				Nonce:   "abc123",
				ChainID: 1,
			},
			wantErr: true,
		},
		{
			name: "missing domain",
			opts: SIWEOptions{
				Address: validAddr,
				URI:     "https://fleming.local/auth",
				Nonce:   "abc123",
				ChainID: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid chain ID",
			opts: SIWEOptions{
				Address: validAddr,
				Domain:  "fleming.local",
				URI:     "https://fleming.local/auth",
				Nonce:   "abc123",
				ChainID: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
