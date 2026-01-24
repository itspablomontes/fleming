package timeline

import (
	"slices"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type RelationshipType string

const (
	RelResultedIn RelationshipType = "resulted_in"

	RelLeadTo RelationshipType = "lead_to"

	RelRequestedBy RelationshipType = "requested_by"

	RelSupports RelationshipType = "supports"

	RelFollowsUp RelationshipType = "follows_up"

	RelContradicts RelationshipType = "contradicts"

	RelAttachedTo RelationshipType = "attached_to"

	RelReplaces RelationshipType = "replaces"

	RelCausedBy RelationshipType = "caused_by"
)

func ValidRelationshipTypes() []RelationshipType {
	return []RelationshipType{
		RelResultedIn, RelLeadTo, RelRequestedBy,
		RelSupports, RelFollowsUp, RelContradicts, RelAttachedTo,
		RelReplaces, RelCausedBy,
	}
}

func (rt RelationshipType) IsValid() bool {
	return slices.Contains(ValidRelationshipTypes(), rt)
}

func (rt RelationshipType) Description() string {
	switch rt {
	case RelResultedIn:
		return "resulted in"
	case RelLeadTo:
		return "lead to"
	case RelRequestedBy:
		return "was requested by"
	case RelSupports:
		return "supports"
	case RelFollowsUp:
		return "follows up on"
	case RelContradicts:
		return "contradicts"
	case RelAttachedTo:
		return "is attached to"
	case RelReplaces:
		return "replaces"
	case RelCausedBy:
		return "was caused by"
	default:
		return "relates to"
	}
}

type Edge struct {
	ID     types.ID `json:"id"`
	FromID types.ID `json:"fromEventId"`

	ToID types.ID `json:"toEventId"`

	Type RelationshipType `json:"relationshipType"`

	Metadata types.Metadata `json:"metadata,omitempty"`
}

func (e *Edge) Validate() error {
	var errs types.ValidationErrors

	if e.FromID.IsEmpty() {
		errs.Add("fromEventId", "source event ID is required")
	}

	if e.ToID.IsEmpty() {
		errs.Add("toEventId", "target event ID is required")
	}

	if e.FromID == e.ToID && !e.FromID.IsEmpty() {
		errs.Add("toEventId", "cannot link event to itself")
	}

	if !e.Type.IsValid() {
		errs.Add("relationshipType", "invalid relationship type")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func (e *Edge) Reverse() Edge {
	return Edge{
		ID:       e.ID,
		FromID:   e.ToID,
		ToID:     e.FromID,
		Type:     e.Type,
		Metadata: e.Metadata,
	}
}
