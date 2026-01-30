import type { EthAddress } from "@/types/ethereum";

/**
 * Consent Types
 */

/**
 * Lifecycle state of a consent request.
 */
export const ConsentState = {
	Requested: "requested",
	Approved: "approved",
	Denied: "denied",
	Revoked: "revoked",
	Expired: "expired",
	Suspended: "suspended",
} as const;

export type ConsentState = (typeof ConsentState)[keyof typeof ConsentState];

/**
 * User-facing labels for consent states.
 */
export const CONSENT_STATE_LABELS: Record<ConsentState, string> = {
	requested: "Pending Review",
	approved: "Access Granted",
	denied: "Denied",
	revoked: "Revoked",
	expired: "Expired",
	suspended: "Suspended",
};

/**
 * UI variants for different consent states.
 */
export const CONSENT_STATE_VARIANTS: Record<
	ConsentState,
	"warning" | "success" | "destructive" | "secondary"
> = {
	requested: "warning",
	approved: "success",
	denied: "secondary",
	revoked: "destructive",
	expired: "secondary",
	suspended: "warning",
};

/**
 * Permission types for consent grants.
 */
export const ConsentPermission = {
	Read: "read",
	Write: "write",
	Share: "share",
} as const;

export type ConsentPermission =
	(typeof ConsentPermission)[keyof typeof ConsentPermission];

/**
 * Represents a consent grant stored in the backend.
 */
export interface ConsentGrant {
	readonly id: string;
	readonly grantor: EthAddress;
	readonly grantee: EthAddress;
	readonly scope?: readonly string[];
	readonly permissions: readonly ConsentPermission[];
	readonly state: ConsentState;
	readonly reason?: string;
	readonly expiresAt?: Date;
	readonly createdAt: Date;
	readonly updatedAt: Date;
}

export interface ConsentRequestPayload {
	readonly grantor: EthAddress;
	readonly permissions: readonly ConsentPermission[];
	readonly reason?: string;
	readonly durationDays?: number;
}
