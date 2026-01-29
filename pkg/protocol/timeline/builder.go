package timeline

import (
	"fmt"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

const SchemaVersionTimeline = protocol.SchemaVersionTimeline

// EventBuilder provides a fluent API for building Event instances with validation.
type EventBuilder struct {
	event *Event
	errs  types.ValidationErrors
}

// NewEventBuilder creates a new EventBuilder instance.
func NewEventBuilder() *EventBuilder {
	return &EventBuilder{
		event: &Event{
			Metadata:      types.NewMetadata(),
			Codes:         make(types.Codes, 0),
			SchemaVersion: SchemaVersionTimeline,
		},
	}
}

// WithID sets the event ID.
func (b *EventBuilder) WithID(id types.ID) *EventBuilder {
	b.event.ID = id
	return b
}

// WithPatientID sets the patient wallet address.
func (b *EventBuilder) WithPatientID(patientID types.WalletAddress) *EventBuilder {
	b.event.PatientID = patientID
	return b
}

// WithType sets the event type.
func (b *EventBuilder) WithType(eventType EventType) *EventBuilder {
	b.event.Type = eventType
	return b
}

// WithTitle sets the event title.
func (b *EventBuilder) WithTitle(title string) *EventBuilder {
	b.event.Title = title
	return b
}

// WithDescription sets the event description.
func (b *EventBuilder) WithDescription(description string) *EventBuilder {
	b.event.Description = description
	return b
}

// WithProvider sets the provider name.
func (b *EventBuilder) WithProvider(provider string) *EventBuilder {
	b.event.Provider = provider
	return b
}

// WithTimestamp sets the event timestamp.
func (b *EventBuilder) WithTimestamp(timestamp time.Time) *EventBuilder {
	b.event.Timestamp = timestamp
	return b
}

// WithCodes sets the medical codes.
func (b *EventBuilder) WithCodes(codes types.Codes) *EventBuilder {
	b.event.Codes = codes
	return b
}

// AddCode adds a single medical code.
func (b *EventBuilder) AddCode(code types.Code) *EventBuilder {
	if err := code.Validate(); err != nil {
		b.errs.Add("codes", fmt.Sprintf("invalid code: %v", err))
		return b
	}
	b.event.Codes = append(b.event.Codes, code)
	return b
}

// WithMetadata sets the metadata map.
func (b *EventBuilder) WithMetadata(metadata types.Metadata) *EventBuilder {
	b.event.Metadata = metadata
	return b
}

// SetMetadata sets a metadata key-value pair.
func (b *EventBuilder) SetMetadata(key string, value any) *EventBuilder {
	b.event.Metadata = b.event.Metadata.Set(key, value)
	return b
}

// WithCreatedAt sets the creation timestamp.
func (b *EventBuilder) WithCreatedAt(createdAt time.Time) *EventBuilder {
	b.event.CreatedAt = createdAt
	return b
}

// WithUpdatedAt sets the update timestamp.
func (b *EventBuilder) WithUpdatedAt(updatedAt time.Time) *EventBuilder {
	b.event.UpdatedAt = updatedAt
	return b
}

// Build validates and returns the Event, or returns an error if validation fails.
func (b *EventBuilder) Build() (*Event, error) {
	if b.errs.HasErrors() {
		return nil, b.errs
	}

	if err := b.event.Validate(); err != nil {
		return nil, fmt.Errorf("event validation failed: %w", err)
	}

	// Set timestamps if not provided
	if b.event.CreatedAt.IsZero() {
		b.event.CreatedAt = time.Now().UTC()
	}
	if b.event.UpdatedAt.IsZero() {
		b.event.UpdatedAt = b.event.CreatedAt
	}

	return b.event, nil
}
