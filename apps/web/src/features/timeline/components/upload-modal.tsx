import { Loader2, Lock, Plus, Trash2, Upload } from "lucide-react";
import { useEffect, useState } from "react";
import { VaultUnlockDialog } from "@/components/common/vault-unlock-dialog";
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
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { useVault } from "@/features/auth/contexts/vault-context";
import { encryptFile, wrapKey } from "@/lib/crypto/encryption";
import { generateDEK } from "@/lib/crypto/keys";
import { addEvent, correctEvent } from "../api";
import {
	EVENT_TYPE_LABELS,
	TimelineEventType as EventTypes,
	type TimelineEvent,
} from "../types";

interface UploadModalProps {
	isOpen: boolean;
	onClose: () => void;
	onSuccess?: () => void;
	editEvent?: TimelineEvent | null;
	onReset?: () => void;
}

export function UploadModal({
	isOpen,
	onClose,
	onSuccess,
	editEvent,
	onReset,
}: UploadModalProps) {
	const { isUnlocked, masterKey } = useVault();
	const [showUnlockDialog, setShowUnlockDialog] = useState(false);

	const [file, setFile] = useState<File | null>(null);
	const [eventType, setEventType] = useState<EventTypes>(EventTypes.VISIT_NOTE);
	const [title, setTitle] = useState("");
	const [description, setDescription] = useState("");
	const [provider, setProvider] = useState("");
	const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
	const [metadata, setMetadata] = useState<{ key: string; value: string }[]>(
		[],
	);
	const [isUploading, setIsUploading] = useState(false);
	const [uploadStatus, setUploadStatus] = useState<string>("");
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		if (editEvent) {
			setEventType(editEvent.type);
			setTitle(editEvent.title || "");
			setDescription(editEvent.description || "");
			setProvider(editEvent.provider || "");
			setDate(new Date(editEvent.timestamp).toISOString().split("T")[0]);
			// Convert object metadata to array for editing
			if (editEvent.metadata) {
				setMetadata(
					Object.entries(editEvent.metadata).map(([key, value]) => ({
						key,
						value: String(value),
					})),
				);
			} else {
				setMetadata([]);
			}
		} else {
			setEventType(EventTypes.VISIT_NOTE);
			setTitle("");
			setDescription("");
			setProvider("");
			setDate(new Date().toISOString().split("T")[0]);
			setMetadata([]);
		}
		setError(null);
	}, [editEvent]);

	const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		if (e.target.files?.[0]) {
			setFile(e.target.files[0]);
			// Auto-fill title if empty
			if (!title) {
				const name = e.target.files[0].name.replace(/\.[^/.]+$/, ""); // remove extension
				setTitle(
					name
						.split(/[-_]/)
						.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
						.join(" "),
				);
			}
		}
	};

	const handleUpload = async () => {
		setError(null);

		// 1. Vault Check
		if (!isUnlocked || !masterKey) {
			setShowUnlockDialog(true);
			return;
		}

		if (!editEvent && !file) {
			setError("Please select a file to upload.");
			return;
		}

		if (!title.trim()) {
			setError("Please provide a title for this event.");
			return;
		}

		setIsUploading(true);
		setUploadStatus("Starting...");

		try {
			// Convert metadata array back to object
			const metadataObj: Record<string, unknown> = {};
			for (const item of metadata) {
				if (item.key.trim()) {
					metadataObj[item.key.trim()] = item.value;
				}
			}

			let payloadFile: File | Blob | undefined = file || undefined;
			let isEncrypted = false;
			let wrappedKeyHex: string | undefined;

			// 2. Encryption Step (if new file present)
			if (file) {
				setUploadStatus("Generating keys...");
				const dek = await generateDEK(); // 256-bit AES-GCM

				setUploadStatus("Encrypting file...");
				const fileData = await file.arrayBuffer();
				const { ciphertext, iv } = await encryptFile(fileData, dek);

				setUploadStatus("Securing keys...");
				const wrappedKey = await wrapKey(dek, masterKey);

				// Combine IV + Ciphertext for storage
				// Format: [IV (12 bytes)] [Ciphertext (N bytes)]
				const encryptedBlob = new Blob([iv.buffer as ArrayBuffer, ciphertext], {
					type: "application/octet-stream",
				});
				payloadFile = encryptedBlob;
				isEncrypted = true;

				// Convert wrapped key to hex for transport
				wrappedKeyHex = Array.from(new Uint8Array(wrappedKey))
					.map((b) => b.toString(16).padStart(2, "0"))
					.join("");
			}

			setUploadStatus("Uploading to vault...");

			if (editEvent) {
				await correctEvent({
					id: editEvent.id,
					eventType,
					title,
					description,
					provider,
					date: new Date(date).toISOString(),
					metadata: metadataObj,
					file: payloadFile,
					isEncrypted,
					wrappedKey: wrappedKeyHex,
				});
			} else {
				await addEvent({
					file: payloadFile as File | Blob,
					eventType,
					title,
					description,
					provider,
					date: new Date(date).toISOString(),
					metadata: metadataObj,
					isEncrypted,
					wrappedKey: wrappedKeyHex,
				});
			}
			onSuccess?.();
			onClose();
			onReset?.();

			// Reset form
			setFile(null);
			setTitle("");
			setDescription("");
			setProvider("");
			setUploadStatus("");
			setMetadata([]);
		} catch (err) {
			console.error("Operation failed:", err);
			setError(
				`Failed to ${editEvent ? "correct" : "upload"} document. Please try again.`,
			);
		} finally {
			setIsUploading(false);
		}
	};

	return (
		<>
			<VaultUnlockDialog
				isOpen={showUnlockDialog}
				onOpenChange={setShowUnlockDialog}
				onSuccess={handleUpload} // Retry upload after unlock
			/>

			<Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
				<DialogContent className="sm:max-w-[500px] bg-white dark:bg-gray-950 border-cyan-200 dark:border-cyan-900 max-h-[90vh] overflow-hidden flex flex-col">
					<DialogHeader>
						<DialogTitle className="text-2xl font-bold text-cyan-900 dark:text-cyan-50 flex items-center gap-2">
							{editEvent ? "Correct Medical Entry" : "Upload Medical Document"}
							{isUnlocked ? (
								<Lock className="h-5 w-5 text-emerald-500" />
							) : (
								<Lock className="h-5 w-5 text-amber-500" />
							)}
						</DialogTitle>
					</DialogHeader>

					<div className="flex-1 overflow-y-auto pr-4 -mr-4">
						<div className="grid gap-4 py-4 px-1">
							<div className="grid gap-2">
								<Label
									htmlFor="file"
									className="text-cyan-900 dark:text-cyan-50"
								>
									{editEvent ? "Updated Document (Optional)" : "Document File"}
								</Label>
								<label
									htmlFor="file-input"
									className={`block border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition-colors ${
										file
											? "border-emerald-500 bg-emerald-50/50 dark:bg-emerald-900/10"
											: "border-cyan-200 dark:border-cyan-800 hover:border-cyan-400 dark:hover:border-cyan-600"
									}`}
								>
									<input
										id="file-input"
										type="file"
										onChange={handleFileChange}
										className="hidden"
									/>
									{file ? (
										<div className="flex flex-col items-center">
											<Upload className="h-8 w-8 text-emerald-600 mb-2" />
											<p className="text-sm font-medium text-emerald-900 dark:text-emerald-50 truncate max-w-full">
												{file.name}
											</p>
											<p className="text-xs text-emerald-600 mt-1">
												{(file.size / 1024 / 1024).toFixed(2)} MB
											</p>
											{isUnlocked && (
												<div className="mt-2 flex items-center gap-1 text-xs text-emerald-600 bg-emerald-100 dark:bg-emerald-900/30 px-2 py-1 rounded-full">
													<Lock className="h-3 w-3" />
													<span>Will be encrypted locally</span>
												</div>
											)}
										</div>
									) : (
										<div className="flex flex-col items-center">
											<Upload className="h-8 w-8 text-cyan-500 mb-2" />
											<p className="text-sm font-medium text-cyan-900 dark:text-cyan-50">
												Click to select or drag and drop
											</p>
											<p className="text-xs text-cyan-600 mt-1">
												PDF, JPG, PNG (up to 10MB)
											</p>
										</div>
									)}
								</label>
							</div>

							<div className="grid grid-cols-2 gap-4">
								<div className="grid gap-2">
									<Label htmlFor="type">Event Type</Label>
									<Select
										value={eventType}
										onValueChange={(value) => setEventType(value as EventTypes)}
									>
										<SelectTrigger className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800">
											<SelectValue placeholder="Select type" />
										</SelectTrigger>
										<SelectContent className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800">
											{(
												Object.entries(EVENT_TYPE_LABELS) as [
													EventTypes,
													string,
												][]
											).map(([value, label]) => (
												<SelectItem key={value} value={value}>
													{label}
												</SelectItem>
											))}
										</SelectContent>
									</Select>
								</div>
								<div className="grid gap-2">
									<Label htmlFor="date">Service Date</Label>
									<Input
										id="date"
										type="date"
										value={date}
										onChange={(e) => setDate(e.target.value)}
										className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800"
									/>
								</div>
							</div>

							<div className="grid gap-2">
								<Label htmlFor="provider">Provider / Facility</Label>
								<Input
									id="provider"
									placeholder="e.g. Dr. Alice Smith or City Hospital"
									value={provider}
									onChange={(e) => setProvider(e.target.value)}
									className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800"
								/>
							</div>

							<div className="grid gap-2">
								<Label htmlFor="title">Title</Label>
								<Input
									id="title"
									placeholder="e.g. Lab Results - Jan 2024"
									value={title}
									onChange={(e) => setTitle(e.target.value)}
									className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800"
								/>
							</div>

							<div className="grid gap-2">
								<Label htmlFor="description">Short Description</Label>
								<Input
									id="description"
									placeholder="e.g. Annual blood work results"
									value={description}
									onChange={(e) => setDescription(e.target.value)}
									className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800"
								/>
							</div>

							<div className="grid gap-2 mt-2">
								<div className="flex items-center justify-between">
									<Label className="text-sm font-medium">
										Additional Details
									</Label>
									<Button
										type="button"
										variant="ghost"
										size="sm"
										onClick={() =>
											setMetadata([...metadata, { key: "", value: "" }])
										}
										className="h-8 text-cyan-600 hover:text-cyan-700"
									>
										<Plus className="h-4 w-4 mr-1" /> Add field
									</Button>
								</div>
								<div className="space-y-2">
									{metadata.map((item, index) => (
										<div key={`${index}-${item.key}`} className="flex gap-2">
											<Input
												placeholder="Key"
												value={item.key}
												onChange={(e) => {
													const newMetadata = [...metadata];
													newMetadata[index] = {
														...newMetadata[index],
														key: e.target.value,
													};
													setMetadata(newMetadata);
												}}
												className="flex-1 bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800 h-8 text-xs"
											/>
											<Input
												placeholder="Value"
												value={item.value}
												onChange={(e) => {
													const newMetadata = [...metadata];
													newMetadata[index] = {
														...newMetadata[index],
														value: e.target.value,
													};
													setMetadata(newMetadata);
												}}
												className="flex-1 bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800 h-8 text-xs"
											/>
											<Button
												type="button"
												variant="ghost"
												size="sm"
												onClick={() =>
													setMetadata(metadata.filter((_, i) => i !== index))
												}
												className="h-8 w-8 p-0 text-red-500 hover:text-red-600 hover:bg-neutral-100 dark:hover:bg-neutral-800"
											>
												<Trash2 className="h-4 w-4" />
											</Button>
										</div>
									))}
								</div>
							</div>
						</div>
					</div>

					<DialogFooter className="mt-4 flex-col gap-2">
						{error && (
							<div className="w-full p-2 mb-2 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded text-sm border border-red-200 dark:border-red-900 text-center">
								{error}
							</div>
						)}
						<div className="flex w-full justify-end gap-2">
							<Button
								variant="outline"
								onClick={onClose}
								disabled={isUploading}
								className="border-cyan-200 dark:border-cyan-800"
							>
								Cancel
							</Button>
							<Button
								onClick={handleUpload}
								disabled={(!editEvent && !file) || isUploading}
								className="bg-emerald-600 hover:bg-emerald-700 text-white min-w-[140px]"
							>
								{isUploading ? (
									<>
										<Loader2 className="mr-2 h-4 w-4 animate-spin" />
										{uploadStatus || "Processing..."}
									</>
								) : // Dynamic label based on lock state
								isUnlocked ? (
									editEvent ? (
										"Confirm Correction"
									) : (
										"Encrypt & Upload"
									)
								) : (
									"Unlock Vault"
								)}
							</Button>
						</div>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</>
	);
}
