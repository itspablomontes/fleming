import { useNavigate } from "@tanstack/react-router";
import {
	createContext,
	type ReactNode,
	useCallback,
	useEffect,
	useState,
} from "react";
import { toast } from "sonner";
import {
	useConnect,
	useConnection,
	useDisconnect,
	useSignMessage,
} from "wagmi";
import { injected } from "wagmi/connectors";
import { deleteCookie, setCookie } from "@/lib/cookie-utils";
import {
	login as apiLogin,
	logout as apiLogout,
	checkAuth,
	getChallenge,
} from "../api";

export const AuthStatus = {
	Idle: "idle",
	Initializing: "initializing",
	Connecting: "connecting",
	Signing: "signing",
	Authenticated: "authenticated",
	Error: "error",
} as const;

export type AuthStatus = (typeof AuthStatus)[keyof typeof AuthStatus];

interface AuthState {
	status: AuthStatus;
	address: string | null;
	error: string | null;
}

interface AuthContextValue extends AuthState {
	isAuthenticated: boolean;
	isLoading: boolean;
	login: () => Promise<void>;
	logout: () => void;
	walletAddress: `0x${string}` | undefined;
	isWalletConnected: boolean;
}

const initialState: AuthState = {
	status: AuthStatus.Initializing,
	address: null,
	error: null,
};

export const AuthContext = createContext<AuthContextValue | null>(null);

interface AuthProviderProps {
	children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
	const navigate = useNavigate();
	const {
		address: walletAddress,
		isConnected: isWalletConnected,
		chainId,
	} = useConnection();
	const { connectAsync } = useConnect();
	const { disconnect: disconnectWallet } = useDisconnect();
	const { signMessageAsync } = useSignMessage();

	const [state, setState] = useState<AuthState>(initialState);

	useEffect(() => {
		const initAuth = async () => {
			console.debug("[Auth] Initializing...");
			try {
				const userData = await checkAuth();
				if (userData) {
					console.debug("[Auth] User data found:", userData.address);
					setState({
						status: AuthStatus.Authenticated,
						address: userData.address,
						error: null,
					});
					setCookie("fleming_has_session", "true", 7);
				} else {
					console.debug("[Auth] No session found on server");
					setState((prev) => ({ ...prev, status: AuthStatus.Idle }));
					deleteCookie("fleming_has_session");
				}
			} catch (error) {
				console.error("[Auth] Initialization error:", error);
				setState((prev) => ({ ...prev, status: AuthStatus.Idle }));
				deleteCookie("fleming_has_session");
			}
		};
		initAuth();
	}, []);

	const login = useCallback(async () => {
		try {
			let currentAddress = walletAddress;
			if (!isWalletConnected || !currentAddress) {
				setState((prev) => ({
					...prev,
					status: AuthStatus.Connecting,
					error: null,
				}));

				const result = await connectAsync({ connector: injected() });
				currentAddress = result.accounts[0];

				if (!currentAddress) {
					throw new Error("No account returned from wallet");
				}
			}

			setState((prev) => ({ ...prev, status: AuthStatus.Signing }));

			const siweMessage = await getChallenge({
				address: currentAddress,
				domain: window.location.host,
				uri: window.location.origin,
				chainId: chainId ?? 1,
			});

			const signature = await signMessageAsync({ message: siweMessage });

			await apiLogin(currentAddress, signature);

			setState({
				status: AuthStatus.Authenticated,
				address: currentAddress,
				error: null,
			});

			console.debug("[Auth] Login successful for:", currentAddress);
			setCookie("fleming_has_session", "true", 7);
			toast.success("Successfully signed in!");

			navigate({ to: "/" });
		} catch (error) {
			const message =
				error instanceof Error ? error.message : "Authentication failed";
			setState((prev) => ({
				...prev,
				status: AuthStatus.Error,
				error: message,
			}));

			if (
				message.includes("User rejected") ||
				message.includes("user rejected")
			) {
				return;
			}

			toast.error(message);
		}
	}, [
		walletAddress,
		isWalletConnected,
		connectAsync,
		signMessageAsync,
		navigate,
		chainId,
	]);

	const logout = useCallback(async () => {
		try {
			await apiLogout();
		} catch (error) {
			console.error("Logout error:", error);
		}
		disconnectWallet();
		deleteCookie("fleming_has_session");
		setState({ ...initialState, status: AuthStatus.Idle });
		navigate({ to: "/" });
		toast.success("Logged out");
	}, [disconnectWallet, navigate]);

	const isAuthenticated =
		state.status === AuthStatus.Authenticated && !!state.address;
	const isLoading =
		state.status === AuthStatus.Initializing ||
		state.status === AuthStatus.Connecting ||
		state.status === AuthStatus.Signing;

	const value: AuthContextValue = {
		...state,
		isAuthenticated,
		isLoading,
		login,
		logout,
		walletAddress,
		isWalletConnected,
	};

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
