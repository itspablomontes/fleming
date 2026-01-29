package attestation

import (
	"time"

	"github.com/google/uuid"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// AttestationBuilder provides a fluent interface for building Attestations.
// Follows the Builder pattern used throughout the protocol layer.
type AttestationBuilder struct {
	att  *Attestation
	errs types.ValidationErrors
}

// NewAttestationBuilder creates a new AttestationBuilder with default values.
func NewAttestationBuilder() *AttestationBuilder {
	return &AttestationBuilder{
		att: &Attestation{
			ID:            types.ID(uuid.New().String()),
			Status:        StatusPendingAttestation,
			Metadata:      types.NewMetadata(),
			SchemaVersion: SchemaVersionAttestation,
		},
		errs: types.ValidationErrors{},
	}
}

// WithID sets the attestation ID (optional - auto-generated if not set).
func (b *AttestationBuilder) WithID(id types.ID) *AttestationBuilder {
	if id.IsEmpty() {
		b.errs.Add("id", "ID cannot be empty")
	}
	b.att.ID = id
	return b
}

// WithEventID sets the event being attested.
func (b *AttestationBuilder) WithEventID(eventID types.ID) *AttestationBuilder {
	if eventID.IsEmpty() {
		b.errs.Add("eventId", "event ID is required")
	}
	b.att.EventID = eventID
	return b
}

// WithEventHash sets the hash of the event being attested.
func (b *AttestationBuilder) WithEventHash(hash string) *AttestationBuilder {
	if hash == "" {
		b.errs.Add("eventHash", "event hash is required")
	}
	b.att.EventHash = hash
	return b
}

// WithAttester sets the attesting provider's wallet address.
func (b *AttestationBuilder) WithAttester(addr types.WalletAddress) *AttestationBuilder {
	if addr.IsEmpty() {
		b.errs.Add("attester", "attester address is required")
	}
	b.att.Attester = addr
	return b
}

// WithAttesterCredentials sets the provider's credentials.
func (b *AttestationBuilder) WithAttesterCredentials(creds *ProviderCredentials) *AttestationBuilder {
	b.att.AttesterCredentials = creds
	return b
}

// WithType sets the attestation type.
func (b *AttestationBuilder) WithType(at AttestationType) *AttestationBuilder {
	if !at.IsValid() {
		b.errs.Add("type", "invalid attestation type")
	}
	b.att.Type = at
	return b
}

// WithNotes sets optional notes from the provider.
func (b *AttestationBuilder) WithNotes(notes string) *AttestationBuilder {
	b.att.Notes = notes
	return b
}

// WithTimestamp sets the attestation timestamp.
func (b *AttestationBuilder) WithTimestamp(t time.Time) *AttestationBuilder {
	if t.IsZero() {
		b.errs.Add("timestamp", "timestamp cannot be zero")
	}
	b.att.Timestamp = t
	return b
}

// WithExpiresAt sets the expiration timestamp (optional).
func (b *AttestationBuilder) WithExpiresAt(t time.Time) *AttestationBuilder {
	b.att.ExpiresAt = &t
	return b
}

// WithTTL sets the expiration as a duration from now.
func (b *AttestationBuilder) WithTTL(duration time.Duration) *AttestationBuilder {
	expiry := time.Now().Add(duration)
	b.att.ExpiresAt = &expiry
	return b
}

// WithSignature sets the ECDSA signature.
func (b *AttestationBuilder) WithSignature(signature string, algorithm string) *AttestationBuilder {
	b.att.Signature = signature
	b.att.SignatureAlgorithm = algorithm
	return b
}

// WithMetadata adds metadata to the attestation.
func (b *AttestationBuilder) WithMetadata(key string, value any) *AttestationBuilder {
	b.att.Metadata = b.att.Metadata.Set(key, value)
	return b
}

// Build validates and returns the attestation.
// The attestation is built in Pending status until signed.
func (b *AttestationBuilder) Build() (*Attestation, error) {
	// Set default timestamp if not set
	if b.att.Timestamp.IsZero() {
		b.att.Timestamp = time.Now().UTC()
	}

	// Check for accumulated errors
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	// Validate the final attestation
	if err := b.att.Validate(); err != nil {
		return nil, err
	}

	return b.att, nil
}

// BuildSigned validates, signs, and returns the attestation.
// Sets the status to Active after signing.
func (b *AttestationBuilder) BuildSigned(signature string, algorithm string) (*Attestation, error) {
	b.att.Signature = signature
	b.att.SignatureAlgorithm = algorithm
	b.att.Status = StatusActiveAttestation

	return b.Build()
}

// MustBuild is like Build but panics on error.
// Use only in tests or when errors are pre-validated.
func (b *AttestationBuilder) MustBuild() *Attestation {
	att, err := b.Build()
	if err != nil {
		panic("AttestationBuilder.MustBuild: " + err.Error())
	}
	return att
}

// ProviderCredentialsBuilder provides a fluent interface for building ProviderCredentials.
type ProviderCredentialsBuilder struct {
	creds *ProviderCredentials
}

// NewProviderCredentialsBuilder creates a new ProviderCredentialsBuilder.
func NewProviderCredentialsBuilder() *ProviderCredentialsBuilder {
	return &ProviderCredentialsBuilder{
		creds: &ProviderCredentials{},
	}
}

// WithName sets the provider's name.
func (b *ProviderCredentialsBuilder) WithName(name string) *ProviderCredentialsBuilder {
	b.creds.Name = name
	return b
}

// WithLicense sets the provider's license information.
func (b *ProviderCredentialsBuilder) WithLicense(number, licenseType string) *ProviderCredentialsBuilder {
	b.creds.LicenseNumber = number
	b.creds.LicenseType = licenseType
	return b
}

// WithSpecialty sets the provider's specialty.
func (b *ProviderCredentialsBuilder) WithSpecialty(specialty string) *ProviderCredentialsBuilder {
	b.creds.Specialty = specialty
	return b
}

// WithOrganization sets the provider's organization.
func (b *ProviderCredentialsBuilder) WithOrganization(org string) *ProviderCredentialsBuilder {
	b.creds.Organization = org
	return b
}

// WithNPI sets the National Provider Identifier.
func (b *ProviderCredentialsBuilder) WithNPI(npi string) *ProviderCredentialsBuilder {
	b.creds.NPI = npi
	return b
}

// Build returns the provider credentials.
func (b *ProviderCredentialsBuilder) Build() *ProviderCredentials {
	return b.creds
}

// AttestationRequestBuilder provides a fluent interface for building attestation requests.
type AttestationRequestBuilder struct {
	req  *AttestationRequest
	errs types.ValidationErrors
}

// NewAttestationRequestBuilder creates a new AttestationRequestBuilder.
func NewAttestationRequestBuilder() *AttestationRequestBuilder {
	return &AttestationRequestBuilder{
		req: &AttestationRequest{
			RequestID:   types.ID(uuid.New().String()),
			RequestedAt: time.Now().UTC(),
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour).UTC(), // Default 7 days
		},
	}
}

// WithEventID sets the event to be attested.
func (b *AttestationRequestBuilder) WithEventID(eventID types.ID) *AttestationRequestBuilder {
	if eventID.IsEmpty() {
		b.errs.Add("eventId", "event ID is required")
	}
	b.req.EventID = eventID
	return b
}

// WithRequester sets the requester's wallet address.
func (b *AttestationRequestBuilder) WithRequester(addr types.WalletAddress) *AttestationRequestBuilder {
	if addr.IsEmpty() {
		b.errs.Add("requester", "requester is required")
	}
	b.req.Requester = addr
	return b
}

// WithTargetAttester sets the target provider (optional).
func (b *AttestationRequestBuilder) WithTargetAttester(addr types.WalletAddress) *AttestationRequestBuilder {
	b.req.TargetAttester = addr
	return b
}

// WithRequestedType sets the type of attestation being requested.
func (b *AttestationRequestBuilder) WithRequestedType(at AttestationType) *AttestationRequestBuilder {
	if !at.IsValid() {
		b.errs.Add("requestedType", "invalid attestation type")
	}
	b.req.RequestedType = at
	return b
}

// WithExpiresAt sets when the request expires.
func (b *AttestationRequestBuilder) WithExpiresAt(t time.Time) *AttestationRequestBuilder {
	b.req.ExpiresAt = t
	return b
}

// WithMessage sets an optional message to the provider.
func (b *AttestationRequestBuilder) WithMessage(msg string) *AttestationRequestBuilder {
	b.req.Message = msg
	return b
}

// Build validates and returns the attestation request.
func (b *AttestationRequestBuilder) Build() (*AttestationRequest, error) {
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	if err := b.req.Validate(); err != nil {
		return nil, err
	}

	return b.req, nil
}
