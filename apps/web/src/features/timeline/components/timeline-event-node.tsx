/**
 * Timeline Event Node Component
 *
 * Individual event marker on the horizontal timeline.
 * Shows as a colored circle with type icon, expands on hover/selection.
 */

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
import { memo, useState } from "react";
import {
	EVENT_TYPE_COLORS,
	EVENT_TYPE_LABELS,
	type TimelineEvent,
	type TimelineEventType,
} from "../types";

// Icon mapping
const ICON_MAP: Record<
	TimelineEventType,
	React.ComponentType<React.SVGProps<SVGSVGElement>>
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
	tombstone: File,
};

interface TimelineEventNodeProps {
	event: TimelineEvent;
	x: number;
	y: number;
	isSelected: boolean;
	onClick: () => void;
	onMouseEnter?: () => void;
	onMouseLeave?: () => void;
	zoomLevel: number;
}

function TimelineEventNodeComponent({
	event,
	x,
	y,
	isSelected,
	onClick,
	onMouseEnter,
	onMouseLeave,
	zoomLevel,
}: TimelineEventNodeProps) {
	const [isHovered, setIsHovered] = useState(false);
	const colors = EVENT_TYPE_COLORS[event.type] || EVENT_TYPE_COLORS.other;
	const Icon = ICON_MAP[event.type] || File;
	const label = EVENT_TYPE_LABELS[event.type] || "Event";

	// Icons should scale with zoom but stay within a healthy range
	// Using sqrt(zoomLevel) makes the scaling less volatile
	const scaleFactor = Math.sqrt(zoomLevel);
	const baseSize = isSelected ? 44 : isHovered ? 36 : 28;
	const size = baseSize * Math.max(scaleFactor, 0.6);
	const showTooltip = isHovered && !isSelected;

	const formattedDate = new Date(event.timestamp).toLocaleDateString("en-US", {
		month: "short",
		day: "numeric",
		year: "numeric",
	});

	return (
		<div
			style={{
				position: "absolute",
				left: x,
				top: y,
				transform: "translate(-50%, -50%)",
				zIndex: isSelected ? 20 : isHovered ? 10 : 1,
			}}
		>
			{/* Event node circle */}
			<button
				type="button"
				onClick={onClick}
				onMouseEnter={() => {
					setIsHovered(true);
					onMouseEnter?.();
				}}
				onMouseLeave={() => {
					setIsHovered(false);
					onMouseLeave?.();
				}}
				style={{
					width: size,
					height: size,
					borderRadius: "50%",
					backgroundColor: isSelected ? colors.border : colors.bg,
					border: `2px solid ${colors.border}`,
					cursor: "pointer",
					display: "flex",
					alignItems: "center",
					justifyContent: "center",
					boxShadow: isSelected
						? `0 0 20px ${colors.glow}, 0 0 40px ${colors.glow}`
						: isHovered
							? `0 0 15px ${colors.glow}`
							: "none",
					transition: "all 0.2s ease",
					animation: isSelected ? "pulse 2s infinite" : "none",
				}}
			>
				<Icon
					className={`${isSelected ? "text-background" : ""}`}
					style={{
						width: size * 0.5,
						height: size * 0.5,
						color: isSelected ? "var(--background)" : colors.border,
					}}
				/>
			</button>

			{/* Hover tooltip */}
			{showTooltip && (
				<div
					style={{
						position: "absolute",
						top: -80,
						left: "50%",
						transform: "translateX(-50%)",
						backgroundColor: "var(--popover)",
						border: `1px solid ${colors.border}`,
						borderRadius: 10,
						padding: "10px 14px",
						minWidth: 160,
						boxShadow: "0 8px 24px rgba(0,0,0,0.2)",
						zIndex: 100,
						pointerEvents: "none",
					}}
				>
					<div
						style={{
							fontSize: 10,
							fontWeight: 700,
							color: colors.text,
							textTransform: "uppercase",
							letterSpacing: "0.5px",
							marginBottom: 4,
						}}
					>
						{label}
					</div>
					<div
						style={{
							fontSize: 13,
							fontWeight: 600,
							color: "var(--popover-foreground)",
							marginBottom: 4,
							whiteSpace: "nowrap",
							overflow: "hidden",
							textOverflow: "ellipsis",
							maxWidth: 180,
						}}
					>
						{event.title}
					</div>
					<div
						style={{
							fontSize: 11,
							color: "var(--muted-foreground)",
						}}
					>
						{formattedDate}
					</div>

					{/* Arrow pointer */}
					<div
						style={{
							position: "absolute",
							bottom: -6,
							left: "50%",
							transform: "translateX(-50%) rotate(45deg)",
							width: 12,
							height: 12,
							backgroundColor: "var(--popover)",
							borderRight: `1px solid ${colors.border}`,
							borderBottom: `1px solid ${colors.border}`,
						}}
					/>
				</div>
			)}

			{/* Selected label below */}
			{isSelected && (
				<div
					style={{
						position: "absolute",
						top: 24,
						left: "50%",
						transform: "translateX(-50%)",
						textAlign: "center",
						whiteSpace: "nowrap",
					}}
				>
					<div
						style={{
							fontSize: 12,
							fontWeight: 600,
							color: "var(--foreground)",
							marginTop: 8,
						}}
					>
						{event.title}
					</div>
					<div
						style={{
							fontSize: 10,
							color: colors.text,
							marginTop: 2,
						}}
					>
						{formattedDate}
					</div>
				</div>
			)}
		</div>
	);
}

export const TimelineEventNode = memo(TimelineEventNodeComponent);
