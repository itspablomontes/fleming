import { AuthButton } from "@/features/auth/components/auth-button";
import { FloatingNav } from "@/features/landing/components/floating-nav";
import { ForDeSciSection } from "@/features/landing/components/for-desci-section";
import { HowItWorks } from "@/features/landing/components/how-it-works";
import { LandingFooter } from "@/features/landing/components/landing-footer";
import { LandingHero } from "@/features/landing/components/landing-hero";
import { MotivationSection } from "@/features/landing/components/motivation-section";
import { OpenSourceBanner } from "@/features/landing/components/open-source-banner";
import { RoadmapSection } from "@/features/landing/components/roadmap-section";
import { TrustSignalsExpanded } from "@/features/landing/components/trust-signals";
import { UseCasesSection } from "@/features/landing/components/use-cases-section";

/**
 * LandingPage
 * Full landing page for unauthenticated users.
 * DeSci-inspired design with expressive animations.
 */
export function LandingPage() {
	return (
		<div className="min-h-screen bg-background overflow-x-hidden">
			{/* Floating navigation */}
			<FloatingNav />

			{/* Hero section */}
			<LandingHero className="mb-10 md:mb-12" actionButton={<AuthButton />} />

			{/* Trust signals expanded (bento grid) */}
			<div className="mb-16 md:mb-24">
				<TrustSignalsExpanded />
			</div>

			{/* Motivation - The Philosophy */}
			<MotivationSection />

			{/* How it works - 3 step flow */}
			<div className="mb-12 md:mb-16">
				<HowItWorks />
			</div>

			{/* Use Cases - Utility First */}
			<UseCasesSection />

			{/* Roadmap - Progress & Future */}
			<RoadmapSection />

			{/* For DeSci - ecosystem integration */}
			<div className="mb-12 md:mb-16">
				<ForDeSciSection />
			</div>

			{/* Open source CTA */}
			<div className="mb-12 md:mb-16">
				<OpenSourceBanner />
			</div>

			{/* Footer */}
			<LandingFooter />
		</div>
	);
}
