import { apiClient } from "@/lib/api-client";

export const checkAuth = async (): Promise<boolean> => {
    try {
        await apiClient("/auth/me");
        return true;
    } catch {
        return false;
    }
};
