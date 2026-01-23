import { FileText, Network } from "lucide-react";
import { useEffect, useState } from "react";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { EventEdge, TimelineEvent } from "../types";
import { RelationshipCluster } from "./relationship-cluster";

interface EventDrawerProps {
	event: TimelineEvent | null;
	relatedEvents: Array<{
		event: TimelineEvent;
		edge: EventEdge;
		direction: "incoming" | "outgoing";
	}>;
	onClose: () => void;
	onEventClick: (event: TimelineEvent) => void;
}

export function EventDrawer({
	event,
	relatedEvents,
	onClose,
	onEventClick,
}: EventDrawerProps) {
	const [viewMode, setViewMode] = useState<"cluster" | "details">(() => {
		if (typeof window !== "undefined") {
			const stored = localStorage.getItem("fleming-drawer-mode");
			return stored === "cluster" || stored === "details" ? stored : "cluster";
		}
		return "cluster";
	});

	useEffect(() => {
		localStorage.setItem("fleming-drawer-mode", viewMode);
	}, [viewMode]);

	if (!event) return null;

	return (
		<Sheet open={!!event} onOpenChange={(open) => !open && onClose()}>
			<SheetContent
				side="bottom"
				className="h-[80vh] sm:h-[600px] p-0 gap-0 border-t-2 border-primary/20 bg-background/95 backdrop-blur-xl"
			>
				<div className="flex h-full flex-col">
					{/* Header with Switch */}
					<div className="flex items-center justify-between px-6 py-4 border-b border-border/50">
						<div className="flex flex-col">
							<SheetTitle className="text-xl font-bold">
								{event.title}
							</SheetTitle>
							<SheetDescription className="text-xs font-mono text-muted-foreground mt-1">
								{new Date(event.timestamp).toLocaleDateString(undefined, {
									weekday: "long",
									year: "numeric",
									month: "long",
									day: "numeric",
								})}
							</SheetDescription>
						</div>

						<div className="flex items-center gap-4">
							<Tabs
								value={viewMode}
								onValueChange={(v: string) =>
									setViewMode(v as "cluster" | "details")
								}
								className="w-[200px]"
							>
								<TabsList className="grid w-full grid-cols-2">
									<TabsTrigger
										value="cluster"
										className="flex items-center gap-2"
									>
										<Network className="w-3.5 h-3.5" />
										Cluster
									</TabsTrigger>
									<TabsTrigger
										value="details"
										className="flex items-center gap-2"
									>
										<FileText className="w-3.5 h-3.5" />
										Details
									</TabsTrigger>
								</TabsList>
							</Tabs>
						</div>
					</div>

					{/* Content Area */}
					<div className="flex-1 overflow-hidden relative">
						{viewMode === "cluster" ? (
							<div className="h-full w-full bg-black/20">
								<RelationshipCluster
									centerEvent={event}
									relatedEvents={relatedEvents}
									onEventClick={onEventClick}
								/>
							</div>
						) : (
							<div className="h-full overflow-y-auto p-6">
								<div className="max-w-2xl mx-auto space-y-8">
									<section>
										<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider mb-3">
											Description
										</h3>
										<div className="prose prose-invert max-w-none">
											<p className="text-base leading-relaxed text-foreground/90">
												{event.description || "No description provided."}
											</p>
										</div>
									</section>

									{event.provider && (
										<section>
											<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider mb-3">
												Provider
											</h3>
											<div className="flex items-center gap-3 p-3 rounded-lg border border-border bg-card/50">
												<div className="h-10 w-10 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold">
													{event.provider[0]}
												</div>
												<div>
													<div className="font-medium">{event.provider}</div>
													<div className="text-xs text-muted-foreground">
														Healthcare Provider
													</div>
												</div>
											</div>
										</section>
									)}

									{event.metadata && (
										<section>
											<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider mb-3">
												Metadata
											</h3>
											<pre className="p-4 rounded-lg bg-black/40 font-mono text-xs overflow-x-auto text-muted-foreground">
												{JSON.stringify(event.metadata, null, 2)}
											</pre>
										</section>
									)}
								</div>
							</div>
						)}
					</div>
				</div>
			</SheetContent>
		</Sheet>
	);
}
