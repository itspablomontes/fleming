import { GitBranch, ShieldCheck, Upload } from "lucide-react";
import { cn } from "@/lib/utils";

/**
 * HowItWorks
 * 3-step visual flow explaining the Fleming protocol.
 * Expressive animations with staggered entrance effects.
 */
interface HowItWorksProps {
	className?: string;
}

const steps = [
	{
		icon: Upload,
		number: "01",
		title: "Upload & Encrypt",
		description:
			"Drag-drop your medical records. They're encrypted client-side before leaving your browser. We never see your data.",
		accent: "from-primary/20 to-primary/5",
	},
	{
		icon: GitBranch,
		number: "02",
		title: "Build Your Graph",
		description:
			"Connect events with relationships. Lab results linked to medications, diagnoses to treatments. Your health story, structured.",
		accent: "from-accent/20 to-accent/5",
	},
	{
		icon: ShieldCheck,
		number: "03",
		title: "Prove Claims",
		description:
			"Generate verifiable credentials that prove facts about your health without revealing the underlying data. ZK-powered privacy.",
		accent: "from-success/20 to-success/5",
	},
];

export function HowItWorks({ className }: HowItWorksProps) {
	return (
		<section
			id="how-it-works"
			className={cn(
				"relative py-16 md:py-24",
				"scroll-mt-20",
				className,
			)}
		>
			{/* Background gradient */}
			<div
				className="pointer-events-none absolute inset-0 overflow-hidden"
				aria-hidden="true"
			>
				<div className="absolute left-1/4 top-0 h-[500px] w-[500px] -translate-x-1/2 rounded-full bg-primary/5 blur-3xl" />
				<div className="absolute right-1/4 bottom-0 h-[400px] w-[400px] translate-x-1/2 rounded-full bg-accent/5 blur-3xl" />
			</div>

			<div className="relative z-10 mx-auto max-w-6xl px-4">
				{/* Section header */}
				<div className="mb-10 md:mb-12 text-center">
					<p className="mb-4 text-sm font-medium uppercase tracking-widest text-primary">
						How It Works
					</p>
					<h2 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl md:text-5xl">
						From data chaos to{" "}
						<span className="gradient-text">verifiable truth</span>
					</h2>
					<p className="mx-auto mt-4 max-w-2xl text-lg text-muted-foreground">
						Three simple steps to take control of your health data and prove
						claims without compromising privacy.
					</p>
				</div>

				{/* Steps grid */}
				<div className="grid gap-8 md:grid-cols-3">
					{steps.map((step, index) => (
						<div
							key={step.number}
							className={cn(
								"group relative",
								"rounded-2xl border border-border bg-card/50 p-8",
								"backdrop-blur-sm",
								"transition-all duration-500 ease-out-expo",
								"hover:border-primary/50 hover:shadow-lg hover:shadow-primary/10",
								"hover:-translate-y-2",
							)}
							style={{
								animationDelay: `${index * 150}ms`,
							}}
						>
							{/* Step number badge */}
							<div
								className={cn(
									"absolute -top-4 left-6",
									"flex h-8 items-center justify-center rounded-full",
									"bg-linear-to-r px-4 text-sm font-bold",
									step.accent,
									"border border-border",
								)}
							>
								<span className="text-foreground">{step.number}</span>
							</div>

							{/* Icon with glow */}
							<div className="mb-6 mt-4">
								<div
									className={cn(
										"inline-flex h-14 w-14 items-center justify-center",
										"rounded-xl bg-linear-to-br",
										step.accent,
										"transition-transform duration-300",
										"group-hover:scale-110",
									)}
								>
									<step.icon className="h-7 w-7 text-foreground" />
								</div>
							</div>

							{/* Content */}
							<h3 className="mb-3 text-xl font-semibold text-foreground">
								{step.title}
							</h3>
							<p className="text-muted-foreground leading-relaxed">
								{step.description}
							</p>

							{/* Decorative connector line (except last) */}
							{index < steps.length - 1 && (
								<div
									className={cn(
										"hidden md:block",
										"absolute right-0 top-1/2 w-8",
										"translate-x-full -translate-y-1/2",
									)}
									aria-hidden="true"
								>
									<div className="h-px w-full bg-linear-to-r from-border to-transparent" />
								</div>
							)}
						</div>
					))}
				</div>
			</div>
		</section>
	);
}
