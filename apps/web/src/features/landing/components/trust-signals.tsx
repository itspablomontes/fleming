import { Link, Lock, ShieldCheck } from "lucide-react";
import { cn } from "@/lib/utils";

/**
 * TrustSignals
 * Feature badges that build trust with first-time users.
 * Displayed below the fold on landing page.
 */
export const TrustSignalsVariant = {
	Horizontal: "horizontal",
	Vertical: "vertical",
} as const;

export type TrustSignalsVariant =
	(typeof TrustSignalsVariant)[keyof typeof TrustSignalsVariant];

interface TrustSignalsProps {
	className?: string;
	variant?: TrustSignalsVariant;
}

const signals = [
	{
		icon: Lock,
		label: "Self-Custodial",
		description: "You control your keys",
	},
	{
		icon: Link,
		label: "Auditable",
		description: "Blockchain-anchored",
	},
	{
		icon: ShieldCheck,
		label: "Encrypted",
		description: "AES-256 protection",
	},
];

export function TrustSignals({
	className,
	variant = "horizontal",
}: TrustSignalsProps) {
	return (
		<div
			className={cn(
				"flex gap-6",
				variant === "vertical"
					? "flex-col"
					: "flex-row flex-wrap justify-center",
				className,
			)}
		>
			{signals.map((signal) => (
				<div
					key={signal.label}
					className="flex items-center gap-3 text-muted-foreground"
				>
					<div className="flex h-10 w-10 items-center justify-center rounded-lg border border-border bg-card">
						<signal.icon className="h-5 w-5 text-primary" aria-hidden="true" />
					</div>
					<div>
						<p className="text-sm font-medium text-foreground">
							{signal.label}
						</p>
						<p className="text-xs text-muted-foreground">
							{signal.description}
						</p>
					</div>
				</div>
			))}
		</div>
	);
}
