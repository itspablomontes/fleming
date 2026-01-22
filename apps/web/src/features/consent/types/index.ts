import type { TimelineEventType } from "@/features/timeline/types";
import type { EthAddress } from "@/types/ethereum";

/**
 * Consent Types
 */

/**
 * Lifecycle state of a consent request.
 */
export const ConsentState = {
	Requested: "requested",
	Pending: "pending",
	Granted: "granted",
	Denied: "denied",
	Revoked: "revoked",
	Expired: "expired",
} as const;

export type ConsentState = (typeof ConsentState)[keyof typeof ConsentState];

/**
 * User-facing labels for consent states.
 */
export const CONSENT_STATE_LABELS: Record<ConsentState, string> = {
	requested: "Pending Your Review",
	pending: "Processing",
	granted: "Access Granted",
	denied: "Denied",
	revoked: "Revoked",
	expired: "Expired",
};

/**
 * UI variants for different consent states.
 */
export const CONSENT_STATE_VARIANTS: Record<
	ConsentState,
	"warning" | "success" | "destructive" | "secondary"
> = {
	requested: "warning",
	pending: "warning",
	granted: "success",
	denied: "secondary",
	revoked: "destructive",
	expired: "secondary",
};

/**
 * Defines the scope of data access granted.
 */
export interface ConsentScope {
	readonly eventTypes?: readonly TimelineEventType[];
	readonly dateFrom?: Date;
	readonly dateTo?: Date;
}

/**
 * Represents a consent record on the chain/system.
 * All fields are readonly to ensure immutability.
 */
export interface Consent {
	readonly id: string;
	readonly grantorId: string;
	readonly grantorAddress: EthAddress;
	readonly granteeId: string;
	readonly granteeAddress: EthAddress;
	readonly granteeDisplayName?: string;
	readonly state: ConsentState;
	readonly scope: ConsentScope;
	readonly reason?: string;
	readonly expiresAt?: Date;
	readonly createdAt: Date;
	readonly updatedAt: Date;
	readonly revokedAt?: Date;
	readonly anchoredTxHash?: string;
}
