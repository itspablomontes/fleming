package vc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// GenerateSalt generates a cryptographically random salt for SD-JWT disclosures.
// Returns a base64url-encoded 16-byte random value.
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(salt), nil
}

// EncodeDisclosure encodes a disclosure in SD-JWT format.
// Format: base64url([salt, claim_name, claim_value])
func EncodeDisclosure(d *Disclosure) (string, error) {
	// Generate salt if not already set
	if d.Salt == "" {
		salt, err := GenerateSalt()
		if err != nil {
			return "", err
		}
		d.Salt = salt
	}

	// Create the disclosure array: [salt, claim_name, claim_value]
	disclosureArray := []any{d.Salt, d.Key, d.Value}

	// JSON encode
	jsonBytes, err := json.Marshal(disclosureArray)
	if err != nil {
		return "", fmt.Errorf("failed to encode disclosure: %w", err)
	}

	// Base64url encode
	encoded := base64.RawURLEncoding.EncodeToString(jsonBytes)
	d.Encoded = encoded

	return encoded, nil
}

// DecodeDisclosure decodes a base64url-encoded SD-JWT disclosure.
func DecodeDisclosure(encoded string) (*Disclosure, error) {
	// Base64url decode
	jsonBytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode disclosure: %w", err)
	}

	// Parse JSON array
	var disclosureArray []any
	if err := json.Unmarshal(jsonBytes, &disclosureArray); err != nil {
		return nil, fmt.Errorf("failed to parse disclosure: %w", err)
	}

	if len(disclosureArray) != 3 {
		return nil, fmt.Errorf("invalid disclosure format: expected 3 elements, got %d", len(disclosureArray))
	}

	salt, ok := disclosureArray[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid disclosure: salt must be string")
	}

	key, ok := disclosureArray[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid disclosure: key must be string")
	}

	return &Disclosure{
		Salt:    salt,
		Key:     key,
		Value:   disclosureArray[2],
		Encoded: encoded,
	}, nil
}

// ComputeDisclosureDigest computes the SHA-256 hash of an encoded disclosure.
// This is used in the SD-JWT payload to reference disclosed claims.
func ComputeDisclosureDigest(encoded string) string {
	hash := sha256.Sum256([]byte(encoded))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// ComputeDisclosureDigestHex computes the SHA-256 hash as hex string.
func ComputeDisclosureDigestHex(encoded string) string {
	hash := sha256.Sum256([]byte(encoded))
	return hex.EncodeToString(hash[:])
}

// DisclosureSet manages a collection of disclosures.
type DisclosureSet struct {
	disclosures map[string]*Disclosure // keyed by claim name
	digests     map[string]string      // digest -> claim name mapping
}

// NewDisclosureSet creates a new disclosure set.
func NewDisclosureSet() *DisclosureSet {
	return &DisclosureSet{
		disclosures: make(map[string]*Disclosure),
		digests:     make(map[string]string),
	}
}

// Add adds a disclosure to the set.
func (ds *DisclosureSet) Add(d *Disclosure) error {
	// Encode if not already encoded
	if d.Encoded == "" {
		if _, err := EncodeDisclosure(d); err != nil {
			return err
		}
	}

	// Compute and store digest
	digest := ComputeDisclosureDigest(d.Encoded)

	ds.disclosures[d.Key] = d
	ds.digests[digest] = d.Key

	return nil
}

// Get retrieves a disclosure by claim name.
func (ds *DisclosureSet) Get(key string) (*Disclosure, bool) {
	d, ok := ds.disclosures[key]
	return d, ok
}

// GetByDigest retrieves a disclosure by its digest.
func (ds *DisclosureSet) GetByDigest(digest string) (*Disclosure, bool) {
	key, ok := ds.digests[digest]
	if !ok {
		return nil, false
	}
	return ds.disclosures[key], true
}

// All returns all disclosures.
func (ds *DisclosureSet) All() []*Disclosure {
	result := make([]*Disclosure, 0, len(ds.disclosures))
	for _, d := range ds.disclosures {
		result = append(result, d)
	}
	return result
}

// EncodedStrings returns all encoded disclosure strings.
func (ds *DisclosureSet) EncodedStrings() []string {
	result := make([]string, 0, len(ds.disclosures))
	for _, d := range ds.disclosures {
		result = append(result, d.Encoded)
	}
	return result
}

// Digests returns all disclosure digests.
func (ds *DisclosureSet) Digests() []string {
	result := make([]string, 0, len(ds.digests))
	for digest := range ds.digests {
		result = append(result, digest)
	}
	return result
}

// SelectDisclosures selects a subset of disclosures by claim name.
// Returns only the encoded strings for the selected claims.
func (ds *DisclosureSet) SelectDisclosures(keys []string) ([]string, error) {
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		d, ok := ds.disclosures[key]
		if !ok {
			return nil, fmt.Errorf("claim not found: %s", key)
		}
		result = append(result, d.Encoded)
	}
	return result, nil
}

// VerifyDisclosure verifies that a disclosure matches the expected digest.
func VerifyDisclosure(encoded, expectedDigest string) bool {
	actualDigest := ComputeDisclosureDigest(encoded)
	return actualDigest == expectedDigest
}
