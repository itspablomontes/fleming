import type { EthAddress } from "@/types/ethereum";

import {
	type AuditAction,
	AuditAnchorStatus,
	type AuditLogEntry,
	type AuditTargetType,
} from "../types";

export interface AuditEntryResponse {
	id: string;
	actor: string;
	action: string;
	resourceType: string;
	resourceId: string;
	timestamp: string;
	metadata?: Record<string, unknown>;
	hash?: string;
	previousHash?: string;
}

export const mapAuditEntry = (entry: AuditEntryResponse): AuditLogEntry => ({
	id: entry.id,
	actor: entry.actor as EthAddress,
	action: entry.action as AuditAction,
	resourceType: entry.resourceType as AuditTargetType,
	resourceId: entry.resourceId,
	timestamp: new Date(entry.timestamp),
	metadata: entry.metadata,
	hash: entry.hash,
	previousHash: entry.previousHash,
	anchorStatus: AuditAnchorStatus.Pending,
});
