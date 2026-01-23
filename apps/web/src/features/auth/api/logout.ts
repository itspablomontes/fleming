import { apiClient } from "@/lib/api-client";

export const logout = async (): Promise<void> => {
	await apiClient("/auth/logout", { method: "POST" });
};
