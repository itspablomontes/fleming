import { format } from "date-fns";
import {
	AlertTriangle,
	ClipboardList,
	Download,
	Eye,
	File,
	FileCheck,
	FileText,
	HeartPulse,
	Lock,
	type LucideIcon,
	Pill,
	ScanLine,
	Stethoscope,
	Syringe,
	TestTube2,
	UserPlus,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import type { TimelineEvent, TimelineEventType } from "../types";
import { EVENT_TYPE_LABELS, TimelineEventType as EventTypes } from "../types";

// Map event types to their corresponding icons
const EVENT_ICONS: Record<TimelineEventType, LucideIcon> = {
	[EventTypes.LAB_RESULT]: TestTube2,
	[EventTypes.PRESCRIPTION]: Pill,
	[EventTypes.IMAGING]: ScanLine,
	[EventTypes.VISIT_NOTE]: FileText,
	[EventTypes.VACCINATION]: Syringe,
	[EventTypes.PROCEDURE]: Stethoscope,
	[EventTypes.ALLERGY]: AlertTriangle,
	[EventTypes.VITAL_SIGNS]: HeartPulse,
	[EventTypes.REFERRAL]: UserPlus,
	[EventTypes.INSURANCE_CLAIM]: FileCheck,
	[EventTypes.OTHER]: File,
	[EventTypes.NOTE]: FileText,
	[EventTypes.CONSULTATION]: Stethoscope,
	[EventTypes.DIAGNOSIS]: ClipboardList,
	[EventTypes.TOMBSTONE]: File,
};

interface TimelineCardProps {
	event: TimelineEvent;
	onClick?: () => void;
}

export function TimelineCard({ event, onClick }: TimelineCardProps) {
	const Icon = EVENT_ICONS[event.type] || File;

	return (
		<Card
			className="p-4 hover:shadow-md transition-all cursor-pointer border-cyan-200 dark:border-cyan-800 hover:border-cyan-400 dark:hover:border-cyan-600"
			onClick={onClick}
		>
			<div className="flex items-start gap-4">
				{/* Icon */}
				<div className="shrink-0 w-12 h-12 rounded-full bg-cyan-100 dark:bg-cyan-900/30 flex items-center justify-center">
					<Icon className="w-6 h-6 text-cyan-600 dark:text-cyan-400" />
				</div>

				{/* Content */}
				<div className="flex-1 min-w-0">
					<div className="flex items-start justify-between gap-2 mb-1">
						<div className="flex items-center gap-2">
							<h3 className="font-semibold text-cyan-900 dark:text-cyan-50">
								{event.title}
							</h3>
							{event.isEncrypted && (
								<Lock
									className="w-4 h-4 text-emerald-600 dark:text-emerald-400"
									aria-label="Encrypted"
								/>
							)}
						</div>
						<time className="text-sm text-cyan-700 dark:text-cyan-300 whitespace-nowrap">
							{format(new Date(event.timestamp), "MMM dd, yyyy")}
						</time>
					</div>

					<p className="text-sm text-cyan-600 dark:text-cyan-400 mb-2">
						{EVENT_TYPE_LABELS[event.type]}
					</p>

					{event.provider && (
						<p className="text-sm text-cyan-700 dark:text-cyan-300">
							{event.provider}
						</p>
					)}

					{/* Actions */}
					{event.blobRef && (
						<div className="flex gap-2 mt-3">
							<Button
								variant="outline"
								size="sm"
								className="text-cyan-700 border-cyan-300 hover:bg-cyan-50 dark:text-cyan-300 dark:border-cyan-700 dark:hover:bg-cyan-900/20"
								onClick={(e) => {
									e.stopPropagation();
									if (event.blobRef) {
										window.open(event.blobRef, "_blank", "noopener,noreferrer");
									}
								}}
							>
								<Eye className="w-4 h-4 mr-1" />
								View
							</Button>
							<Button
								variant="outline"
								size="sm"
								className="text-cyan-700 border-cyan-300 hover:bg-cyan-50 dark:text-cyan-300 dark:border-cyan-700 dark:hover:bg-cyan-900/20"
								onClick={(e) => {
									e.stopPropagation();
									if (event.blobRef) {
										const link = document.createElement("a");
										link.href = event.blobRef;
										link.download = `fleming-doc-${event.id}`;
										document.body.appendChild(link);
										link.click();
										document.body.removeChild(link);
									}
								}}
							>
								<Download className="w-4 h-4 mr-1" />
								Download
							</Button>
						</div>
					)}
				</div>
			</div>
		</Card>
	);
}
