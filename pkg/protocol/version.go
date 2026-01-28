package protocol

// Version constants for the Fleming Protocol.
const (
	VersionMajor = 1
	VersionMinor = 0
	VersionPatch = 0
	VersionString = "1.0.0"
)

// Version represents a protocol version with major, minor, patch, and schema version.
type Version struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Patch int    `json:"patch"`
	Schema string `json:"schema,omitempty"` // Schema version (e.g., "timeline.v1", "consent.v1")
}

// CurrentVersion returns the current protocol version.
func CurrentVersion() Version {
	return Version{
		Major: VersionMajor,
		Minor: VersionMinor,
		Patch: VersionPatch,
	}
}

// String returns the version as a semantic version string (e.g., "1.0.0").
func (v Version) String() string {
	return VersionString
}

// SchemaVersion returns the schema version string for a given component.
func SchemaVersion(component string) string {
	return component + ".v1"
}

// Schema versions for each protocol component.
const (
	SchemaVersionTimeline    = "timeline.v1"
	SchemaVersionConsent     = "consent.v1"
	SchemaVersionAudit       = "audit.v1"
	SchemaVersionIdentity    = "identity.v1"
	SchemaVersionVC          = "vc.v1"
	SchemaVersionAttestation = "attestation.v1"
)
