import { apiClient } from "@/lib/api-client";

import type { AuditAction, AuditLogEntry, AuditTargetType } from "../types";
import { type AuditEntryResponse, mapAuditEntry } from "./mappers";

interface AuditEntriesResponse {
	entries: AuditEntryResponse[];
}

export interface AuditQueryParams {
	actor?: string;
	resourceId?: string;
	resourceType?: AuditTargetType;
	action?: AuditAction;
	startTime?: string;
	endTime?: string;
	limit?: number;
	offset?: number;
}

export const queryAuditEntries = async (
	params: AuditQueryParams,
): Promise<AuditLogEntry[]> => {
	const searchParams = new URLSearchParams();

	if (params.actor) {
		searchParams.set("actor", params.actor);
	}
	if (params.resourceId) {
		searchParams.set("resourceId", params.resourceId);
	}
	if (params.resourceType) {
		searchParams.set("resourceType", params.resourceType);
	}
	if (params.action) {
		searchParams.set("action", params.action);
	}
	if (params.startTime) {
		searchParams.set("startTime", params.startTime);
	}
	if (params.endTime) {
		searchParams.set("endTime", params.endTime);
	}
	if (params.limit !== undefined) {
		searchParams.set("limit", String(params.limit));
	}
	if (params.offset !== undefined) {
		searchParams.set("offset", String(params.offset));
	}

	const query = searchParams.toString();
	const response = await apiClient(
		`/api/audit/query${query ? `?${query}` : ""}`,
	);
	return (response as AuditEntriesResponse).entries.map(mapAuditEntry);
};
