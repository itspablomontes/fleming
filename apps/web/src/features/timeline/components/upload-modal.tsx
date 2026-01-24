import { Loader2, Upload } from "lucide-react";
import { useEffect, useState } from "react";
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
import { addEvent, correctEvent } from "../api";
import { EVENT_TYPE_LABELS, type TimelineEvent, TimelineEventType as EventTypes } from "../types";

interface UploadModalProps {
	isOpen: boolean;
	onClose: () => void;
	onSuccess?: () => void;
	editEvent?: TimelineEvent | null;
	onReset?: () => void;
}

export function UploadModal({ isOpen, onClose, onSuccess, editEvent, onReset }: UploadModalProps) {
	const [file, setFile] = useState<File | null>(null);
	const [eventType, setEventType] = useState<string>(EventTypes.VISIT_NOTE);
	const [description, setDescription] = useState("");
	const [provider, setProvider] = useState("");
	const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
	const [isUploading, setIsUploading] = useState(false);

	useEffect(() => {
		if (editEvent) {
			setEventType(editEvent.type);
			setDescription(editEvent.description || "");
			setProvider(editEvent.provider || "");
			setDate(new Date(editEvent.timestamp).toISOString().split("T")[0]);
		} else {
			setEventType(EventTypes.VISIT_NOTE);
			setDescription("");
			setProvider("");
			setDate(new Date().toISOString().split("T")[0]);
		}
	}, [editEvent]);

	const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		if (e.target.files?.[0]) {
			setFile(e.target.files[0]);
		}
	};

	const handleUpload = async () => {
		if (!editEvent && !file) return;

		setIsUploading(true);
		try {
			if (editEvent) {
				await correctEvent({
					id: editEvent.id,
					eventType,
					description,
					provider,
					date: new Date(date).toISOString(),
					file: file || undefined,
				});
			} else {
				await addEvent({
					file: file!,
					eventType,
					description,
					provider,
					date: new Date(date).toISOString(),
				});
			}
			onSuccess?.();
			onClose();
			onReset?.();
			// Reset form
			setFile(null);
			setDescription("");
			setProvider("");
		} catch (error) {
			console.error("Operation failed:", error);
			alert(`Failed to ${editEvent ? "correct" : "upload"} document. Please try again.`);
		} finally {
			setIsUploading(false);
		}
	};

	return (
		<Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
			<DialogContent className="sm:max-w-[425px] bg-white dark:bg-gray-950 border-cyan-200 dark:border-cyan-900">
				<DialogHeader>
					<DialogTitle className="text-2xl font-bold text-cyan-900 dark:text-cyan-50">
						{editEvent ? "Correct Medical Entry" : "Upload Medical Document"}
					</DialogTitle>
				</DialogHeader>

				<div className="grid gap-4 py-4">
					<div className="grid gap-2">
						<Label htmlFor="file" className="text-cyan-900 dark:text-cyan-50">
							{editEvent ? "Updated Document (Optional)" : "Document File"}
						</Label>
						<label
							htmlFor="file-input"
							className={`block border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition-colors ${file
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
							<Select value={eventType} onValueChange={setEventType}>
								<SelectTrigger className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800">
									<SelectValue placeholder="Select type" />
								</SelectTrigger>
								<SelectContent className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800">
									{Object.entries(EVENT_TYPE_LABELS).map(([value, label]) => (
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
						<Label htmlFor="description">Short Description</Label>
						<Input
							id="description"
							placeholder="e.g. Annual blood work results"
							value={description}
							onChange={(e) => setDescription(e.target.value)}
							className="bg-white dark:bg-gray-900 border-cyan-200 dark:border-cyan-800"
						/>
					</div>
				</div>

				<DialogFooter>
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
						className="bg-emerald-600 hover:bg-emerald-700 text-white"
					>
						{isUploading ? (
							<>
								<Loader2 className="mr-2 h-4 w-4 animate-spin" />
								{editEvent ? "Saving Correction..." : "Encrypting & Uploading..."}
							</>
						) : (
							editEvent ? "Confirm Correction" : "Encrypt & Upload"
						)}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
