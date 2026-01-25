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
import { type JSX, useEffect, useMemo, useState } from "react";
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
	isCompact?: boolean,
): { nodes: Node<EventNodeData>[]; edges: Edge<RelationshipEdgeData>[] } {
	// Sort events by timestamp for layout
	const sortedEvents = [...data.events].sort(
		(a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
	);

	// Create nodes with automatic layout
	const nodes: Node<EventNodeData>[] = sortedEvents.map((event, index) => {
		// Compact layout for smaller screens to reduce panning
		const columns = isCompact ? 1 : 4;
		const col = index % columns;
		const row = Math.floor(index / columns);
		const columnWidth = isCompact ? 280 : 400;
		const rowHeight = isCompact ? 220 : 280;

		return {
			id: event.id,
			type: "eventNode",
			position: {
				x: col * columnWidth + 80,
				y: row * rowHeight + 80,
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
			color: "var(--primary)",
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
}: TimelineGraphProps): JSX.Element {
	const [isCompact, setIsCompact] = useState(false);

	useEffect(() => {
		const media = window.matchMedia("(max-width: 768px)");
		const update = () => setIsCompact(media.matches);
		update();
		media.addEventListener("change", update);
		return () => media.removeEventListener("change", update);
	}, []);

	const { nodes: initialNodes, edges: initialEdges } = useMemo(
		() => graphDataToFlow(data, onEventClick, selectedEventId, isCompact),
		[data, onEventClick, selectedEventId, isCompact],
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
			isCompact,
		);
		setNodes(newNodes);
		setEdges(newEdges);
	}, [data, onEventClick, selectedEventId, isCompact, setNodes, setEdges]);

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
				fitViewOptions={{ padding: isCompact ? 0.1 : 0.2 }}
				minZoom={isCompact ? 0.5 : 0.3}
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
					// Use a very light opacity primary for background pattern
					color="var(--muted-foreground)"
					className="opacity-20"
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
							<polygon points="0 0, 15 7.5, 0 15" fill="var(--primary)" />
						</marker>
					</defs>
				</svg>
			</ReactFlow>
		</div>
	);
}
