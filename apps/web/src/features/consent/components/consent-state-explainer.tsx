import { ChevronDown, ChevronUp, HelpCircle, Info } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { CONSENT_STATE_LABELS, CONSENT_STATE_VARIANTS } from "../types";

/**
 * ConsentStateExplainer
 * Collapsible section that explains the consent lifecycle with a visual state machine.
 * User preferred option B: Collapsible "Understanding consent states" section.
 */
interface ConsentStateExplainerProps {
	className?: string;
	defaultOpen?: boolean;
}

// State metadata for the visual diagram
const stateNodes: Record<
	string,
	{
		label: string;
		variant: "warning" | "success" | "destructive" | "secondary";
		icon: string;
		description: string;
		isFinal: boolean;
	}
> = {
	requested: {
		label: CONSENT_STATE_LABELS.requested,
		variant: CONSENT_STATE_VARIANTS.requested,
		icon: "⏳",
		description: "Awaiting your decision",
		isFinal: false,
	},
	approved: {
		label: CONSENT_STATE_LABELS.approved,
		variant: CONSENT_STATE_VARIANTS.approved,
		icon: "✅",
		description: "Grantee can access your data",
		isFinal: false,
	},
	denied: {
		label: CONSENT_STATE_LABELS.denied,
		variant: CONSENT_STATE_VARIANTS.denied,
		icon: "❌",
		description: "Request was rejected",
		isFinal: true,
	},
	revoked: {
		label: CONSENT_STATE_LABELS.revoked,
		variant: CONSENT_STATE_VARIANTS.revoked,
		icon: "🚫",
		description: "Access was withdrawn",
		isFinal: true,
	},
	expired: {
		label: CONSENT_STATE_LABELS.expired,
		variant: CONSENT_STATE_VARIANTS.expired,
		icon: "⌛",
		description: "Time-limited access ended",
		isFinal: true,
	},
	suspended: {
		label: CONSENT_STATE_LABELS.suspended,
		variant: CONSENT_STATE_VARIANTS.suspended,
		icon: "⏸️",
		description: "Temporarily paused",
		isFinal: false,
	},
};

const variantStyles: Record<string, string> = {
	warning: "bg-warning/20 border-warning/50 text-warning-foreground",
	success: "bg-success/20 border-success/50 text-success",
	destructive: "bg-destructive/20 border-destructive/50 text-destructive",
	secondary: "bg-muted/50 border-border text-muted-foreground",
};

export function ConsentStateExplainer({
	className,
	defaultOpen = false,
}: ConsentStateExplainerProps) {
	const [isOpen, setIsOpen] = useState(defaultOpen);

	return (
		<div
			className={cn(
				"rounded-xl border border-border bg-card/50",
				className,
			)}
		>
			{/* Trigger button */}
			<Button
				variant="ghost"
				onClick={() => setIsOpen(!isOpen)}
				className={cn(
					"w-full justify-between gap-2 p-4 h-auto",
					"text-muted-foreground hover:text-foreground",
					"transition-colors duration-200",
				)}
			>
				<div className="flex items-center gap-2">
					<HelpCircle className="h-4 w-4" />
					<span className="text-sm font-medium">
						Understanding consent states
					</span>
				</div>
				{isOpen ? (
					<ChevronUp className="h-4 w-4" />
				) : (
					<ChevronDown className="h-4 w-4" />
				)}
			</Button>

			{/* Collapsible content */}
			{isOpen && (
				<div className="px-4 pb-4 space-y-6 animate-in fade-in slide-in-from-top-2 duration-200">
					{/* Visual state diagram */}
					<div className="relative">
						{/* State flow diagram - horizontal on desktop, vertical on mobile */}
						<div className="flex flex-col md:flex-row items-center md:items-start gap-4 md:gap-2 justify-center">
							{/* Requested state */}
							<StateNode state="requested" />

							{/* Arrow to approved/denied */}
							<div className="flex flex-col md:flex-row items-center gap-2">
								<Arrow direction="horizontal" label="Approve" className="hidden md:flex" />
								<Arrow direction="vertical" label="Approve" className="md:hidden" />
							</div>

							{/* Approved state with branches */}
							<div className="flex flex-col items-center gap-2">
								<StateNode state="approved" />
								
								{/* Branches from approved */}
								<div className="flex gap-4 mt-2">
									<div className="flex flex-col items-center gap-1">
										<span className="text-xs text-muted-foreground">Revoke</span>
										<StateNode state="revoked" size="sm" />
									</div>
									<div className="flex flex-col items-center gap-1">
										<span className="text-xs text-muted-foreground">Expire</span>
										<StateNode state="expired" size="sm" />
									</div>
									<div className="flex flex-col items-center gap-1">
										<span className="text-xs text-muted-foreground">Suspend</span>
										<StateNode state="suspended" size="sm" />
									</div>
								</div>
							</div>
						</div>

						{/* Denied branch (separate) */}
						<div className="flex items-center gap-2 mt-4 justify-center opacity-60">
							<span className="text-xs text-muted-foreground">or Deny →</span>
							<StateNode state="denied" size="sm" />
						</div>
					</div>

					{/* Legend */}
					<div className="border-t border-border pt-4">
						<p className="text-xs text-muted-foreground mb-3">
							<Info className="inline h-3 w-3 mr-1" />
							Solid borders indicate active states. Faded states are terminal (no further actions).
						</p>
						
						<div className="grid grid-cols-2 md:grid-cols-3 gap-2">
							{Object.entries(stateNodes).map(([key, node]) => (
								<div
									key={key}
									className={cn(
										"flex items-center gap-2 rounded-lg p-2 text-xs",
										node.isFinal && "opacity-60",
									)}
								>
									<span>{node.icon}</span>
									<span className="font-medium">{node.label}</span>
								</div>
							))}
						</div>
					</div>
				</div>
			)}
		</div>
	);
}

/**
 * StateNode - Individual state in the diagram
 */
interface StateNodeProps {
	state: keyof typeof stateNodes;
	size?: "sm" | "md";
	className?: string;
}

function StateNode({ state, size = "md", className }: StateNodeProps) {
	const node = stateNodes[state];
	
	return (
		<Tooltip>
			<TooltipTrigger asChild>
				<div
					className={cn(
						"flex items-center gap-2 rounded-lg border-2",
						"transition-all duration-200",
						"hover:scale-105 cursor-help",
						variantStyles[node.variant],
						size === "sm" ? "px-2 py-1 text-xs" : "px-3 py-2 text-sm",
						node.isFinal && "opacity-70 border-dashed",
						className,
					)}
				>
					<span>{node.icon}</span>
					<span className="font-medium">{node.label}</span>
				</div>
			</TooltipTrigger>
			<TooltipContent side="top" className="max-w-xs">
				<p>{node.description}</p>
			</TooltipContent>
		</Tooltip>
	);
}

/**
 * Arrow - Transition arrow between states
 */
interface ArrowProps {
	direction: "horizontal" | "vertical";
	label?: string;
	className?: string;
}

function Arrow({ direction, label, className }: ArrowProps) {
	return (
		<div
			className={cn(
				"flex items-center gap-1",
				direction === "vertical" && "flex-col",
				className,
			)}
		>
			{label && (
				<span className="text-xs text-muted-foreground">{label}</span>
			)}
			<div
				className={cn(
					"bg-border",
					direction === "horizontal" ? "h-px w-8" : "w-px h-8",
				)}
			/>
			<div
				className={cn(
					"border-border",
					direction === "horizontal"
						? "border-r-2 border-t-2 w-2 h-2 rotate-45 -ml-1"
						: "border-b-2 border-r-2 w-2 h-2 rotate-45 -mt-1",
				)}
			/>
		</div>
	);
}
