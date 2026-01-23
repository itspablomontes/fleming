/**
 * Relationship Cluster Component
 *
 * A mini graph view that appears when an event is selected.
 * Shows the selected event and all its related events in a radial cluster.
 */

import {
	Background,
	BackgroundVariant,
	Controls,
	type Edge,
	MarkerType,
	type Node,
	ReactFlow,
} from "@xyflow/react";
import { useMemo } from "react";
import "@xyflow/react/dist/style.css";

import type { EventEdge, TimelineEvent } from "../types";
import { EventNode } from "./event-node";
import { RelationshipEdge } from "./relationship-edge";

// Register custom node and edge types
const nodeTypes = {
	eventNode: EventNode,
};

const edgeTypes = {
	relationshipEdge: RelationshipEdge,
};

interface RelationshipClusterProps {
	centerEvent: TimelineEvent;
	relatedEvents: Array<{
		event: TimelineEvent;
		edge: EventEdge;
		direction: "incoming" | "outgoing";
	}>;
	onEventClick: (event: TimelineEvent) => void;
}

/**
 * Calculate radial positions for related events around the center.
 */
function getRadialPosition(
	index: number,
	total: number,
	radius: number,
): { x: number; y: number } {
	const angle = (2 * Math.PI * index) / total - Math.PI / 2; // Start from top
	return {
		x: Math.cos(angle) * radius + 300, // Center X
		y: Math.sin(angle) * radius + 150, // Center Y
	};
}

export function RelationshipCluster({
	centerEvent,
	relatedEvents,
	onEventClick,
}: RelationshipClusterProps) {
	// Find all related events and edges
	const { nodes, edges } = useMemo(() => {
		// Create center node (selected event)
		const centerNode: Node = {
			id: centerEvent.id,
			type: "eventNode",
			position: { x: 0, y: 0 },
			data: {
				event: centerEvent,
				isSelected: true,
				onClick: onEventClick,
			},
		};

		// Create surrounding nodes
		const surroundingNodes: Node[] = relatedEvents.map(({ event }, index) => ({
			id: event.id,
			type: "eventNode",
			position: getRadialPosition(index, relatedEvents.length, 400),
			data: {
				event,
				isSelected: false,
				onClick: onEventClick,
			},
		}));

		// Create edges
		const flowEdges: Edge[] = relatedEvents.map(({ edge, event }) => {
			// Determine color based on relationship type
			let edgeColor = "#22d3ee"; // Default Cyan
			if (edge.relationshipType === "contradicts")
				edgeColor = "#ef4444"; // Red
			else if (edge.relationshipType === "resulted_in")
				edgeColor = "#22c55e"; // Green
			else if (edge.relationshipType === "supports") edgeColor = "#3b82f6"; // Blue

			return {
				id: `${centerEvent.id}-${event.id}`,
				source: centerEvent.id,
				target: event.id,
				type: "relationshipEdge",
				animated: true,
				markerEnd: {
					type: MarkerType.ArrowClosed,
					width: 12,
					height: 12,
					color: edgeColor,
				},
				style: {
					stroke: edgeColor,
				},
				data: {
					relationshipType: edge.relationshipType,
				},
			};
		});

		return {
			nodes: [centerNode, ...surroundingNodes],
			edges: flowEdges,
		};
	}, [centerEvent, relatedEvents, onEventClick]);

	return (
		<div style={{ width: "100%", height: "100%" }}>
			<ReactFlow
				nodes={nodes}
				edges={edges}
				nodeTypes={nodeTypes}
				edgeTypes={edgeTypes}
				fitView
				fitViewOptions={{ padding: 0.3 }}
				minZoom={0.1}
				maxZoom={2}
				nodesDraggable={true}
				nodesConnectable={false}
				elementsSelectable={true}
				panOnDrag={true} // Now fully interactive by default as per request
				zoomOnScroll={true}
				zoomOnPinch={true}
				zoomOnDoubleClick={true}
			>
				<Background
					variant={BackgroundVariant.Dots}
					gap={30} // Increased gap for airy feel
					size={1}
					color="rgba(34, 211, 238, 0.1)"
				/>
				<Controls
					showInteractive={false}
					className="bg-background/50 backdrop-blur border border-border"
				/>
			</ReactFlow>
		</div>
	);
}
