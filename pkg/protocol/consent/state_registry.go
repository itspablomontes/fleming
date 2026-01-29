package consent

import (
	"sync"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

var (
	// defaultStateRegistry is the default registry for consent state types.
	defaultStateRegistry types.TypeRegistry[State]

	// stateRegistryOnce ensures the registry is initialized only once.
	stateRegistryOnce sync.Once
)

func init() {
	// Initialize default registry on package load
	stateRegistryOnce.Do(func() {
		defaultStateRegistry = types.NewTypeRegistry[State]()
		RegisterDefaultStates()
	})
}

// GetStateRegistry returns the default state type registry.
func GetStateRegistry() types.TypeRegistry[State] {
	return defaultStateRegistry
}

// RegisterState registers a custom state type at runtime.
// Note: State machine transitions remain explicit and must be updated separately.
func RegisterState(state State, metadata types.TypeMetadata) error {
	return defaultStateRegistry.Register(state, metadata)
}

// ValidStates returns all valid state types (backward compatibility).
func ValidStates() []State {
	return defaultStateRegistry.ValidTypes()
}

// RegisterDefaultStates registers all built-in state types.
func RegisterDefaultStates() {
	reg := defaultStateRegistry
	types.RegisterBatch(reg, map[State]types.TypeMetadata{
		StateRequested: {
			Name:        "Requested",
			Description: "Consent request pending approval",
			Since:       "0.1.0",
		},
		StateApproved: {
			Name:        "Approved",
			Description: "Consent grant approved and active",
			Since:       "0.1.0",
		},
		StateDenied: {
			Name:        "Denied",
			Description: "Consent request denied (terminal)",
			Since:       "0.1.0",
		},
		StateRevoked: {
			Name:        "Revoked",
			Description: "Consent grant revoked by grantor (terminal)",
			Since:       "0.1.0",
		},
		StateExpired: {
			Name:        "Expired",
			Description: "Consent grant expired (terminal)",
			Since:       "0.1.0",
		},
		StateSuspended: {
			Name:        "Suspended",
			Description: "Consent grant temporarily suspended (can be resumed)",
			Since:       "0.1.0",
		},
	})
}
