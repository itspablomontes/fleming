import type { EthAddress } from "@/types/ethereum";

/**
 * Auth & Identity Types
 */

/**
 * Defines the role of a user in the system.
 */
export const UserRole = {
	Patient: "patient",
	Doctor: "doctor",
	Researcher: "researcher",
} as const;

export type UserRole = (typeof UserRole)[keyof typeof UserRole];

/**
 * Represents an authenticated user in the system.
 * All fields are readonly to ensure immutability.
 */
export interface User {
	readonly id: string;
	readonly address: EthAddress;
	readonly role: UserRole;
	readonly createdAt: Date;
	readonly displayName?: string;
}

/**
 * Represents an active user session.
 */
export interface Session {
	readonly userId: string;
	readonly address: EthAddress;
	readonly role: UserRole;
	readonly issuedAt: Date;
	readonly expiresAt: Date;
}
