import { useRef } from "react";
import type { GraphData } from "../types";

interface TimelineMinimapProps {
	data: GraphData;
	minDate: Date;
	maxDate: Date;
	viewportStart: number;
	viewportWidth: number;
	totalWidth: number;
	onScroll: (percentage: number) => void;
}

export function TimelineMinimap({
	data,
	minDate,
	maxDate,
	viewportStart,
	viewportWidth,
	totalWidth,
	onScroll,
}: TimelineMinimapProps) {
	const trackRef = useRef<HTMLDivElement>(null);

	// Calculate normalized positions (0-100%)
	const viewportLeftPct = Math.max(
		0,
		Math.min(100, (viewportStart / totalWidth) * 100),
	);
	const viewportWidthPct = Math.max(
		5,
		Math.min(100, (viewportWidth / totalWidth) * 100),
	);

	// Drag handler: pass center position (0–1); clamp to track so drag-to-end works.
	const handleMouseDown = (e: React.MouseEvent) => {
		const track = trackRef.current;
		if (!track) return;

		const updatePosition = (clientX: number) => {
			const rect = track.getBoundingClientRect();
			const clickX = Math.max(0, Math.min(rect.width, clientX - rect.left));
			const clickPct = (clickX / rect.width) * 100;
			onScroll(clickPct / 100);
		};

		updatePosition(e.clientX);

		const handleMouseMove = (moveEvent: MouseEvent) => {
			updatePosition(moveEvent.clientX);
		};

		const handleMouseUp = () => {
			document.removeEventListener("mousemove", handleMouseMove);
			document.removeEventListener("mouseup", handleMouseUp);
		};

		document.addEventListener("mousemove", handleMouseMove);
		document.addEventListener("mouseup", handleMouseUp);
	};

	const handleTrackClick = (e: React.MouseEvent) => {
		if (!trackRef.current) return;
		const rect = trackRef.current.getBoundingClientRect();
		const clickX = Math.max(0, Math.min(rect.width, e.clientX - rect.left));
		const clickPct = (clickX / rect.width) * 100;
		onScroll(clickPct / 100);
	};

	return (
		<div
			className="fixed bottom-6 left-1/2 -translate-x-1/2 z-40 bg-background/90 backdrop-blur border border-border rounded-lg shadow-lg p-2 flex flex-col gap-2 transition-all hover:opacity-100 opacity-90"
			style={{ width: "min(600px, 80vw)" }}
		>
			{/* biome-ignore lint/a11y/noStaticElementInteractions: Custom interactive slider track */}
			{/* biome-ignore lint/a11y/useKeyWithClickEvents: Mouse-centric navigation control */}
			<div
				ref={trackRef}
				className="relative h-12 w-full bg-muted/30 rounded cursor-pointer overflow-hidden border border-border/50"
				onClick={handleTrackClick}
			>
				{/* Event Dots */}
				{data.events.map((event) => {
					const eventDate = new Date(event.timestamp);
					const totalMs = maxDate.getTime() - minDate.getTime();
					const eventMs = eventDate.getTime() - minDate.getTime();
					const leftPct = (eventMs / totalMs) * 100;

					return (
						<div
							key={event.id}
							className="absolute top-1/2 -translate-y-1/2 w-1.5 h-1.5 rounded-full bg-primary/40 pointer-events-none"
							style={{ left: `${leftPct}%` }}
						/>
					);
				})}

				{/* Viewport Indicator */}
				{/* biome-ignore lint/a11y/noStaticElementInteractions: Viewport indicator is intentionally interactive via dragging */}
				<div
					className="absolute top-0 bottom-0 border-2 border-primary/50 bg-primary/10 rounded cursor-grab active:cursor-grabbing hover:bg-primary/20 transition-colors"
					style={{
						left: `${viewportLeftPct}%`,
						width: `${viewportWidthPct}%`,
					}}
					onMouseDown={(e) => {
						e.stopPropagation(); // Prevent track click
						handleMouseDown(e);
					}}
				/>
			</div>
			<div className="flex justify-between text-[10px] text-muted-foreground px-1 uppercase tracking-wider font-medium">
				<span>
					{minDate.toLocaleDateString(undefined, {
						month: "short",
						year: "numeric",
					})}
				</span>
				<span>
					{maxDate.toLocaleDateString(undefined, {
						month: "short",
						year: "numeric",
					})}
				</span>
			</div>
		</div>
	);
}
