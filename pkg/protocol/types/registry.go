package types

import (
	"fmt"
	"sync"
)

// TypeMetadata contains metadata about a registered type.
type TypeMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Deprecated  bool   `json:"deprecated"`
	Since       string `json:"since"` // Version when added (e.g., "1.0.0")
}

// TypeRegistry is a generic interface for type registries that allow runtime registration
// and validation of enum-like types.
type TypeRegistry[T comparable] interface {
	// Register adds a new type value with metadata to the registry.
	Register(value T, metadata TypeMetadata) error

	// IsValid checks if a type value is registered and valid.
	IsValid(value T) bool

	// ValidTypes returns all registered type values.
	ValidTypes() []T

	// GetMetadata retrieves metadata for a type value.
	GetMetadata(value T) (TypeMetadata, bool)
}

// registry is a thread-safe implementation of TypeRegistry.
type registry[T comparable] struct {
	mu        sync.RWMutex
	types     map[T]TypeMetadata
	typeOrder []T // Preserve registration order
}

// NewTypeRegistry creates a new thread-safe type registry.
func NewTypeRegistry[T comparable]() TypeRegistry[T] {
	return &registry[T]{
		types:     make(map[T]TypeMetadata),
		typeOrder: make([]T, 0),
	}
}

// Register adds a type value with metadata.
func (r *registry[T]) Register(value T, metadata TypeMetadata) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.types[value]; exists {
		return fmt.Errorf("type %v is already registered", value)
	}

	if metadata.Name == "" {
		metadata.Name = fmt.Sprintf("%v", value)
	}

	r.types[value] = metadata
	r.typeOrder = append(r.typeOrder, value)
	return nil
}

// IsValid checks if a type value is registered.
func (r *registry[T]) IsValid(value T) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, exists := r.types[value]
	return exists && !meta.Deprecated
}

// ValidTypes returns all registered type values in registration order.
func (r *registry[T]) ValidTypes() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]T, len(r.typeOrder))
	copy(result, r.typeOrder)
	return result
}

// GetMetadata retrieves metadata for a type value.
func (r *registry[T]) GetMetadata(value T) (TypeMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, exists := r.types[value]
	return meta, exists
}

// RegisterBatch registers multiple types at once.
func RegisterBatch[T comparable](reg TypeRegistry[T], types map[T]TypeMetadata) error {
	for value, metadata := range types {
		if err := reg.Register(value, metadata); err != nil {
			return fmt.Errorf("failed to register %v: %w", value, err)
		}
	}
	return nil
}
