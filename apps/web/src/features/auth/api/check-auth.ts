import { apiClient } from "@/lib/api-client";

export interface AuthResponse {
	address: string;
	status: string;
}

export const checkAuth = async (): Promise<AuthResponse | null> => {
	try {
		const data = await apiClient("/auth/me");
		return data as AuthResponse;
	} catch {
		return null;
	}
};
