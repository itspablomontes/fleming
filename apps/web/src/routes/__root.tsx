import { createRootRoute, Outlet } from "@tanstack/react-router";
import { useEffect } from "react";
import { Toaster } from "sonner";
import { AuthProvider } from "@/features/auth/context/auth-context";
import { VaultProvider } from "@/features/auth/contexts/vault-context";

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
			<VaultProvider>
				<Outlet />
				<Toaster richColors position="top-right" />
			</VaultProvider>
		</AuthProvider>
	);
}
