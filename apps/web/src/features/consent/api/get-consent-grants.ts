import { apiClient } from "@/lib/api-client";
import type { ConsentGrant } from "../types";
import { type ConsentGrantResponse, mapConsentGrant } from "./mappers";

interface ConsentGrantsResponse {
	grants: ConsentGrantResponse[];
}

export const getActiveGrants = async (): Promise<ConsentGrant[]> => {
	const response = await apiClient("/api/consent/active");
	const payload = response as ConsentGrantsResponse;
	return payload.grants.map(mapConsentGrant);
};

export const getMyGrants = async (): Promise<ConsentGrant[]> => {
	const response = await apiClient("/api/consent/grants");
	const payload = response as ConsentGrantsResponse;
	return payload.grants.map(mapConsentGrant);
};
