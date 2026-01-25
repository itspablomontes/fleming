import {
	Download,
	Edit,
	File,
	FileText,
	Loader2,
	Lock,
	Network,
	Share2,
	ShieldCheck,
	Trash2,
} from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { ConfirmationModal } from "@/components/common/confirmation-modal";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useVault } from "@/features/auth/contexts/vault-context";
import {
	decryptChunkedBuffer,
	decryptFile,
	unwrapKey,
	wrapKey,
} from "@/lib/crypto/encryption";
import { deriveMasterKey } from "@/lib/crypto/keys";
import { getFileKey, shareFileKey } from "../api";
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
	const [shareFile, setShareFile] = useState<EventFile | null>(null);
	const [shareAddress, setShareAddress] = useState("");
	const [shareSignature, setShareSignature] = useState("");
	const [shareSalt, setShareSalt] = useState("");
	const [isSharing, setIsSharing] = useState(false);
	const [shareError, setShareError] = useState<string | null>(null);
	const [downloadError, setDownloadError] = useState<string | null>(null);

	const toBytes = useCallback((hex: string) => {
		const normalized = hex.startsWith("0x") ? hex.slice(2) : hex;
		const bytes =
			normalized.match(/.{1,2}/g)?.map((byte) => Number.parseInt(byte, 16)) ??
			[];
		return new Uint8Array(bytes);
	}, []);

	const toHex = useCallback((buffer: ArrayBuffer) => {
		return Array.from(new Uint8Array(buffer))
			.map((b) => b.toString(16).padStart(2, "0"))
			.join("");
	}, []);

	const downloadFile = useCallback(
		async (file: EventFile) => {
			if (!event || !masterKey) return;

			setIsDownloading(file.id);
			try {
				// 1. Fetch wrapped key for current actor
				const keyResponse = await getFileKey({
					eventId: event.id,
					fileId: file.id,
					patientId: event.patientId,
				});

				// 2. Fetch encrypted file content
				const response = await fetch(
					`/api/timeline/events/${event.id}/files/${file.id}?patientId=${event.patientId}`,
					{ credentials: "include" },
				);
				if (!response.ok) {
					throw new Error("Failed to fetch encrypted file");
				}
				const fullBuffer = await response.arrayBuffer();

				// 3. Unwrap DEK
				const wrappedKeyBytes = toBytes(keyResponse.wrappedKey);
				const dek = await unwrapKey(wrappedKeyBytes, masterKey);

				const metadata = (file.metadata ?? {}) as {
					isMultipart?: boolean;
					chunkSize?: number;
					totalSize?: number;
					ivLength?: number;
				};
				const isMultipart = metadata.isMultipart === true;

				// 4. Decrypt file
				const decryptedBuffer = isMultipart
					? await decryptChunkedBuffer(fullBuffer, dek, {
							chunkSize: Number(metadata.chunkSize),
							totalSize: Number(metadata.totalSize),
							ivLength: Number(metadata.ivLength ?? 12),
						})
					: await decryptFile(
							fullBuffer.slice(12),
							new Uint8Array(fullBuffer.slice(0, 12)),
							dek,
						);

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
				setDownloadError("Failed to decrypt and download file.");
			} finally {
				setIsDownloading(null);
			}
		},
		[event, masterKey, toBytes],
	);

	const handleShare = useCallback(async () => {
		if (!event || !shareFile || !masterKey) return;
		setShareError(null);
		setIsSharing(true);
		try {
			if (!shareAddress.trim() || !shareSignature.trim() || !shareSalt.trim()) {
				throw new Error("Recipient address, signature, and salt are required.");
			}

			const keyResponse = await getFileKey({
				eventId: event.id,
				fileId: shareFile.id,
				patientId: event.patientId,
			});
			const wrappedKeyBytes = toBytes(keyResponse.wrappedKey);
			const dek = await unwrapKey(wrappedKeyBytes, masterKey);

			const saltBytes = new TextEncoder().encode(shareSalt.trim());
			const recipientKey = await deriveMasterKey(
				shareSignature.trim(),
				saltBytes,
			);
			const recipientWrapped = await wrapKey(dek, recipientKey);
			const recipientWrappedHex = toHex(recipientWrapped);

			await shareFileKey({
				eventId: event.id,
				fileId: shareFile.id,
				grantee: shareAddress.trim(),
				wrappedKey: recipientWrappedHex,
			});

			setShareFile(null);
			setShareAddress("");
			setShareSignature("");
			setShareSalt("");
		} catch (error) {
			const message =
				error instanceof Error ? error.message : "Failed to share file.";
			setShareError(message);
		} finally {
			setIsSharing(false);
		}
	}, [
		event,
		masterKey,
		shareAddress,
		shareFile,
		shareSalt,
		shareSignature,
		toBytes,
		toHex,
	]);

	useEffect(() => {
		localStorage.setItem("fleming-drawer-mode", viewMode);
	}, [viewMode]);

	if (!event) return null;

	return (
		<>
			<Sheet open={!!event} onOpenChange={(open) => !open && onClose()}>
				<SheetContent
					side="bottom"
					className="h-[80vh] sm:h-[600px] p-0 gap-0 border-t-2 border-primary/20 bg-background/95 backdrop-blur-xl"
				>
					<div className="flex h-full flex-col">
						{/* Header with Switch */}
						<div className="flex items-center justify-between px-6 py-4 pr-12 border-b border-border/50">
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

										{event.metadata &&
											Object.keys(event.metadata).length > 0 && (
												<section>
													<h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider mb-3">
														Metadata
													</h3>
													<div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
														{Object.entries(event.metadata).map(
															([key, value]) => (
																<div
																	key={key}
																	className="flex flex-col p-2 rounded border border-border/50 bg-black/10"
																>
																	<span className="text-[10px] text-muted-foreground uppercase">
																		{key}
																	</span>
																	<span className="text-sm truncate">
																		{String(value)}
																	</span>
																</div>
															),
														)}
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
																	<div
																		className="font-semibold text-foreground truncate max-w-[200px]"
																		title={file.fileName}
																	>
																		{file.fileName}
																	</div>
																	<div className="text-xs text-muted-foreground flex items-center gap-2">
																		<span>{file.mimeType}</span>
																		<span>•</span>
																		<span>
																			{(file.fileSize / 1024).toFixed(1)} KB
																		</span>
																	</div>
																</div>
															</div>
															<div className="flex items-center gap-2">
																<button
																	type="button"
																	onClick={() => setShareFile(file)}
																	disabled={!isUnlocked}
																	className="p-2.5 rounded-full bg-secondary/30 text-secondary-foreground hover:bg-secondary disabled:opacity-50 disabled:cursor-not-allowed transition-all"
																	title={
																		isUnlocked
																			? "Share encrypted file"
																			: "Unlock Vault to Share"
																	}
																	aria-label="Share encrypted file"
																>
																	<Share2 className="w-5 h-5" />
																</button>
																<button
																	type="button"
																	onClick={() => downloadFile(file)}
																	disabled={!!isDownloading || !isUnlocked}
																	className="p-2.5 rounded-full bg-primary/10 text-primary hover:bg-primary hover:text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all"
																	title={
																		isUnlocked
																			? "Decrypt & Download"
																			: "Unlock Vault to Download"
																	}
																	aria-label="Decrypt & Download"
																>
																	{isDownloading === file.id ? (
																		<Loader2 className="w-5 h-5 animate-spin" />
																	) : (
																		<Download className="w-5 h-5" />
																	)}
																</button>
															</div>
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

			<Dialog
				open={!!shareFile}
				onOpenChange={(open) => !open && setShareFile(null)}
			>
				<DialogContent className="max-w-md">
					<DialogHeader>
						<DialogTitle>Share encrypted file</DialogTitle>
					</DialogHeader>
					<div className="grid gap-4">
						<div className="grid gap-2">
							<Label htmlFor="share-address">Recipient address</Label>
							<Input
								id="share-address"
								value={shareAddress}
								onChange={(e) => setShareAddress(e.target.value)}
								placeholder="0x..."
							/>
						</div>
						<div className="grid gap-2">
							<Label htmlFor="share-signature">Recipient signature</Label>
							<Input
								id="share-signature"
								value={shareSignature}
								onChange={(e) => setShareSignature(e.target.value)}
								placeholder="Signature used to derive recipient key"
							/>
						</div>
						<div className="grid gap-2">
							<Label htmlFor="share-salt">Recipient encryption salt</Label>
							<Input
								id="share-salt"
								value={shareSalt}
								onChange={(e) => setShareSalt(e.target.value)}
								placeholder="Recipient salt from /api/auth/me"
							/>
						</div>
						{shareError && (
							<div className="text-sm text-destructive">{shareError}</div>
						)}
					</div>
					<DialogFooter className="mt-4">
						<Button variant="outline" onClick={() => setShareFile(null)}>
							Cancel
						</Button>
						<Button onClick={handleShare} disabled={isSharing}>
							{isSharing ? "Sharing..." : "Share"}
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>

			<ConfirmationModal
				isOpen={!!downloadError}
				onOpenChange={(open) => !open && setDownloadError(null)}
				onConfirm={() => setDownloadError(null)}
				title="Download failed"
				message={downloadError ?? "Unable to download file."}
				confirmLabel="OK"
				cancelLabel="Close"
				variant="warning"
			/>
		</>
	);
}
