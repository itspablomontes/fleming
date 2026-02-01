/**
 * Relationship Type Picker Modal
 *
 * Modal that appears after selecting two events in link mode.
 * Shows a visual preview of the connection and relationship type options.
 */

import {
	AlertTriangle,
	ArrowRight,
	Check,
	FileCheck,
	Link2,
	Loader2,
	RefreshCw,
	Undo2,
	Zap,
} from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { linkEvents } from "../api";
import { useTimelineCoordinator } from "../stores/timeline-coordinator";
import {
	EVENT_TYPE_COLORS,
	EVENT_TYPE_LABELS,
	RELATIONSHIP_LABELS,
	type RelationshipType,
	type TimelineEvent,
} from "../types";

// Icons for each relationship type
const RELATIONSHIP_ICONS: Record<RelationshipType, typeof ArrowRight> = {
	resulted_in: ArrowRight,
	lead_to: ArrowRight,
	requested_by: Undo2,
	supports: FileCheck,
	follows_up: RefreshCw,
	contradicts: AlertTriangle,
	attached_to: Link2,
	replaces: RefreshCw,
	caused_by: Zap,
};

interface RelationshipTypePickerProps {
	isOpen: boolean;
	source: TimelineEvent | null;
	target: TimelineEvent | null;
	onClose: () => void;
	onSuccess: () => void;
}

export function RelationshipTypePicker({
	isOpen,
	source,
	target,
	onClose,
	onSuccess,
}: RelationshipTypePickerProps) {
	const [selectedType, setSelectedType] = useState<RelationshipType | null>(
		null,
	);
	const [isCreating, setIsCreating] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const completeLinking = useTimelineCoordinator(
		(state) => state.completeLinking,
	);

	const handleCreate = async () => {
		if (!source || !target || !selectedType) return;

		setIsCreating(true);
		setError(null);

		try {
			await linkEvents({
				fromEventId: source.id,
				toEventId: target.id,
				relationshipType: selectedType,
			});
			completeLinking();
			onSuccess();
			onClose();
		} catch (err) {
			console.error("Failed to link events:", err);
			setError("Failed to create relationship. Please try again.");
		} finally {
			setIsCreating(false);
		}
	};

	const handleClose = () => {
		setSelectedType(null);
		setError(null);
		onClose();
	};

	if (!source || !target) return null;

	const sourceColors = EVENT_TYPE_COLORS[source.type];
	const targetColors = EVENT_TYPE_COLORS[target.type];

	return (
		<Dialog open={isOpen} onOpenChange={(open) => !open && handleClose()}>
			<DialogContent className="sm:max-w-[500px] bg-background border-border">
				<DialogHeader>
					<DialogTitle className="text-xl font-bold flex items-center gap-2">
						<Link2 className="h-5 w-5 text-primary" />
						Create Relationship
					</DialogTitle>
				</DialogHeader>

				{/* Connection Preview */}
				<div className="flex items-center justify-center gap-4 py-6 px-4 bg-muted/30 rounded-lg border border-border/50">
					{/* Source Event */}
					<div
						className="flex flex-col items-center gap-2 p-3 rounded-lg border-2 max-w-[140px]"
						style={{
							borderColor: sourceColors.border,
							backgroundColor: sourceColors.bg,
						}}
					>
						<span
							className="text-xs font-medium uppercase tracking-wide"
							style={{ color: sourceColors.text }}
						>
							{EVENT_TYPE_LABELS[source.type]}
						</span>
						<span className="text-sm font-semibold text-center line-clamp-2">
							{source.title}
						</span>
						<span className="text-xs text-muted-foreground">
							{new Date(source.timestamp).toLocaleDateString()}
						</span>
					</div>

					{/* Arrow */}
					<div className="flex flex-col items-center gap-1">
						<ArrowRight className="h-6 w-6 text-primary" />
						{selectedType && (
							<span className="text-xs text-primary font-medium">
								{RELATIONSHIP_LABELS[selectedType]}
							</span>
						)}
					</div>

					{/* Target Event */}
					<div
						className="flex flex-col items-center gap-2 p-3 rounded-lg border-2 max-w-[140px]"
						style={{
							borderColor: targetColors.border,
							backgroundColor: targetColors.bg,
						}}
					>
						<span
							className="text-xs font-medium uppercase tracking-wide"
							style={{ color: targetColors.text }}
						>
							{EVENT_TYPE_LABELS[target.type]}
						</span>
						<span className="text-sm font-semibold text-center line-clamp-2">
							{target.title}
						</span>
						<span className="text-xs text-muted-foreground">
							{new Date(target.timestamp).toLocaleDateString()}
						</span>
					</div>
				</div>

				{/* Relationship Type Grid */}
				<div className="grid grid-cols-3 gap-2 py-2">
					{(
						Object.entries(RELATIONSHIP_LABELS) as [RelationshipType, string][]
					).map(([type, label]) => {
						const Icon = RELATIONSHIP_ICONS[type];
						const isSelected = selectedType === type;

						return (
							<button
								key={type}
								type="button"
								onClick={() => setSelectedType(type)}
								className={`flex flex-col items-center gap-1.5 p-3 rounded-lg border-2 transition-all cursor-pointer ${
									isSelected
										? "border-primary bg-primary/10 text-primary"
										: "border-border hover:border-primary/50 hover:bg-muted/50"
								}`}
							>
								<Icon className="h-4 w-4" />
								<span className="text-xs font-medium text-center leading-tight">
									{label}
								</span>
								{isSelected && (
									<Check className="h-3 w-3 text-primary absolute top-1 right-1" />
								)}
							</button>
						);
					})}
				</div>

				{error && (
					<div className="p-2 bg-destructive/10 text-destructive text-sm rounded border border-destructive/20 text-center">
						{error}
					</div>
				)}

				<DialogFooter className="gap-2">
					<Button variant="outline" onClick={handleClose} disabled={isCreating}>
						Cancel
					</Button>
					<Button
						onClick={handleCreate}
						disabled={!selectedType || isCreating}
						className="min-w-[120px]"
					>
						{isCreating ? (
							<>
								<Loader2 className="h-4 w-4 mr-2 animate-spin" />
								Creating...
							</>
						) : (
							<>
								<Link2 className="h-4 w-4 mr-2" />
								Create Link
							</>
						)}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
