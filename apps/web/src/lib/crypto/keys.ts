export const ALGORITHM_AES_GCM = "AES-GCM";
export const ALGORITHM_AES_KW = "AES-KW";
export const HASH_SHA_256 = "SHA-256";
export const KEY_LENGTH_256 = 256;

/**
 * Derives a Master Key (KEK) from a wallet signature using HKDF.
 * This function is deterministic: same signature + same salt = same key.
 *
 * @param signature - The signature string from the wallet (hex).
 * @param salt - The user's unique salt (Uint8Array).
 * @returns CryptoKey - The derived master key (used for wrapping DEKs).
 */
export async function deriveMasterKey(
  signature: string,
  salt: Uint8Array,
): Promise<CryptoKey> {
  const enc = new TextEncoder();
  const signatureKeyMaterial = await window.crypto.subtle.importKey(
    "raw",
    enc.encode(signature),
    { name: "HKDF" },
    false,
    ["deriveKey"],
  );

  return window.crypto.subtle.deriveKey(
    {
      name: "HKDF",
      hash: HASH_SHA_256,
      salt: salt as BufferSource,
      info: enc.encode("Fleming-E2EE-Master-Key-v1"), // Context separation
    },
    signatureKeyMaterial,
    {
      name: ALGORITHM_AES_KW, // Used to wrap/unwrap DEKs
      length: KEY_LENGTH_256,
    },
    false, // Master Key is never exportable
    ["wrapKey", "unwrapKey"],
  );
}

/**
 * Generates a random Data Encryption Key (DEK).
 * Used to encrypt a single file.
 */
export async function generateDEK(): Promise<CryptoKey> {
  return window.crypto.subtle.generateKey(
    {
      name: ALGORITHM_AES_GCM,
      length: KEY_LENGTH_256,
    },
    true, // DEKs must be exportable to be wrapped
    ["encrypt", "decrypt"],
  );
}

/**
 * Helper to import a raw key (for testing or recovery).
 */
export async function importKey(rawKey: BufferSource): Promise<CryptoKey> {
  return window.crypto.subtle.importKey(
    "raw",
    rawKey,
    { name: ALGORITHM_AES_GCM },
    true,
    ["encrypt", "decrypt"],
  );
}

/**
 * Helper to export a key to raw bytes.
 */
export async function exportKey(key: CryptoKey): Promise<ArrayBuffer> {
  return window.crypto.subtle.exportKey("raw", key);
}
