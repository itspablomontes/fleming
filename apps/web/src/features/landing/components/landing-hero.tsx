import { Logo } from "@/components/common/logo";
import { cn } from "@/lib/utils";
import { TrustSignals } from "./trust-signals";

/**
 * LandingHero
 * Primary hero section for unauthenticated users.
 * DeSci-inspired dark aesthetic with cyan/lime accents.
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
				"px-4 py-12 text-center",
				className,
			)}
		>
			<div
				className="pointer-events-none absolute inset-0 overflow-hidden"
				aria-hidden="true"
			>
				<div className="absolute left-1/2 top-1/2 h-[600px] w-[600px] -translate-x-1/2 -translate-y-1/2 rounded-full bg-primary/5 blur-3xl" />
				<div className="absolute right-1/4 top-1/4 h-[400px] w-[400px] rounded-full bg-accent/5 blur-3xl" />
			</div>

			<div className="relative z-10 flex max-w-2xl flex-col items-center gap-8">
				<Logo size="lg" />

				<div className="space-y-4">
					<h1 className="text-4xl font-bold tracking-tight text-foreground sm:text-5xl md:text-6xl">
						Your medical records,{" "}
						<span className="bg-linear-to-r from-primary to-accent bg-clip-text text-transparent">
							your control.
						</span>
					</h1>

					<p className="mx-auto max-w-md text-lg text-muted-foreground">
						Store, encrypt, and manage who can access your health data.
						Self-sovereign. Auditable. Secure.
					</p>
				</div>

				<div className="flex flex-col items-center gap-4">
					{actionButton}

					<p className="text-xs text-muted-foreground">
						We never have access to your funds. Just your identity.
					</p>
				</div>

				<div className="mt-8 pt-8 border-t border-border">
					<TrustSignals />
				</div>
			</div>
		</div>
	);
}
