import { ArrowRight, Sparkles } from "lucide-react";
import { Logo } from "@/components/common/logo";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { TrustSignals } from "./trust-signals";

/**
 * LandingHero
 * Primary hero section with animated gradients and health passport concept.
 * Expressive animations with staggered entrance effects.
 */
interface LandingHeroProps {
	className?: string;
	actionButton?: React.ReactNode;
}

export function LandingHero({ className, actionButton }: LandingHeroProps) {
	return (
		<div
			className={cn(
				"relative min-h-screen flex flex-col items-center justify-center",
				"px-4 py-24 pt-32 text-center", // Extra top padding for floating nav
				className,
			)}
		>
			{/* Animated background gradients */}
			<div
				className="pointer-events-none absolute inset-0 overflow-hidden"
				aria-hidden="true"
			>
				{/* Primary gradient blob - animates */}
				<div
					className={cn(
						"absolute left-1/2 top-1/2 h-[700px] w-[700px]",
						"-translate-x-1/2 -translate-y-1/2",
						"rounded-full bg-primary/8 blur-3xl",
						"animate-pulse",
					)}
					style={{ animationDuration: "4s" }}
				/>
				{/* Secondary accent blob */}
				<div
					className={cn(
						"absolute right-1/4 top-1/3 h-[500px] w-[500px]",
						"rounded-full bg-accent/6 blur-3xl",
						"animate-pulse",
					)}
					style={{ animationDuration: "6s", animationDelay: "1s" }}
				/>
				{/* Tertiary ambient blob */}
				<div
					className={cn(
						"absolute left-1/4 bottom-1/4 h-[400px] w-[400px]",
						"rounded-full bg-success/5 blur-3xl",
						"animate-pulse",
					)}
					style={{ animationDuration: "5s", animationDelay: "2s" }}
				/>

				{/* Grid pattern overlay */}
				<div
					className={cn(
						"absolute inset-0",
						"bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)]",
						"bg-size-[60px_60px]",
						"mask-[radial-gradient(ellipse_80%_50%_at_50%_50%,black_40%,transparent_100%)]",
					)}
				/>
			</div>

			{/* Hero content */}
			<div
				className={cn(
					"relative z-10 flex max-w-3xl flex-col items-center gap-8",
					"animate-in fade-in slide-in-from-bottom-8 duration-1000",
				)}
			>
				{/* Announcement badge */}
				<div
					className={cn(
						"inline-flex items-center gap-2",
						"rounded-full border border-primary/30 bg-primary/10 px-4 py-1.5",
						"text-sm text-primary",
						"animate-in fade-in zoom-in-95 duration-500",
					)}
				>
					<Sparkles className="h-4 w-4" />
					<span>Open Source Health Protocol</span>
				</div>

				{/* Logo */}
				<Logo size="lg" />

				{/* Main headline */}
				<div className="space-y-6">
					<h1
						className={cn(
							"text-4xl font-bold tracking-tight text-foreground",
							"sm:text-5xl md:text-6xl lg:text-7xl",
							"animate-in fade-in slide-in-from-bottom-4 duration-700",
						)}
						style={{ animationDelay: "200ms" }}
					>
						Prove what matters,{" "}
						<span className="gradient-text">reveal nothing else</span>
					</h1>

					<p
						className={cn(
							"mx-auto max-w-xl text-lg text-muted-foreground",
							"sm:text-xl",
							"animate-in fade-in slide-in-from-bottom-4 duration-700",
						)}
						style={{ animationDelay: "400ms" }}
					>
						Self-sovereign health identity for the longevity community.
						Upload, encrypt, and generate ZK-powered proofs—without revealing your data.
					</p>
				</div>

				{/* CTAs */}
				<div
					className={cn(
						"flex flex-col items-center gap-4 sm:flex-row sm:gap-6",
						"animate-in fade-in slide-in-from-bottom-4 duration-700",
					)}
					style={{ animationDelay: "600ms" }}
				>
					{actionButton}

					<Button
						variant="ghost"
						size="lg"
						asChild
						className="gap-2 text-muted-foreground hover:text-foreground"
					>
						<a href="#how-it-works">
							Learn more
							<ArrowRight className="h-4 w-4" />
						</a>
					</Button>
				</div>

				{/* Subtle helper text */}
				<p
					className={cn(
						"text-xs text-muted-foreground/60",
						"animate-in fade-in duration-1000",
					)}
					style={{ animationDelay: "800ms" }}
				>
					We never have access to your funds. Connect your wallet to authenticate.
				</p>

				{/* Trust signals */}
				<div
					className={cn(
						"mt-8 pt-8 border-t border-border/50",
						"animate-in fade-in slide-in-from-bottom-4 duration-700",
					)}
					style={{ animationDelay: "1000ms" }}
				>
					<TrustSignals variant="horizontal" />
				</div>
			</div>

			{/* Scroll indicator */}
			<div
				className={cn(
					"absolute bottom-8 left-1/2 -translate-x-1/2",
					"flex flex-col items-center gap-2",
					"text-muted-foreground/50",
					"animate-bounce",
				)}
				style={{ animationDuration: "2s" }}
			>
				<span className="text-xs uppercase tracking-widest">Scroll</span>
				<div className="h-8 w-px bg-linear-to-b from-muted-foreground/50 to-transparent" />
			</div>
		</div>
	);
}
