import { Link, useLocation } from "@tanstack/react-router";
import { Clock, Menu, NotebookText, ShieldCheck } from "lucide-react";
import type { JSX } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet";
import { useMyConsentGrants } from "@/features/consent/hooks/use-consent-grants";
import { cn } from "@/lib/utils";

const navItems = [
	{
		label: "Timeline",
		description: "Your medical history overview",
		to: "/",
		icon: Clock,
	},
	{
		label: "Consent",
		description: "Manage access requests",
		to: "/consent",
		icon: ShieldCheck,
	},
	{
		label: "Audit Log",
		description: "Security & activity trail",
		to: "/audit",
		icon: NotebookText,
	},
] as const;

interface AppNavigationProps {
	triggerClassName?: string;
}

export function AppNavigation({
	triggerClassName,
}: AppNavigationProps): JSX.Element {
	const location = useLocation();
	const { data: grants } = useMyConsentGrants();
	const pendingCount =
		grants?.filter((grant) => grant.state === "requested").length ?? 0;

	return (
		<Sheet>
			<SheetTrigger asChild>
				<Button
					size="icon"
					className={cn("rounded-full", triggerClassName)}
					aria-label="Open navigation menu"
				>
					<Menu className="h-5 w-5" />
				</Button>
			</SheetTrigger>
			<SheetContent side="right" className="w-80">
				<SheetHeader className="space-y-2">
					<SheetTitle>Navigate</SheetTitle>
					<SheetDescription>
						Jump to core areas of your Fleming timeline.
					</SheetDescription>
				</SheetHeader>
				<div className="space-y-2 px-4 pb-6">
					{navItems.map((item) => {
						const isActive = location.pathname === item.to;
						const Icon = item.icon;

						return (
							<Button
								key={item.to}
								variant={isActive ? "secondary" : "ghost"}
								className={cn(
									"flex h-auto w-full items-start justify-between gap-3 rounded-lg px-4 py-3 text-left",
									isActive && "border border-border",
								)}
								asChild
							>
								<Link to={item.to}>
									<span className="flex items-start gap-3">
										<Icon className="mt-0.5 h-4 w-4" aria-hidden="true" />
										<span className="space-y-1">
											<span className="block text-sm font-medium text-foreground">
												{item.label}
											</span>
											<span className="block text-xs text-muted-foreground">
												{item.description}
											</span>
										</span>
									</span>
									{item.to === "/consent" && pendingCount > 0 && (
										<Badge variant="secondary" className="ml-auto">
											{pendingCount}
										</Badge>
									)}
								</Link>
							</Button>
						);
					})}
				</div>
			</SheetContent>
		</Sheet>
	);
}
