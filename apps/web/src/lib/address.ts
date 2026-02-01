import { getAddress, isAddress } from "viem";
import type { EthAddress } from "@/types/ethereum";

export type AddressValidationResult =
	| { ok: true; checksum: EthAddress; normalized: EthAddress }
	| { ok: false; error: string };

/**
 * validateUserAddressInput
 *
 * Best-practice validation behavior:
 * - Accept all-lower or all-upper addresses (cannot checksum-validate those).
 * - If the user provides mixed-case, reject invalid EIP-55 checksum casing.
 * - Always return a checksummed address for display.
 */
export function validateUserAddressInput(input: string): AddressValidationResult {
	const raw = input.trim();
	if (!raw) {
		return { ok: false, error: "Address is required" };
	}

	// Fast reject: avoid throwing for obviously invalid strings.
	if (!isAddress(raw as `0x${string}`, { strict: false })) {
		return { ok: false, error: "Enter a valid Ethereum address" };
	}

	try {
		const checksum = getAddress(raw as `0x${string}`) as EthAddress
		const normalized = checksum.toLowerCase() as EthAddress;
		return { ok: true, checksum, normalized };
	} catch {
		// getAddress() rejects invalid-checksum mixed-case addresses
		return { ok: false, error: "Invalid address checksum" };
	}
}

/** Canonical storage form used by backend: lowercased. */
export function normalizeAddress(input: string): EthAddress {
	const res = validateUserAddressInput(input);
	if (!res.ok) {
		throw new Error(res.error);
	}
	return res.normalized;
}

/** Preferred display form: EIP-55 checksummed. */
export function toChecksumAddress(input: string | EthAddress): EthAddress {
	const res = validateUserAddressInput(String(input));
	if (!res.ok) {
		throw new Error(res.error);
	}
	return res.checksum;
}

