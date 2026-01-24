package types

import "errors"

var (
	ErrInvalidID = errors.New("invalid ID: cannot be empty")

	ErrInvalidAddress = errors.New("invalid address: must be 0x followed by 40 hex characters")

	ErrValidationFailed = errors.New("validation failed")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{Field: field, Message: message}
}

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
