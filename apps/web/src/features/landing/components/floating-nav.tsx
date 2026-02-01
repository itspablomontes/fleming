import { Github } from "lucide-react";
import { Logo } from "@/components/common/logo";
import { ThemeToggle } from "@/components/common/theme-toggle";
import { Button } from "@/components/ui/button";
import { AuthButton } from "@/features/auth/components/auth-button";
import { cn } from "@/lib/utils";

/**
 * FloatingNav
 * Glassmorphism navigation bar with logo, nav links, and wallet connect.
 * Floats above content with subtle backdrop blur.
 */
interface FloatingNavProps {
	className?: string;
}

const navLinks = [
	{ href: "#how-it-works", label: "How It Works" },
	{ href: "#roadmap", label: "Roadmap" },
	{ href: "#for-desci", label: "For DeSci" },
	{ href: "#open-source", label: "Open Source" },
];

export function FloatingNav({ className }: FloatingNavProps) {
	return (
		<header
			className={cn(
				"fixed top-4 left-4 right-4 z-50",
				"mx-auto max-w-6xl",
				"glass rounded-2xl px-4 py-3",
				"flex items-center justify-between gap-4",
				"animate-in fade-in slide-in-from-top-4 duration-700",
				className,
			)}
		>
			<Logo size="sm" />

			<nav className="hidden md:flex items-center gap-6">
				{navLinks.map((link) => (
					<a
						key={link.href}
						href={link.href}
						className={cn(
							"text-sm font-medium text-muted-foreground",
							"hover:text-foreground transition-colors duration-200",
							"relative after:absolute after:bottom-0 after:left-0 after:h-0.5 after:w-0",
							"after:bg-primary after:transition-all after:duration-300",
							"hover:after:w-full",
						)}
					>
						{link.label}
					</a>
				))}
			</nav>

			<div className="flex items-center gap-3">
				<Button
					variant="ghost"
					size="sm"
					asChild
					className="hidden sm:flex gap-2 text-muted-foreground hover:text-foreground"
				>
					<a
						href="https://github.com/itspablomontes/fleming"
						target="_blank"
						rel="noopener noreferrer"
					>
						<Github className="h-4 w-4" />
						<span className="hidden lg:inline">GitHub</span>
					</a>
				</Button>
				<ThemeToggle />
				<AuthButton />
			</div>
		</header>
	);
}
