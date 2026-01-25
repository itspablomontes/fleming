import type { JSX, ReactNode } from "react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
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
import { Label } from "@/components/ui/label";

import type { ConsentGrant } from "../types";

interface RevokeConsentDialogProps {
	grant: ConsentGrant;
	onRevoke: (grantId: string) => void;
	isRevoking?: boolean;
	children: ReactNode;
}

export function RevokeConsentDialog({
	grant,
	onRevoke,
	isRevoking,
	children,
}: RevokeConsentDialogProps): JSX.Element {
	const [isOpen, setIsOpen] = useState(false);
	const [confirmed, setConfirmed] = useState(false);

	return (
		<Dialog
			open={isOpen}
			onOpenChange={(next) => {
				setIsOpen(next);
				if (!next) {
					setConfirmed(false);
				}
			}}
		>
			<DialogTrigger asChild>{children}</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Revoke access?</DialogTitle>
					<DialogDescription>
						Revoking access immediately blocks this provider from viewing or
						updating your timeline.
					</DialogDescription>
				</DialogHeader>

				<div className="space-y-3">
					<div className="flex items-start gap-3 rounded-md border border-border p-3">
						<input
							id={`revoke-confirm-${grant.id}`}
							type="checkbox"
							className="mt-1 h-4 w-4 rounded border-gray-300 accent-rose-600 focus-visible:ring-2 focus-visible:ring-rose-500 focus-visible:ring-offset-2"
							checked={confirmed}
							onChange={(event) => setConfirmed(event.target.checked)}
						/>
						<Label
							htmlFor={`revoke-confirm-${grant.id}`}
							className="cursor-pointer text-sm text-foreground"
						>
							I understand this action cannot be undone.
						</Label>
					</div>
				</div>

				<DialogFooter>
					<DialogClose asChild>
						<Button variant="outline">Cancel</Button>
					</DialogClose>
					<Button
						variant="destructive"
						disabled={!confirmed || isRevoking}
						onClick={() => onRevoke(grant.id)}
					>
						{isRevoking ? "Revoking..." : "Revoke access"}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
