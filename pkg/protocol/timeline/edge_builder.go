package timeline

import (
	"fmt"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// EdgeBuilder provides a fluent API for building Edge instances with validation.
type EdgeBuilder struct {
	edge *Edge
	errs types.ValidationErrors
}

// NewEdgeBuilder creates a new EdgeBuilder instance.
func NewEdgeBuilder() *EdgeBuilder {
	return &EdgeBuilder{
		edge: &Edge{
			Metadata: types.NewMetadata(),
		},
	}
}

// WithID sets the edge ID.
func (b *EdgeBuilder) WithID(id types.ID) *EdgeBuilder {
	b.edge.ID = id
	return b
}

// WithFromID sets the source event ID.
func (b *EdgeBuilder) WithFromID(fromID types.ID) *EdgeBuilder {
	b.edge.FromID = fromID
	return b
}

// WithToID sets the target event ID.
func (b *EdgeBuilder) WithToID(toID types.ID) *EdgeBuilder {
	b.edge.ToID = toID
	return b
}

// WithType sets the relationship type.
func (b *EdgeBuilder) WithType(relType RelationshipType) *EdgeBuilder {
	b.edge.Type = relType
	return b
}

// WithMetadata sets the metadata map.
func (b *EdgeBuilder) WithMetadata(metadata types.Metadata) *EdgeBuilder {
	b.edge.Metadata = metadata
	return b
}

// SetMetadata sets a metadata key-value pair.
func (b *EdgeBuilder) SetMetadata(key string, value any) *EdgeBuilder {
	b.edge.Metadata = b.edge.Metadata.Set(key, value)
	return b
}

// Build validates and returns the Edge, or returns an error if validation fails.
func (b *EdgeBuilder) Build() (*Edge, error) {
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	if err := b.edge.Validate(); err != nil {
		return nil, fmt.Errorf("edge validation failed: %w", err)
	}

	return b.edge, nil
}
