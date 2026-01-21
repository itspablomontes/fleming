import { Loader2, Wallet } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

/**
 * ConnectWalletButton
 * Primary CTA for wallet connection.
 * This is a presentation component — wallet logic will be added via Wagmi hooks.
 */
interface ConnectWalletButtonProps {
	className?: string;
	isConnecting?: boolean;
	isConnected?: boolean;
	onClick?: () => void;
}

export function ConnectWalletButton({
	className,
	isConnecting = false,
	isConnected = false,
	onClick,
}: ConnectWalletButtonProps) {
	if (isConnected) {
		return (
			<Button
				variant="secondary"
				className={cn("gap-2", className)}
				onClick={onClick}
			>
				<Wallet className="h-4 w-4" aria-hidden="true" />
				Sign In
			</Button>
		);
	}

	return (
		<Button
			size="lg"
			className={cn(
				"gap-2 glow-primary",
				"bg-linear-to-r from-primary to-primary/80",
				"hover:from-primary/90 hover:to-primary/70",
				"transition-all duration-300",
				className,
			)}
			onClick={onClick}
			disabled={isConnecting}
		>
			{isConnecting ? (
				<>
					<Loader2 className="h-5 w-5 animate-spin" aria-hidden="true" />
					Connecting...
				</>
			) : (
				<>
					<Wallet className="h-5 w-5" aria-hidden="true" />
					Connect Wallet
				</>
			)}
		</Button>
	);
}
