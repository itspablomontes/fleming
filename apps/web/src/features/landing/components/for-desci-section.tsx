import {
	Bot,
	Brain,
	Dna,
	FlaskConical,
	HeartPulse,
	Users,
} from "lucide-react";
import { cn } from "@/lib/utils";

/**
 * ForDeSciSection
 * Showcases how Fleming fits into the DeSci ecosystem.
 * Focus on longevity research DAOs and AI agents without naming specific projects.
 */
interface ForDeSciSectionProps {
	className?: string;
}

const useCases = [
	{
		icon: HeartPulse,
		title: "Longevity Research DAOs",
		description:
			"Enable biomarker range proofs for cohort eligibility. Users prove they meet trial criteria without revealing exact values.",
		features: [
			"Cohort eligibility verification",
			"Consented data feeds for research",
			"RWE (Real-World Evidence) contribution",
		],
		gradient: "from-rose-500/30 to-rose-500/10",
		iconColor: "text-rose-400",
	},
	{
		icon: Dna,
		title: "Decentralized Biotech",
		description:
			"Provide verifiable health data for IP-NFT validation and patient-centric clinical trials. Neutral infrastructure for all.",
		features: [
			"IP-NFT patient data validation",
			"Decentralized biobanking",
			"Cross-DAO data portability",
		],
		gradient: "from-emerald-500/30 to-emerald-500/10",
		iconColor: "text-emerald-400",
	},
	{
		icon: Bot,
		title: "AI Health Agents",
		description:
			"Export consented graph data for AI agents analyzing biomarker correlations and generating research hypotheses.",
		features: [
			"JSON-LD/RDF data export",
			"Consent-scoped data feeds",
			"BioAgent integration ready",
		],
		gradient: "from-violet-500/30 to-violet-500/10",
		iconColor: "text-violet-400",
	},
];

const stats = [
	{ value: "23+", label: "Event Types", icon: FlaskConical },
	{ value: "17+", label: "Relationship Types", icon: Brain },
	{ value: "∞", label: "Community Driven", icon: Users },
];

export function ForDeSciSection({ className }: ForDeSciSectionProps) {
	return (
		<section
			id="for-desci"
			className={cn(
				"relative py-16 md:py-24",
				"scroll-mt-20",
				className,
			)}
		>
			{/* Background */}
			<div
				className="pointer-events-none absolute inset-0 overflow-hidden"
				aria-hidden="true"
			>
				<div className="absolute inset-0 bg-linear-to-b from-transparent via-primary/5 to-transparent" />
			</div>

			<div className="relative z-10 mx-auto max-w-6xl px-4">
				{/* Section header */}
				<div className="mb-10 md:mb-12 text-center">
					<p className="mb-4 text-sm font-medium uppercase tracking-widest text-primary">
						For DeSci
					</p>
					<h2 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl md:text-5xl">
						Built for{" "}
						<span className="gradient-text">decentralized science</span>
					</h2>
					<p className="mx-auto mt-4 max-w-2xl text-lg text-muted-foreground">
						Token-free neutral infrastructure designed to integrate with
						research DAOs, biotech communities, and AI agents.
					</p>
				</div>

				{/* Use cases grid */}
				<div className="grid gap-6 lg:grid-cols-3">
					{useCases.map((useCase, index) => (
						<div
							key={useCase.title}
							className={cn(
								"group relative overflow-hidden",
								"rounded-2xl border border-border",
								"bg-linear-to-br from-card/80 to-card/40",
								"backdrop-blur-sm p-8",
								"transition-all duration-500 ease-out-expo",
								"hover:border-primary/30 hover:shadow-xl hover:shadow-primary/5",
							)}
							style={{
								animationDelay: `${index * 100}ms`,
							}}
						>
							{/* Gradient background on hover */}
							<div
								className={cn(
									"absolute inset-0 opacity-0 transition-opacity duration-500",
									"bg-linear-to-br",
									useCase.gradient,
									"group-hover:opacity-100",
								)}
								aria-hidden="true"
							/>

							{/* Content */}
							<div className="relative z-10">
								{/* Icon */}
								<div
									className={cn(
										"mb-6 inline-flex h-12 w-12 items-center justify-center",
										"rounded-xl bg-background/80 border border-border",
										"transition-transform duration-300",
										"group-hover:scale-110 group-hover:rotate-3",
									)}
								>
									<useCase.icon className={cn("h-6 w-6", useCase.iconColor)} />
								</div>

								<h3 className="mb-3 text-xl font-semibold text-foreground">
									{useCase.title}
								</h3>
								<p className="mb-6 text-muted-foreground leading-relaxed">
									{useCase.description}
								</p>

								{/* Feature list */}
								<ul className="space-y-2">
									{useCase.features.map((feature) => (
										<li
											key={feature}
											className="flex items-center gap-2 text-sm text-muted-foreground"
										>
											<span className="h-1.5 w-1.5 rounded-full bg-primary" />
											{feature}
										</li>
									))}
								</ul>
							</div>
						</div>
					))}
				</div>

				{/* Stats bar */}
				<div
					className={cn(
						"mt-16 rounded-2xl border border-border",
						"bg-card/50 backdrop-blur-sm",
						"p-8",
					)}
				>
					<div className="grid grid-cols-3 gap-8 text-center">
						{stats.map((stat) => (
							<div key={stat.label} className="space-y-2">
								<div className="flex items-center justify-center gap-2">
									<stat.icon className="h-5 w-5 text-primary" />
									<span className="text-3xl font-bold text-foreground">
										{stat.value}
									</span>
								</div>
								<p className="text-sm text-muted-foreground">{stat.label}</p>
							</div>
						))}
					</div>
				</div>
			</div>
		</section>
	);
}
