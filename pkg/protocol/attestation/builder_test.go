package attestation

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestAttestationBuilder_WithEventID(t *testing.T) {
	validID, _ := types.NewID("event-1")
	builder := NewAttestationBuilder()

	builder.WithEventID(validID)
	if builder.att.EventID != validID {
		t.Error("WithEventID() did not set event ID")
	}

	emptyID := types.ID("")
	builder2 := NewAttestationBuilder()
	builder2.WithEventID(emptyID)
	if !builder2.errs.HasErrors() {
		t.Error("WithEventID() with empty ID should add error")
	}
}

func TestAttestationBuilder_WithEventHash(t *testing.T) {
	builder := NewAttestationBuilder()

	builder.WithEventHash("hash123")
	if builder.att.EventHash != "hash123" {
		t.Error("WithEventHash() did not set hash")
	}

	// Empty hash should add error
	builder2 := NewAttestationBuilder()
	builder2.WithEventHash("")
	if !builder2.errs.HasErrors() {
		t.Error("WithEventHash() with empty hash should add error")
	}
}

func TestAttestationBuilder_WithAttester(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	builder := NewAttestationBuilder()

	builder.WithAttester(validAddr)
	if builder.att.Attester != validAddr {
		t.Error("WithAttester() did not set attester")
	}

	// Empty address should add error
	emptyAddr := types.WalletAddress("")
	builder2 := NewAttestationBuilder()
	builder2.WithAttester(emptyAddr)
	if !builder2.errs.HasErrors() {
		t.Error("WithAttester() with empty address should add error")
	}
}

func TestAttestationBuilder_WithType(t *testing.T) {
	builder := NewAttestationBuilder()

	builder.WithType(AttestVerified)
	if builder.att.Type != AttestVerified {
		t.Error("WithType() did not set type")
	}

	// Invalid type should add error
	builder2 := NewAttestationBuilder()
	builder2.WithType("invalid")
	if !builder2.errs.HasErrors() {
		t.Error("WithType() with invalid type should add error")
	}
}

func TestAttestationBuilder_WithTTL(t *testing.T) {
	builder := NewAttestationBuilder()
	ttl := 30 * 24 * time.Hour

	builder.WithTTL(ttl)
	if builder.att.ExpiresAt == nil {
		t.Error("WithTTL() should set ExpiresAt")
	}
	if builder.att.ExpiresAt.Before(time.Now()) {
		t.Error("WithTTL() should set future expiration")
	}
}

func TestAttestationBuilder_Build(t *testing.T) {
	validEventID, _ := types.NewID("event-1")
	validAttester, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name    string
		builder func() *AttestationBuilder
		wantErr bool
	}{
		{
			name: "valid attestation",
			builder: func() *AttestationBuilder {
				return NewAttestationBuilder().
					WithEventID(validEventID).
					WithEventHash("hash123").
					WithAttester(validAttester).
					WithType(AttestVerified)
			},
			wantErr: false,
		},
		{
			name: "missing event ID",
			builder: func() *AttestationBuilder {
				return NewAttestationBuilder().
					WithEventHash("hash123").
					WithAttester(validAttester).
					WithType(AttestVerified)
			},
			wantErr: true,
		},
		{
			name: "missing event hash",
			builder: func() *AttestationBuilder {
				return NewAttestationBuilder().
					WithEventID(validEventID).
					WithAttester(validAttester).
					WithType(AttestVerified)
			},
			wantErr: true,
		},
		{
			name: "missing attester",
			builder: func() *AttestationBuilder {
				return NewAttestationBuilder().
					WithEventID(validEventID).
					WithEventHash("hash123").
					WithType(AttestVerified)
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			builder: func() *AttestationBuilder {
				return NewAttestationBuilder().
					WithEventID(validEventID).
					WithEventHash("hash123").
					WithAttester(validAttester).
					WithType("invalid")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.builder()
			att, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if att == nil {
					t.Error("Build() returned nil for valid attestation")
				}
				if att.Status != StatusPendingAttestation {
					t.Errorf("Build() status = %v, want %v", att.Status, StatusPendingAttestation)
				}
			}
		})
	}
}

func TestAttestationBuilder_BuildSigned(t *testing.T) {
	validEventID, _ := types.NewID("event-1")
	validAttester, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	builder := NewAttestationBuilder().
		WithEventID(validEventID).
		WithEventHash("hash123").
		WithAttester(validAttester).
		WithType(AttestVerified)

	att, err := builder.BuildSigned("sig123", "ES256K")
	if err != nil {
		t.Errorf("BuildSigned() error = %v", err)
		return
	}

	if att.Status != StatusActiveAttestation {
		t.Errorf("BuildSigned() status = %v, want %v", att.Status, StatusActiveAttestation)
	}
	if att.Signature != "sig123" {
		t.Errorf("BuildSigned() signature = %v, want sig123", att.Signature)
	}
	if att.SignatureAlgorithm != "ES256K" {
		t.Errorf("BuildSigned() algorithm = %v, want ES256K", att.SignatureAlgorithm)
	}
}

func TestProviderCredentialsBuilder(t *testing.T) {
	builder := NewProviderCredentialsBuilder().
		WithName("Dr. Smith").
		WithLicense("MD12345", "MD").
		WithSpecialty("Cardiology").
		WithOrganization("Hospital A").
		WithNPI("1234567890")

	creds := builder.Build()

	if creds.Name != "Dr. Smith" {
		t.Errorf("WithName() name = %v, want Dr. Smith", creds.Name)
	}
	if creds.LicenseNumber != "MD12345" {
		t.Errorf("WithLicense() number = %v, want MD12345", creds.LicenseNumber)
	}
	if creds.LicenseType != "MD" {
		t.Errorf("WithLicense() type = %v, want MD", creds.LicenseType)
	}
	if creds.Specialty != "Cardiology" {
		t.Errorf("WithSpecialty() = %v, want Cardiology", creds.Specialty)
	}
	if creds.Organization != "Hospital A" {
		t.Errorf("WithOrganization() = %v, want Hospital A", creds.Organization)
	}
	if creds.NPI != "1234567890" {
		t.Errorf("WithNPI() = %v, want 1234567890", creds.NPI)
	}
}

func TestAttestationRequestBuilder(t *testing.T) {
	validEventID, _ := types.NewID("event-1")
	validRequester, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name    string
		builder func() *AttestationRequestBuilder
		wantErr bool
	}{
		{
			name: "valid request",
			builder: func() *AttestationRequestBuilder {
				return NewAttestationRequestBuilder().
					WithEventID(validEventID).
					WithRequester(validRequester).
					WithRequestedType(AttestVerified)
			},
			wantErr: false,
		},
		{
			name: "missing event ID",
			builder: func() *AttestationRequestBuilder {
				return NewAttestationRequestBuilder().
					WithRequester(validRequester).
					WithRequestedType(AttestVerified)
			},
			wantErr: true,
		},
		{
			name: "missing requester",
			builder: func() *AttestationRequestBuilder {
				return NewAttestationRequestBuilder().
					WithEventID(validEventID).
					WithRequestedType(AttestVerified)
			},
			wantErr: true,
		},
		{
			name: "invalid requested type",
			builder: func() *AttestationRequestBuilder {
				return NewAttestationRequestBuilder().
					WithEventID(validEventID).
					WithRequester(validRequester).
					WithRequestedType("invalid")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.builder()
			req, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && req == nil {
				t.Error("Build() returned nil for valid request")
			}
		})
	}
}
