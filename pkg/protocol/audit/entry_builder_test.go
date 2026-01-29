package audit

import (
	"testing"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestEntryBuilder_WithActor(t *testing.T) {
	validActor, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	builder := NewEntryBuilder()

	builder.WithActor(validActor)
	if builder.entry.Actor != validActor {
		t.Error("WithActor() did not set actor")
	}
}

func TestEntryBuilder_WithAction(t *testing.T) {
	builder := NewEntryBuilder()

	builder.WithAction(ActionCreate)
	if builder.entry.Action != ActionCreate {
		t.Error("WithAction() did not set action")
	}
}

func TestEntryBuilder_WithResourceType(t *testing.T) {
	builder := NewEntryBuilder()

	builder.WithResourceType(ResourceEvent)
	if builder.entry.ResourceType != ResourceEvent {
		t.Error("WithResourceType() did not set resource type")
	}
}

func TestEntryBuilder_WithResourceID(t *testing.T) {
	validID, _ := types.NewID("event-1")
	builder := NewEntryBuilder()

	builder.WithResourceID(validID)
	if builder.entry.ResourceID != validID {
		t.Error("WithResourceID() did not set resource ID")
	}
}

func TestEntryBuilder_WithPreviousHash(t *testing.T) {
	builder := NewEntryBuilder()

	builder.WithPreviousHash("hash123")
	if builder.entry.PreviousHash != "hash123" {
		t.Error("WithPreviousHash() did not set previous hash")
	}
}

func TestEntryBuilder_Build(t *testing.T) {
	validActor, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validID, _ := types.NewID("event-1")

	tests := []struct {
		name    string
		builder func() *EntryBuilder
		wantErr bool
	}{
		{
			name: "valid entry",
			builder: func() *EntryBuilder {
				return NewEntryBuilder().
					WithActor(validActor).
					WithAction(ActionCreate).
					WithResourceType(ResourceEvent).
					WithResourceID(validID)
			},
			wantErr: false,
		},
		{
			name: "missing actor",
			builder: func() *EntryBuilder {
				return NewEntryBuilder().
					WithAction(ActionCreate).
					WithResourceType(ResourceEvent).
					WithResourceID(validID)
			},
			wantErr: true,
		},
		{
			name: "missing action",
			builder: func() *EntryBuilder {
				return NewEntryBuilder().
					WithActor(validActor).
					WithResourceType(ResourceEvent).
					WithResourceID(validID)
			},
			wantErr: true,
		},
		{
			name: "missing resource ID",
			builder: func() *EntryBuilder {
				return NewEntryBuilder().
					WithActor(validActor).
					WithAction(ActionCreate).
					WithResourceType(ResourceEvent)
			},
			wantErr: true,
		},
		{
			name: "invalid action",
			builder: func() *EntryBuilder {
				return NewEntryBuilder().
					WithActor(validActor).
					WithAction("invalid").
					WithResourceType(ResourceEvent).
					WithResourceID(validID)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.builder()
			entry, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if entry == nil {
					t.Error("Build() returned nil for valid entry")
				}
				if entry.Hash == "" {
					t.Error("Build() should compute hash")
				}
				if entry.Timestamp.IsZero() {
					t.Error("Build() should set timestamp if not provided")
				}
			}
		})
	}
}

func TestEntryBuilder_SetMetadata(t *testing.T) {
	builder := NewEntryBuilder()
	builder.SetMetadata("key", "value")

	val := builder.entry.Metadata.GetString("key")
	if val != "value" {
		t.Errorf("SetMetadata() value = %v, want value", val)
	}
}
