import { Loader2 } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { ConfirmationModal } from "@/components/common/confirmation-modal";
import { Logo } from "@/components/common/logo";
import { ThemeToggle } from "@/components/common/theme-toggle";
import { deleteEvent, getGraphData } from "../api";
import { EventDrawer } from "../components/event-drawer";
import { HorizontalTimeline } from "../components/horizontal-timeline";
import { TimelineItem } from "../components/timeline-item";
import { UploadFAB } from "../components/upload-fab";
import { UploadModal } from "../components/upload-modal";
import { useEditStore } from "@/features/timeline/stores/edit-store";
import { useTimelineCoordinator } from "@/features/timeline/stores/timeline-coordinator";
import type { EventEdge, GraphData, TimelineEvent } from "../types";

export function TimelineViewPage(): JSX.Element {
	const [selectedEvent, setSelectedEvent] = useState<TimelineEvent | null>(
		null,
	);
	const [isLoading, setIsLoading] = useState(true);
	const [data, setData] = useState<GraphData>({ events: [], edges: [] });
	const [uploadModalOpen, setUploadModalOpen] = useState(false);
	const [archiveTarget, setArchiveTarget] = useState<TimelineEvent | null>(null);
	const [archiveError, setArchiveError] = useState<string | null>(null);
	const [isArchiving, setIsArchiving] = useState(false);
	const cancelEdit = useEditStore((state) => state.cancelEdit);
	const startEdit = useTimelineCoordinator((state) => state.startEdit);
	const startUpload = useTimelineCoordinator((state) => state.startUpload);
	const resetAll = useTimelineCoordinator((state) => state.resetAll);

	const refreshData = useCallback(async () => {
		try {
			setIsLoading(true);
			const graphData = await getGraphData();
			setData(graphData);
		} catch (error) {
			console.error("Failed to load timeline data", error);
		} finally {
			setIsLoading(false);
		}
	}, []);

	useEffect(() => {
		refreshData();
	}, [refreshData]);

	const handleEventSelect = useCallback((event: TimelineEvent | null) => {
		setSelectedEvent(event);
	}, []);

	const handleEventClick = useCallback((event: TimelineEvent) => {
		setSelectedEvent(event);
	}, []);

	const handleEdit = useCallback((event: TimelineEvent) => {
		setSelectedEvent(null);
		startEdit(event);
		setUploadModalOpen(true);
	}, [startEdit]);

	const handleArchive = useCallback((event: TimelineEvent) => {
		setArchiveTarget(event);
	}, []);

	const confirmArchive = useCallback(async () => {
		if (!archiveTarget) return;
		setIsArchiving(true);
		try {
			await deleteEvent(archiveTarget.id);
			await refreshData();
			setSelectedEvent(null);
			setArchiveTarget(null);
		} catch (error) {
			console.error("Failed to archive event", error);
			setArchiveError("Failed to archive event. Please try again.");
		} finally {
			setIsArchiving(false);
		}
	}, [archiveTarget, refreshData]);

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
		<div className="flex flex-col h-screen bg-background text-foreground overflow-hidden relative">
			{/* Header */}
			<header className="flex items-center justify-between px-4 py-4 sm:px-6 bg-background/70 backdrop-blur-md border-b border-border/10 shrink-0 z-10">
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

			{/* Desktop Timeline */}
			<div
				className={`absolute inset-0 top-[72px] z-0 hidden md:flex items-center transition-all duration-700 ease-in-out pointer-events-none ${selectedEvent
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

			{/* Mobile List */}
			<main className="flex-1 md:hidden overflow-y-auto px-4 py-6 space-y-4 safe-area-bottom">
				{isLoading ? (
					<div className="flex items-center justify-center h-full">
						<Loader2 className="w-8 h-8 text-primary animate-spin" />
					</div>
				) : (
					<>
						{data.events.length === 0 ? (
							<div className="text-sm text-muted-foreground text-center py-12">
								No timeline events yet. Upload your first record to begin.
							</div>
						) : (
							data.events.map((event) => (
								<TimelineItem
									key={event.id}
									event={event}
									onView={handleEventClick}
								/>
							))
						)}
					</>
				)}
			</main>

			{/* Event Drawer (Cluster/Details) */}
			<EventDrawer
				event={selectedEvent}
				relatedEvents={relatedEvents}
				onClose={handleCloseDrawer}
				onEventClick={handleEventClick}
				onEdit={handleEdit}
				onArchive={handleArchive}
			/>

			{/* Upload FAB */}
			<UploadFAB
				onClick={() => {
					startUpload();
					setUploadModalOpen(true);
				}}
			/>

			{/* Upload Modal */}
			<UploadModal
				isOpen={uploadModalOpen}
				onClose={() => {
					setUploadModalOpen(false);
					resetAll();
				}}
				onSuccess={() => {
					refreshData();
					console.log("Operation successful");
					cancelEdit();
				}}
			/>

			<ConfirmationModal
				isOpen={!!archiveTarget}
				onOpenChange={(open) => !open && setArchiveTarget(null)}
				onConfirm={confirmArchive}
				title="Archive this record?"
				message="This action is immutable and will be recorded in your audit trail."
				confirmLabel="Archive"
				cancelLabel="Cancel"
				variant="warning"
				isConfirming={isArchiving}
			/>

			<ConfirmationModal
				isOpen={!!archiveError}
				onOpenChange={(open) => !open && setArchiveError(null)}
				onConfirm={() => setArchiveError(null)}
				title="Archive failed"
				message={archiveError ?? "Unable to archive record."}
				confirmLabel="OK"
				cancelLabel="Close"
				variant="warning"
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
