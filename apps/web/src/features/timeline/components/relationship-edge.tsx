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
						? "#ffffff"
						: style?.stroke || "rgba(34, 211, 238, 0.4)", // White on select, else type color
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
							color: selected ? "#ffffff" : "#22d3ee",
							backgroundColor: "#0b1221", // --card color from DeSci theme
							padding: "4px 10px",
							borderRadius: 6,
							border: `1px solid ${selected ? "#22d3ee" : "rgba(34, 211, 238, 0.3)"}`,
							pointerEvents: "all",
							backdropFilter: "blur(4px)",
							boxShadow: selected ? "0 0 15px rgba(34, 211, 238, 0.4)" : "none",
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
