import type { EthAddress } from "@/types/ethereum";

/**
 * Audit Log Types
 */

export const AuditAction = {
	ConsentRequested: "consent.request",
	ConsentGranted: "consent.approve",
	ConsentDenied: "consent.deny",
	ConsentRevoked: "consent.revoke",
	ConsentExpired: "consent.expire",
	RecordCreated: "create",
	RecordRead: "read",
	RecordUpdated: "update",
	RecordDeleted: "delete",
	FileUploaded: "file.upload",
	FileDownloaded: "file.download",
	FileShared: "file.share",
	UserAuthenticated: "auth.login",
	UserLoggedOut: "auth.logout",
} as const;

export type AuditAction = (typeof AuditAction)[keyof typeof AuditAction];

export const AuditTargetType = {
	Consent: "consent",
	Event: "event",
	File: "file",
	Session: "session",
} as const;

export type AuditTargetType =
	(typeof AuditTargetType)[keyof typeof AuditTargetType];

export const AuditAnchorStatus = {
	Pending: "pending",
	Anchored: "anchored",
	Failed: "failed",
} as const;

export type AuditAnchorStatus =
	(typeof AuditAnchorStatus)[keyof typeof AuditAnchorStatus];

export interface AuditLogEntry {
	id: string;
	action: AuditAction;
	actor: EthAddress;
	resourceId: string;
	resourceType: AuditTargetType;
	timestamp: Date;
	metadata?: Record<string, unknown>;
	hash?: string;
	previousHash?: string;
	anchoredTxHash?: string;
	anchorStatus?: AuditAnchorStatus;
}

export interface MerkleProofStep {
	hash: string;
	isLeft: boolean;
}

export interface MerkleProof {
	entryHash: string;
	steps: MerkleProofStep[];
}
