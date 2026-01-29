package vc

import (
	"time"

	"github.com/google/uuid"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// CredentialBuilder provides a fluent interface for building Verifiable Credentials.
// Follows the Builder pattern used throughout the protocol layer.
type CredentialBuilder struct {
	cred *Credential
	errs types.ValidationErrors
}

// NewCredentialBuilder creates a new CredentialBuilder with default values.
func NewCredentialBuilder() *CredentialBuilder {
	return &CredentialBuilder{
		cred: &Credential{
			ID:            types.ID(uuid.New().String()),
			Claims:        make(map[string]any),
			Disclosures:   make([]Disclosure, 0),
			SourceEventIDs: make([]types.ID, 0),
			Status:        StatusPending,
			SchemaVersion: SchemaVersionVC,
		},
		errs: types.ValidationErrors{},
	}
}

// WithID sets the credential ID (optional - auto-generated if not set).
func (b *CredentialBuilder) WithID(id types.ID) *CredentialBuilder {
	if id.IsEmpty() {
		b.errs.Add("id", "ID cannot be empty")
	}
	b.cred.ID = id
	return b
}

// WithIssuer sets the issuer wallet address.
func (b *CredentialBuilder) WithIssuer(addr types.WalletAddress) *CredentialBuilder {
	if addr.IsEmpty() {
		b.errs.Add("issuer", "issuer address is required")
	}
	b.cred.Issuer = addr
	return b
}

// WithSubject sets the subject wallet address.
func (b *CredentialBuilder) WithSubject(addr types.WalletAddress) *CredentialBuilder {
	if addr.IsEmpty() {
		b.errs.Add("subject", "subject address is required")
	}
	b.cred.Subject = addr
	return b
}

// WithClaimType sets the claim type.
func (b *CredentialBuilder) WithClaimType(ct ClaimType) *CredentialBuilder {
	if !ct.IsValid() {
		b.errs.Add("claimType", "invalid claim type")
	}
	b.cred.ClaimType = ct
	return b
}

// AddClaim adds a claim key-value pair.
// If disclosed is true, the claim will be included in the SD-JWT disclosures.
func (b *CredentialBuilder) AddClaim(key string, value any, disclosed bool) *CredentialBuilder {
	if key == "" {
		b.errs.Add("claims", "claim key cannot be empty")
		return b
	}
	b.cred.Claims[key] = value

	// If this is a selective disclosure claim, add to disclosures
	if disclosed {
		d := Disclosure{
			Key:   key,
			Value: value,
			// Salt will be generated during SD-JWT encoding
		}
		b.cred.Disclosures = append(b.cred.Disclosures, d)
	}

	return b
}

// AddBloodworkClaim adds a bloodwork range claim.
func (b *CredentialBuilder) AddBloodworkClaim(claim *BloodworkRangeClaim) *CredentialBuilder {
	if err := claim.Validate(); err != nil {
		b.errs.Add("claims", "invalid bloodwork claim: "+err.Error())
		return b
	}
	for k, v := range claim.ToMap() {
		b.cred.Claims[k] = v
	}
	b.cred.ClaimType = ClaimBloodworkRange
	return b
}

// AddProtocolAdherenceClaim adds a protocol adherence claim.
func (b *CredentialBuilder) AddProtocolAdherenceClaim(claim *ProtocolAdherenceClaim) *CredentialBuilder {
	if err := claim.Validate(); err != nil {
		b.errs.Add("claims", "invalid protocol adherence claim: "+err.Error())
		return b
	}
	for k, v := range claim.ToMap() {
		b.cred.Claims[k] = v
	}
	b.cred.ClaimType = ClaimProtocolAdherence
	return b
}

// AddBiometricPercentileClaim adds a biometric percentile claim.
func (b *CredentialBuilder) AddBiometricPercentileClaim(claim *BiometricPercentileClaim) *CredentialBuilder {
	if err := claim.Validate(); err != nil {
		b.errs.Add("claims", "invalid biometric percentile claim: "+err.Error())
		return b
	}
	for k, v := range claim.ToMap() {
		b.cred.Claims[k] = v
	}
	b.cred.ClaimType = ClaimBiometricPercentile
	return b
}

// WithSourceEvents sets the source event IDs that back this credential.
func (b *CredentialBuilder) WithSourceEvents(eventIDs ...types.ID) *CredentialBuilder {
	if len(eventIDs) == 0 {
		b.errs.Add("sourceEventIds", "at least one source event is required")
		return b
	}
	b.cred.SourceEventIDs = eventIDs
	return b
}

// WithIssuedAt sets the issuance timestamp.
func (b *CredentialBuilder) WithIssuedAt(t time.Time) *CredentialBuilder {
	if t.IsZero() {
		b.errs.Add("issuedAt", "issuedAt cannot be zero")
	}
	b.cred.IssuedAt = t
	return b
}

// WithExpiresAt sets the expiration timestamp (optional).
func (b *CredentialBuilder) WithExpiresAt(t time.Time) *CredentialBuilder {
	b.cred.ExpiresAt = &t
	return b
}

// WithTTL sets the expiration as a duration from now.
func (b *CredentialBuilder) WithTTL(duration time.Duration) *CredentialBuilder {
	expiry := time.Now().Add(duration)
	b.cred.ExpiresAt = &expiry
	return b
}

// WithRevocationIndex sets the revocation list index for this credential.
func (b *CredentialBuilder) WithRevocationIndex(index uint64) *CredentialBuilder {
	b.cred.RevocationIndex = &index
	return b
}

// Build validates and returns the credential.
func (b *CredentialBuilder) Build() (*Credential, error) {
	// Set default issuedAt if not set
	if b.cred.IssuedAt.IsZero() {
		b.cred.IssuedAt = time.Now().UTC()
	}

	// Set status to active
	b.cred.Status = StatusActive

	// Check for accumulated errors
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	// Validate the final credential
	if err := b.cred.Validate(); err != nil {
		return nil, err
	}

	return b.cred, nil
}

// MustBuild is like Build but panics on error.
// Use only in tests or when errors are pre-validated.
func (b *CredentialBuilder) MustBuild() *Credential {
	cred, err := b.Build()
	if err != nil {
		panic("CredentialBuilder.MustBuild: " + err.Error())
	}
	return cred
}

// PresentationBuilder builds credential presentations with selective disclosure.
type PresentationBuilder struct {
	credential     *Credential
	disclosedKeys  map[string]bool
	errs           types.ValidationErrors
}

// NewPresentationBuilder creates a builder for presenting a credential.
func NewPresentationBuilder(cred *Credential) *PresentationBuilder {
	return &PresentationBuilder{
		credential:    cred,
		disclosedKeys: make(map[string]bool),
	}
}

// DiscloseKey marks a claim key to be disclosed in the presentation.
func (b *PresentationBuilder) DiscloseKey(key string) *PresentationBuilder {
	if _, exists := b.credential.Claims[key]; !exists {
		b.errs.Add("disclosedKeys", "claim key not found: "+key)
		return b
	}
	b.disclosedKeys[key] = true
	return b
}

// DiscloseAll marks all claims to be disclosed.
func (b *PresentationBuilder) DiscloseAll() *PresentationBuilder {
	for key := range b.credential.Claims {
		b.disclosedKeys[key] = true
	}
	return b
}

// Build creates a presentation with only the disclosed claims.
func (b *PresentationBuilder) Build() (*Credential, error) {
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	if !b.credential.IsUsable() {
		return nil, types.NewValidationError("credential", "credential is not usable (status: "+string(b.credential.Status)+")")
	}

	// Create a copy with only disclosed claims
	presentation := &Credential{
		ID:            b.credential.ID,
		Issuer:        b.credential.Issuer,
		Subject:       b.credential.Subject,
		ClaimType:     b.credential.ClaimType,
		Claims:        make(map[string]any),
		Disclosures:   make([]Disclosure, 0),
		IssuedAt:      b.credential.IssuedAt,
		ExpiresAt:     b.credential.ExpiresAt,
		Status:        b.credential.Status,
		SchemaVersion: b.credential.SchemaVersion,
	}

	// Only include disclosed claims
	for key, value := range b.credential.Claims {
		if b.disclosedKeys[key] {
			presentation.Claims[key] = value
			presentation.Disclosures = append(presentation.Disclosures, Disclosure{
				Key:   key,
				Value: value,
			})
		}
	}

	return presentation, nil
}
