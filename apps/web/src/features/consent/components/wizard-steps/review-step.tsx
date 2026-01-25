import type { JSX } from "react";

import { AddressDisplay } from "@/components/common/address-display";
import { Badge } from "@/components/ui/badge";

import type { ConsentPermission } from "../../types";
import type { ConsentRequestFormValues } from "../consent-request-wizard-types";

interface ReviewStepProps {
	values: ConsentRequestFormValues;
}

const permissionLabels: Record<ConsentPermission, string> = {
	read: "Read",
	write: "Write",
	share: "Share",
};

export function ReviewStep({ values }: ReviewStepProps): JSX.Element {
	return (
		<div className="space-y-4">
			<div className="space-y-2">
				<p className="text-xs uppercase tracking-wide text-muted-foreground">
					Patient
				</p>
				{values.grantor ? (
					<AddressDisplay address={values.grantor as `0x${string}`} showCopy />
				) : (
					<p className="text-sm text-muted-foreground">No address provided.</p>
				)}
			</div>

			<div className="space-y-2">
				<p className="text-xs uppercase tracking-wide text-muted-foreground">
					Permissions
				</p>
				{values.permissions.length > 0 ? (
					<div className="flex flex-wrap gap-2">
						{values.permissions.map((permission) => (
							<Badge key={permission} variant="secondary">
								{permissionLabels[permission]}
							</Badge>
						))}
					</div>
				) : (
					<p className="text-sm text-muted-foreground">
						No permissions selected.
					</p>
				)}
			</div>

			<div className="space-y-2">
				<p className="text-xs uppercase tracking-wide text-muted-foreground">
					Duration
				</p>
				<p className="text-sm text-foreground">
					{values.durationDays ? `${values.durationDays} days` : "Not set"}
				</p>
			</div>

			<div className="space-y-2">
				<p className="text-xs uppercase tracking-wide text-muted-foreground">
					Reason
				</p>
				<p className="text-sm text-foreground">
					{values.reason?.trim() || "No reason provided."}
				</p>
			</div>
		</div>
	);
}
