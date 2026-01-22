import { Loader2, LogOut, Wallet } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { AuthStatus } from "../context/auth-context";
import { useAuth } from "../hooks/use-auth";
import { SiweModal } from "./siwe-modal";

export function AuthButton() {
	const {
		isAuthenticated,
		isLoading,
		status,
		login,
		logout,
		walletAddress,
		isWalletConnected,
	} = useAuth();
	const [showModal, setShowModal] = useState(false);

	const handleClick = () => {
		if (isAuthenticated) {
			logout();
			return;
		}

		setShowModal(true);
	};

	const handleConfirmSign = async () => {
		await login();
		setShowModal(false);
	};

	if (isAuthenticated) {
		return (
			<Button variant="outline" onClick={logout} className="gap-2">
				<LogOut className="h-4 w-4" />
				Disconnect {walletAddress?.slice(0, 6)}...
			</Button>
		);
	}

	const getButtonContent = () => {
		if (status === AuthStatus.Connecting) {
			return (
				<>
					<Loader2 className="mr-2 h-4 w-4 animate-spin" />
					Connecting...
				</>
			);
		}
		if (status === AuthStatus.Signing) {
			return (
				<>
					<Loader2 className="mr-2 h-4 w-4 animate-spin" />
					Signing...
				</>
			);
		}
		if (isWalletConnected) {
			return (
				<>
					<Wallet className="mr-2 h-4 w-4" />
					Sign In
				</>
			);
		}
		return (
			<>
				<Wallet className="mr-2 h-4 w-4" />
				Connect Wallet
			</>
		);
	};

	return (
		<>
			<Button onClick={handleClick} disabled={isLoading}>
				{getButtonContent()}
			</Button>

			<SiweModal
				open={showModal}
				onOpenChange={setShowModal}
				address={walletAddress || null}
				isSigning={status === AuthStatus.Signing}
				onSign={handleConfirmSign}
				onCancel={() => setShowModal(false)}
			/>
		</>
	);
}
