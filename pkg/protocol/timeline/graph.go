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

type Graph interface {
	GraphReader
	GraphWriter
}

type GraphData struct {
	Events []Event `json:"events"`
	Edges  []Edge  `json:"edges"`
}

func NewGraphData() GraphData {
	return GraphData{
		Events: make([]Event, 0),
		Edges:  make([]Edge, 0),
	}
}

func (g *GraphData) AddEvent(event Event) {
	g.Events = append(g.Events, event)
}

func (g *GraphData) AddEdge(edge Edge) {
	g.Edges = append(g.Edges, edge)
}

func (g *GraphData) FindEvent(id types.ID) *Event {
	for i := range g.Events {
		if g.Events[i].ID == id {
			return &g.Events[i]
		}
	}
	return nil
}

func (g *GraphData) GetOutgoingEdges(eventID types.ID) []Edge {
	var edges []Edge
	for _, edge := range g.Edges {
		if edge.FromID == eventID {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (g *GraphData) GetIncomingEdges(eventID types.ID) []Edge {
	var edges []Edge
	for _, edge := range g.Edges {
		if edge.ToID == eventID {
			edges = append(edges, edge)
		}
	}
	return edges
}
