import { useMutation } from "@tanstack/react-query";
import {
	CheckCircle,
	Loader2,
	ShieldAlert,
	ShieldCheck,
	ShieldQuestion,
} from "lucide-react";
import type { JSX } from "react";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { verifyIntegrity } from "../api";

interface AuditStatusBarProps {
	className?: string;
}

export function AuditStatusBar({ className }: AuditStatusBarProps): JSX.Element {
	const [integrityState, setIntegrityState] = useState<{
		status: "idle" | "valid" | "invalid";
		message?: string;
		lastVerified?: Date;
	}>({ status: "idle" });

	const verifyMutation = useMutation({
		mutationFn: verifyIntegrity,
		onSuccess: (data) => {
			setIntegrityState({
				status: data.valid ? "valid" : "invalid",
				message: data.message,
				lastVerified: new Date(),
			});
		},
		onError: (error) => {
			setIntegrityState({
				status: "invalid",
				message: error instanceof Error ? error.message : "Verification failed",
				lastVerified: new Date(),
			});
		},
	});

	useEffect(() => {
		verifyMutation.mutate();
	}, [verifyMutation.mutate]);



	const statusConfig = {
		idle: {
			icon: ShieldQuestion,
			bgClass: "bg-muted/50 border-border",
			iconClass: "text-muted-foreground",
			label: "Not verified",
			description: "Run verification to confirm chain integrity",
		},
		valid: {
			icon: ShieldCheck,
			bgClass: "bg-success/10 border-success/30",
			iconClass: "text-success",
			label: "Chain Valid",
			description: integrityState.message || "All entries verified",
		},
		invalid: {
			icon: ShieldAlert,
			bgClass: "bg-destructive/10 border-destructive/30",
			iconClass: "text-destructive",
			label: "Chain Invalid",
			description: integrityState.message || "Verification failed",
		},
	};

	const config = statusConfig[integrityState.status];
	const Icon = config.icon;

	return (
		<div
			className={cn(
				"flex flex-col sm:flex-row sm:items-center gap-3 p-4 rounded-xl border transition-colors duration-300",
				config.bgClass,
				className,
			)}
		>
			<div className="flex items-center gap-3 flex-1 min-w-0">
				<div
					className={cn(
						"flex h-10 w-10 shrink-0 items-center justify-center rounded-lg",
						integrityState.status === "valid" && "bg-success/20",
						integrityState.status === "invalid" && "bg-destructive/20",
						integrityState.status === "idle" && "bg-muted",
					)}
				>
					<Icon className={cn("h-5 w-5", config.iconClass)} />
				</div>
				<div className="min-w-0 flex-1">
					<div className="flex items-center gap-2">
						<span className="font-semibold text-foreground text-sm">
							{config.label}
						</span>
						{integrityState.status === "valid" && (
							<CheckCircle className="h-4 w-4 text-success" />
						)}
					</div>
					<p className="text-xs text-muted-foreground truncate">
						{config.description}
					</p>
					{integrityState.lastVerified && (
						<p className="text-xs text-muted-foreground/70 mt-0.5">
							Last verified: {integrityState.lastVerified.toLocaleTimeString()}
						</p>
					)}
				</div>
			</div>

			<Button
				onClick={() => verifyMutation.mutate()}
				disabled={verifyMutation.isPending}
				variant={integrityState.status === "valid" ? "outline" : "default"}
				size="sm"
				className="gap-2 shrink-0"
			>
				{verifyMutation.isPending ? (
					<>
						<Loader2 className="h-4 w-4 animate-spin" />
						Verifying...
					</>
				) : (
					<>
						<ShieldCheck className="h-4 w-4" />
						Verify Chain
					</>
				)}
			</Button>
		</div>
	);
}
