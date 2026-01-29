package timeline

import (
	"sync"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

var (
	// defaultRelationshipTypeRegistry is the default registry for relationship types.
	defaultRelationshipTypeRegistry types.TypeRegistry[RelationshipType]

	// relationshipTypeRegistryOnce ensures the registry is initialized only once.
	relationshipTypeRegistryOnce sync.Once
)

func init() {
	// Initialize default registry on package load
	relationshipTypeRegistryOnce.Do(func() {
		defaultRelationshipTypeRegistry = types.NewTypeRegistry[RelationshipType]()
		RegisterDefaultRelationshipTypes()
	})
}

// GetRelationshipTypeRegistry returns the default relationship type registry.
func GetRelationshipTypeRegistry() types.TypeRegistry[RelationshipType] {
	return defaultRelationshipTypeRegistry
}

// RegisterRelationshipType registers a custom relationship type at runtime.
func RegisterRelationshipType(relType RelationshipType, metadata types.TypeMetadata) error {
	return defaultRelationshipTypeRegistry.Register(relType, metadata)
}

// ValidRelationshipTypes returns all valid relationship types (backward compatibility).
func ValidRelationshipTypes() []RelationshipType {
	return defaultRelationshipTypeRegistry.ValidTypes()
}

// RegisterDefaultRelationshipTypes registers all built-in relationship types.
func RegisterDefaultRelationshipTypes() {
	reg := defaultRelationshipTypeRegistry
	types.RegisterBatch(reg, map[RelationshipType]types.TypeMetadata{
		// Core relationship types
		RelResultedIn: {
			Name:        "Resulted In",
			Description: "Event A resulted in event B",
			Since:       "0.1.0",
		},
		RelLeadTo: {
			Name:        "Lead To",
			Description: "Event A led to event B",
			Since:       "0.1.0",
		},
		RelRequestedBy: {
			Name:        "Requested By",
			Description: "Event A was requested by event B",
			Since:       "0.1.0",
		},
		RelSupports: {
			Name:        "Supports",
			Description: "Event A supports event B",
			Since:       "0.1.0",
		},
		RelFollowsUp: {
			Name:        "Follows Up",
			Description: "Event A follows up on event B",
			Since:       "0.1.0",
		},
		RelContradicts: {
			Name:        "Contradicts",
			Description: "Event A contradicts event B",
			Since:       "0.1.0",
		},
		RelAttachedTo: {
			Name:        "Attached To",
			Description: "Event A is attached to event B",
			Since:       "0.1.0",
		},
		RelReplaces: {
			Name:        "Replaces",
			Description: "Event A replaces event B (append-only correction)",
			Since:       "0.1.0",
		},
		RelCausedBy: {
			Name:        "Caused By",
			Description: "Event A was caused by event B",
			Since:       "0.1.0",
		},

		// Provider attestation (CRITICAL for cosigning feature)
		RelCosignedBy: {
			Name:        "Co-signed By",
			Description: "Event was co-signed by a healthcare provider",
			Since:       "0.1.0",
		},
		RelAttestedBy: {
			Name:        "Attested By",
			Description: "Event accuracy was attested by a provider",
			Since:       "0.1.0",
		},

		// Medical relationships
		RelTreats: {
			Name:        "Treats",
			Description: "Treatment relationship (e.g., medication treats condition)",
			Since:       "0.1.0",
		},
		RelMonitors: {
			Name:        "Monitors",
			Description: "Monitoring relationship (e.g., lab test monitors medication effect)",
			Since:       "0.1.0",
		},
		RelContraindicated: {
			Name:        "Contraindicated",
			Description: "Contraindication relationship between events",
			Since:       "0.1.0",
		},
		RelDerivedFrom: {
			Name:        "Derived From",
			Description: "Data derived from another event",
			Since:       "0.1.0",
		},
		RelPartOf: {
			Name:        "Part Of",
			Description: "Event is part of a larger entity or protocol",
			Since:       "0.1.0",
		},

		// AI/Suggestions
		RelSuggestedBy: {
			Name:        "Suggested By",
			Description: "Relationship was suggested by AI or rule engine",
			Since:       "0.1.0",
		},
	})
}
