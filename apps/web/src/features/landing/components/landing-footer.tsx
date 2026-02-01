import { BookOpen, Github, Twitter } from "lucide-react";
import { Logo } from "@/components/common/logo";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

/**
 * LandingFooter
 * Simple footer with logo, navigation links, and social icons.
 */
interface LandingFooterProps {
	className?: string;
}

const footerLinks = {
	product: [
		{ label: "How It Works", href: "#how-it-works" },
		{ label: "For DeSci", href: "#for-desci" },
		{ label: "Features", href: "#features" },
	],
	resources: [
		{
			label: "Documentation",
			href: "https://github.com/itspablomontes/fleming/tree/main/docs",
			external: true,
			comingSoon: true,
		},
		{ label: "Architecture", href: "https://github.com/itspablomontes/fleming/blob/main/docs/ARCHITECTURE.md", external: true },
		{ label: "Roadmap", href: "https://github.com/itspablomontes/fleming/blob/main/docs/ROADMAP.md", external: true },
	],
	community: [
		{ label: "GitHub", href: "https://github.com/itspablomontes/fleming", external: true },
		{ label: "Contributing", href: "https://github.com/itspablomontes/fleming/blob/main/CONTRIBUTING.md", external: true },
		{ label: "Discord", href: "#", external: true, comingSoon: true },
	],
};

const socialLinks = [
	{ icon: Github, href: "https://github.com/itspablomontes/fleming", label: "GitHub" },
	{ icon: Twitter, href: "#", label: "Twitter" },
	{
		icon: BookOpen,
		href: "https://github.com/itspablomontes/fleming/tree/main/docs",
		label: "Docs",
		comingSoon: true,
	},
];

export function LandingFooter({ className }: LandingFooterProps) {
	return (
		<footer
			className={cn(
				"relative border-t border-border",
				"bg-card/30 backdrop-blur-sm",
				"py-12 md:py-16",
				className,
			)}
		>
			<div className="mx-auto max-w-6xl px-4">
				<div className="grid gap-12 md:grid-cols-4">
					{/* Brand column */}
					<div className="md:col-span-1">
						<Logo size="sm" />
						<p className="mt-4 text-sm text-muted-foreground leading-relaxed">
							Self-sovereign health identity for the longevity community.
							ZK-powered proofs for patients, providers, and research.
						</p>

						{/* Social links */}
						<div className="mt-6 flex items-center gap-4">
							{socialLinks.map((link) => {
								const icon = <link.icon className="h-4 w-4" />;

								if (link.comingSoon) {
									return (
										<div
											key={link.label}
											aria-disabled="true"
											className={cn(
												"flex h-9 w-9 items-center justify-center",
												"rounded-lg border border-border bg-card/50",
												"text-muted-foreground opacity-60",
											)}
										>
											{icon}
										</div>
									);
								}

								return (
									<a
										key={link.label}
										href={link.href}
										target="_blank"
										rel="noopener noreferrer"
										aria-label={link.label}
										className={cn(
											"flex h-9 w-9 items-center justify-center",
											"rounded-lg border border-border bg-card/50",
											"text-muted-foreground transition-all duration-200",
											"hover:border-primary/50 hover:text-foreground hover:bg-primary/10",
										)}
									>
										{icon}
									</a>
								);
							})}
						</div>
					</div>

					{/* Links columns */}
					<div className="md:col-span-3 grid grid-cols-2 gap-8 sm:grid-cols-3">
						{/* Product */}
						<div>
							<h4 className="mb-4 text-sm font-semibold text-foreground">
								Product
							</h4>
							<ul className="space-y-3">
								{footerLinks.product.map((link) => (
									<li key={link.label}>
										<a
											href={link.href}
											className="text-sm text-muted-foreground hover:text-foreground transition-colors"
										>
											{link.label}
										</a>
									</li>
								))}
							</ul>
						</div>

						{/* Resources */}
						<div>
							<h4 className="mb-4 text-sm font-semibold text-foreground">
								Resources
							</h4>
							<ul className="space-y-3">
								{footerLinks.resources.map((link) => (
									<li key={link.label}>
										{link.comingSoon ? (
											<div className="flex items-center gap-2 text-sm text-muted-foreground">
												<span>{link.label}</span>
												<Badge
													variant="secondary"
													className="bg-muted text-muted-foreground"
												>
													Coming soon
												</Badge>
											</div>
										) : (
											<a
												href={link.href}
												target={link.external ? "_blank" : undefined}
												rel={link.external ? "noopener noreferrer" : undefined}
												className="text-sm text-muted-foreground hover:text-foreground transition-colors"
											>
												{link.label}
												{link.external && (
													<span className="ml-1 text-xs">↗</span>
												)}
											</a>
										)}
									</li>
								))}
							</ul>
						</div>

						{/* Community */}
						<div>
							<h4 className="mb-4 text-sm font-semibold text-foreground">
								Community
							</h4>
							<ul className="space-y-3">
								{footerLinks.community.map((link) => (
									<li key={link.label}>
										{link.comingSoon ? (
											<div className="flex items-center gap-2 text-sm text-muted-foreground">
												<span>{link.label}</span>
												<Badge
													variant="secondary"
													className="bg-muted text-muted-foreground"
												>
													Coming soon
												</Badge>
											</div>
										) : (
											<a
												href={link.href}
												target={link.external ? "_blank" : undefined}
												rel={link.external ? "noopener noreferrer" : undefined}
												className="text-sm text-muted-foreground hover:text-foreground transition-colors"
											>
												{link.label}
												{link.external && (
													<span className="ml-1 text-xs">↗</span>
												)}
											</a>
										)}
									</li>
								))}
							</ul>
						</div>
					</div>
				</div>

				{/* Bottom bar */}
				<div className="mt-12 flex flex-col items-center justify-between gap-4 border-t border-border pt-8 sm:flex-row">
					<p className="text-sm text-muted-foreground">
						© {new Date().getFullYear()} Fleming Protocol. MIT Licensed.
					</p>
					<p className="text-sm text-muted-foreground">
						Built with 💚 for the DeSci community
					</p>
				</div>
			</div>
		</footer>
	);
}
