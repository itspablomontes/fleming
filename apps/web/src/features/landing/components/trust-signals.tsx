import { Database, Link, Lock, ShieldCheck } from "lucide-react";
import { cn } from "@/lib/utils";

/**
 * TrustSignals
 * Feature badges that build trust with first-time users.
 * Expanded version with more detail and polish.
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
		description: "Your keys, your data",
		detail: "Wallet-derived encryption",
	},
	{
		icon: Link,
		label: "On-Chain Anchored",
		description: "Base L2 Merkle roots",
		detail: "Tamper-evident proofs",
	},
	{
		icon: ShieldCheck,
		label: "E2E Encrypted",
		description: "AES-256-GCM",
		detail: "Zero backend access",
	},
	{
		icon: Database,
		label: "Auditable",
		description: "Hash-chained logs",
		detail: "Verifiable integrity",
	},
];

export function TrustSignals({
	className,
	variant = "horizontal",
}: TrustSignalsProps) {
	return (
		<div
			className={cn(
				"flex gap-6 md:gap-8",
				variant === "vertical"
					? "flex-col"
					: "flex-row flex-wrap justify-center",
				className,
			)}
		>
			{signals.map((signal, index) => (
				<div
					key={signal.label}
					className={cn(
						"group flex items-center gap-3",
						"transition-all duration-300",
						"hover:scale-105",
					)}
					style={{
						animationDelay: `${index * 100}ms`,
					}}
				>
					{/* Icon container with glow on hover */}
					<div
						className={cn(
							"relative flex h-12 w-12 items-center justify-center",
							"rounded-xl border border-border bg-card/50",
							"transition-all duration-300",
							"group-hover:border-primary/50 group-hover:shadow-lg group-hover:shadow-primary/20",
						)}
					>
						<signal.icon
							className="h-5 w-5 text-primary transition-transform duration-300 group-hover:scale-110"
							aria-hidden="true"
						/>
					</div>

					{/* Text content */}
					<div className="text-left">
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

/**
 * TrustSignalsExpanded
 * Full-width bento grid version for below-the-fold placement.
 */
interface TrustSignalsExpandedProps {
	className?: string;
}

export function TrustSignalsExpanded({ className }: TrustSignalsExpandedProps) {
	return (
		<section
			id="features"
			className={cn("relative py-16 md:py-24 scroll-mt-20", className)}
		>
			<div className="mx-auto max-w-6xl px-4">
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
					{signals.map((signal, index) => (
						<div
							key={signal.label}
							className={cn(
								"group relative overflow-hidden",
								"rounded-2xl border border-border",
								"bg-card/50 backdrop-blur-sm p-6",
								"transition-all duration-500 ease-out-expo",
								"hover:border-primary/30 hover:bg-card/80",
								"hover:-translate-y-1 hover:shadow-xl hover:shadow-primary/10",
							)}
							style={{
								animationDelay: `${index * 100}ms`,
							}}
						>
							{/* Gradient overlay on hover */}
							<div
								className={cn(
									"absolute inset-0 opacity-0 transition-opacity duration-500",
									"bg-linear-to-br from-primary/10 to-transparent",
									"group-hover:opacity-100",
								)}
								aria-hidden="true"
							/>

							{/* Content */}
							<div className="relative z-10">
								<div
									className={cn(
										"mb-4 inline-flex h-12 w-12 items-center justify-center",
										"rounded-xl bg-primary/10",
										"transition-transform duration-300",
										"group-hover:scale-110",
									)}
								>
									<signal.icon className="h-6 w-6 text-primary" />
								</div>

								<h3 className="mb-1 text-lg font-semibold text-foreground">
									{signal.label}
								</h3>
								<p className="mb-2 text-sm text-muted-foreground">
									{signal.description}
								</p>
								<p className="text-xs text-primary/80">{signal.detail}</p>
							</div>
						</div>
					))}
				</div>
			</div>
		</section>
	);
}
