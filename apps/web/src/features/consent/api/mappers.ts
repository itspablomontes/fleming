import type { EthAddress } from "@/types/ethereum";

import type { ConsentGrant, ConsentPermission, ConsentState } from "../types";

export interface ConsentGrantResponse {
	id: string;
	grantor: string;
	grantee: string;
	scope?: string[] | null;
	permissions: string[];
	state: string;
	reason?: string | null;
	expiresAt?: string | null;
	createdAt: string;
	updatedAt: string;
}

const zeroTimePrefix = "0001-01-01";

const parseOptionalDate = (value?: string | null): Date | undefined => {
	if (!value) {
		return undefined;
	}
	if (value.startsWith(zeroTimePrefix)) {
		return undefined;
	}
	const parsed = new Date(value);
	if (Number.isNaN(parsed.getTime())) {
		return undefined;
	}
	return parsed;
};

const parseRequiredDate = (value: string): Date => {
	const parsed = new Date(value);
	if (Number.isNaN(parsed.getTime())) {
		throw new Error("Invalid date value in consent payload");
	}
	return parsed;
};

export const mapConsentGrant = (grant: ConsentGrantResponse): ConsentGrant => ({
	id: grant.id,
	grantor: grant.grantor as EthAddress,
	grantee: grant.grantee as EthAddress,
	scope: grant.scope ?? undefined,
	permissions: grant.permissions as ConsentPermission[],
	state: grant.state as ConsentState,
	reason: grant.reason ?? undefined,
	expiresAt: parseOptionalDate(grant.expiresAt),
	createdAt: parseRequiredDate(grant.createdAt),
	updatedAt: parseRequiredDate(grant.updatedAt),
});
