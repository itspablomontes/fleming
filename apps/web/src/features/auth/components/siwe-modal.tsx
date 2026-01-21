import { CheckCircle, Loader2, PenLine } from "lucide-react";
import { AddressDisplay } from "@/components/common/address-display";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import type { EthAddress } from "@/types/ethereum";

/**
 * SiweModal
 * Explains SIWE signing before triggering wallet popup.
 * Builds trust by explicitly stating what the signature does NOT do.
 */
interface SiweModalProps {
	open: boolean;
	onOpenChange: (open: boolean) => void;
	address: EthAddress | null;
	isSigning?: boolean;
	onSign?: () => void;
	onCancel?: () => void;
}

const trustPoints = [
	"Proves you own this address",
	"Does NOT authorize transactions",
	"Does NOT cost any gas",
];

export function SiweModal({
	open,
	onOpenChange,
	address,
	isSigning = false,
	onSign,
	onCancel,
}: SiweModalProps) {
	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className="sm:max-w-md">
				<DialogHeader>
					<DialogTitle className="text-xl">Sign in to Fleming</DialogTitle>
					<DialogDescription className="text-muted-foreground">
						Your wallet will ask you to sign a message. This:
					</DialogDescription>
				</DialogHeader>

				<div className="space-y-6 py-4">
					{/* Trust points */}
					<ul className="space-y-3">
						{trustPoints.map((point) => (
							<li
								key={point}
								className="flex items-center gap-3 text-sm text-foreground"
							>
								<CheckCircle className="h-4 w-4 text-success" />
								{point}
							</li>
						))}
					</ul>

					{/* Address display */}
					{address && (
						<div className="rounded-lg border border-border bg-card p-4">
							<p className="mb-2 text-xs text-muted-foreground">
								Signing with:
							</p>
							<AddressDisplay address={address} showCopy={false} />
						</div>
					)}

					{/* Actions */}
					<div className="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
						<Button variant="ghost" onClick={onCancel} disabled={isSigning}>
							Cancel
						</Button>
						<Button
							onClick={onSign}
							disabled={isSigning}
							className={cn(
								"gap-2",
								"bg-linear-to-r from-primary to-primary/80",
							)}
						>
							{isSigning ? (
								<>
									<Loader2 className="h-4 w-4 animate-spin" />
									Waiting for signature...
								</>
							) : (
								<>
									<PenLine className="h-4 w-4" />
									Sign & Continue
								</>
							)}
						</Button>
					</div>
				</div>
			</DialogContent>
		</Dialog>
	);
}
