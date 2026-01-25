import { apiClient } from "@/lib/api-client";

export const revokeConsent = async (grantId: string): Promise<void> => {
	if (!grantId) {
		throw new Error("Consent grant id is required");
	}
	await apiClient(`/api/consent/${grantId}/revoke`, {
		method: "POST",
	});
};
