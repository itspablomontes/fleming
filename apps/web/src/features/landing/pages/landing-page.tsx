import { AuthButton } from "@/features/auth/components/auth-button";
import { LandingHero } from "@/features/landing/components/landing-hero";

export function LandingPage() {
	return (
		<div className="min-h-screen bg-background">
			<LandingHero actionButton={<AuthButton />} />
		</div>
	);
}
