package vc

import (
	"encoding/base64"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Errorf("GenerateSalt() error = %v", err)
		return
	}

	if salt == "" {
		t.Error("GenerateSalt() returned empty string")
	}

	// Should be base64url encoded (no padding)
	decoded, err := base64.RawURLEncoding.DecodeString(salt)
	if err != nil {
		t.Errorf("GenerateSalt() returned invalid base64url: %v", err)
		return
	}

	if len(decoded) != 16 {
		t.Errorf("GenerateSalt() decoded length = %d, want 16", len(decoded))
	}

	// Should generate different salts
	salt2, _ := GenerateSalt()
	if salt == salt2 {
		t.Error("GenerateSalt() should generate unique salts")
	}
}

func TestEncodeDisclosure(t *testing.T) {
	d := &Disclosure{
		Key:   "marker",
		Value: "718-7",
	}

	encoded, err := EncodeDisclosure(d)
	if err != nil {
		t.Errorf("EncodeDisclosure() error = %v", err)
		return
	}

	if encoded == "" {
		t.Error("EncodeDisclosure() returned empty string")
	}

	if d.Encoded != encoded {
		t.Error("EncodeDisclosure() should set Encoded field")
	}

	if d.Salt == "" {
		t.Error("EncodeDisclosure() should generate salt if missing")
	}

	// Encoding same disclosure again should use existing salt
	originalSalt := d.Salt
	encoded2, err := EncodeDisclosure(d)
	if err != nil {
		t.Errorf("EncodeDisclosure() error = %v", err)
		return
	}
	if d.Salt != originalSalt {
		t.Error("EncodeDisclosure() should preserve existing salt")
	}
	if encoded2 != encoded {
		t.Error("EncodeDisclosure() with same salt should produce same encoding")
	}
}

func TestDecodeDisclosure(t *testing.T) {
	// First encode a disclosure
	d := &Disclosure{
		Key:   "marker",
		Value: "718-7",
	}
	encoded, err := EncodeDisclosure(d)
	if err != nil {
		t.Fatalf("EncodeDisclosure() error = %v", err)
	}

	// Decode it
	decoded, err := DecodeDisclosure(encoded)
	if err != nil {
		t.Errorf("DecodeDisclosure() error = %v", err)
		return
	}

	if decoded.Salt != d.Salt {
		t.Errorf("DecodeDisclosure() salt = %v, want %v", decoded.Salt, d.Salt)
	}
	if decoded.Key != d.Key {
		t.Errorf("DecodeDisclosure() key = %v, want %v", decoded.Key, d.Key)
	}
	if decoded.Value != d.Value {
		t.Errorf("DecodeDisclosure() value = %v, want %v", decoded.Value, d.Value)
	}
	if decoded.Encoded != encoded {
		t.Error("DecodeDisclosure() should set Encoded field")
	}

	// Invalid base64
	_, err = DecodeDisclosure("invalid-base64!")
	if err == nil {
		t.Error("DecodeDisclosure() with invalid base64 should error")
	}

	// Invalid JSON
	invalidJSON := base64.RawURLEncoding.EncodeToString([]byte("not json"))
	_, err = DecodeDisclosure(invalidJSON)
	if err == nil {
		t.Error("DecodeDisclosure() with invalid JSON should error")
	}

	// Wrong array length
	wrongLength := base64.RawURLEncoding.EncodeToString([]byte(`["salt","key"]`))
	_, err = DecodeDisclosure(wrongLength)
	if err == nil {
		t.Error("DecodeDisclosure() with wrong array length should error")
	}
}

func TestComputeDisclosureDigest(t *testing.T) {
	encoded := "eyJzYWx0IjoiMTIzNCIsImtleSI6Im1hcmtlciIsInZhbHVlIjoiNzE4LTcifQ"
	digest := ComputeDisclosureDigest(encoded)

	if digest == "" {
		t.Error("ComputeDisclosureDigest() returned empty string")
	}

	// Should be base64url encoded
	_, err := base64.RawURLEncoding.DecodeString(digest)
	if err != nil {
		t.Errorf("ComputeDisclosureDigest() returned invalid base64url: %v", err)
	}

	// Same input should produce same digest
	digest2 := ComputeDisclosureDigest(encoded)
	if digest != digest2 {
		t.Error("ComputeDisclosureDigest() should be deterministic")
	}
}

func TestDisclosureSet_Add(t *testing.T) {
	ds := NewDisclosureSet()

	d := &Disclosure{
		Key:   "marker",
		Value: "718-7",
	}

	err := ds.Add(d)
	if err != nil {
		t.Errorf("Add() error = %v", err)
		return
	}

	if d.Encoded == "" {
		t.Error("Add() should encode disclosure")
	}

	retrieved, ok := ds.Get("marker")
	if !ok {
		t.Error("Get() should find added disclosure")
	}
	if retrieved.Key != "marker" {
		t.Errorf("Get() key = %v, want marker", retrieved.Key)
	}
}

func TestDisclosureSet_GetByDigest(t *testing.T) {
	ds := NewDisclosureSet()

	d := &Disclosure{
		Key:   "marker",
		Value: "718-7",
	}
	ds.Add(d)

	digest := ComputeDisclosureDigest(d.Encoded)
	retrieved, ok := ds.GetByDigest(digest)
	if !ok {
		t.Error("GetByDigest() should find disclosure by digest")
	}
	if retrieved.Key != "marker" {
		t.Errorf("GetByDigest() key = %v, want marker", retrieved.Key)
	}

	_, ok = ds.GetByDigest("nonexistent")
	if ok {
		t.Error("GetByDigest() should not find nonexistent digest")
	}
}

func TestDisclosureSet_SelectDisclosures(t *testing.T) {
	ds := NewDisclosureSet()

	ds.Add(&Disclosure{Key: "marker", Value: "718-7"})
	ds.Add(&Disclosure{Key: "value", Value: 15.0})
	ds.Add(&Disclosure{Key: "range", Value: "13.5-17.5"})

	selected, err := ds.SelectDisclosures([]string{"marker", "value"})
	if err != nil {
		t.Errorf("SelectDisclosures() error = %v", err)
		return
	}

	if len(selected) != 2 {
		t.Errorf("SelectDisclosures() returned %d disclosures, want 2", len(selected))
	}

	// Missing key should error
	_, err = ds.SelectDisclosures([]string{"nonexistent"})
	if err == nil {
		t.Error("SelectDisclosures() with missing key should error")
	}
}

func TestVerifyDisclosure(t *testing.T) {
	encoded := "eyJzYWx0IjoiMTIzNCIsImtleSI6Im1hcmtlciIsInZhbHVlIjoiNzE4LTcifQ"
	digest := ComputeDisclosureDigest(encoded)

	if !VerifyDisclosure(encoded, digest) {
		t.Error("VerifyDisclosure() should verify correct digest")
	}

	if VerifyDisclosure(encoded, "wrong-digest") {
		t.Error("VerifyDisclosure() should reject wrong digest")
	}
}
