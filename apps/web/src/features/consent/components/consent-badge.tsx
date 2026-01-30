import { cva, type VariantProps } from "class-variance-authority";
import { Ban, CheckCircle, Clock, PauseCircle, Timer, XCircle } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import type { ConsentState } from "@/features/consent/types";
import { CONSENT_STATE_LABELS } from "@/features/consent/types";
import { cn } from "@/lib/utils";

const consentBadgeVariants = cva("inline-flex items-center gap-1.5", {
	variants: {
		state: {
			requested: "bg-warning/15 text-(--warning-text) border-warning/30",
			approved: "bg-success/15 text-(--success-text) border-success/30",
			denied: "bg-muted text-muted-foreground border-border",
			revoked: "bg-destructive/15 text-destructive border-destructive/30",
			expired: "bg-muted text-muted-foreground border-border",
			suspended: "bg-warning/15 text-(--warning-text) border-warning/30",
		},
	},
	defaultVariants: {
		state: "requested",
	},
});

const stateIcons: Record<ConsentState, React.ElementType> = {
	requested: Clock,
	approved: CheckCircle,
	denied: XCircle,
	revoked: Ban,
	expired: Timer,
	suspended: PauseCircle,
};

interface ConsentBadgeProps extends VariantProps<typeof consentBadgeVariants> {
	state: ConsentState;
	showLabel?: boolean;
	className?: string;
}

/**
 * ConsentBadge
 * Displays the current state of a consent with appropriate color and icon.
 */
export function ConsentBadge({
	state,
	showLabel = true,
	className,
}: ConsentBadgeProps) {
	const Icon = stateIcons[state];
	const label = CONSENT_STATE_LABELS[state];

	return (
		<Badge
			variant="outline"
			className={cn(consentBadgeVariants({ state }), "dark", className)}
		>
			<Icon className="h-3.5 w-3.5" aria-hidden="true" />
			{showLabel && <span>{label}</span>}
			<span className="sr-only">{label}</span>
		</Badge>
	);
}
