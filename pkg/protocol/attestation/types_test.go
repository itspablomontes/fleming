package attestation

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestAttestationType_IsValid(t *testing.T) {
	tests := []struct {
		at   AttestationType
		want bool
	}{
		{AttestAccurate, true},
		{AttestVerified, true},
		{AttestGenerated, true},
		{AttestReviewed, true},
		{AttestAmended, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.at), func(t *testing.T) {
			if got := tt.at.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttestationStatus_IsValid(t *testing.T) {
	tests := []struct {
		status AttestationStatus
		want   bool
	}{
		{StatusPendingAttestation, true},
		{StatusActiveAttestation, true},
		{StatusRevokedAttestation, true},
		{StatusExpiredAttestation, true},
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

func TestAttestationStatus_IsActive(t *testing.T) {
	tests := []struct {
		status AttestationStatus
		want   bool
	}{
		{StatusActiveAttestation, true},
		{StatusPendingAttestation, false},
		{StatusRevokedAttestation, false},
		{StatusExpiredAttestation, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttestation_Validate(t *testing.T) {
	validID, _ := types.NewID("att-1")
	validEventID, _ := types.NewID("event-1")
	validAttester, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name    string
		att     Attestation
		wantErr bool
	}{
		{
			name: "valid attestation",
			att: Attestation{
				ID:                 validID,
				EventID:            validEventID,
				EventHash:          "hash123",
				Attester:           validAttester,
				Type:               AttestVerified,
				Status:             StatusActiveAttestation,
				Signature:          "sig123",
				SignatureAlgorithm: "ES256K",
				Timestamp:          time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			att: Attestation{
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing event ID",
			att: Attestation{
				ID:        validID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing event hash",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing attester",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      "invalid",
				Status:    StatusActiveAttestation,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    "invalid",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing signature for active status",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing timestamp",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Signature: "sig123",
			},
			wantErr: true,
		},
		{
			name: "pending status without signature",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusPendingAttestation,
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.att.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttestation_IsExpired(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name string
		att  Attestation
		want bool
	}{
		{
			name: "no expiration",
			att: Attestation{
				ExpiresAt: nil,
			},
			want: false,
		},
		{
			name: "expired",
			att: Attestation{
				ExpiresAt: &past,
			},
			want: true,
		},
		{
			name: "not expired",
			att: Attestation{
				ExpiresAt: &future,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.att.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttestation_IsValid(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	past := now.Add(-time.Hour)

	validID, _ := types.NewID("att-1")
	validEventID, _ := types.NewID("event-1")
	validAttester, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name string
		att  Attestation
		want bool
	}{
		{
			name: "active and not expired",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Signature: "sig123",
				Timestamp: now,
				ExpiresAt: &future,
			},
			want: true,
		},
		{
			name: "active but expired",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusActiveAttestation,
				Signature: "sig123",
				Timestamp: now,
				ExpiresAt: &past,
			},
			want: false,
		},
		{
			name: "revoked",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusRevokedAttestation,
				Signature: "sig123",
				Timestamp: now,
			},
			want: false,
		},
		{
			name: "pending",
			att: Attestation{
				ID:        validID,
				EventID:   validEventID,
				EventHash: "hash123",
				Attester:  validAttester,
				Type:      AttestVerified,
				Status:    StatusPendingAttestation,
				Timestamp: now,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.att.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttestationRequest_Validate(t *testing.T) {
	validRequester, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validID, _ := types.NewID("req-1")
	validEventID, _ := types.NewID("event-1")

	tests := []struct {
		name    string
		req     AttestationRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: AttestationRequest{
				RequestID:     validID,
				EventID:       validEventID,
				Requester:     validRequester,
				RequestedType: AttestVerified,
				RequestedAt:   time.Now(),
				ExpiresAt:     time.Now().Add(7 * 24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "missing request ID",
			req: AttestationRequest{
				EventID:       validEventID,
				Requester:     validRequester,
				RequestedType: AttestVerified,
			},
			wantErr: true,
		},
		{
			name: "missing event ID",
			req: AttestationRequest{
				RequestID:     validID,
				Requester:     validRequester,
				RequestedType: AttestVerified,
			},
			wantErr: true,
		},
		{
			name: "missing requester",
			req: AttestationRequest{
				RequestID:     validID,
				EventID:       validEventID,
				RequestedType: AttestVerified,
			},
			wantErr: true,
		},
		{
			name: "invalid requested type",
			req: AttestationRequest{
				RequestID:     validID,
				EventID:       validEventID,
				Requester:     validRequester,
				RequestedType: "invalid",
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
