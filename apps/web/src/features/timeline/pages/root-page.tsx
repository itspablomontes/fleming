import { LoadingScreen } from "@/components/common/loading-screen";
import { useAuth } from "@/features/auth/hooks/use-auth";
import { LandingPage } from "@/features/landing/pages/landing-page";
import { TimelineViewPage } from "@/features/timeline/pages/timeline-view-page";

/**
 * RootPage component handles the conditional rendering between
 * LandingPage and TimelineViewPage based on authentication status.
 *
 * This component is used as the primary content for the root path (/).
 */
export function RootPage() {
	const { isAuthenticated, isLoading } = useAuth();

	console.debug("[RootPage] Render:", { isAuthenticated, isLoading });

	if (isLoading) {
		return <LoadingScreen />;
	}

	return isAuthenticated ? <TimelineViewPage /> : <LandingPage />;
}
