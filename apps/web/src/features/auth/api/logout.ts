import { apiClient } from "@/lib/api-client";

export const logout = async (): Promise<void> => {
	await apiClient("/api/auth/logout", { method: "POST" });
};
