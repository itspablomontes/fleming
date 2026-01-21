import { cn } from "@/lib/utils";

/**
 * Fleming Logo
 * Minimal medical-tech branding mark with DeSci aesthetic.
 */
interface LogoProps {
	className?: string;
	size?: "sm" | "md" | "lg";
	showText?: boolean;
}

const sizes = {
	sm: "h-6 w-6",
	md: "h-8 w-8",
	lg: "h-10 w-10",
};

const textSizes = {
	sm: "text-lg",
	md: "text-xl",
	lg: "text-2xl",
};

export function Logo({ className, size = "md", showText = true }: LogoProps) {
	return (
		<div className={cn("flex items-center gap-2", className)}>
			{/* DNA Helix inspired mark */}
			<div
				className={cn(
					"relative flex items-center justify-center rounded-lg bg-linear-to-br from-primary to-accent",
					sizes[size],
				)}
			>
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					strokeWidth="2"
					strokeLinecap="round"
					strokeLinejoin="round"
					className="h-[60%] w-[60%] text-primary-foreground"
					aria-hidden="true"
				>
					{/* Simplified DNA helix / medical cross hybrid */}
					<path d="M12 2v20" />
					<path d="M2 12h20" />
					<circle cx="12" cy="12" r="3" fill="currentColor" opacity="0.5" />
				</svg>
			</div>

			{showText && (
				<span
					className={cn(
						"font-semibold tracking-tight text-foreground",
						textSizes[size],
					)}
				>
					Fleming
				</span>
			)}
		</div>
	);
}
