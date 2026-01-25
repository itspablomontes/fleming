import { ALGORITHM_AES_GCM, ALGORITHM_AES_KW } from "./keys";

// 12 bytes IV is standard for AES-GCM
const IV_LENGTH = 12;

interface EncryptedFile {
  ciphertext: ArrayBuffer;
  iv: Uint8Array;
}

/**
 * Encrypts file data using a Data Encryption Key (DEK).
 * AES-GCM provides both confidentiality and integrity.
 */
export async function encryptFile(
  data: BufferSource,
  dek: CryptoKey,
): Promise<EncryptedFile> {
  const iv = window.crypto.getRandomValues(new Uint8Array(IV_LENGTH));

  const ciphertext = await window.crypto.subtle.encrypt(
    {
      name: ALGORITHM_AES_GCM,
      iv: iv,
    },
    dek,
    data,
  );

  return {
    ciphertext,
    iv,
  };
}

/**
 * Decrypts file data using a Data Encryption Key (DEK).
 * Throws error if integrity check fails.
 */
export async function decryptFile(
  ciphertext: BufferSource,
  iv: Uint8Array,
  dek: CryptoKey,
): Promise<ArrayBuffer> {
  return window.crypto.subtle.decrypt(
    {
      name: ALGORITHM_AES_GCM,
      iv: iv as BufferSource,
    },
    dek,
    ciphertext,
  );
}

/**
 * Wraps (encrypts) a DEK using the Master Key (KEK).
 * Uses AES-KW (Key Wrap) algorithm.
 */
export async function wrapKey(
  dek: CryptoKey,
  kek: CryptoKey,
): Promise<ArrayBuffer> {
  return window.crypto.subtle.wrapKey("raw", dek, kek, ALGORITHM_AES_KW);
}

/**
 * Unwraps (decrypts) a wrapped DEK using the Master Key (KEK).
 */
export async function unwrapKey(
  wrappedKey: BufferSource,
  kek: CryptoKey,
): Promise<CryptoKey> {
  return window.crypto.subtle.unwrapKey(
    "raw",
    wrappedKey,
    kek,
    ALGORITHM_AES_KW,
    { name: ALGORITHM_AES_GCM }, // The algorithm of the key being unwrapped
    true, // Unwrapped DEK is exportable
    ["encrypt", "decrypt"],
  );
}
