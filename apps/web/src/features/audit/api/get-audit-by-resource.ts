import { apiClient } from "@/lib/api-client";

import type { AuditLogEntry } from "../types";
import { mapAuditEntry, type AuditEntryResponse } from "./mappers";

interface AuditEntriesResponse {
	entries: AuditEntryResponse[];
}

export const getAuditEntriesByResource = async (
	resourceId: string,
): Promise<AuditLogEntry[]> => {
	const response = await apiClient(`/api/audit/resource/${resourceId}`);
	return (response as AuditEntriesResponse).entries.map(mapAuditEntry);
};
