import { apiClient } from "@/lib/api-client";

import type { AuditLogEntry } from "../types";
import { mapAuditEntry, type AuditEntryResponse } from "./mappers";

interface AuditLogsResponse {
	entries: AuditEntryResponse[];
}

export const getAuditLogs = async (): Promise<AuditLogEntry[]> => {
	const response = await apiClient("/api/audit");
	return (response as AuditLogsResponse).entries.map(mapAuditEntry);
};
