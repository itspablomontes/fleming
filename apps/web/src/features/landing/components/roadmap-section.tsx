import {
	Activity,
	CheckCircle,
	Circle,
	Clock,
	Code2,
	Database,
	Rocket,
	Sparkles,
} from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useState } from "react";

import { cn } from "@/lib/utils";

/**
 * RoadmapSection
 * Interactive roadmap displaying the project's evolution.
 * Uses tabs to switch between phases and animated layout for cards.
 */

// Data structure
type PhaseStatus = "active" | "complete" | "planned";

interface RoadmapPhase {
	id: string;
	title: string;
	subtitle: string;
	status: PhaseStatus;
	icon: React.ElementType;
	description: string;
	items: {
		label: string;
		status: "done" | "progress" | "pending";
	}[];
	color: string;
}

const phases: RoadmapPhase[] = [
	{
		id: "foundation",
		title: "Foundation",
		subtitle: "Protocol Core (Complete)",
		status: "complete",
		icon: Database,
		color: "from-teal-500 to-emerald-500",
		description:
			"The bedrock of Fleming. Self-sovereign medical data vault with true end-to-end encryption and identity management.",
		items: [
			{ label: "Identity & Wallet Auth (SIWE)", status: "done" },
			{ label: "Medical Timeline Protocol", status: "done" },
			{ label: "End-to-End Encryption (E2EE)", status: "done" },
			{ label: "Consent State Machine", status: "done" },
			{ label: "Cryptographic Audit Log", status: "done" },
			{ label: "Doctor/Patient Workflows", status: "done" },
		],
	},
	{
		id: "integration",
		title: "Integration",
		subtitle: "Anchoring (Current)",
		status: "active",
		icon: Activity,
		color: "from-cyan-500 to-blue-500",
		description:
			"Mandatory integration phase. Bridging the gap between private data and public trust with on-chain anchoring.",
		items: [
			{ label: "Merkle Tree Construction", status: "done" },
			{ label: "Base L2 Anchoring Contract", status: "done" },
			{ label: "Batch Verification API", status: "progress" },
			{ label: "Automated Cron Anchoring", status: "progress" },
			{ label: "Data Export (Vault ZIP)", status: "pending" },
		],
	},
	{
		id: "evolution",
		title: "Evolution",
		subtitle: "Health Passport (Next)",
		status: "planned",
		icon: Sparkles,
		color: "from-violet-500 to-fuchsia-500",
		description:
			"Transforming into a ZK-powered Health Passport. Prove what matters, reveal nothing else.",
		items: [
			{ label: "Verifiable Credential (VC) Builder", status: "pending" },
			{ label: "Ranges Proof (Zero-Knowledge)", status: "pending" },
			{ label: "Smart Ingestion (OCR + AI)", status: "pending" },
			{ label: "Provider Cosigning", status: "pending" },
			{ label: "Trust Ecosystem", status: "pending" },
		],
	},
];

export function RoadmapSection() {
	const [activePhaseId, setActivePhaseId] = useState<string>("integration");

	const activePhase = phases.find((p) => p.id === activePhaseId) || phases[0];

	return (
		<section
			id="roadmap"
			className="relative py-24 bg-muted/30 overflow-hidden"
		>
			{/* Background decorations */}
			<div className="absolute top-0 inset-x-0 h-px bg-linear-to-r from-transparent via-border to-transparent opacity-50" />
			<div className="absolute -left-20 top-20 w-72 h-72 bg-primary/5 rounded-full blur-3xl" />
			<div className="absolute -right-20 bottom-20 w-96 h-96 bg-accent/5 rounded-full blur-3xl" />

			<div className="container relative mx-auto px-4 max-w-6xl">
				{/* Header */}
				<div className="text-center mb-16 space-y-4">
					<h2
						className={cn(
							"text-3xl md:text-4xl font-bold tracking-tight inline-block gradient-text",
							"animate-in fade-in zoom-in duration-700",
						)}
					>
						Development Roadmap
					</h2>
					<p
						className={cn(
							"text-muted-foreground text-lg max-w-2xl mx-auto",
							"animate-in fade-in slide-in-from-bottom-4 duration-700 delay-100",
						)}
					>
						Building to learn, learning to build. From a secure data vault to a
						sovereign health passport.
					</p>
				</div>

				{/* Tabs Navigation */}
				<div className="flex flex-wrap justify-center gap-4 mb-12">
					{phases.map((phase) => {
						const isActive = activePhaseId === phase.id;
						const Icon = phase.icon;

						return (
							<button
								type="button"
								key={phase.id}
								onClick={() => setActivePhaseId(phase.id)}
								className={cn(
									"relative group flex items-center gap-3 px-6 py-3 rounded-full border transition-all duration-300",
									isActive
										? "bg-background border-primary/50 shadow-lg scale-105"
										: "bg-background/50 border-transparent hover:border-border hover:bg-background/80",
								)}
							>
								{/* Active Pill Glow */}
								{isActive && (
									<div
										className={cn(
											"absolute inset-0 rounded-full opacity-20 bg-linear-to-r",
											phase.color,
										)}
									/>
								)}

								<div
									className={cn(
										"p-2 rounded-full transition-colors",
										isActive
											? "bg-primary/10 text-primary"
											: "bg-muted text-muted-foreground group-hover:text-foreground",
									)}
								>
									<Icon className="w-5 h-5" />
								</div>

								<div className="text-left">
									<div
										className={cn(
											"font-semibold text-sm",
											isActive ? "text-foreground" : "text-muted-foreground",
										)}
									>
										{phase.title}
									</div>
									<div className="text-xs text-muted-foreground">
										{phase.status === "complete" && "Completed"}
										{phase.status === "active" && "In Progress"}
										{phase.status === "planned" && "Coming Soon"}
									</div>
								</div>

								{isActive && (
									<motion.div
										layoutId="activeTabIndicator"
										className="absolute -bottom-2 left-1/2 -translate-x-1/2 w-1 h-1 bg-primary rounded-full hidden md:block"
									/>
								)}
							</button>
						);
					})}
				</div>

				{/* Content Card */}
				<div className="relative min-h-[400px]">
					<AnimatePresence mode="wait">
						<motion.div
							key={activePhase.id}
							initial={{ opacity: 0, y: 20, scale: 0.98 }}
							animate={{ opacity: 1, y: 0, scale: 1 }}
							exit={{ opacity: 0, y: -20, scale: 0.98 }}
							transition={{ duration: 0.3 }}
							className="grid md:grid-cols-5 gap-8 bg-card/50 backdrop-blur-sm border rounded-3xl p-8 shadow-sm"
						>
							{/* Left: Summary */}
							<div className="md:col-span-2 space-y-6">
								<div className="space-y-2">
									<div
										className={cn(
											"text-xs font-mono font-bold uppercase tracking-wider text-transparent bg-clip-text bg-linear-to-r",
											activePhase.color,
										)}
									>
										{activePhase.subtitle}
									</div>
									<h3 className="text-2xl font-bold">{activePhase.title}</h3>
								</div>
								<p className="text-muted-foreground leading-relaxed">
									{activePhase.description}
								</p>

								{/* Status Badge */}
								<div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-background border text-sm font-medium">
									{activePhase.status === "complete" && (
										<>
											<CheckCircle className="w-4 h-4 text-emerald-500" />
											<span>Alpha Complete</span>
										</>
									)}
									{activePhase.status === "active" && (
										<>
											<Clock className="w-4 h-4 text-blue-500 animate-pulse" />
											<span>Active Development</span>
										</>
									)}
									{activePhase.status === "planned" && (
										<>
											<Rocket className="w-4 h-4 text-purple-500" />
											<span>Future Release</span>
										</>
									)}
								</div>
							</div>

							{/* Right: Checklist */}
							<div className="md:col-span-3 bg-background/50 rounded-2xl p-6 border">
								<h4 className="font-semibold mb-4 flex items-center gap-2">
									<Code2 className="w-4 h-4 text-muted-foreground" />
									Key Deliverables
								</h4>
								<div className="grid gap-3">
									{activePhase.items.map((item) => (
										<div
											key={item.label}
											className="flex items-start gap-3 p-3 rounded-lg hover:bg-background transition-colors"
										>
											<div className="mt-0.5">
												{item.status === "done" && (
													<div className="w-5 h-5 rounded-full bg-emerald-500/10 flex items-center justify-center">
														<CheckCircle className="w-3.5 h-3.5 text-emerald-500" />
													</div>
												)}
												{item.status === "progress" && (
													<div className="w-5 h-5 rounded-full bg-blue-500/10 flex items-center justify-center">
														<div className="w-2 h-2 rounded-full bg-blue-500 animate-ping absolute opacity-75" />
														<Circle className="w-3.5 h-3.5 text-blue-500 fill-blue-500/50" />
													</div>
												)}
												{item.status === "pending" && (
													<div className="w-5 h-5 rounded-full bg-muted flex items-center justify-center">
														<Circle className="w-3.5 h-3.5 text-muted-foreground" />
													</div>
												)}
											</div>
											<div className="flex-1">
												<div
													className={cn(
														"text-sm font-medium",
														item.status === "pending"
															? "text-muted-foreground"
															: "text-foreground",
													)}
												>
													{item.label}
												</div>
											</div>
										</div>
									))}
								</div>
							</div>
						</motion.div>
					</AnimatePresence>
				</div>
			</div>
		</section>
	);
}
