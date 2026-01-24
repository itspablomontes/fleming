package types

import (
	"fmt"
	"regexp"
	"strings"
)

type CodingSystem string

const (
	CodingICD10  CodingSystem = "ICD-10"
	CodingLOINC  CodingSystem = "LOINC"
	CodingCustom CodingSystem = "custom"
)

func ValidCodingSystems() []CodingSystem {
	return []CodingSystem{CodingICD10, CodingLOINC, CodingCustom}
}

func (cs CodingSystem) IsValid() bool {
	switch cs {
	case CodingICD10, CodingLOINC, CodingCustom:
		return true
	}
	return false
}

type Code struct {
	System CodingSystem `json:"system"`

	Value string `json:"code"`

	Display string `json:"display,omitempty"`
}

func NewCode(system CodingSystem, value string) (Code, error) {
	c := Code{System: system, Value: value}
	if err := c.Validate(); err != nil {
		return Code{}, err
	}
	return c, nil
}

func NewCodeWithDisplay(system CodingSystem, value, display string) (Code, error) {
	c := Code{System: system, Value: value, Display: display}
	if err := c.Validate(); err != nil {
		return Code{}, err
	}
	return c, nil
}

var (
	icd10Regex = regexp.MustCompile(`^[A-Z][0-9]{2}(\.[0-9A-Z]{1,4})?$`)
	loincRegex = regexp.MustCompile(`^[0-9]{1,5}-[0-9]$`)
)

func (c Code) Validate() error {
	if c.Value == "" {
		return NewValidationError("code", "value cannot be empty")
	}

	value := strings.TrimSpace(c.Value)

	switch c.System {
	case CodingICD10:
		if !icd10Regex.MatchString(strings.ToUpper(value)) {
			return NewValidationError("code", fmt.Sprintf("invalid ICD-10 format: %s", value))
		}
	case CodingLOINC:
		if !loincRegex.MatchString(value) {
			return NewValidationError("code", fmt.Sprintf("invalid LOINC format: %s", value))
		}
	case CodingCustom:
	default:
		return NewValidationError("system", fmt.Sprintf("unsupported coding system: %s", c.System))
	}

	return nil
}

func (c Code) IsEmpty() bool {
	return c.Value == ""
}

func (c Code) String() string {
	if c.Display != "" {
		return fmt.Sprintf("%s|%s (%s)", c.System, c.Value, c.Display)
	}
	return fmt.Sprintf("%s|%s", c.System, c.Value)
}

func (c Code) Equals(other Code) bool {
	return c.System == other.System && strings.EqualFold(c.Value, other.Value)
}

type Codes []Code

func (codes Codes) HasSystem(system CodingSystem) bool {
	for _, c := range codes {
		if c.System == system {
			return true
		}
	}
	return false
}

func (codes Codes) BySystem(system CodingSystem) (Code, bool) {
	for _, c := range codes {
		if c.System == system {
			return c, true
		}
	}
	return Code{}, false
}
