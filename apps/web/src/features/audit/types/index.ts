import type { EthAddress } from "@/types/ethereum";

/**
 * Audit Log Types
 */

export const AuditAction = {
	ConsentRequested: "consent_requested",
	ConsentGranted: "consent_granted",
	ConsentDenied: "consent_denied",
	ConsentRevoked: "consent_revoked",
	ConsentExpired: "consent_expired",
	RecordUploaded: "record_uploaded",
	RecordViewed: "record_viewed",
	RecordDownloaded: "record_downloaded",
	UserAuthenticated: "user_authenticated",
	UserLoggedOut: "user_logged_out",
} as const;

export type AuditAction = (typeof AuditAction)[keyof typeof AuditAction];

export const AuditTargetType = {
	Consent: "consent",
	Event: "event",
	User: "user",
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
	actorId: string;
	actorAddress: EthAddress;
	targetId?: string;
	targetType?: AuditTargetType;
	timestamp: Date;
	metadata?: Record<string, unknown>;
	anchoredTxHash?: string;
	anchorStatus: AuditAnchorStatus;
}
