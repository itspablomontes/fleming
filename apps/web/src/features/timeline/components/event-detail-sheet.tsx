/**
 * Event Detail Sheet Component
 *
 * Side panel showing detailed information about a selected timeline event.
 */

import {
	AlertTriangle,
	ArrowLeft,
	ArrowRight,
	Calendar,
	ClipboardList,
	File,
	FileCheck,
	FileText,
	HeartPulse,
	type LucideIcon,
	Pill,
	ScanLine,
	Scissors,
	Stethoscope,
	Syringe,
	TestTube2,
	User,
	UserPlus,
	X,
} from "lucide-react";
import {
	EVENT_TYPE_COLORS,
	EVENT_TYPE_LABELS,
	type EventEdge,
	RELATIONSHIP_LABELS,
	type TimelineEvent,
	type TimelineEventType,
} from "../types";

// Icon mapping
const ICON_MAP: Record<TimelineEventType, LucideIcon> = {
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

interface RelatedEventInfo {
	event: TimelineEvent;
	edge: EventEdge;
	direction: "incoming" | "outgoing";
}

interface EventDetailSheetProps {
	event: TimelineEvent | null;
	relatedEvents?: RelatedEventInfo[];
	onClose: () => void;
	onEventClick?: (event: TimelineEvent) => void;
}

export function EventDetailSheet({
	event,
	relatedEvents = [],
	onClose,
	onEventClick,
}: EventDetailSheetProps) {
	if (!event) return null;

	const colors = EVENT_TYPE_COLORS[event.type] || EVENT_TYPE_COLORS.other;
	const Icon = ICON_MAP[event.type] || File;
	const label = EVENT_TYPE_LABELS[event.type] || "Event";

	const formattedDate = new Date(event.timestamp).toLocaleDateString("en-US", {
		weekday: "long",
		month: "long",
		day: "numeric",
		year: "numeric",
	});

	const formattedTime = new Date(event.timestamp).toLocaleTimeString("en-US", {
		hour: "numeric",
		minute: "2-digit",
	});

	return (
		<div
			style={{
				position: "absolute",
				right: 0,
				top: 0,
				bottom: 0,
				width: 400,
				backgroundColor: "var(--background)",
				borderLeft: "1px solid var(--border)",
				boxShadow: "-4px 0 20px rgba(0,0,0,0.1)",
				display: "flex",
				flexDirection: "column",
				zIndex: 100,
			}}
		>
			{/* Header */}
			<div
				style={{
					padding: 20,
					borderBottom: "1px solid var(--border)",
					display: "flex",
					alignItems: "flex-start",
					gap: 16,
				}}
			>
				<div
					className="dark"
					style={{
						backgroundColor: colors.border,
						borderRadius: 12,
						padding: 12,
						display: "flex",
						alignItems: "center",
						justifyContent: "center",
					}}
				>
					<Icon style={{ width: 24, height: 24, color: "#ffffff" }} />
				</div>
				<div style={{ flex: 1 }}>
					<div
						style={{
							fontSize: 11,
							fontWeight: 600,
							color: "var(--primary)",
							textTransform: "uppercase",
							letterSpacing: "0.5px",
							marginBottom: 4,
						}}
					>
						{label}
					</div>
					<h2
						style={{
							fontSize: 18,
							fontWeight: 600,
							color: "var(--foreground)",
							margin: 0,
						}}
					>
						{event.title}
					</h2>
				</div>
				<button
					type="button"
					onClick={onClose}
					style={{
						background: "none",
						border: "none",
						cursor: "pointer",
						padding: 8,
						borderRadius: 8,
						display: "flex",
						alignItems: "center",
						justifyContent: "center",
					}}
					aria-label="Close"
				>
					<X style={{ width: 20, height: 20, color: "var(--muted-foreground)" }} />
				</button>
			</div>

			{/* Content */}
			<div
				style={{
					flex: 1,
					overflow: "auto",
					padding: 20,
				}}
			>
				{/* Date and Provider */}
				<div style={{ marginBottom: 24 }}>
					<div
						style={{
							display: "flex",
							alignItems: "center",
							gap: 8,
							marginBottom: 8,
						}}
					>
						<Calendar style={{ width: 16, height: 16, color: "var(--muted-foreground)" }} />
						<span style={{ fontSize: 14, color: "var(--foreground)" }}>
							{formattedDate} at {formattedTime}
						</span>
					</div>
					{event.provider && (
						<div
							style={{
								display: "flex",
								alignItems: "center",
								gap: 8,
							}}
						>
							<User style={{ width: 16, height: 16, color: "var(--muted-foreground)" }} />
							<span style={{ fontSize: 14, color: "var(--foreground)" }}>
								{event.provider}
							</span>
						</div>
					)}
				</div>

				{/* Description */}
				{event.description && (
					<div style={{ marginBottom: 24 }}>
						<h3
							style={{
								fontSize: 12,
								fontWeight: 600,
								color: "var(--muted-foreground)",
								textTransform: "uppercase",
								letterSpacing: "0.5px",
								marginBottom: 8,
							}}
						>
							Description
						</h3>
						<p
							style={{
								fontSize: 14,
								color: "var(--foreground)",
								lineHeight: 1.6,
								margin: 0,
							}}
						>
							{event.description}
						</p>
					</div>
				)}

				{/* Metadata */}
				{event.metadata && Object.keys(event.metadata).length > 0 && (
					<div style={{ marginBottom: 24 }}>
						<h3
							style={{
								fontSize: 12,
								fontWeight: 600,
								color: "var(--muted-foreground)",
								textTransform: "uppercase",
								letterSpacing: "0.5px",
								marginBottom: 8,
							}}
						>
							Details
						</h3>
						<div
							style={{
								backgroundColor: "var(--muted)",
								borderRadius: 8,
								padding: 12,
							}}
						>
							{Object.entries(event.metadata).map(([key, value]) => (
								<div
									key={key}
									style={{
										display: "flex",
										justifyContent: "space-between",
										padding: "6px 0",
										borderBottom: "1px solid var(--border)",
									}}
								>
									<span
										style={{
											fontSize: 13,
											color: "var(--muted-foreground)",
											textTransform: "capitalize",
										}}
									>
										{key.replace(/_/g, " ")}
									</span>
									<span
										style={{ fontSize: 13, color: "var(--foreground)", fontWeight: 500 }}
									>
										{String(value)}
									</span>
								</div>
							))}
						</div>
					</div>
				)}

				{/* Related Events */}
				{relatedEvents.length > 0 && (
					<div>
						<h3
							style={{
								fontSize: 12,
								fontWeight: 600,
								color: "var(--muted-foreground)",
								textTransform: "uppercase",
								letterSpacing: "0.5px",
								marginBottom: 8,
							}}
						>
							Related Events ({relatedEvents.length})
						</h3>
						<div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
							{relatedEvents.map(({ event: relatedEvent, edge, direction }) => {
								const relColors =
									EVENT_TYPE_COLORS[relatedEvent.type] ||
									EVENT_TYPE_COLORS.other;
								const RelIcon = ICON_MAP[relatedEvent.type] || File;

								return (
									<button
										key={relatedEvent.id}
										type="button"
										onClick={() => onEventClick?.(relatedEvent)}
										style={{
											display: "flex",
											alignItems: "center",
											gap: 12,
											padding: 12,
											backgroundColor: "var(--card)",
											border: "1px solid var(--border)",
											borderRadius: 8,
											cursor: "pointer",
											textAlign: "left",
											width: "100%",
										}}
									>
										{direction === "outgoing" ? (
											<ArrowRight
												style={{ width: 14, height: 14, color: "var(--muted-foreground)" }}
											/>
										) : (
											<ArrowLeft
												style={{ width: 14, height: 14, color: "var(--muted-foreground)" }}
											/>
										)}
										<div
											className="dark"
											style={{
												backgroundColor: relColors.border,
												borderRadius: 6,
												padding: 6,
											}}
										>
											<RelIcon
												style={{
													width: 14,
													height: 14,
													color: "#ffffff",
												}}
											/>
										</div>
										<div style={{ flex: 1 }}>
											<div
												style={{
													fontSize: 13,
													fontWeight: 500,
													color: "var(--foreground)",
												}}
											>
												{relatedEvent.title}
											</div>
											<div
												style={{
													fontSize: 11,
													color: "var(--muted-foreground)",
												}}
											>
												{RELATIONSHIP_LABELS[edge.relationshipType]}
											</div>
										</div>
									</button>
								);
							})}
						</div>
					</div>
				)}
			</div>
		</div>
	);
}
