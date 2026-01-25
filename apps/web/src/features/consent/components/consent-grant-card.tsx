import { format, formatDistanceToNow } from "date-fns";
import type { JSX } from "react";
import { useMemo, useState } from "react";
import { AddressDisplay } from "@/components/common/address-display";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ConsentBadge } from "@/features/consent/components/consent-badge";
import { cn } from "@/lib/utils";

import type { ConsentGrant, ConsentPermission, ConsentState } from "../types";
import { RevokeConsentDialog } from "./revoke-consent-dialog";

const permissionLabels: Record<ConsentPermission, string> = {
	read: "Read",
	write: "Write",
	share: "Share",
};

interface ConsentGrantCardProps {
	grant: ConsentGrant;
	onRevoke?: (grantId: string) => void;
	isRevoking?: boolean;
	className?: string;
}

const formatDate = (value?: Date): string => {
	if (!value) {
		return "No expiration";
	}
	return format(value, "MMM d, yyyy");
};

const getEffectiveState = (grant: ConsentGrant): ConsentState => {
	if (grant.state === "approved" && grant.expiresAt) {
		if (grant.expiresAt.getTime() <= Date.now()) {
			return "expired";
		}
	}
	return grant.state;
};

export function ConsentGrantCard({
	grant,
	onRevoke,
	isRevoking,
	className,
}: ConsentGrantCardProps): JSX.Element {
	const [detailsOpen, setDetailsOpen] = useState(false);

	const effectiveState = useMemo(() => getEffectiveState(grant), [grant]);
	const expiresIn =
		grant.expiresAt && effectiveState === "approved"
			? formatDistanceToNow(grant.expiresAt, { addSuffix: true })
			: undefined;

	return (
		<Card className={cn("border-border bg-white dark:bg-gray-900", className)}>
			<CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
				<div className="space-y-2">
					<CardTitle className="text-base font-semibold text-foreground">
						Consent grant
					</CardTitle>
					<div className="grid gap-2 md:grid-cols-2">
						<div className="space-y-1">
							<p className="text-xs text-muted-foreground">Grantor</p>
							<AddressDisplay address={grant.grantor} showCopy />
						</div>
						<div className="space-y-1">
							<p className="text-xs text-muted-foreground">Grantee</p>
							<AddressDisplay address={grant.grantee} showCopy />
						</div>
					</div>
				</div>
				<ConsentBadge state={effectiveState} />
			</CardHeader>
			<CardContent className="space-y-4 text-sm text-muted-foreground">
				<div className="grid gap-3 md:grid-cols-3">
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
						{expiresIn && (
							<p className="text-xs text-muted-foreground">{expiresIn}</p>
						)}
					</div>
					<div className="space-y-1">
						<p className="text-xs uppercase tracking-wide text-muted-foreground">
							Updated
						</p>
						<p className="text-sm text-foreground">
							{format(grant.updatedAt, "MMM d, yyyy")}
						</p>
					</div>
				</div>

				{detailsOpen && (
					<div className="space-y-2">
						{grant.reason && (
							<div className="space-y-1">
								<p className="text-xs uppercase tracking-wide text-muted-foreground">
									Reason
								</p>
								<p className="text-sm text-foreground">{grant.reason}</p>
							</div>
						)}
						{grant.scope && grant.scope.length > 0 && (
							<div className="space-y-1">
								<p className="text-xs uppercase tracking-wide text-muted-foreground">
									Scope
								</p>
								<p className="text-sm text-foreground">
									{grant.scope.join(", ")}
								</p>
							</div>
						)}
					</div>
				)}

				<div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
					<Button
						variant="ghost"
						size="sm"
						type="button"
						onClick={() => setDetailsOpen((prev) => !prev)}
						className="cursor-pointer text-xs text-muted-foreground hover:text-foreground"
					>
						{detailsOpen ? "Hide details" : "View details"}
					</Button>

					{onRevoke && effectiveState === "approved" && (
						<RevokeConsentDialog
							grant={grant}
							onRevoke={onRevoke}
							isRevoking={isRevoking}
						>
							<Button
								variant="destructive"
								size="sm"
								className="cursor-pointer"
							>
								Revoke access
							</Button>
						</RevokeConsentDialog>
					)}
				</div>
			</CardContent>
		</Card>
	);
}
