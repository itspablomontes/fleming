// Package vc provides Verifiable Credential types and utilities for the Protocol layer.
// Implements SD-JWT (Selective Disclosure JSON Web Tokens) for privacy-preserving claims.
package vc

import (
	"sync"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// SchemaVersionVC is the schema version for verifiable credentials.
const SchemaVersionVC = protocol.SchemaVersionVC

// ClaimType represents the type of claim in a verifiable credential.
type ClaimType string

const (
	// Health claim types
	ClaimBloodworkRange      ClaimType = "BloodworkRange"      // Biomarkers within optimal ranges
	ClaimProtocolAdherence   ClaimType = "ProtocolAdherence"   // Intervention duration compliance
	ClaimBiometricPercentile ClaimType = "BiometricPercentile" // Biometric ranking (HRV, VO2max)
	ClaimStackValidation     ClaimType = "StackValidation"     // Supplement/medication stack validation

	// Provider attestation types
	ClaimProviderAttestation ClaimType = "ProviderAttestation" // Provider confirmed accuracy
	ClaimLabVerification     ClaimType = "LabVerification"     // Lab results verified

	// Identity claims
	ClaimAgeOver ClaimType = "AgeOver" // Age is over threshold (without revealing exact age)
)

var (
	defaultClaimTypeRegistry types.TypeRegistry[ClaimType]
	claimTypeRegistryOnce    sync.Once
)

func init() {
	claimTypeRegistryOnce.Do(func() {
		defaultClaimTypeRegistry = types.NewTypeRegistry[ClaimType]()
		RegisterDefaultClaimTypes()
	})
}

// GetClaimTypeRegistry returns the default claim type registry.
func GetClaimTypeRegistry() types.TypeRegistry[ClaimType] {
	return defaultClaimTypeRegistry
}

// RegisterClaimType registers a custom claim type at runtime.
func RegisterClaimType(ct ClaimType, metadata types.TypeMetadata) error {
	return defaultClaimTypeRegistry.Register(ct, metadata)
}

// IsValid checks if the claim type is registered.
func (ct ClaimType) IsValid() bool {
	return defaultClaimTypeRegistry.IsValid(ct)
}

// RegisterDefaultClaimTypes registers all built-in claim types.
func RegisterDefaultClaimTypes() {
	reg := defaultClaimTypeRegistry
	types.RegisterBatch(reg, map[ClaimType]types.TypeMetadata{
		ClaimBloodworkRange: {
			Name:        "Bloodwork Range",
			Description: "Proves biomarkers are within specified ranges",
			Since:       "0.1.0",
		},
		ClaimProtocolAdherence: {
			Name:        "Protocol Adherence",
			Description: "Proves adherence to intervention protocol for minimum duration",
			Since:       "0.1.0",
		},
		ClaimBiometricPercentile: {
			Name:        "Biometric Percentile",
			Description: "Proves biometric values rank above specified percentile",
			Since:       "0.1.0",
		},
		ClaimStackValidation: {
			Name:        "Stack Validation",
			Description: "Proves supplement/medication stack meets criteria",
			Since:       "0.1.0",
		},
		ClaimProviderAttestation: {
			Name:        "Provider Attestation",
			Description: "Provider confirmed the accuracy of health data",
			Since:       "0.1.0",
		},
		ClaimLabVerification: {
			Name:        "Lab Verification",
			Description: "Lab results have been verified by authorized lab",
			Since:       "0.1.0",
		},
		ClaimAgeOver: {
			Name:        "Age Over",
			Description: "Proves subject is over specified age without revealing exact age",
			Since:       "0.1.0",
		},
	})
}

// CredentialStatus represents the status of a verifiable credential.
type CredentialStatus string

const (
	StatusActive  CredentialStatus = "active"  // Credential is valid and active
	StatusRevoked CredentialStatus = "revoked" // Credential has been revoked
	StatusExpired CredentialStatus = "expired" // Credential has expired
	StatusPending CredentialStatus = "pending" // Credential is pending issuance
)

// IsValid checks if the status is a known status.
func (s CredentialStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusRevoked, StatusExpired, StatusPending:
		return true
	}
	return false
}

// IsUsable returns true if the credential can be presented.
func (s CredentialStatus) IsUsable() bool {
	return s == StatusActive
}

// Credential represents a Verifiable Credential using SD-JWT format.
// This is the protocol-level representation - the actual SD-JWT encoding
// is handled by the builder.
type Credential struct {
	// ID is the unique identifier for the credential
	ID types.ID `json:"id"`

	// Issuer is the wallet address of the credential issuer
	Issuer types.WalletAddress `json:"issuer"`

	// Subject is the wallet address of the credential subject (patient)
	Subject types.WalletAddress `json:"subject"`

	// ClaimType identifies the type of claim
	ClaimType ClaimType `json:"claimType"`

	// Claims contains the actual claim data (key-value pairs)
	Claims map[string]any `json:"claims"`

	// Disclosures contains selective disclosure information
	// Only populated when presenting with disclosures
	Disclosures []Disclosure `json:"disclosures,omitempty"`

	// SourceEventIDs are the timeline event IDs that back this credential
	SourceEventIDs []types.ID `json:"sourceEventIds,omitempty"`

	// IssuedAt is when the credential was issued
	IssuedAt time.Time `json:"issuedAt"`

	// ExpiresAt is when the credential expires (optional)
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// Status is the current status of the credential
	Status CredentialStatus `json:"status"`

	// RevocationIndex is the index in the revocation list (if revocable)
	RevocationIndex *uint64 `json:"revocationIndex,omitempty"`

	// SchemaVersion is the protocol schema version
	SchemaVersion string `json:"schemaVersion"`
}

// Validate validates the credential structure.
func (c *Credential) Validate() error {
	var errs types.ValidationErrors

	if c.ID.IsEmpty() {
		errs.Add("id", "credential ID is required")
	}

	if c.Issuer.IsEmpty() {
		errs.Add("issuer", "issuer is required")
	}

	if c.Subject.IsEmpty() {
		errs.Add("subject", "subject is required")
	}

	if !c.ClaimType.IsValid() {
		errs.Add("claimType", "invalid claim type")
	}

	if len(c.Claims) == 0 {
		errs.Add("claims", "at least one claim is required")
	}

	if c.IssuedAt.IsZero() {
		errs.Add("issuedAt", "issuedAt is required")
	}

	if !c.Status.IsValid() {
		errs.Add("status", "invalid status")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// IsExpired checks if the credential has expired.
func (c *Credential) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// IsUsable checks if the credential can be presented.
func (c *Credential) IsUsable() bool {
	return c.Status.IsUsable() && !c.IsExpired()
}

// Disclosure represents a selective disclosure element.
// In SD-JWT, each disclosure is a base64-encoded JSON array: [salt, claim_name, claim_value]
type Disclosure struct {
	// Salt is the random salt used for this disclosure
	Salt string `json:"salt"`

	// Key is the claim name being disclosed
	Key string `json:"key"`

	// Value is the claim value being disclosed
	Value any `json:"value"`

	// Encoded is the base64url-encoded disclosure string
	Encoded string `json:"encoded,omitempty"`
}

// CredentialRequest represents a request for a verifiable credential.
type CredentialRequest struct {
	// RequestID is the unique identifier for this request
	RequestID types.ID `json:"requestId"`

	// Requester is the wallet address of the requester (usually the subject)
	Requester types.WalletAddress `json:"requester"`

	// ClaimType is the type of claim being requested
	ClaimType ClaimType `json:"claimType"`

	// ClaimCriteria contains the criteria for the claim
	ClaimCriteria map[string]any `json:"claimCriteria"`

	// SourceEventIDs are the event IDs to use as evidence
	SourceEventIDs []types.ID `json:"sourceEventIds"`

	// RequestedAt is when the request was made
	RequestedAt time.Time `json:"requestedAt"`
}

// Validate validates the credential request.
func (r *CredentialRequest) Validate() error {
	var errs types.ValidationErrors

	if r.RequestID.IsEmpty() {
		errs.Add("requestId", "request ID is required")
	}

	if r.Requester.IsEmpty() {
		errs.Add("requester", "requester is required")
	}

	if !r.ClaimType.IsValid() {
		errs.Add("claimType", "invalid claim type")
	}

	if len(r.SourceEventIDs) == 0 {
		errs.Add("sourceEventIds", "at least one source event is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
