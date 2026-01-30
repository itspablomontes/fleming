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
	ConsentSuspend: "consent.suspend",
	ConsentResume: "consent.resume",
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
	VC: "vc",
	ZKProof: "zk_proof",
	Attestation: "attestation",
} as const;

export type AuditTargetType =
	(typeof AuditTargetType)[keyof typeof AuditTargetType];

/**
 * Human-readable label for audit action (protocol-aligned; unknown actions get formatted).
 */
export function formatAuditAction(action: string): string {
	const known: Record<string, string> = {
		[AuditAction.ConsentRequested]: "Consent requested",
		[AuditAction.ConsentGranted]: "Consent approved",
		[AuditAction.ConsentDenied]: "Consent denied",
		[AuditAction.ConsentRevoked]: "Consent revoked",
		[AuditAction.ConsentExpired]: "Consent expired",
		[AuditAction.ConsentSuspend]: "Consent suspended",
		[AuditAction.ConsentResume]: "Consent resumed",
		[AuditAction.RecordCreated]: "Create",
		[AuditAction.RecordRead]: "Read",
		[AuditAction.RecordUpdated]: "Update",
		[AuditAction.RecordDeleted]: "Delete",
		[AuditAction.FileUploaded]: "File upload",
		[AuditAction.FileDownloaded]: "File download",
		[AuditAction.FileShared]: "File share",
		[AuditAction.UserAuthenticated]: "Login",
		[AuditAction.UserLoggedOut]: "Logout",
	};
	if (known[action]) return known[action];
	return action.replace(/\./g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

/**
 * Human-readable label for resource type (protocol-aligned; unknown types get formatted).
 */
export function formatAuditResourceType(resourceType: string): string {
	const known: Record<string, string> = {
		[AuditTargetType.Consent]: "Consent",
		[AuditTargetType.Event]: "Event",
		[AuditTargetType.File]: "File",
		[AuditTargetType.Session]: "Session",
		[AuditTargetType.VC]: "Verifiable credential",
		[AuditTargetType.ZKProof]: "ZK proof",
		[AuditTargetType.Attestation]: "Attestation",
	};
	if (known[resourceType]) return known[resourceType];
	return resourceType.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

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
