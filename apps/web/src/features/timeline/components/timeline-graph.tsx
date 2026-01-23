/**
 * Timeline Graph Component
 *
 * Main React Flow canvas for displaying the patient timeline as an interactive graph.
 */

import {
	Background,
	BackgroundVariant,
	Controls,
	type Edge,
	MarkerType,
	type Node,
	ReactFlow,
	useEdgesState,
	useNodesState,
} from "@xyflow/react";
import { useEffect, useMemo } from "react";
import "@xyflow/react/dist/style.css";

import type { GraphData, TimelineEvent } from "../types";
import { EventNode, type EventNodeData } from "./event-node";
import {
	RelationshipEdge,
	type RelationshipEdgeData,
} from "./relationship-edge";

// Register custom node and edge types
const nodeTypes = {
	eventNode: EventNode,
};

const edgeTypes = {
	relationshipEdge: RelationshipEdge,
};

interface TimelineGraphProps {
	data: GraphData;
	onEventClick?: (event: TimelineEvent) => void;
	selectedEventId?: string | null;
}

/**
 * Convert GraphData to React Flow nodes and edges.
 */
function graphDataToFlow(
	data: GraphData,
	onEventClick?: (event: TimelineEvent) => void,
	selectedEventId?: string | null,
): { nodes: Node<EventNodeData>[]; edges: Edge<RelationshipEdgeData>[] } {
	// Sort events by timestamp for layout
	const sortedEvents = [...data.events].sort(
		(a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
	);

	// Create nodes with automatic layout
	const nodes: Node<EventNodeData>[] = sortedEvents.map((event, index) => {
		// Simple grid layout: events spread horizontally with some vertical variation
		const col = index % 4;
		const row = Math.floor(index / 4);

		return {
			id: event.id,
			type: "eventNode",
			position: {
				x: col * 400 + 80,
				y: row * 280 + 80,
			},
			data: {
				event,
				isSelected: event.id === selectedEventId,
				onClick: onEventClick,
			},
		};
	});

	// Create edges
	const edges: Edge<RelationshipEdgeData>[] = data.edges.map((edge) => ({
		id: edge.id,
		source: edge.fromEventId,
		target: edge.toEventId,
		type: "relationshipEdge",
		animated: false,
		markerEnd: {
			type: MarkerType.ArrowClosed,
			width: 15,
			height: 15,
			color: "#22d3ee", // cyan-400 for dark theme
		},
		data: {
			relationshipType: edge.relationshipType,
		},
	}));

	return { nodes, edges };
}

export function TimelineGraph({
	data,
	onEventClick,
	selectedEventId,
}: TimelineGraphProps) {
	const { nodes: initialNodes, edges: initialEdges } = useMemo(
		() => graphDataToFlow(data, onEventClick, selectedEventId),
		[data, onEventClick, selectedEventId],
	);

	const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
	const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

	// Update nodes when selection changes
	useEffect(() => {
		setNodes((nds) =>
			nds.map((node) => ({
				...node,
				data: {
					...node.data,
					isSelected: node.id === selectedEventId,
				},
			})),
		);
	}, [selectedEventId, setNodes]);

	// Update when data changes
	useEffect(() => {
		const { nodes: newNodes, edges: newEdges } = graphDataToFlow(
			data,
			onEventClick,
			selectedEventId,
		);
		setNodes(newNodes);
		setEdges(newEdges);
	}, [data, onEventClick, selectedEventId, setNodes, setEdges]);

	return (
		<div style={{ width: "100%", height: "100%" }}>
			<ReactFlow
				nodes={nodes}
				edges={edges}
				onNodesChange={onNodesChange}
				onEdgesChange={onEdgesChange}
				nodeTypes={nodeTypes}
				edgeTypes={edgeTypes}
				fitView
				fitViewOptions={{ padding: 0.2 }}
				minZoom={0.3}
				maxZoom={2}
				defaultEdgeOptions={{
					type: "relationshipEdge",
				}}
			>
				<Controls className="bg-card! border-border!" />
				<Background
					variant={BackgroundVariant.Dots}
					gap={24}
					size={1}
					color="rgba(34, 211, 238, 0.15)" // cyan dots
				/>

				{/* Custom arrowhead marker */}
				<svg aria-label="Timeline arrows" role="img">
					<defs>
						<marker
							id="arrowhead"
							markerWidth="15"
							markerHeight="15"
							refX="13"
							refY="7.5"
							orient="auto"
						>
							<polygon
								points="0 0, 15 7.5, 0 15"
								fill="#22d3ee" // cyan-400
							/>
						</marker>
					</defs>
				</svg>
			</ReactFlow>
		</div>
	);
}
