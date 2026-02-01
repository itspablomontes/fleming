import { CheckCircle2, Clock, XCircle } from "lucide-react";
import type { JSX } from "react";
import { cn } from "@/lib/utils";

interface ConsentStatsBarProps {
	pendingCount: number;
	activeCount: number;
	expiredCount: number;
	onFilterClick?: (filter: "pending" | "active" | "expired") => void;
	activeFilter?: "pending" | "active" | "expired" | null;
	isLoading?: boolean;
}

export function ConsentStatsBar({
	pendingCount,
	activeCount,
	expiredCount,
	onFilterClick,
	activeFilter,
	isLoading,
}: ConsentStatsBarProps): JSX.Element {
	const stats = [
		{
			id: "pending" as const,
			label: "Pending",
			count: pendingCount,
			icon: Clock,
			colorClass: "text-amber-500 dark:text-amber-400",
            bgClass: "bg-amber-500/10 border-amber-500/20",
            activeClass: "ring-amber-500/50 border-amber-500",
		},
		{
			id: "active" as const,
			label: "Active",
			count: activeCount,
			icon: CheckCircle2,
			colorClass: "text-emerald-500 dark:text-emerald-400",
            bgClass: "bg-emerald-500/10 border-emerald-500/20",
            activeClass: "ring-emerald-500/50 border-emerald-500",
		},
		{
			id: "expired" as const,
			label: "Expired",
			count: expiredCount,
			icon: XCircle,
			colorClass: "text-slate-500 dark:text-slate-400",
            bgClass: "bg-slate-500/10 border-slate-500/20",
            activeClass: "ring-slate-500/50 border-slate-500",
		},
	];

	return (
		<div className="grid grid-cols-3 gap-3 w-full">
			{stats.map((stat) => {
				const Icon = stat.icon;
				const isActive = activeFilter === stat.id;

				return (
					<button
						key={stat.id}
						type="button"
						onClick={() => onFilterClick?.(stat.id)}
                        disabled={isLoading}
						className={cn(
							"relative flex items-center justify-between p-3 rounded-lg border transition-all duration-200 group",
                            "hover:bg-accent/50 hover:border-accent-foreground/20",
                            stat.bgClass,
                            isActive && cn("ring-1 ring-offset-0", stat.activeClass),
                            isLoading && "opacity-50 cursor-wait"
						)}
					>
                        <div className="flex items-center gap-3">
                            <div className={cn("p-1.5 rounded-md bg-background/50 backdrop-blur-sm", stat.colorClass)}>
                                <Icon className="w-4 h-4" />
                            </div>
                            <span className="text-sm font-medium text-foreground/80 group-hover:text-foreground transition-colors">
                                {stat.label}
                            </span>
                        </div>
                        
						<span className={cn("text-lg font-bold tabular-nums tracking-tight", stat.colorClass)}>
							{isLoading ? "-" : stat.count}
						</span>
					</button>
				);
			})}
		</div>
	);
}
