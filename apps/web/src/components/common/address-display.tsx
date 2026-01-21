import { Check, Copy } from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import type { EthAddress } from "@/types";

interface AddressDisplayProps {
	address: EthAddress;
	displayName?: string;
	showAvatar?: boolean;
	showCopy?: boolean;
	truncate?: boolean;
	className?: string;
}

/**
 * Generates a deterministic color from an Ethereum address.
 * Simple implementation inspired by Jazzicon.
 */
function addressToColor(address: string): string {
	const hash = address.slice(2, 10);
	const hue = Number.parseInt(hash, 16) % 360;
	return `hsl(${hue}, 65%, 55%)`;
}

/**
 * Generates a simple gradient background from address.
 */
function addressToGradient(address: string): string {
	const color1 = addressToColor(address);
	const hash2 = address.slice(10, 18);
	const hue2 = Number.parseInt(hash2, 16) % 360;
	const color2 = `hsl(${hue2}, 60%, 50%)`;
	return `linear-gradient(135deg, ${color1}, ${color2})`;
}

/**
 * Truncates an Ethereum address for display.
 * E.g., 0x1234...5678
 */
function truncateAddress(
	address: string,
	startChars = 6,
	endChars = 4,
): string {
	if (address.length <= startChars + endChars + 3) return address;
	return `${address.slice(0, startChars)}...${address.slice(-endChars)}`;
}

/**
 * AddressDisplay
 * Displays an Ethereum address with optional avatar, ENS name, and copy button.
 * Follows Web3 UX best practices for address display.
 */
export function AddressDisplay({
	address,
	displayName,
	showAvatar = true,
	showCopy = true,
	truncate = true,
	className,
}: AddressDisplayProps) {
	const [copied, setCopied] = useState(false);

	const gradient = useMemo(() => addressToGradient(address), [address]);
	const displayAddress = truncate ? truncateAddress(address) : address;

	const handleCopy = useCallback(async () => {
		try {
			await navigator.clipboard.writeText(address);
			setCopied(true);
			setTimeout(() => setCopied(false), 2000);
		} catch (err) {
			console.error("Failed to copy address:", err);
		}
	}, [address]);

	return (
		<div className={cn("inline-flex items-center gap-2", className)}>
			{showAvatar && (
				<Avatar className="h-6 w-6">
					<div
						className="h-full w-full rounded-full"
						style={{ background: gradient }}
						aria-hidden="true"
					/>
				</Avatar>
			)}

			<div className="flex flex-col">
				{displayName && (
					<span className="text-sm font-medium text-foreground">
						{displayName}
					</span>
				)}
				<code className="font-mono text-xs text-muted-foreground">
					{displayAddress}
				</code>
			</div>

			{showCopy && (
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<Button
								variant="ghost"
								size="icon"
								className="h-7 w-7"
								onClick={handleCopy}
								aria-label={copied ? "Copied!" : "Copy address"}
							>
								{copied ? (
									<Check className="h-3.5 w-3.5 text-success" />
								) : (
									<Copy className="h-3.5 w-3.5" />
								)}
							</Button>
						</TooltipTrigger>
						<TooltipContent>
							<p>{copied ? "Copied!" : "Copy address"}</p>
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			)}
		</div>
	);
}
