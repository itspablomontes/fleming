import {
	AlertTriangle,
	ClipboardList,
	Eye,
	File,
	FileCheck,
	FileText,
	HeartPulse,
	Lock,
	Pill,
	ScanLine,
	Stethoscope,
	Syringe,
	TestTube2,
	UserPlus,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import {
	EVENT_TYPE_LABELS,
	type TimelineEvent,
	type TimelineEventType,
} from "../types";

const iconComponents: Record<TimelineEventType, React.ElementType> = {
	consultation: Stethoscope,
	diagnosis: ClipboardList,
	prescription: Pill,
	procedure: Stethoscope,
	lab_result: TestTube2,
	imaging: ScanLine,
	note: FileText,
	visit_note: FileText,
	vaccination: Syringe,
	allergy: AlertTriangle,
	vital_signs: HeartPulse,
	referral: UserPlus,
	insurance_claim: FileCheck,
	other: File,
	tombstone: File,
};

interface TimelineItemProps {
	event: TimelineEvent;
	onView?: (event: TimelineEvent) => void;
	className?: string;
}

/**
 * TimelineItem
 * Displays a single medical event in the patient's timeline.
 * Shows event type icon, title, provider, date, and encryption status.
 */
export function TimelineItem({ event, onView, className }: TimelineItemProps) {
	const Icon = iconComponents[event.type];
	const typeLabel = EVENT_TYPE_LABELS[event.type];

	const timestamp = new Date(event.timestamp);
	const formattedDate = new Intl.DateTimeFormat("en-US", {
		month: "short",
		day: "numeric",
		year: "numeric",
	}).format(timestamp);

	return (
		<Card
			className={cn(
				"cursor-pointer transition-colors hover:bg-accent/50",
				className,
			)}
			onClick={() => onView?.(event)}
		>
			<CardHeader className="flex flex-row items-start gap-4 p-4 pb-2">
				<div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
					<Icon className="h-5 w-5" aria-hidden="true" />
				</div>
				<div className="flex-1 space-y-1">
					<div className="flex items-center justify-between">
						<h3 className="font-semibold leading-none">{event.title}</h3>
						<time
							dateTime={timestamp.toISOString()}
							className="text-xs text-muted-foreground"
						>
							{formattedDate}
						</time>
					</div>
					<Badge variant="secondary" className="text-xs dark">
						{typeLabel}
					</Badge>
				</div>
			</CardHeader>
			<CardContent className="p-4 pt-0">
				<div className="flex items-center justify-between">
					<div className="space-y-1">
						{event.description && (
							<p className="text-sm text-muted-foreground line-clamp-2">
								{event.description}
							</p>
						)}
						{event.provider && (
							<p className="text-xs text-muted-foreground">{event.provider}</p>
						)}
					</div>
					<div className="flex items-center gap-2">
						{event.isEncrypted && (
							<Badge
								variant="outline"
								className="gap-1 text-xs text-success border-success/30 dark"
							>
								<Lock className="h-3 w-3" />
								Encrypted
							</Badge>
						)}
						<Button
							variant="ghost"
							size="sm"
							onClick={(e) => {
								e.stopPropagation();
								onView?.(event);
							}}
						>
							<Eye className="mr-1 h-4 w-4" />
							View
						</Button>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
