import { Loader2 } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { Logo } from "@/components/common/logo";
import { ThemeToggle } from "@/components/common/theme-toggle";
import { getGraphData } from "../api/get-graph";
import { EventDrawer } from "../components/event-drawer";
import { HorizontalTimeline } from "../components/horizontal-timeline";
import { UploadFAB } from "../components/upload-fab";
import { UploadModal } from "../components/upload-modal";
import type { EventEdge, GraphData, TimelineEvent } from "../types";

export function TimelineViewPage() {
	const [selectedEvent, setSelectedEvent] = useState<TimelineEvent | null>(
		null,
	);
	const [isLoading, setIsLoading] = useState(true);
	const [data, setData] = useState<GraphData>({ events: [], edges: [] });
	const [uploadModalOpen, setUploadModalOpen] = useState(false);

	useEffect(() => {
		const loadaData = async () => {
			try {
				const graphData = await getGraphData();
				setData(graphData);
			} catch (error) {
				console.error("Failed to load timeline data", error);
			} finally {
				setIsLoading(false);
			}
		};
		loadaData();
	}, []);

	const handleEventSelect = useCallback((event: TimelineEvent | null) => {
		setSelectedEvent(event);
	}, []);

	const handleEventClick = useCallback((event: TimelineEvent) => {
		setSelectedEvent(event);
	}, []);

	const handleCloseDrawer = useCallback(() => {
		setSelectedEvent(null);
	}, []);

	// Get related events for drawer
	const relatedEvents = useMemo(() => {
		if (!selectedEvent || !data) return [];

		return data.edges
			.filter(
				(edge) =>
					edge.fromEventId === selectedEvent.id ||
					edge.toEventId === selectedEvent.id,
			)
			.map((edge) => {
				const relatedId =
					edge.fromEventId === selectedEvent.id
						? edge.toEventId
						: edge.fromEventId;
				const relatedEvent = data.events.find((e) => e.id === relatedId);
				if (!relatedEvent) return null;

				return {
					event: relatedEvent,
					edge,
					direction:
						edge.fromEventId === selectedEvent.id
							? ("outgoing" as const)
							: ("incoming" as const),
				};
			})
			.filter(Boolean) as {
				event: TimelineEvent;
				edge: EventEdge;
				direction: "incoming" | "outgoing";
			}[];
	}, [selectedEvent, data]);

	return (
		<div
			className="flex flex-col h-screen bg-background text-foreground overflow-hidden"
		>
			{/* Background Timeline Layer */}
			<div
				className={`absolute inset-0 z-0 flex items-center transition-all duration-700 ease-in-out pointer-events-none ${selectedEvent
					? "opacity-30 -translate-y-1/3 blur-[1.4px]"
					: "opacity-100 translate-y-0"
					}`}
			>
				{isLoading ? (
					<div className="flex items-center justify-center w-full h-full">
						<Loader2 className="w-8 h-8 text-primary animate-spin" />
					</div>
				) : (
					<div className="h-fit w-full pointer-events-auto">
						<HorizontalTimeline
							data={data}
							selectedEventId={selectedEvent?.id || null}
							onEventSelect={handleEventSelect}
						/>
					</div>
				)}
			</div>

			{/* Content Overlays (Header, FABs, etc) */}
			<main className="relative z-10 flex flex-col h-full pointer-events-none">
				{/* Header */}
				<header className="flex items-center justify-between px-6 py-4 bg-background/20 backdrop-blur-md border-b border-border/10 shrink-0 pointer-events-auto">
					<div className="flex items-center gap-4">
						<Logo size="sm" />
						<h1 className="text-lg font-bold tracking-tight text-foreground">
							Medical Timeline
						</h1>
					</div>

					<div className="flex items-center gap-4">
						<div className="text-xs text-muted-foreground hidden sm:block">
							Scroll to explore • Click for details
						</div>
						<ThemeToggle />
					</div>
				</header>
				<div className="flex-1" />{" "}
				{/* Spacer for background timeline visibility */}
			</main>

			{/* Event Drawer (Cluster/Details) */}
			<EventDrawer
				event={selectedEvent}
				relatedEvents={relatedEvents}
				onClose={handleCloseDrawer}
				onEventClick={handleEventClick}
			/>

			{/* Upload FAB */}
			<UploadFAB onClick={() => setUploadModalOpen(true)} />

			{/* Upload Modal */}
			<UploadModal
				isOpen={uploadModalOpen}
				onClose={() => setUploadModalOpen(false)}
				onSuccess={() => {
					// TODO: Refresh timeline data when API is connected
					console.log("Upload successful");
				}}
			/>

			{/* CSS Animation */}
			<style>{`
				@keyframes pulse {
					0%, 100% {
						box-shadow: 0 0 20px var(--glow-primary);
					}
					50% {
						box-shadow: 0 0 35px var(--glow-primary);
					}
				}
			`}</style>
		</div>
	);
}
