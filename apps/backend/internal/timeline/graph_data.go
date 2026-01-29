package timeline

// GraphData is an in-memory representation of the graph for visualization/export.
// This is an implementation detail and should not be part of the protocol layer.
type GraphData struct {
	Events []TimelineEvent `json:"events"`
	Edges  []EventEdge     `json:"edges"`
}

func NewGraphData() GraphData {
	return GraphData{
		Events: make([]TimelineEvent, 0),
		Edges:  make([]EventEdge, 0),
	}
}

func (g *GraphData) AddEvent(event TimelineEvent) {
	g.Events = append(g.Events, event)
}

func (g *GraphData) AddEdge(edge EventEdge) {
	g.Edges = append(g.Edges, edge)
}

func (g *GraphData) FindEvent(id string) *TimelineEvent {
	for i := range g.Events {
		if g.Events[i].ID == id {
			return &g.Events[i]
		}
	}
	return nil
}

func (g *GraphData) GetOutgoingEdges(eventID string) []EventEdge {
	var edges []EventEdge
	for _, edge := range g.Edges {
		if edge.FromEventID == eventID {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (g *GraphData) GetIncomingEdges(eventID string) []EventEdge {
	var edges []EventEdge
	for _, edge := range g.Edges {
		if edge.ToEventID == eventID {
			edges = append(edges, edge)
		}
	}
	return edges
}
