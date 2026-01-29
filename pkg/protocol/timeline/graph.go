package timeline

import (
	"context"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type GraphReader interface {
	GetEvent(ctx context.Context, id types.ID) (*Event, error)

	GetTimeline(ctx context.Context, patientID types.WalletAddress) ([]Event, error)

	GetRelated(ctx context.Context, eventID types.ID, depth int) ([]Event, []Edge, error)
}

type GraphWriter interface {
	CreateEvent(ctx context.Context, event *Event) error

	UpdateEvent(ctx context.Context, event *Event) error

	DeleteEvent(ctx context.Context, id types.ID) error

	CreateEdge(ctx context.Context, edge *Edge) error

	DeleteEdge(ctx context.Context, id types.ID) error
}

// Graph is a convenience interface that combines GraphReader and GraphWriter.
// Consumers should prefer composing GraphReader and GraphWriter separately when possible
// to follow Interface Segregation Principle.
type Graph interface {
	GraphReader
	GraphWriter
}
