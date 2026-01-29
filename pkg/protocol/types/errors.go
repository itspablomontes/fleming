package types

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidID = errors.New("invalid ID: cannot be empty")

	ErrInvalidAddress = errors.New("invalid address: must be 0x followed by 40 hex characters")

	ErrValidationFailed = errors.New("validation failed")
)

// ProtocolError is the base interface for all protocol errors.
type ProtocolError interface {
	error
	Code() string
	Field() string
}

// ValidationError represents a field validation error.
type ValidationError struct {
	field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.Message)
}

func (e ValidationError) Code() string {
	return "VALIDATION_ERROR"
}

func (e ValidationError) Field() string {
	return e.field
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{field: field, Message: message}
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation succeeded"
	}
	if len(ve) == 1 {
		return ve[0].Error()
	}
	msg := "multiple validation errors: "
	for i, e := range ve {
		if i > 0 {
			msg += "; "
		}
		msg += e.Error()
	}
	return msg
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

func (ve *ValidationErrors) Add(field, message string) {
	*ve = append(*ve, NewValidationError(field, message))
}

// DomainError represents a domain logic error.
type DomainError struct {
	code    string
	Message string
}

func (e DomainError) Error() string {
	return e.Message
}

func (e DomainError) Code() string {
	return e.code
}

func (e DomainError) Field() string {
	return ""
}

func NewDomainError(code, message string) DomainError {
	return DomainError{code: code, Message: message}
}

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	resourceType string
	ResourceID   string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.resourceType, e.ResourceID)
}

func (e NotFoundError) Code() string {
	return "NOT_FOUND"
}

func (e NotFoundError) Field() string {
	return "id"
}

func NewNotFoundError(resourceType, resourceID string) NotFoundError {
	return NotFoundError{resourceType: resourceType, ResourceID: resourceID}
}
