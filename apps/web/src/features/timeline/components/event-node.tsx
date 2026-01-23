/**
 * Event Node Component
 *
 * Custom React Flow node representing a timeline event.
 * Displays event type icon, title, date, and type-based styling.
 */

import { Handle, Position } from "@xyflow/react";
import {
	AlertTriangle,
	ClipboardList,
	File,
	FileCheck,
	FileText,
	HeartPulse,
	Pill,
	ScanLine,
	Scissors,
	Stethoscope,
	Syringe,
	TestTube2,
	UserPlus,
} from "lucide-react";
import { memo } from "react";
import {
	EVENT_TYPE_COLORS,
	EVENT_TYPE_LABELS,
	type TimelineEvent,
	type TimelineEventType,
} from "../types";

// Icon mapping
const ICON_MAP: Record<
	TimelineEventType,
	React.ComponentType<{ className?: string }>
> = {
	consultation: Stethoscope,
	diagnosis: ClipboardList,
	prescription: Pill,
	procedure: Scissors,
	lab_result: TestTube2,
	imaging: ScanLine,
	note: FileText,
	vaccination: Syringe,
	allergy: AlertTriangle,
	visit_note: FileText,
	vital_signs: HeartPulse,
	referral: UserPlus,
	insurance_claim: FileCheck,
	other: File,
};

export interface EventNodeData {
	event: TimelineEvent;
	isSelected?: boolean;
	onClick?: (event: TimelineEvent) => void;
	[key: string]: unknown; // Index signature for React Flow compatibility
}

function EventNodeComponent({ data }: { data: EventNodeData }) {
	const { event, isSelected, onClick } = data;
	const colors = EVENT_TYPE_COLORS[event.type] || EVENT_TYPE_COLORS.other;
	const Icon = ICON_MAP[event.type] || File;
	const label = EVENT_TYPE_LABELS[event.type] || "Event";

	const formattedDate = new Date(event.timestamp).toLocaleDateString("en-US", {
		month: "short",
		day: "numeric",
		year: "numeric",
	});

	return (
		// biome-ignore lint/a11y/useSemanticElements: Cannot use <button> as it contains <Handle> divs
		<div
			className="event-node group"
			style={{
				backgroundColor: "var(--card)",
				borderColor: isSelected ? colors.border : "var(--border)",
				borderWidth: 1,
				borderStyle: "solid",
				borderRadius: 16,
				padding: "16px 20px",
				minWidth: 240,
				maxWidth: 320,
				cursor: "pointer",
				backdropFilter: "blur(12px)",
				boxShadow: isSelected
					? `0 0 20px var(--glow-primary), 0 0 40px var(--glow-primary)`
					: `0 4px 20px rgba(0,0,0,0.1)`,
				transition: "all 0.2s ease",
			}}
			onClick={() => onClick?.(event)}
			onKeyDown={(e) => {
				if (e.key === "Enter" || e.key === " ") {
					onClick?.(event);
				}
			}}
			role="button"
			tabIndex={0}
		>
			{/* Input handle (left) */}
			<Handle
				type="target"
				position={Position.Left}
				style={{
					background: "var(--background)",
					width: 12,
					height: 12,
					border: `2px solid ${colors.border}`,
				}}
			/>

			{/* Header with icon and type */}
			<div
				style={{
					display: "flex",
					alignItems: "center",
					gap: 10,
					marginBottom: 10,
				}}
			>
				<div
					className="dark"
					style={{
						backgroundColor: colors.border,
						borderRadius: 8,
						padding: 6,
						display: "flex",
						alignItems: "center",
						justifyContent: "center",
						boxShadow: `0 0 10px ${colors.glow}`,
					}}
				>
					<Icon className="w-4 h-4 text-white" />
				</div>
				<span
					style={{
						fontSize: 10,
						fontWeight: 700,
						color: colors.text,
						textTransform: "uppercase",
						letterSpacing: "1px",
					}}
				>
					{label}
				</span>
			</div>

			{/* Title */}
			<div
				style={{
					fontSize: 15,
					fontWeight: 600,
					color: "var(--foreground)",
					marginBottom: 6,
					overflow: "hidden",
					textOverflow: "ellipsis",
					whiteSpace: "nowrap",
				}}
			>
				{event.title}
			</div>

			{/* Date and provider */}
			<div
				style={{
					fontSize: 12,
					color: "var(--muted-foreground)",
					display: "flex",
					justifyContent: "space-between",
					alignItems: "center",
					fontWeight: 500,
				}}
			>
				<span>{formattedDate}</span>
				{event.provider && (
					<span
						style={{
							overflow: "hidden",
							textOverflow: "ellipsis",
							whiteSpace: "nowrap",
							maxWidth: 120,
							color: "var(--primary)", // Use primary for provider info
							opacity: 0.9,
						}}
					>
						{event.provider}
					</span>
				)}
			</div>

			{/* Output handle (right) */}
			<Handle
				type="source"
				position={Position.Right}
				style={{
					background: colors.border,
					width: 10,
					height: 10,
					border: `2px solid ${colors.bg}`,
				}}
			/>
		</div>
	);
}

export const EventNode = memo(EventNodeComponent);
