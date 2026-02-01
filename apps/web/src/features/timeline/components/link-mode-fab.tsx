import { Link2, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { useTimelineCoordinator } from "../stores/timeline-coordinator";

export function LinkModeFAB() {
	const isLinkMode = useTimelineCoordinator((state) => state.isLinkMode);
	const linkSource = useTimelineCoordinator((state) => state.linkSource);
	const startLinkMode = useTimelineCoordinator((state) => state.startLinkMode);
	const cancelLinkMode = useTimelineCoordinator((state) => state.cancelLinkMode);

	const handleClick = () => {
		if (isLinkMode) {
			cancelLinkMode();
		} else {
			startLinkMode();
		}
	};

	return (
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger asChild>
					<Button
						onClick={handleClick}
						size="lg"
						className={`fixed bottom-6 left-24 h-14 w-14 rounded-full shadow-lg transition-all hover:scale-110 focus:ring-4 z-50 ${
							isLinkMode
								? "bg-amber-600 hover:bg-amber-700 text-white focus:ring-amber-500/50 animate-pulse"
								: "bg-cyan-600 hover:bg-cyan-700 text-white focus:ring-cyan-500/50"
						}`}
						aria-label={isLinkMode ? "Cancel linking" : "Link events"}
					>
						{isLinkMode ? (
							<X className="h-6 w-6" />
						) : (
							<Link2 className="h-6 w-6" />
						)}
					</Button>
				</TooltipTrigger>
				<TooltipContent side="top" className="bg-background border-border">
					{isLinkMode ? (
						<div className="text-center">
							<p className="font-medium">
								{linkSource ? "Click target event" : "Click source event"}
							</p>
							<p className="text-xs text-muted-foreground">
								Press Esc or click to cancel
							</p>
						</div>
					) : (
						<p>Link events together</p>
					)}
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>
	);
}
