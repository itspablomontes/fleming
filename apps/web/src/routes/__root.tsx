import { createRootRoute, Outlet } from "@tanstack/react-router";
import { useEffect } from "react";
import { Toaster } from "sonner";
import { AuthProvider } from "@/features/auth/context/auth-context";

export const Route = createRootRoute({
	component: RootComponent,
});

function RootComponent() {
	// Initialize theme on mount
	useEffect(() => {
		const stored = localStorage.getItem("fleming-theme");
		const theme = stored === "light" ? "light" : "dark";
		document.documentElement.classList.add(theme);
	}, []);

	return (
		<AuthProvider>
			{/* Skip link for keyboard users */}
			<a
				href="#main-content"
				className="sr-only focus:not-sr-only focus:absolute focus:z-100 focus:top-4 focus:left-4 focus:px-4 focus:py-2 focus:bg-primary focus:text-primary-foreground focus:rounded-md"
			>
				Skip to main content
			</a>
			<Outlet />
			<Toaster richColors position="top-right" />
		</AuthProvider>
	);
}
