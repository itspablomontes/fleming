import { createRootRoute, Outlet } from "@tanstack/react-router";
import { Toaster } from "sonner";
import { AuthProvider } from "@/features/auth/context/auth-context";

export const Route = createRootRoute({
	component: () => (
		<AuthProvider>
			<Outlet />
			<Toaster richColors position="top-right" />
		</AuthProvider>
	),
});
