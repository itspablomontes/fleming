import { Download, Edit, File, FileText, Loader2, Lock, Network, ShieldCheck, Trash2 } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useVault } from "@/features/auth/contexts/vault-context";
import { decryptFile, unwrapKey } from "@/lib/crypto/encryption";
import type { EventEdge, EventFile, TimelineEvent } from "../types";
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
	onEdit?: (event: TimelineEvent) => void;
	onArchive?: (event: TimelineEvent) => void;
}

export function EventDrawer({
	event,
	relatedEvents,
	onClose,
	onEventClick,
	onEdit,
	onArchive,
}: EventDrawerProps) {
	const [viewMode, setViewMode] = useState<"cluster" | "details">(() => {
		if (typeof window !== "undefined") {
			const stored = localStorage.getItem("fleming-drawer-mode");
			return stored === "cluster" || stored === "details" ? stored : "cluster";
		}
		return "cluster";
	});
	const { masterKey, isUnlocked } = useVault();
	const [isDownloading, setIsDownloading] = useState<string | null>(null);

	const downloadFile = useCallback(
		async (file: EventFile) => {
			if (!event || !masterKey) return;

			setIsDownloading(file.id);
			try {
				// 1. Fetch encrypted file content
				const response = await fetch(
					`/api/timeline/events/${event.id}/files/${file.id}`,
					{ credentials: "include" }
				);
				if (!response.ok) {
					throw new Error("Failed to fetch encrypted file");
				}
				const fullBuffer = await response.arrayBuffer();

				// 2. Extract IV (12 bytes) and Ciphertext
				const iv = new Uint8Array(fullBuffer.slice(0, 12));
				const ciphertext = fullBuffer.slice(12);

				// 3. Unwrap DEK
				const wrappedKeyHex = file.wrappedDek;
				if (!wrappedKeyHex) throw new Error("File is missing wrapped key");

				const wrappedKeyBytes = new Uint8Array(
					wrappedKeyHex.match(/.{1,2}/g)?.map((byte) => parseInt(byte, 16)) || []
				);

				const dek = await unwrapKey(wrappedKeyBytes, masterKey);

				// 4. Decrypt file
				const decryptedBuffer = await decryptFile(ciphertext, iv, dek);

				// 5. Create a download link
				const decryptedBlob = new Blob([decryptedBuffer], {
					type: file.mimeType,
				});
				const url = URL.createObjectURL(decryptedBlob);
				const a = document.createElement("a");
				a.href = url;
				a.download = file.fileName;
				document.body.appendChild(a);
				a.click();
				document.body.removeChild(a);
				URL.revokeObjectURL(url);
			} catch (error) {
				console.error("Error downloading file:", error);
				alert("Failed to decrypt and download file.");
			} finally {
				setIsDownloading(null);
			}
		},
		[event, masterKey]
	);

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
							{event.type !== "tombstone" && (
								<div className="flex items-center gap-1 bg-muted/50 rounded-lg p-1 border border-border/50">
									<button
										type="button"
										onClick={() => onEdit?.(event)}
										className="p-1.5 hover:bg-background rounded-md transition-colors text-muted-foreground hover:text-foreground"
										title="Correct/Edit Entry"
									>
										<Edit className="w-4 h-4" />
									</button>
									<button
										type="button"
										onClick={() => onArchive?.(event)}
										className="p-1.5 hover:bg-background rounded-md transition-colors text-muted-foreground hover:text-destructive"
										title="Archive Entry"
									>
										<Trash2 className="w-4 h-4" />
									</button>
								</div>
							)}
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

									{event.metadata && Object.keys(event.metadata).length > 0 && (
										<section>
											<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider mb-3">
												Metadata
											</h3>
											<div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
												{Object.entries(event.metadata).map(([key, value]) => (
													<div key={key} className="flex flex-col p-2 rounded border border-border/50 bg-black/10">
														<span className="text-[10px] text-muted-foreground uppercase">{key}</span>
														<span className="text-sm truncate">{String(value)}</span>
													</div>
												))}
											</div>
										</section>
									)}

									{event.files && event.files.length > 0 && (
										<section>
											<div className="flex items-center justify-between mb-3">
												<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
													Secure Attachments
												</h3>
												{isUnlocked ? (
													<div className="flex items-center gap-1 text-[10px] text-emerald-500 font-medium bg-emerald-500/10 px-2 py-0.5 rounded-full border border-emerald-500/20">
														<ShieldCheck className="w-3 h-3" />
														END-TO-END ENCRYPTED
													</div>
												) : (
													<div className="flex items-center gap-1 text-[10px] text-amber-500 font-medium bg-amber-500/10 px-2 py-0.5 rounded-full border border-amber-500/20">
														<Lock className="w-3 h-3" />
														VAULT LOCKED
													</div>
												)}
											</div>
											<div className="space-y-3">
												{event.files.map((file) => (
													<div
														key={file.id}
														className="group flex items-center justify-between p-4 rounded-xl border border-border bg-card/50 hover:bg-card hover:border-primary/30 transition-all"
													>
														<div className="flex items-center gap-4 overflow-hidden">
															<div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center text-primary shrink-0">
																<File className="w-6 h-6" />
															</div>
															<div className="overflow-hidden">
																<div className="font-semibold text-foreground truncate">
																	{file.fileName}
																</div>
																<div className="text-xs text-muted-foreground flex items-center gap-2">
																	<span>{file.mimeType}</span>
																	<span>•</span>
																	<span>{(file.fileSize / 1024).toFixed(1)} KB</span>
																</div>
															</div>
														</div>
														<button
															type="button"
															onClick={() => downloadFile(file)}
															disabled={!!isDownloading || !isUnlocked}
															className="p-2.5 rounded-full bg-primary/10 text-primary hover:bg-primary hover:text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all"
															title={isUnlocked ? "Decrypt & Download" : "Unlock Vault to Download"}
														>
															{isDownloading === file.id ? (
																<Loader2 className="w-5 h-5 animate-spin" />
															) : (
																<Download className="w-5 h-5" />
															)}
														</button>
													</div>
												))}
											</div>
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
