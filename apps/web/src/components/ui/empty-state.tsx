import type { LucideIcon } from "lucide-react";
import type { JSX, ReactNode } from "react";
import { cn } from "@/lib/utils";

interface EmptyStateProps {
	icon?: LucideIcon;
	title: string;
	description?: string;
	action?: ReactNode;
	className?: string;
}

export function EmptyState({
	icon: Icon,
	title,
	description,
	action,
	className,
}: EmptyStateProps): JSX.Element {
	return (
		<div
			className={cn(
				"flex min-h-[300px] flex-col items-center justify-center gap-4 rounded-xl border border-dashed border-border bg-muted/30 p-8 text-center animate-in fade-in zoom-in-95 duration-200",
				className,
			)}
		>
			{Icon && (
				<div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted shadow-sm">
					<Icon className="h-6 w-6 text-muted-foreground" />
				</div>
			)}
			<div className="max-w-xs space-y-1">
				<h3 className="text-sm font-semibold text-foreground">{title}</h3>
				{description && (
					<p className="text-sm text-muted-foreground">{description}</p>
				)}
			</div>
			{action && <div className="mt-2">{action}</div>}
		</div>
	);
}
