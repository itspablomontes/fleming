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
 * JWT payload from backend.
 * This represents what's actually in the JWT claims.
 */
export interface JWTPayload {
	readonly sub: EthAddress; // Subject = wallet address
	readonly exp: number; // Expiration timestamp (Unix seconds)
	readonly iat: number; // Issued at timestamp (Unix seconds)
}

/**
 * Represents an authenticated user in the system.
 * Note: Currently backend only provides address via JWT.
 * Other fields are populated client-side or will be added later.
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
