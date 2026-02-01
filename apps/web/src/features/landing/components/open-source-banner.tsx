import { BookOpen, Github, Heart, MessageCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

/**
 * OpenSourceBanner
 * Community-focused open source contribution CTA.
 * "Built by the community, for the community" tone.
 */
interface OpenSourceBannerProps {
	className?: string;
}

const contributionLinks = [
	{
		icon: Github,
		label: "Fork on GitHub",
		href: "https://github.com/itspablomontes/fleming",
		description: "Star, fork, and contribute code",
	},
	{
		icon: BookOpen,
		label: "Read the Docs",
		href: "https://github.com/itspablomontes/fleming/tree/main/docs",
		description: "Architecture, roadmap, and guides",
	},
	{
		icon: MessageCircle,
		label: "Join Discord",
		href: "#", // TODO: Add Discord link
		description: "Chat with the community",
	},
];

export function OpenSourceBanner({ className }: OpenSourceBannerProps) {
	return (
		<section
			id="open-source"
			className={cn(
				"relative py-24 md:py-32",
				"scroll-mt-20",
				className,
			)}
		>
			{/* Animated background gradient */}
			<div
				className="pointer-events-none absolute inset-0 overflow-hidden"
				aria-hidden="true"
			>
				<div
					className={cn(
						"absolute inset-0",
						"bg-linear-to-br from-primary/10 via-transparent to-accent/10",
						"animate-pulse",
					)}
					style={{ animationDuration: "4s" }}
				/>
			</div>

			<div className="relative z-10 mx-auto max-w-4xl px-4 text-center">
				{/* Community heart icon */}
				<div
					className={cn(
						"mx-auto mb-8 inline-flex h-20 w-20 items-center justify-center",
						"rounded-full bg-linear-to-br from-primary/20 to-accent/20",
						"border border-primary/20",
						"animate-pulse",
					)}
					style={{ animationDuration: "2s" }}
				>
					<Heart className="h-10 w-10 text-primary" />
				</div>

				{/* Headline */}
				<h2 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl md:text-5xl">
					Built by the community,{" "}
					<span className="gradient-text">for the community</span>
				</h2>

				<p className="mx-auto mt-6 max-w-2xl text-lg text-muted-foreground">
					Fleming is MIT-licensed open source software. We believe health data
					sovereignty should be accessible to everyone. Join us in building the
					future of self-sovereign health identity.
				</p>

				{/* Contribution cards */}
				<div className="mt-12 grid gap-4 sm:grid-cols-3">
					{contributionLinks.map((link, index) => (
						<a
							key={link.label}
							href={link.href}
							target={link.href.startsWith("http") ? "_blank" : undefined}
							rel={link.href.startsWith("http") ? "noopener noreferrer" : undefined}
							className={cn(
								"group relative overflow-hidden",
								"rounded-xl border border-border",
								"bg-card/50 backdrop-blur-sm p-6",
								"transition-all duration-300 ease-out-expo",
								"hover:border-primary/50 hover:bg-card/80",
								"hover:-translate-y-1 hover:shadow-lg hover:shadow-primary/10",
							)}
							style={{
								animationDelay: `${index * 100}ms`,
							}}
						>
							{/* Icon */}
							<div
								className={cn(
									"mx-auto mb-4 flex h-12 w-12 items-center justify-center",
									"rounded-lg bg-primary/10",
									"transition-transform duration-300",
									"group-hover:scale-110",
								)}
							>
								<link.icon className="h-6 w-6 text-primary" />
							</div>

							<h3 className="mb-1 font-semibold text-foreground">
								{link.label}
							</h3>
							<p className="text-sm text-muted-foreground">
								{link.description}
							</p>
						</a>
					))}
				</div>

				{/* MIT License badge */}
				<div className="mt-12 inline-flex items-center gap-2 rounded-full border border-border bg-card/50 px-4 py-2 text-sm text-muted-foreground">
					<span className="h-2 w-2 rounded-full bg-success animate-pulse" />
					MIT Licensed • Free Forever
				</div>

				{/* CTA */}
				<div className="mt-8">
					<Button
						size="lg"
						variant="outline"
						asChild
						className={cn(
							"gap-2 text-base",
							"border-primary/50 hover:bg-primary/10",
						)}
					>
						<a
							href="https://github.com/itspablomontes/fleming"
							target="_blank"
							rel="noopener noreferrer"
						>
							<Github className="h-5 w-5" />
							View on GitHub
						</a>
					</Button>
				</div>
			</div>
		</section>
	);
}
