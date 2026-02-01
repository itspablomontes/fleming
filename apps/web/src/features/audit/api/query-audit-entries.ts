import { apiClient } from "@/lib/api-client";

import type { AuditAction, AuditLogEntry, AuditTargetType } from "../types";
import { type AuditEntryResponse, mapAuditEntry } from "./mappers";

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

export interface AuditQueryResult {
	entries: AuditLogEntry[];
	total: number;
}

export const queryAuditEntries = async (
	params: AuditQueryParams,
): Promise<AuditQueryResult> => {
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
	
	// Backend returns { entries: [], total: number }
	// We cast to any because the apiClient return type is generic/unknown by default
	// but mapping assumes specific structure.
	const result = response as { entries: AuditEntryResponse[]; total: number };
	
	return {
		entries: result.entries.map(mapAuditEntry),
		total: result.total,
	};
};
