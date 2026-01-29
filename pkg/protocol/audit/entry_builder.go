package audit

import (
	"fmt"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// EntryBuilder provides a fluent API for building Entry instances with validation.
type EntryBuilder struct {
	entry *Entry
	errs  types.ValidationErrors
}

// NewEntryBuilder creates a new EntryBuilder instance.
func NewEntryBuilder() *EntryBuilder {
	return &EntryBuilder{
		entry: &Entry{
			Metadata:      types.NewMetadata(),
			SchemaVersion: SchemaVersionAudit,
		},
	}
}

// WithID sets the entry ID.
func (b *EntryBuilder) WithID(id types.ID) *EntryBuilder {
	b.entry.ID = id
	return b
}

// WithActor sets the actor wallet address.
func (b *EntryBuilder) WithActor(actor types.WalletAddress) *EntryBuilder {
	b.entry.Actor = actor
	return b
}

// WithAction sets the action.
func (b *EntryBuilder) WithAction(action Action) *EntryBuilder {
	b.entry.Action = action
	return b
}

// WithResourceType sets the resource type.
func (b *EntryBuilder) WithResourceType(resourceType ResourceType) *EntryBuilder {
	b.entry.ResourceType = resourceType
	return b
}

// WithResourceID sets the resource ID.
func (b *EntryBuilder) WithResourceID(resourceID types.ID) *EntryBuilder {
	b.entry.ResourceID = resourceID
	return b
}

// WithTimestamp sets the timestamp.
func (b *EntryBuilder) WithTimestamp(timestamp time.Time) *EntryBuilder {
	b.entry.Timestamp = timestamp
	return b
}

// WithPreviousHash sets the previous hash for chaining.
func (b *EntryBuilder) WithPreviousHash(previousHash string) *EntryBuilder {
	b.entry.PreviousHash = previousHash
	return b
}

// WithMetadata sets the metadata map.
func (b *EntryBuilder) WithMetadata(metadata types.Metadata) *EntryBuilder {
	b.entry.Metadata = metadata
	return b
}

// SetMetadata sets a metadata key-value pair.
func (b *EntryBuilder) SetMetadata(key string, value any) *EntryBuilder {
	b.entry.Metadata = b.entry.Metadata.Set(key, value)
	return b
}

// Build validates and returns the Entry, computing the hash automatically.
func (b *EntryBuilder) Build() (*Entry, error) {
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	if b.entry.Timestamp.IsZero() {
		b.entry.Timestamp = time.Now().UTC()
	}

	if err := b.entry.Validate(); err != nil {
		return nil, fmt.Errorf("entry validation failed: %w", err)
	}

	b.entry.SetHash()

	return b.entry, nil
}
