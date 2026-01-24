package timeline

import (
	"testing"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestRelationshipType_IsValid(t *testing.T) {
	tests := []struct {
		rt   RelationshipType
		want bool
	}{
		{RelResultedIn, true},
		{RelSupports, true},
		{RelFollowsUp, true},
		{RelReplaces, true},
		{RelCausedBy, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.rt), func(t *testing.T) {
			if got := tt.rt.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEdge_Validate(t *testing.T) {
	tests := []struct {
		name    string
		edge    Edge
		wantErr bool
	}{
		{
			name: "valid edge",
			edge: Edge{
				FromID: "event-1",
				ToID:   "event-2",
				Type:   RelResultedIn,
			},
			wantErr: false,
		},
		{
			name: "missing fromID",
			edge: Edge{
				ToID: "event-2",
				Type: RelResultedIn,
			},
			wantErr: true,
		},
		{
			name: "missing toID",
			edge: Edge{
				FromID: "event-1",
				Type:   RelResultedIn,
			},
			wantErr: true,
		},
		{
			name: "self-reference",
			edge: Edge{
				FromID: "event-1",
				ToID:   "event-1",
				Type:   RelResultedIn,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			edge: Edge{
				FromID: "event-1",
				ToID:   "event-2",
				Type:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.edge.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEdge_Reverse(t *testing.T) {
	edge := Edge{
		ID:     "edge-1",
		FromID: "event-1",
		ToID:   "event-2",
		Type:   RelResultedIn,
	}

	reversed := edge.Reverse()

	if reversed.FromID != "event-2" {
		t.Errorf("Expected reversed FromID to be event-2, got %s", reversed.FromID)
	}
	if reversed.ToID != "event-1" {
		t.Errorf("Expected reversed ToID to be event-1, got %s", reversed.ToID)
	}
	if reversed.ID != edge.ID {
		t.Error("Reverse should preserve ID")
	}
}

func TestGraphData(t *testing.T) {
	g := NewGraphData()

	event := Event{ID: types.ID("evt-1"), Title: "Test Event"}
	edge := Edge{ID: types.ID("edge-1"), FromID: "evt-1", ToID: "evt-2", Type: RelResultedIn}

	g.AddEvent(event)
	g.AddEdge(edge)

	if len(g.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(g.Events))
	}

	found := g.FindEvent("evt-1")
	if found == nil {
		t.Error("Expected to find event")
	}

	outgoing := g.GetOutgoingEdges("evt-1")
	if len(outgoing) != 1 {
		t.Errorf("Expected 1 outgoing edge, got %d", len(outgoing))
	}

	incoming := g.GetIncomingEdges("evt-2")
	if len(incoming) != 1 {
		t.Errorf("Expected 1 incoming edge, got %d", len(incoming))
	}
}
