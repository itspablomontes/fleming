package vc

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestClaimType_IsValid(t *testing.T) {
	tests := []struct {
		ct   ClaimType
		want bool
	}{
		{ClaimBloodworkRange, true},
		{ClaimProtocolAdherence, true},
		{ClaimBiometricPercentile, true},
		{ClaimStackValidation, true},
		{ClaimProviderAttestation, true},
		{ClaimLabVerification, true},
		{ClaimAgeOver, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.ct), func(t *testing.T) {
			if got := tt.ct.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentialStatus_IsValid(t *testing.T) {
	tests := []struct {
		status CredentialStatus
		want   bool
	}{
		{StatusActive, true},
		{StatusRevoked, true},
		{StatusExpired, true},
		{StatusPending, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentialStatus_IsUsable(t *testing.T) {
	tests := []struct {
		status CredentialStatus
		want   bool
	}{
		{StatusActive, true},
		{StatusPending, false},
		{StatusRevoked, false},
		{StatusExpired, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsUsable(); got != tt.want {
				t.Errorf("IsUsable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredential_Validate(t *testing.T) {
	validIssuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validSubject, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")
	validID, _ := types.NewID("cred-1")

	tests := []struct {
		name    string
		cred    Credential
		wantErr bool
	}{
		{
			name: "valid credential",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  time.Now(),
				Status:    StatusActive,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			cred: Credential{
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  time.Now(),
				Status:    StatusActive,
			},
			wantErr: true,
		},
		{
			name: "missing issuer",
			cred: Credential{
				ID:        validID,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  time.Now(),
				Status:    StatusActive,
			},
			wantErr: true,
		},
		{
			name: "missing subject",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  time.Now(),
				Status:    StatusActive,
			},
			wantErr: true,
		},
		{
			name: "invalid claim type",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: "invalid",
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  time.Now(),
				Status:    StatusActive,
			},
			wantErr: true,
		},
		{
			name: "empty claims",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    nil,
				IssuedAt:  time.Now(),
				Status:    StatusActive,
			},
			wantErr: true,
		},
		{
			name: "missing issuedAt",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				Status:    StatusActive,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  time.Now(),
				Status:    "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cred.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCredential_IsExpired(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		cred    Credential
		want    bool
	}{
		{
			name: "no expiration",
			cred: Credential{
				ExpiresAt: nil,
			},
			want: false,
		},
		{
			name: "expired",
			cred: Credential{
				ExpiresAt: &past,
			},
			want: true,
		},
		{
			name: "not expired",
			cred: Credential{
				ExpiresAt: &future,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cred.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredential_IsUsable(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	past := now.Add(-time.Hour)

	validID, _ := types.NewID("cred-1")
	validIssuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validSubject, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")

	tests := []struct {
		name string
		cred Credential
		want bool
	}{
		{
			name: "active and not expired",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  now,
				Status:    StatusActive,
				ExpiresAt: &future,
			},
			want: true,
		},
		{
			name: "active but expired",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  now,
				Status:    StatusActive,
				ExpiresAt: &past,
			},
			want: false,
		},
		{
			name: "revoked",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  now,
				Status:    StatusRevoked,
			},
			want: false,
		},
		{
			name: "pending",
			cred: Credential{
				ID:        validID,
				Issuer:    validIssuer,
				Subject:   validSubject,
				ClaimType: ClaimBloodworkRange,
				Claims:    map[string]any{"marker": "718-7"},
				IssuedAt:  now,
				Status:    StatusPending,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cred.IsUsable(); got != tt.want {
				t.Errorf("IsUsable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentialRequest_Validate(t *testing.T) {
	validRequester, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validID, _ := types.NewID("req-1")
	validEventID, _ := types.NewID("event-1")

	tests := []struct {
		name    string
		req     CredentialRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CredentialRequest{
				RequestID:      validID,
				Requester:      validRequester,
				ClaimType:      ClaimBloodworkRange,
				ClaimCriteria:  map[string]any{"marker": "718-7"},
				SourceEventIDs: []types.ID{validEventID},
				RequestedAt:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing request ID",
			req: CredentialRequest{
				Requester:      validRequester,
				ClaimType:      ClaimBloodworkRange,
				SourceEventIDs: []types.ID{validEventID},
			},
			wantErr: true,
		},
		{
			name: "missing requester",
			req: CredentialRequest{
				RequestID:      validID,
				ClaimType:      ClaimBloodworkRange,
				SourceEventIDs: []types.ID{validEventID},
			},
			wantErr: true,
		},
		{
			name: "invalid claim type",
			req: CredentialRequest{
				RequestID:      validID,
				Requester:      validRequester,
				ClaimType:      "invalid",
				SourceEventIDs: []types.ID{validEventID},
			},
			wantErr: true,
		},
		{
			name: "no source events",
			req: CredentialRequest{
				RequestID:      validID,
				Requester:      validRequester,
				ClaimType:      ClaimBloodworkRange,
				SourceEventIDs: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
