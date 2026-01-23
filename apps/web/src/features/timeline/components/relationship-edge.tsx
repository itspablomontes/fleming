/**
 * Relationship Edge Component
 *
 * Custom React Flow edge with relationship type label.
 */

import {
	BaseEdge,
	EdgeLabelRenderer,
	type EdgeProps,
	getBezierPath,
} from "@xyflow/react";
import { memo } from "react";
import { RELATIONSHIP_LABELS, type RelationshipType } from "../types";

export interface RelationshipEdgeData {
	relationshipType: RelationshipType;
	[key: string]: unknown;
}

function RelationshipEdgeComponent({
	id,
	sourceX,
	sourceY,
	targetX,
	targetY,
	sourcePosition,
	targetPosition,
	data,
	selected,
	markerEnd,
	style,
}: EdgeProps) {
	const [edgePath, labelX, labelY] = getBezierPath({
		sourceX,
		sourceY,
		sourcePosition,
		targetX,
		targetY,
		targetPosition,
	});

	const edgeData = data as RelationshipEdgeData;
	const label = edgeData?.relationshipType
		? RELATIONSHIP_LABELS[edgeData.relationshipType]
		: "";


	return (
		<>
			<BaseEdge
				id={id}
				path={edgePath}
				style={{
					...style, // Use dynamic color from parent
					strokeWidth: selected ? 2.5 : 1.5,
					stroke: selected
						? "var(--foreground)"
						: style?.stroke || "var(--primary)", // White on select, else theme primary
					opacity: selected ? 1 : 0.4,
				}}
				markerEnd={markerEnd} // Pass markerEnd so arrowheads match color
			/>
			{label && (
				<EdgeLabelRenderer>
					<div
						style={{
							position: "absolute",
							transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
							fontSize: 10,
							fontWeight: 700,
							color: selected ? "var(--foreground)" : "var(--primary)",
							backgroundColor: "var(--card)",
							padding: "4px 10px",
							borderRadius: 6,
							border: `1px solid ${selected ? "var(--primary)" : "var(--border)"}`,
							pointerEvents: "all",
							backdropFilter: "blur(4px)",
							boxShadow: selected ? "0 0 15px var(--glow-primary)" : "none",
						}}
						className="nodrag nopan"
					>
						{label}
					</div>
				</EdgeLabelRenderer>
			)}
		</>
	);
}

export const RelationshipEdge = memo(RelationshipEdgeComponent);
