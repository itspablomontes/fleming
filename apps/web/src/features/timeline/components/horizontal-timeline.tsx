/**
 * Horizontal Timeline Component
 *
 * Main chronological timeline view with events positioned left-to-right by date.
 * Clicking an event reveals its relationship cluster below.
 */

import {
	ChevronLeft,
	ChevronRight,
	Map as MapIcon,
	ZoomIn,
	ZoomOut,
} from "lucide-react";
import {
	useCallback,
	useEffect,
	useLayoutEffect,
	useMemo,
	useRef,
	useState,
} from "react";
import { Button } from "@/components/ui/button";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import type { GraphData, TimelineEvent } from "../types";
import { TimelineEventNode } from "./timeline-event-node";
import { TimelineMinimap } from "./timeline-minimap";

interface HorizontalTimelineProps {
	data: GraphData;
	selectedEventId: string | null;
	onEventSelect: (event: TimelineEvent | null) => void;
}

/**
 * Calculate the horizontal position of an event based on its timestamp.
 */
function getEventPosition(
	timestamp: string,
	minDate: Date,
	maxDate: Date,
	trackWidth: number,
): number {
	const eventDate = new Date(timestamp);
	const totalRange = maxDate.getTime() - minDate.getTime();
	if (totalRange === 0) return trackWidth / 2;

	const eventOffset = eventDate.getTime() - minDate.getTime();
	const padding = 80; // Padding from edges
	const usableWidth = trackWidth - padding * 2;
	return padding + (eventOffset / totalRange) * usableWidth;
}

/**
 * Format month/year for axis labels.
 */
function formatAxisLabel(date: Date): string {
	return date.toLocaleDateString("en-US", { month: "short", year: "numeric" });
}

export function HorizontalTimeline({
	data,
	selectedEventId,
	onEventSelect,
}: HorizontalTimelineProps) {
	const containerRef = useRef<HTMLDivElement>(null);
	const [isDragging, setIsDragging] = useState(false);
	const [startX, setStartX] = useState(0);
	const [scrollLeft, setScrollLeft] = useState(0);
	const [zoomLevel, setZoomLevel] = useState(1); // 0.5 to 3
	const [showMinimap, setShowMinimap] = useState(true);
	const [hoveredEventId, setHoveredEventId] = useState<string | null>(null);
	const [zoomAnchor, setZoomAnchor] = useState<{
		ratio: number;
		x: number;
	} | null>(null);

	// Calculate time range from events
	const { minDate, maxDate, sortedEvents } = useMemo(() => {
		if (!data.events.length) {
			return {
				minDate: new Date(),
				maxDate: new Date(),
				sortedEvents: [],
			};
		}

		const sorted = [...data.events].sort(
			(a, b) =>
				new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
		);

		const min = new Date(sorted[0].timestamp);
		const max = new Date(sorted[sorted.length - 1].timestamp);

		// Add padding to time range (1 month before and after)
		min.setMonth(min.getMonth() - 1);
		max.setMonth(max.getMonth() + 1);

		return { minDate: min, maxDate: max, sortedEvents: sorted };
	}, [data.events]);

	// Generate axis labels (one per month)
	const axisLabels = useMemo(() => {
		const labels: { date: Date; label: string }[] = [];
		const current = new Date(minDate);
		current.setDate(1); // Start of month

		while (current <= maxDate) {
			labels.push({
				date: new Date(current),
				label: formatAxisLabel(current),
			});
			current.setMonth(current.getMonth() + 1);
		}

		return labels;
	}, [minDate, maxDate]);

	// Track width calculation based on zoom - significantly increased for better resolution
	const pixelsPerMonth = 600 * zoomLevel;
	const trackWidth = Math.max(
		containerRef.current?.clientWidth || 1200,
		axisLabels.length * pixelsPerMonth,
	);

	// Calculate layout with vertical staggering to avoid collisions
	const eventLayout = useMemo(() => {
		const layout: Record<string, { x: number; y: number; stemHeight: number }> =
			{};
		const COLLISION_THRESHOLD = 120 * zoomLevel; // Horizontal distance to trigger stagger

		let lastX = -Infinity;
		let level = 0;

		sortedEvents.forEach((event) => {
			const x = getEventPosition(event.timestamp, minDate, maxDate, trackWidth);

			// If too close to previous event, change vertical level
			if (x - lastX < COLLISION_THRESHOLD) {
				level = (level + 1) % 3;
			} else {
				level = 0;
			}

			// Stagger Y positions around the line (which is at 100)
			// Levels: 40 (High), 70 (Mid), 50 (Slightly High)
			const yOffsets = [40, 75, 20];
			const y = yOffsets[level];

			layout[event.id] = {
				x,
				y,
				stemHeight: 100 - y, // Connects down to the line at 100
			};

			lastX = x;
		});

		return layout;
	}, [sortedEvents, minDate, maxDate, trackWidth, zoomLevel]);

	// Zoom handlers
	const handleZoomIn = useCallback(() => {
		if (!containerRef.current) return;
		const container = containerRef.current;
		const viewportWidth = container.clientWidth;
		const centerOffset = viewportWidth / 2;
		const currentWidth = trackWidth;
		const pointerRatio = (container.scrollLeft + centerOffset) / currentWidth;

		setZoomLevel((prev) => {
			const next = Math.min(prev + 0.25, 3);
			if (next !== prev) {
				setZoomAnchor({ ratio: pointerRatio, x: centerOffset });
			}
			return next;
		});
	}, [trackWidth]);

	const handleZoomOut = useCallback(() => {
		if (!containerRef.current) return;
		const container = containerRef.current;
		const viewportWidth = container.clientWidth;
		const centerOffset = viewportWidth / 2;
		const currentWidth = trackWidth;
		const pointerRatio = (container.scrollLeft + centerOffset) / currentWidth;

		setZoomLevel((prev) => {
			const next = Math.max(prev - 0.25, 0.5);
			if (next !== prev) {
				setZoomAnchor({ ratio: pointerRatio, x: centerOffset });
			}
			return next;
		});
	}, [trackWidth]);

	const handleWheel = useCallback(
		(e: React.WheelEvent) => {
			if (!containerRef.current) return;

			// Prevent default vertical scrolling to use wheel for zoom
			e.preventDefault();

			const container = containerRef.current;
			const rect = container.getBoundingClientRect();
			const mouseX = e.clientX - rect.left;

			// Capture ratio relative to current trackWidth for consistency
			const currentWidth = trackWidth;
			const pointerRatio = (container.scrollLeft + mouseX) / currentWidth;

			// Use a smaller delta for smoother wheel zoom
			const delta = e.deltaY > 0 ? -0.1 : 0.1;

			setZoomLevel((prev) => {
				const next = Math.min(Math.max(prev + delta, 0.4), 4);
				if (next !== prev) {
					setZoomAnchor({ ratio: pointerRatio, x: mouseX });
				}
				return next;
			});
		},
		[trackWidth],
	);

	// Use layout effect to adjust scroll position immediately after zoom updates trackWidth
	// to prevent any visible jumping/flicking
	useLayoutEffect(() => {
		if (zoomAnchor && containerRef.current) {
			const container = containerRef.current;
			container.scrollLeft = zoomAnchor.ratio * trackWidth - zoomAnchor.x;
			setZoomAnchor(null);
		}
	}, [zoomAnchor, trackWidth]);

	// Jump handlers
	const scrollToDate = useCallback(
		(date: Date) => {
			if (!containerRef.current) return;

			const position = getEventPosition(
				date.toISOString(),
				minDate,
				maxDate,
				trackWidth,
			);

			containerRef.current.scrollTo({
				left: position - containerRef.current.clientWidth / 2,
				behavior: "smooth",
			});
		},
		[minDate, maxDate, trackWidth],
	);

	const handleJumpPrev = useCallback(() => {
		if (selectedEventId) {
			// Find current index and select previous
			const currentIndex = sortedEvents.findIndex(
				(e) => e.id === selectedEventId,
			);
			if (currentIndex > 0) {
				const prevEvent = sortedEvents[currentIndex - 1];
				onEventSelect(prevEvent);
				scrollToDate(new Date(prevEvent.timestamp));
				return;
			}
		}

		// Default scroll behavior if no selection or at start
		if (containerRef.current) {
			containerRef.current.scrollBy({
				left: -pixelsPerMonth,
				behavior: "smooth",
			});
		}
	}, [
		pixelsPerMonth,
		selectedEventId,
		sortedEvents,
		onEventSelect,
		scrollToDate,
	]);

	const handleJumpNext = useCallback(() => {
		if (selectedEventId) {
			// Find current index and select next
			const currentIndex = sortedEvents.findIndex(
				(e) => e.id === selectedEventId,
			);
			if (currentIndex !== -1 && currentIndex < sortedEvents.length - 1) {
				const nextEvent = sortedEvents[currentIndex + 1];
				onEventSelect(nextEvent);
				scrollToDate(new Date(nextEvent.timestamp));
				return;
			}
		}

		if (containerRef.current) {
			containerRef.current.scrollBy({
				left: pixelsPerMonth,
				behavior: "smooth",
			});
		}
	}, [
		pixelsPerMonth,
		selectedEventId,
		sortedEvents,
		onEventSelect,
		scrollToDate,
	]);

	// Keyboard navigation: first press = one step, then rAF-driven continuous pan/zoom
	const keyRepeatRef = useRef<{
		activeKey: string | null;
		rafId: number | null;
		timeoutId: ReturnType<typeof setTimeout> | null;
		lastZoomTime: number;
	}>({ activeKey: null, rafId: null, timeoutId: null, lastZoomTime: 0 });

	const clearKeyRepeat = useCallback(() => {
		const r = keyRepeatRef.current;
		if (r.timeoutId != null) {
			clearTimeout(r.timeoutId);
			r.timeoutId = null;
		}
		if (r.rafId != null) {
			cancelAnimationFrame(r.rafId);
			r.rafId = null;
		}
		r.activeKey = null;
	}, []);

	useEffect(() => {
		const INITIAL_DELAY_MS = 300;
		const PX_PER_FRAME = 8;
		const ZOOM_THROTTLE_MS = 120;

		const isEditable = (el: Element | null) => {
			if (!el) return false;
			const tag = (el as HTMLElement).tagName;
			if (tag === "INPUT" || tag === "TEXTAREA") return true;
			return (el as HTMLElement).isContentEditable;
		};

		const tick = () => {
			const r = keyRepeatRef.current;
			if (r.activeKey == null) return;
			const container = containerRef.current;
			const now = performance.now();

			if (r.activeKey === "ArrowLeft") {
				if (container)
					container.scrollBy({ left: -PX_PER_FRAME, behavior: "auto" });
			} else if (r.activeKey === "ArrowRight") {
				if (container)
					container.scrollBy({ left: PX_PER_FRAME, behavior: "auto" });
			} else if (r.activeKey === "+" || r.activeKey === "=") {
				if (now - r.lastZoomTime >= ZOOM_THROTTLE_MS) {
					r.lastZoomTime = now;
					handleZoomIn();
				}
			} else if (r.activeKey === "-") {
				if (now - r.lastZoomTime >= ZOOM_THROTTLE_MS) {
					r.lastZoomTime = now;
					handleZoomOut();
				}
			}

			r.rafId = requestAnimationFrame(tick);
		};

		const handleKeyDown = (e: KeyboardEvent) => {
			if (isEditable(document.activeElement)) return;
			const key = e.key;
			if (
				key !== "ArrowLeft" &&
				key !== "ArrowRight" &&
				key !== "+" &&
				key !== "=" &&
				key !== "-"
			)
				return;

			if (!e.repeat) {
				e.preventDefault();
				clearKeyRepeat();
				if (key === "ArrowLeft") handleJumpPrev();
				else if (key === "ArrowRight") handleJumpNext();
				else if (key === "+" || key === "=") handleZoomIn();
				else if (key === "-") handleZoomOut();

				const r = keyRepeatRef.current;
				r.activeKey = key;
				r.lastZoomTime = performance.now();
				r.timeoutId = setTimeout(() => {
					r.timeoutId = null;
					if (r.activeKey == null) return;
					r.rafId = requestAnimationFrame(tick);
				}, INITIAL_DELAY_MS);
			}
		};

		const handleKeyUp = (e: KeyboardEvent) => {
			const key = e.key;
			if (
				key === "ArrowLeft" ||
				key === "ArrowRight" ||
				key === "+" ||
				key === "=" ||
				key === "-"
			) {
				const r = keyRepeatRef.current;
				if (r.activeKey === key) clearKeyRepeat();
			}
		};

		const handleVisibility = () => {
			if (document.visibilityState === "hidden") clearKeyRepeat();
		};

		window.addEventListener("keydown", handleKeyDown);
		window.addEventListener("keyup", handleKeyUp);
		document.addEventListener("visibilitychange", handleVisibility);
		return () => {
			window.removeEventListener("keydown", handleKeyDown);
			window.removeEventListener("keyup", handleKeyUp);
			document.removeEventListener("visibilitychange", handleVisibility);
			clearKeyRepeat();
		};
	}, [handleJumpPrev, handleJumpNext, handleZoomIn, handleZoomOut, clearKeyRepeat]);

	// Drag to scroll handlers
	const handleMouseDown = useCallback((e: React.MouseEvent) => {
		if (!containerRef.current) return;
		setIsDragging(true);
		setStartX(e.pageX - containerRef.current.offsetLeft);
		setScrollLeft(containerRef.current.scrollLeft);
	}, []);

	const handleMouseUp = useCallback(() => {
		setIsDragging(false);
	}, []);

	const handleMouseMove = useCallback(
		(e: React.MouseEvent) => {
			if (!isDragging || !containerRef.current) return;
			e.preventDefault();
			const x = e.pageX - containerRef.current.offsetLeft;
			const walk = (x - startX) * 1.5;
			containerRef.current.scrollLeft = scrollLeft - walk;
		},
		[isDragging, startX, scrollLeft],
	);

	const handleEventClick = useCallback(
		(event: TimelineEvent) => {
			if (selectedEventId === event.id) {
				onEventSelect(null); // Deselect if clicking same event
			} else {
				onEventSelect(event);
			}
		},
		[selectedEventId, onEventSelect],
	);

	const handleMinimapScroll = useCallback((percentage: number) => {
		if (!containerRef.current) return;
		const totalWidth = containerRef.current.scrollWidth;
		const viewportWidth = containerRef.current.clientWidth;
		const maxScroll = totalWidth - viewportWidth;

		containerRef.current.scrollLeft = Math.max(
			0,
			Math.min(maxScroll, totalWidth * percentage - viewportWidth / 2),
		);
	}, []);

	// Update minimap on scroll
	const [viewportStart, setViewportStart] = useState(0);
	const [viewportWidth, setViewportWidth] = useState(100);

	const handleScroll = useCallback(() => {
		if (!containerRef.current) return;
		setViewportStart(containerRef.current.scrollLeft);
		setViewportWidth(containerRef.current.clientWidth);
	}, []);

	// Initial and resize observer for viewport measurement
	useEffect(() => {
		const updateMeasurements = () => {
			if (containerRef.current) {
				setViewportStart(containerRef.current.scrollLeft);
				setViewportWidth(containerRef.current.clientWidth);
			}
		};

		updateMeasurements();
		window.addEventListener("resize", updateMeasurements);
		return () => window.removeEventListener("resize", updateMeasurements);
	}, []); // Re-measure when layout changes

	return (
		// biome-ignore lint/a11y/noStaticElementInteractions: Custom interactive timeline container
		// biome-ignore lint/a11y/useAriaPropsSupportedByRole: Div acts as a complex focusable region
		<div
			ref={containerRef}
			className="no-scrollbar"
			style={{
				width: "100%",
				overflowX: "auto",
				overflowY: "hidden",
				cursor: isDragging ? "grabbing" : "grab",
				userSelect: "none",
				padding: "24px 0",
				outline: "none",
			}}
			aria-label="Interactive timeline, use arrow keys to navigate, scroll or drag to pan"
			onKeyDown={(e) => {
				if (e.key === "ArrowLeft") handleJumpPrev();
				if (e.key === "ArrowRight") handleJumpNext();
			}}
			onMouseDown={handleMouseDown}
			onMouseUp={handleMouseUp}
			onMouseLeave={handleMouseUp}
			onMouseMove={handleMouseMove}
			onWheel={handleWheel}
			onScroll={handleScroll}
		>
			{/* Controls Overlay */}
			<div className="fixed bottom-6 right-6 z-50 flex flex-col gap-2 p-2 bg-background/80 backdrop-blur-md border border-border rounded-xl shadow-lg">
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<Button variant="ghost" size="icon" onClick={handleZoomIn}>
								<ZoomIn className="w-4 h-4" />
							</Button>
						</TooltipTrigger>
						<TooltipContent side="left">Zoom In (+)</TooltipContent>
					</Tooltip>

					<Tooltip>
						<TooltipTrigger asChild>
							<Button variant="ghost" size="icon" onClick={handleZoomOut}>
								<ZoomOut className="w-4 h-4" />
							</Button>
						</TooltipTrigger>
						<TooltipContent side="left">Zoom Out (-)</TooltipContent>
					</Tooltip>

					<div className="h-px bg-border my-1" />

					<Tooltip>
						<TooltipTrigger asChild>
							<Button
								variant="ghost"
								size="icon"
								onClick={() => setShowMinimap(!showMinimap)}
							>
								<MapIcon className="w-4 h-4" />
							</Button>
						</TooltipTrigger>
						<TooltipContent side="left">Toggle Minimap</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			</div>

			{/* Jump Navigation Overlay */}
			<div className="absolute top-1/2 left-4 z-40 -translate-y-1/2 opacity-0 hover:opacity-100 transition-opacity">
				<Button
					variant="secondary"
					size="icon"
					className="rounded-full shadow-lg"
					onClick={handleJumpPrev}
				>
					<ChevronLeft className="w-6 h-6" />
				</Button>
			</div>
			<div className="absolute top-1/2 right-4 z-40 -translate-y-1/2 opacity-0 hover:opacity-100 transition-opacity">
				<Button
					variant="secondary"
					size="icon"
					className="rounded-full shadow-lg"
					onClick={handleJumpNext}
				>
					<ChevronRight className="w-6 h-6" />
				</Button>
			</div>

			{/* Timeline container */}
			<div
				style={{
					position: "relative",
					width: trackWidth,
					height: 350, // Increased height significantly for staggering + tooltips
					minWidth: "100%",
				}}
			>
				{/* Time axis labels */}
				<div
					style={{
						position: "absolute",
						top: 220, // Slightly below the line at 200
						left: 0,
						right: 0,
						height: 24,
						display: "flex",
					}}
				>
					{axisLabels.map(({ date, label }) => {
						const x = getEventPosition(
							date.toISOString(),
							minDate,
							maxDate,
							trackWidth,
						);
						return (
							<button
								key={date.toISOString()}
								type="button"
								onClick={() => scrollToDate(date)}
								className="absolute transform -translate-x-1/2 font-semibold text-muted-foreground tracking-[1px] uppercase hover:text-primary transition-colors cursor-pointer"
								style={{
									left: x,
									fontSize: Math.max(
										8,
										Math.min(12, 10 * Math.sqrt(zoomLevel)),
									),
									opacity: Math.max(0.4, Math.min(1, 0.7 * zoomLevel)),
								}}
							>
								{label}
							</button>
						);
					})}
				</div>

				<div
					style={{
						position: "absolute",
						top: 200, // Centered vertically in the 350px container
						left: 0,
						right: 0,
						height: 1, // Ultra thin sleek line
						background:
							"linear-gradient(90deg, transparent, var(--border) 10%, var(--border) 90%, transparent)", // Fade edges
						boxShadow: "0 0 10px var(--muted-foreground)", // Subtle glow
						borderRadius: 0,
						opacity: 0.5,
						zIndex: 0,
					}}
				/>

				{/* Active range highlight on the track */}
				<div
					style={{
						position: "absolute",
						top: 200,
						left: 40,
						right: 40,
						height: 1,
						background:
							"linear-gradient(90deg, transparent, var(--primary) 10%, var(--primary) 90%, transparent)", // Cyan active range with fade
						opacity: 0.6,
						boxShadow: "0 0 15px var(--glow-primary)", // Cyan glow
						zIndex: 1,
					}}
				/>

				{/* Event nodes */}
				{sortedEvents.map((event) => {
					const layout = eventLayout[event.id];
					if (!layout) return null;

					const { x, y, stemHeight } = layout;
					const isSelected = event.id === selectedEventId;

					return (
						<div key={event.id}>
							{/* Vertical Stem */}
							<div
								style={{
									position: "absolute",
									left: x,
									top: y + 100, // Adjusted for new center
									height: stemHeight,
									width: 1,
									// Gradient stem: Fades from node color down to line
									background:
										isSelected || hoveredEventId === event.id
											? "linear-gradient(to bottom, var(--primary), transparent)"
											: "linear-gradient(to bottom, var(--muted-foreground), transparent)",
									transform: "translateX(-50%)",
									transition: "all 0.3s ease",
									pointerEvents: "none",
									zIndex: 0,
								}}
							/>

							{/* Connection Dot on the line */}
							<div
								style={{
									position: "absolute",
									left: x,
									top: 200,
									width: 4,
									height: 4,
									borderRadius: "50%",
									backgroundColor:
										isSelected || hoveredEventId === event.id
											? "var(--primary)"
											: "var(--muted-foreground)",
									transform: "translate(-50%, -50%)",
									transition: "all 0.3s ease",
									zIndex: 2,
								}}
							/>

							<TimelineEventNode
								event={event}
								x={x}
								y={y + 100} // Offset because we increased container height
								isSelected={isSelected}
								onClick={() => handleEventClick(event)}
								onMouseEnter={() => setHoveredEventId(event.id)}
								onMouseLeave={() => setHoveredEventId(null)}
								zoomLevel={zoomLevel}
							/>
						</div>
					);
				})}

				{/* Connection lines to indicate relationships */}
				{(selectedEventId || hoveredEventId) && (
					<svg
						style={{
							position: "absolute",
							top: 0,
							left: 0,
							width: "100%",
							height: "100%",
							pointerEvents: "none",
						}}
					>
						<title>Relationship connections</title>
						{data.edges
							.filter(
								(edge) =>
									edge.fromEventId === (hoveredEventId || selectedEventId) ||
									edge.toEventId === (hoveredEventId || selectedEventId),
							)
							.map((edge) => {
								const fromEvent = data.events.find(
									(e) => e.id === edge.fromEventId,
								);
								const toEvent = data.events.find(
									(e) => e.id === edge.toEventId,
								);
								if (!fromEvent || !toEvent) return null;

								// Arch over the top of events
								// Draw paths between nodes at their respective Y levels
								const fromLayout = eventLayout[fromEvent.id];
								const toLayout = eventLayout[toEvent.id];
								if (!fromLayout || !toLayout) return null;

								const x1 = fromLayout.x;
								const y1 = fromLayout.y + 100;
								const x2 = toLayout.x;
								const y2 = toLayout.y + 100;

								// Control point Y=0 or higher for a nice arc
								const midX = (x1 + x2) / 2;
								const controlY = Math.min(y1, y2) - 60;

								return (
									<path
										key={edge.id}
										d={`M ${x1} ${y1} Q ${midX} ${controlY} ${x2} ${y2}`}
										fill="none"
										stroke={
											hoveredEventId && !selectedEventId
												? "var(--primary)"
												: "var(--primary)"
										}
										style={{
											opacity: hoveredEventId && !selectedEventId ? 0.4 : 0.6,
										}}
										strokeWidth={2}
										strokeDasharray="4 4"
									/>
								);
							})}
					</svg>
				)}
			</div>
			{/* Minimap */}
			{showMinimap && (
				<TimelineMinimap
					data={data}
					minDate={minDate}
					maxDate={maxDate}
					viewportStart={viewportStart}
					viewportWidth={viewportWidth}
					totalWidth={trackWidth}
					onScroll={handleMinimapScroll}
				/>
			)}
		</div>
	);
}
