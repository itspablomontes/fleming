import { useState } from "react";
import { SiweModal } from "@/features/auth/components/siwe-modal";
import { LandingHero } from "@/features/landing/components/landing-hero";
import type { EthAddress } from "@/types/ethereum";

export function LandingPage() {
	// Mock state for now - will be replaced by Wagmi hooks in next phase
	const [isConnecting, setIsConnecting] = useState(false);
	const [isConnected, setIsConnected] = useState(false);
	const [isSiweOpen, setIsSiweOpen] = useState(false);
	const [isSigning, setIsSigning] = useState(false);
	const [userAddress, setUserAddress] = useState<EthAddress | null>(null);

	const handleConnect = () => {
		setIsConnecting(true);
		// Simulate connection delay
		setTimeout(() => {
			setIsConnecting(false);
			setIsConnected(true);
			setUserAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F");
			setIsSiweOpen(true);
		}, 1000);
	};

	const handleSign = () => {
		setIsSigning(true);
		// Simulate signing delay
		setTimeout(() => {
			setIsSigning(false);
			setIsSiweOpen(false);
			// Would redirect to dashboard here
			console.log("Signed in!");
		}, 1500);
	};

	const handleSignIn = () => {
		if (isConnected) {
			setIsSiweOpen(true);
		} else {
			handleConnect();
		}
	};

	return (
		<>
			<LandingHero
				isConnecting={isConnecting}
				isConnected={isConnected}
				onConnect={handleConnect}
				onSignIn={handleSignIn}
			/>

			<SiweModal
				open={isSiweOpen}
				onOpenChange={setIsSiweOpen}
				address={userAddress}
				isSigning={isSigning}
				onSign={handleSign}
				onCancel={() => setIsSiweOpen(false)}
			/>
		</>
	);
}
