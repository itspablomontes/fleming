import { apiClient } from "@/lib/api-client";

export interface LoginResponse {
	success: boolean;
}

export const login = async (
	address: string,
	signature: string,
): Promise<boolean> => {
	const res = await apiClient("/auth/login", {
		body: { address, signature },
	});
	return res.success;
};
