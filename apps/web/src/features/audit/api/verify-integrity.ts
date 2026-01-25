import { apiClient } from "@/lib/api-client";

interface VerifyIntegrityResponse {
	valid: boolean;
	message: string;
}

export const verifyIntegrity = async (): Promise<VerifyIntegrityResponse> => {
	return apiClient("/api/audit/verify");
};
