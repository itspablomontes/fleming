package types

import (
	"fmt"
	"slices"
)

// Validator interface for types that can validate themselves.
type Validator interface {
	Validate() error
}

// ValidateRequired checks if a string field is non-empty.
func ValidateRequired(field string, value string) error {
	if value == "" {
		return NewValidationError(field, "is required")
	}
	return nil
}

// ValidateNonEmpty checks if a string field is non-empty (alias for ValidateRequired).
func ValidateNonEmpty(field string, value string) error {
	return ValidateRequired(field, value)
}

// ValidateEnum checks if a value is in the list of valid values.
func ValidateEnum(field string, value string, valid []string) error {
	if !slices.Contains(valid, value) {
		return NewValidationError(field, fmt.Sprintf("must be one of: %v", valid))
	}
	return nil
}

// ValidateID checks if an ID is valid (non-empty).
func ValidateID(field string, id ID) error {
	if id.IsEmpty() {
		return NewValidationError(field, "ID is required")
	}
	return nil
}

// ValidateWalletAddress checks if a wallet address is valid.
func ValidateWalletAddress(field string, addr WalletAddress) error {
	if addr.IsEmpty() {
		return NewValidationError(field, "wallet address is required")
	}
	return nil
}
