// Public API for auth feature

export { AuthButton } from "./components/auth-button";
export { AuthContext, AuthProvider, AuthStatus } from "./context/auth-context";
export { type UseAuthReturn, useAuth } from "./hooks/use-auth";
export type { Session, User, UserRole } from "./types";
