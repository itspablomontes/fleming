import { use } from "react";
import { AuthContext, type AuthStatus } from "../context/auth-context";

export interface UseAuthReturn {
	isAuthenticated: boolean;
	isLoading: boolean;
	status: AuthStatus;
	address: string | null;
	walletAddress: `0x${string}` | undefined;
	isWalletConnected: boolean;
	error: string | null;
	login: () => Promise<void>;
	logout: () => void;
}

/**
 * Hook to access authentication state and actions.
 * Must be used within an AuthProvider.
 */
export function useAuth(): UseAuthReturn {
	const context = use(AuthContext);

	if (!context) {
		throw new Error("useAuth must be used within an AuthProvider");
	}

	return context;
}
