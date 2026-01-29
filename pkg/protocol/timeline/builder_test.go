package timeline

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestEventBuilder_WithPatientID(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	builder := NewEventBuilder()

	builder.WithPatientID(validAddr)
	if builder.event.PatientID != validAddr {
		t.Error("WithPatientID() did not set patient ID")
	}
}

func TestEventBuilder_WithType(t *testing.T) {
	builder := NewEventBuilder()

	builder.WithType(EventLabResult)
	if builder.event.Type != EventLabResult {
		t.Error("WithType() did not set event type")
	}
}

func TestEventBuilder_AddCode(t *testing.T) {
	builder := NewEventBuilder()

	validCode, _ := types.NewCode(types.CodingICD10, "E11.9")
	builder.AddCode(validCode)

	if len(builder.event.Codes) != 1 {
		t.Errorf("AddCode() expected 1 code, got %d", len(builder.event.Codes))
	}

	// Invalid code should add error
	invalidCode := types.Code{System: "invalid", Value: "code"}
	builder2 := NewEventBuilder()
	builder2.AddCode(invalidCode)
	if !builder2.errs.HasErrors() {
		t.Error("AddCode() with invalid code should add error")
	}
}

func TestEventBuilder_Build(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name    string
		builder func() *EventBuilder
		wantErr bool
	}{
		{
			name: "valid event",
			builder: func() *EventBuilder {
				return NewEventBuilder().
					WithPatientID(validAddr).
					WithType(EventLabResult).
					WithTitle("Blood Test").
					WithTimestamp(time.Now())
			},
			wantErr: false,
		},
		{
			name: "missing patient ID",
			builder: func() *EventBuilder {
				return NewEventBuilder().
					WithType(EventLabResult).
					WithTitle("Blood Test").
					WithTimestamp(time.Now())
			},
			wantErr: true,
		},
		{
			name: "missing title",
			builder: func() *EventBuilder {
				return NewEventBuilder().
					WithPatientID(validAddr).
					WithType(EventLabResult).
					WithTimestamp(time.Now())
			},
			wantErr: true,
		},
		{
			name: "invalid event type",
			builder: func() *EventBuilder {
				return NewEventBuilder().
					WithPatientID(validAddr).
					WithType("invalid").
					WithTitle("Blood Test").
					WithTimestamp(time.Now())
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.builder()
			event, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if event == nil {
					t.Error("Build() returned nil for valid event")
				}
				if event.CreatedAt.IsZero() {
					t.Error("Build() should set CreatedAt if not provided")
				}
			}
		})
	}
}

func TestEventBuilder_SetMetadata(t *testing.T) {
	builder := NewEventBuilder()
	builder.SetMetadata("key", "value")

	val := builder.event.Metadata.GetString("key")
	if val != "value" {
		t.Errorf("SetMetadata() value = %v, want value", val)
	}
}
