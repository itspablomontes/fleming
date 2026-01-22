/**
 * Notification Types
 */

export type NotificationType =
	| "consent_request_received"
	| "consent_granted"
	| "consent_revoked"
	| "consent_expired"
	| "record_shared";

export interface Notification {
	id: string;
	userId: string;
	type: NotificationType;
	title: string;
	message: string;
	read: boolean;
	relatedId?: string;
	createdAt: Date;
}
