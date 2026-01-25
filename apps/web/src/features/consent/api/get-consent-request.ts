import { apiClient } from "@/lib/api-client";
import type { ConsentGrant } from "../types";
import { type ConsentGrantResponse, mapConsentGrant } from "./mappers";

export const getConsentGrant = async (
	grantId: string,
): Promise<ConsentGrant> => {
	if (!grantId) {
		throw new Error("Consent grant id is required");
	}
	const response = await apiClient(`/api/consent/${grantId}`);
	return mapConsentGrant(response as ConsentGrantResponse);
};
