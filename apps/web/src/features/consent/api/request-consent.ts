import { apiClient } from "@/lib/api-client";
import type { ConsentGrant, ConsentRequestPayload } from "../types";
import { type ConsentGrantResponse, mapConsentGrant } from "./mappers";

export const requestConsent = async (
	payload: ConsentRequestPayload,
): Promise<ConsentGrant> => {
	const response = await apiClient("/api/consent/request", { body: payload });
	return mapConsentGrant(response as ConsentGrantResponse);
};
