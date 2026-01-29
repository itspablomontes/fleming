package timeline

import (
	"testing"
)

func TestRelationshipType_IsValid(t *testing.T) {
	tests := []struct {
		rt   RelationshipType
		want bool
	}{
		// Core relationships
		{RelResultedIn, true},
		{RelLeadTo, true},
		{RelRequestedBy, true},
		{RelSupports, true},
		{RelFollowsUp, true},
		{RelContradicts, true},
		{RelAttachedTo, true},
		{RelReplaces, true},
		{RelCausedBy, true},
		// Provider attestation
		{RelCosignedBy, true},
		{RelAttestedBy, true},
		// Medical relationships
		{RelTreats, true},
		{RelMonitors, true},
		{RelContraindicated, true},
		{RelDerivedFrom, true},
		{RelPartOf, true},
		// AI/Suggestions
		{RelSuggestedBy, true},
		// Invalid
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

func TestRelationshipType_Description(t *testing.T) {
	tests := []struct {
		rt   RelationshipType
		want string
	}{
		{RelResultedIn, "resulted in"},
		{RelCosignedBy, "was co-signed by"},
		{RelAttestedBy, "was attested by"},
		{RelTreats, "treats"},
		{RelMonitors, "monitors"},
		{RelContraindicated, "is contraindicated with"},
		{RelDerivedFrom, "was derived from"},
		{RelPartOf, "is part of"},
		{RelSuggestedBy, "was suggested by"},
		{"unknown", "relates to"},
	}

	for _, tt := range tests {
		t.Run(string(tt.rt), func(t *testing.T) {
			got := tt.rt.Description()
			if got != tt.want {
				t.Errorf("Description() = %v, want %v", got, tt.want)
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
