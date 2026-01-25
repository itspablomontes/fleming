import { format } from "date-fns";
import type { JSX } from "react";
import { AddressDisplay } from "@/components/common/address-display";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { ConsentBadge } from "@/features/consent/components/consent-badge";
import { cn } from "@/lib/utils";

import type { ConsentGrant, ConsentPermission } from "../types";

const permissionLabels: Record<ConsentPermission, string> = {
	read: "Read",
	write: "Write",
	share: "Share",
};

interface ConsentRequestCardProps {
	grant: ConsentGrant;
	onApprove: (grantId: string) => void;
	onDeny: (grantId: string) => void;
	isApproving?: boolean;
	isDenying?: boolean;
	className?: string;
}

const formatDate = (value?: Date): string => {
	if (!value) {
		return "No expiration";
	}
	return format(value, "MMM d, yyyy");
};

export function ConsentRequestCard({
	grant,
	onApprove,
	onDeny,
	isApproving,
	isDenying,
	className,
}: ConsentRequestCardProps): JSX.Element {
	const isBusy = Boolean(isApproving || isDenying);

	return (
		<Card className={cn("border-border bg-white dark:bg-gray-900", className)}>
			<CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
				<div className="space-y-2">
					<CardTitle className="text-base font-semibold text-foreground">
						Access request
					</CardTitle>
					<div className="space-y-1">
						<p className="text-xs text-muted-foreground">Requested by</p>
						<AddressDisplay address={grant.grantee} showCopy />
					</div>
				</div>
				<ConsentBadge state={grant.state} />
			</CardHeader>
			<CardContent className="space-y-4 text-sm text-muted-foreground">
				<div className="grid gap-2 md:grid-cols-2">
					<div className="space-y-1">
						<p className="text-xs uppercase tracking-wide text-muted-foreground">
							Permissions
						</p>
						<div className="flex flex-wrap gap-2">
							{grant.permissions.map((permission) => (
								<Badge key={permission} variant="secondary">
									{permissionLabels[permission]}
								</Badge>
							))}
						</div>
					</div>
					<div className="space-y-1">
						<p className="text-xs uppercase tracking-wide text-muted-foreground">
							Expires
						</p>
						<p className="text-sm text-foreground">
							{formatDate(grant.expiresAt)}
						</p>
					</div>
				</div>

				{grant.reason && (
					<div className="space-y-1">
						<p className="text-xs uppercase tracking-wide text-muted-foreground">
							Reason
						</p>
						<p className="text-sm text-foreground">{grant.reason}</p>
					</div>
				)}

				<div className="flex flex-col gap-2 sm:flex-row sm:justify-end">
					<Button
						variant="outline"
						onClick={() => onApprove(grant.id)}
						disabled={isBusy}
					>
						{isApproving ? "Approving..." : "Approve"}
					</Button>
					<Dialog>
						<DialogTrigger asChild>
							<Button
								variant="destructive"
								disabled={isBusy}
								className="cursor-pointer"
							>
								{isDenying ? "Denying..." : "Deny"}
							</Button>
						</DialogTrigger>
						<DialogContent>
							<DialogHeader>
								<DialogTitle>Deny access request?</DialogTitle>
								<DialogDescription>
									This request will be rejected immediately. The doctor will
									need to request access again.
								</DialogDescription>
							</DialogHeader>
							<DialogFooter>
								<DialogClose asChild>
									<Button variant="outline">Cancel</Button>
								</DialogClose>
								<DialogClose asChild>
									<Button
										variant="destructive"
										onClick={() => onDeny(grant.id)}
									>
										Deny request
									</Button>
								</DialogClose>
							</DialogFooter>
						</DialogContent>
					</Dialog>
				</div>
			</CardContent>
		</Card>
	);
}
