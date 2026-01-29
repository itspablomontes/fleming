// Package attestation provides provider attestation types for the Protocol layer.
// Attestations allow healthcare providers to verify and co-sign patient health data,
// adding trust and credibility to self-reported information.
package attestation

import (
	"sync"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// SchemaVersionAttestation is the schema version for attestations.
const SchemaVersionAttestation = protocol.SchemaVersionAttestation

// AttestationType represents the type of attestation provided by a healthcare provider.
type AttestationType string

const (
	// AttestAccurate indicates the provider confirms the data is accurate.
	AttestAccurate AttestationType = "accurate"

	// AttestVerified indicates the provider has verified the results independently.
	AttestVerified AttestationType = "verified"

	// AttestGenerated indicates the provider generated the data (e.g., ordered the test).
	AttestGenerated AttestationType = "generated"

	// AttestReviewed indicates the provider has reviewed the data.
	AttestReviewed AttestationType = "reviewed"

	// AttestAmended indicates the provider has amended/corrected the data.
	AttestAmended AttestationType = "amended"
)

var (
	defaultAttestationTypeRegistry types.TypeRegistry[AttestationType]
	attestationTypeRegistryOnce    sync.Once
)

func init() {
	attestationTypeRegistryOnce.Do(func() {
		defaultAttestationTypeRegistry = types.NewTypeRegistry[AttestationType]()
		RegisterDefaultAttestationTypes()
	})
}

// GetAttestationTypeRegistry returns the default attestation type registry.
func GetAttestationTypeRegistry() types.TypeRegistry[AttestationType] {
	return defaultAttestationTypeRegistry
}

// RegisterAttestationType registers a custom attestation type at runtime.
func RegisterAttestationType(at AttestationType, metadata types.TypeMetadata) error {
	return defaultAttestationTypeRegistry.Register(at, metadata)
}

// IsValid checks if the attestation type is registered.
func (at AttestationType) IsValid() bool {
	return defaultAttestationTypeRegistry.IsValid(at)
}

// RegisterDefaultAttestationTypes registers all built-in attestation types.
func RegisterDefaultAttestationTypes() {
	reg := defaultAttestationTypeRegistry
	types.RegisterBatch(reg, map[AttestationType]types.TypeMetadata{
		AttestAccurate: {
			Name:        "Accurate",
			Description: "Provider confirms the data is accurate as reported",
			Since:       "0.1.0",
		},
		AttestVerified: {
			Name:        "Verified",
			Description: "Provider has independently verified the data",
			Since:       "0.1.0",
		},
		AttestGenerated: {
			Name:        "Generated",
			Description: "Provider generated this data (ordered test, performed procedure)",
			Since:       "0.1.0",
		},
		AttestReviewed: {
			Name:        "Reviewed",
			Description: "Provider has reviewed the data",
			Since:       "0.1.0",
		},
		AttestAmended: {
			Name:        "Amended",
			Description: "Provider has amended or corrected the data",
			Since:       "0.1.0",
		},
	})
}

// AttestationStatus represents the status of an attestation.
type AttestationStatus string

const (
	// StatusPendingAttestation is when the attestation is awaiting provider action.
	StatusPendingAttestation AttestationStatus = "pending"

	// StatusActiveAttestation is when the attestation is valid and active.
	StatusActiveAttestation AttestationStatus = "active"

	// StatusRevokedAttestation is when the provider has revoked the attestation.
	StatusRevokedAttestation AttestationStatus = "revoked"

	// StatusExpiredAttestation is when the attestation has expired.
	StatusExpiredAttestation AttestationStatus = "expired"
)

// IsValid checks if the attestation status is valid.
func (s AttestationStatus) IsValid() bool {
	switch s {
	case StatusPendingAttestation, StatusActiveAttestation, StatusRevokedAttestation, StatusExpiredAttestation:
		return true
	}
	return false
}

// IsActive returns true if the attestation is currently valid.
func (s AttestationStatus) IsActive() bool {
	return s == StatusActiveAttestation
}

// Attestation represents a provider's attestation of a timeline event.
// An attestation binds a provider's cryptographic signature to an event,
// confirming its accuracy or validity.
type Attestation struct {
	// ID is the unique identifier for this attestation
	ID types.ID `json:"id"`

	// EventID is the timeline event being attested
	EventID types.ID `json:"eventId"`

	// EventHash is the hash of the event at the time of attestation
	// This ensures the attestation is bound to a specific version
	EventHash string `json:"eventHash"`

	// Attester is the wallet address of the attesting provider
	Attester types.WalletAddress `json:"attester"`

	// AttesterCredentials describes the provider's qualifications
	AttesterCredentials *ProviderCredentials `json:"attesterCredentials,omitempty"`

	// Type is the type of attestation
	Type AttestationType `json:"type"`

	// Status is the current status of the attestation
	Status AttestationStatus `json:"status"`

	// Signature is the ECDSA signature of the attestation
	Signature string `json:"signature"`

	// SignatureAlgorithm is the algorithm used (e.g., "ES256K")
	SignatureAlgorithm string `json:"signatureAlgorithm"`

	// Notes are optional notes from the provider
	Notes string `json:"notes,omitempty"`

	// Timestamp is when the attestation was created
	Timestamp time.Time `json:"timestamp"`

	// ExpiresAt is when the attestation expires (optional)
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// Metadata contains additional attestation data
	Metadata types.Metadata `json:"metadata,omitempty"`

	// SchemaVersion is the protocol schema version
	SchemaVersion string `json:"schemaVersion"`
}

// ProviderCredentials represents the attesting provider's qualifications.
type ProviderCredentials struct {
	// Name is the provider's name or organization
	Name string `json:"name"`

	// LicenseNumber is the provider's license number
	LicenseNumber string `json:"licenseNumber,omitempty"`

	// LicenseType is the type of license (e.g., "MD", "DO", "NP")
	LicenseType string `json:"licenseType,omitempty"`

	// Specialty is the provider's medical specialty
	Specialty string `json:"specialty,omitempty"`

	// Organization is the provider's organization
	Organization string `json:"organization,omitempty"`

	// NPI is the National Provider Identifier (US)
	NPI string `json:"npi,omitempty"`
}

// Validate validates the attestation structure.
func (a *Attestation) Validate() error {
	var errs types.ValidationErrors

	if a.ID.IsEmpty() {
		errs.Add("id", "attestation ID is required")
	}

	if a.EventID.IsEmpty() {
		errs.Add("eventId", "event ID is required")
	}

	if a.EventHash == "" {
		errs.Add("eventHash", "event hash is required")
	}

	if a.Attester.IsEmpty() {
		errs.Add("attester", "attester address is required")
	}

	if !a.Type.IsValid() {
		errs.Add("type", "invalid attestation type")
	}

	if !a.Status.IsValid() {
		errs.Add("status", "invalid attestation status")
	}

	if a.Signature == "" && a.Status == StatusActiveAttestation {
		errs.Add("signature", "signature is required for active attestations")
	}

	if a.Timestamp.IsZero() {
		errs.Add("timestamp", "timestamp is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// IsExpired checks if the attestation has expired.
func (a *Attestation) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// IsValid checks if the attestation is currently valid (active and not expired).
func (a *Attestation) IsValid() bool {
	return a.Status.IsActive() && !a.IsExpired()
}

// AttestationRequest represents a request for a provider to attest an event.
type AttestationRequest struct {
	// RequestID is the unique identifier for this request
	RequestID types.ID `json:"requestId"`

	// EventID is the event to be attested
	EventID types.ID `json:"eventId"`

	// Requester is the wallet address requesting the attestation (usually patient)
	Requester types.WalletAddress `json:"requester"`

	// TargetAttester is the provider being asked to attest (optional)
	TargetAttester types.WalletAddress `json:"targetAttester,omitempty"`

	// RequestedType is the type of attestation being requested
	RequestedType AttestationType `json:"requestedType"`

	// RequestedAt is when the request was made
	RequestedAt time.Time `json:"requestedAt"`

	// ExpiresAt is when the request expires
	ExpiresAt time.Time `json:"expiresAt"`

	// Message is an optional message to the provider
	Message string `json:"message,omitempty"`
}

// Validate validates the attestation request.
func (r *AttestationRequest) Validate() error {
	var errs types.ValidationErrors

	if r.RequestID.IsEmpty() {
		errs.Add("requestId", "request ID is required")
	}

	if r.EventID.IsEmpty() {
		errs.Add("eventId", "event ID is required")
	}

	if r.Requester.IsEmpty() {
		errs.Add("requester", "requester is required")
	}

	if !r.RequestedType.IsValid() {
		errs.Add("requestedType", "invalid attestation type")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
