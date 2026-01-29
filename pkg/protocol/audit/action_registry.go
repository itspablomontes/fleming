package audit

import (
	"sync"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

var (
	// defaultActionRegistry is the default registry for audit action types.
	defaultActionRegistry types.TypeRegistry[Action]

	// defaultResourceTypeRegistry is the default registry for resource types.
	defaultResourceTypeRegistry types.TypeRegistry[ResourceType]

	// actionRegistryOnce ensures the registry is initialized only once.
	actionRegistryOnce sync.Once

	// resourceTypeRegistryOnce ensures the registry is initialized only once.
	resourceTypeRegistryOnce sync.Once
)

func init() {
	// Initialize default registries on package load
	actionRegistryOnce.Do(func() {
		defaultActionRegistry = types.NewTypeRegistry[Action]()
		RegisterDefaultActions()
	})

	resourceTypeRegistryOnce.Do(func() {
		defaultResourceTypeRegistry = types.NewTypeRegistry[ResourceType]()
		RegisterDefaultResourceTypes()
	})
}

// GetActionRegistry returns the default action type registry.
func GetActionRegistry() types.TypeRegistry[Action] {
	return defaultActionRegistry
}

// GetResourceTypeRegistry returns the default resource type registry.
func GetResourceTypeRegistry() types.TypeRegistry[ResourceType] {
	return defaultResourceTypeRegistry
}

// RegisterAction registers a custom action type at runtime.
func RegisterAction(action Action, metadata types.TypeMetadata) error {
	return defaultActionRegistry.Register(action, metadata)
}

// RegisterResourceType registers a custom resource type at runtime.
func RegisterResourceType(resourceType ResourceType, metadata types.TypeMetadata) error {
	return defaultResourceTypeRegistry.Register(resourceType, metadata)
}

// ValidActions returns all valid action types (backward compatibility).
func ValidActions() []Action {
	return defaultActionRegistry.ValidTypes()
}

// ValidResourceTypes returns all valid resource types (backward compatibility).
func ValidResourceTypes() []ResourceType {
	return defaultResourceTypeRegistry.ValidTypes()
}

// RegisterDefaultActions registers all built-in action types.
func RegisterDefaultActions() {
	reg := defaultActionRegistry
	types.RegisterBatch(reg, map[Action]types.TypeMetadata{
		// CRUD operations
		ActionCreate: {
			Name:        "Create",
			Description: "Create a resource",
			Since:       "0.1.0",
		},
		ActionRead: {
			Name:        "Read",
			Description: "Read a resource",
			Since:       "0.1.0",
		},
		ActionUpdate: {
			Name:        "Update",
			Description: "Update a resource",
			Since:       "0.1.0",
		},
		ActionDelete: {
			Name:        "Delete",
			Description: "Delete a resource",
			Since:       "0.1.0",
		},

		// Consent operations
		ActionConsentRequest: {
			Name:        "Consent Request",
			Description: "Request consent grant",
			Since:       "0.1.0",
		},
		ActionConsentApprove: {
			Name:        "Consent Approve",
			Description: "Approve consent grant",
			Since:       "0.1.0",
		},
		ActionConsentDeny: {
			Name:        "Consent Deny",
			Description: "Deny consent grant",
			Since:       "0.1.0",
		},
		ActionConsentRevoke: {
			Name:        "Consent Revoke",
			Description: "Revoke consent grant",
			Since:       "0.1.0",
		},
		ActionConsentExpire: {
			Name:        "Consent Expire",
			Description: "Consent grant expired",
			Since:       "0.1.0",
		},
		ActionConsentSuspend: {
			Name:        "Consent Suspend",
			Description: "Temporarily suspend consent grant",
			Since:       "0.1.0",
		},
		ActionConsentResume: {
			Name:        "Consent Resume",
			Description: "Resume suspended consent grant",
			Since:       "0.1.0",
		},

		// Authentication
		ActionLogin: {
			Name:        "Login",
			Description: "User login",
			Since:       "0.1.0",
		},
		ActionLogout: {
			Name:        "Logout",
			Description: "User logout",
			Since:       "0.1.0",
		},

		// File operations
		ActionUpload: {
			Name:        "File Upload",
			Description: "Upload a file",
			Since:       "0.1.0",
		},
		ActionDownload: {
			Name:        "File Download",
			Description: "Download a file",
			Since:       "0.1.0",
		},
		ActionShare: {
			Name:        "File Share",
			Description: "Share file access",
			Since:       "0.1.0",
		},

		// Verifiable Credentials
		ActionVCIssue: {
			Name:        "VC Issue",
			Description: "Issue a verifiable credential",
			Since:       "0.1.0",
		},
		ActionVCRevoke: {
			Name:        "VC Revoke",
			Description: "Revoke a verifiable credential",
			Since:       "0.1.0",
		},
		ActionVCVerify: {
			Name:        "VC Verify",
			Description: "Verify a verifiable credential",
			Since:       "0.1.0",
		},
		ActionVCPresent: {
			Name:        "VC Present",
			Description: "Present a credential with selective disclosure",
			Since:       "0.1.0",
		},

		// Zero-Knowledge Proofs
		ActionZKGenerate: {
			Name:        "ZK Generate",
			Description: "Generate a zero-knowledge proof",
			Since:       "0.1.0",
		},
		ActionZKVerify: {
			Name:        "ZK Verify",
			Description: "Verify a zero-knowledge proof",
			Since:       "0.1.0",
		},

		// Attestation
		ActionCosign: {
			Name:        "Attestation Cosign",
			Description: "Provider co-signs an event",
			Since:       "0.1.0",
		},
		ActionAttest: {
			Name:        "Attestation Attest",
			Description: "Provider attests to accuracy of an event",
			Since:       "0.1.0",
		},
	})
}

// RegisterDefaultResourceTypes registers all built-in resource types.
func RegisterDefaultResourceTypes() {
	reg := defaultResourceTypeRegistry
	types.RegisterBatch(reg, map[ResourceType]types.TypeMetadata{
		// Core resources
		ResourceEvent: {
			Name:        "Event",
			Description: "Timeline event",
			Since:       "0.1.0",
		},
		ResourceFile: {
			Name:        "File",
			Description: "File attachment",
			Since:       "0.1.0",
		},
		ResourceConsent: {
			Name:        "Consent",
			Description: "Consent grant",
			Since:       "0.1.0",
		},
		ResourceSession: {
			Name:        "Session",
			Description: "User session",
			Since:       "0.1.0",
		},

		// Verifiable Credentials
		ResourceVC: {
			Name:        "Verifiable Credential",
			Description: "Verifiable credential (SD-JWT)",
			Since:       "0.1.0",
		},

		// Zero-Knowledge Proofs
		ResourceZKProof: {
			Name:        "ZK Proof",
			Description: "Zero-knowledge proof",
			Since:       "0.1.0",
		},

		// Attestation
		ResourceAttestation: {
			Name:        "Attestation",
			Description: "Provider attestation",
			Since:       "0.1.0",
		},
	})
}
