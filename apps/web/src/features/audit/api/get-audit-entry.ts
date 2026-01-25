import { apiClient } from "@/lib/api-client";

import type { AuditLogEntry } from "../types";
import { mapAuditEntry, type AuditEntryResponse } from "./mappers";

interface AuditEntryResponseWrapper {
	entry: AuditEntryResponse;
}

export const getAuditEntry = async (entryId: string): Promise<AuditLogEntry> => {
	const response = await apiClient(`/api/audit/entries/${entryId}`);
	return mapAuditEntry((response as AuditEntryResponseWrapper).entry);
};
