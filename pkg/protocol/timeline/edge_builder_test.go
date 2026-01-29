package timeline

import (
	"testing"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestEdgeBuilder_WithFromID(t *testing.T) {
	validID, _ := types.NewID("event-1")
	builder := NewEdgeBuilder()

	builder.WithFromID(validID)
	if builder.edge.FromID != validID {
		t.Error("WithFromID() did not set from ID")
	}
}

func TestEdgeBuilder_WithToID(t *testing.T) {
	validID, _ := types.NewID("event-2")
	builder := NewEdgeBuilder()

	builder.WithToID(validID)
	if builder.edge.ToID != validID {
		t.Error("WithToID() did not set to ID")
	}
}

func TestEdgeBuilder_WithType(t *testing.T) {
	builder := NewEdgeBuilder()

	builder.WithType(RelResultedIn)
	if builder.edge.Type != RelResultedIn {
		t.Error("WithType() did not set relationship type")
	}
}

func TestEdgeBuilder_Build(t *testing.T) {
	fromID, _ := types.NewID("event-1")
	toID, _ := types.NewID("event-2")

	tests := []struct {
		name    string
		builder func() *EdgeBuilder
		wantErr bool
	}{
		{
			name: "valid edge",
			builder: func() *EdgeBuilder {
				return NewEdgeBuilder().
					WithFromID(fromID).
					WithToID(toID).
					WithType(RelResultedIn)
			},
			wantErr: false,
		},
		{
			name: "missing from ID",
			builder: func() *EdgeBuilder {
				return NewEdgeBuilder().
					WithToID(toID).
					WithType(RelResultedIn)
			},
			wantErr: true,
		},
		{
			name: "missing to ID",
			builder: func() *EdgeBuilder {
				return NewEdgeBuilder().
					WithFromID(fromID).
					WithType(RelResultedIn)
			},
			wantErr: true,
		},
		{
			name: "self-reference",
			builder: func() *EdgeBuilder {
				return NewEdgeBuilder().
					WithFromID(fromID).
					WithToID(fromID).
					WithType(RelResultedIn)
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			builder: func() *EdgeBuilder {
				return NewEdgeBuilder().
					WithFromID(fromID).
					WithToID(toID).
					WithType("invalid")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.builder()
			edge, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && edge == nil {
				t.Error("Build() returned nil for valid edge")
			}
		})
	}
}
