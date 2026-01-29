package types

import (
	"fmt"
	"regexp"
	"strings"
)

type CodingSystem string

const (
	// Standard medical coding systems
	CodingICD10  CodingSystem = "ICD-10"  // International Classification of Diseases
	CodingLOINC  CodingSystem = "LOINC"   // Logical Observation Identifiers Names and Codes
	CodingSNOMED CodingSystem = "SNOMED"  // SNOMED CT medical terminology
	CodingRxNorm CodingSystem = "RxNorm"  // Medication terminology

	// Longevity/Biohacking namespace
	CodingBIOHACK CodingSystem = "BIOHACK" // Custom namespace for longevity interventions

	// Fallback
	CodingCustom CodingSystem = "custom" // Custom/proprietary codes
)

// BIOHACK namespace codes for longevity interventions
const (
	// Medications/Protocols
	BiohackRapamycin = "BIOHACK:RAPA"      // Rapamycin/Sirolimus protocol
	BiohackMetformin = "BIOHACK:METF"      // Metformin
	BiohackNAD       = "BIOHACK:NAD"       // NAD+ precursors
	BiohackNMN       = "BIOHACK:NMN"       // Nicotinamide mononucleotide
	BiohackNR        = "BIOHACK:NR"        // Nicotinamide riboside
	BiohackPeptides  = "BIOHACK:PEPT"      // Peptides (BPC-157, etc.)
	BiohackResveratrol = "BIOHACK:RESV"    // Resveratrol
	BiohackBerberine = "BIOHACK:BERB"      // Berberine

	// Biometrics/Measurements
	BiohackHRV       = "BIOHACK:HRV"       // Heart Rate Variability
	BiohackVO2Max    = "BIOHACK:VO2MAX"    // VO2 Max
	BiohackDEXA      = "BIOHACK:DEXA"      // DEXA body composition
	BiohackGrip      = "BIOHACK:GRIP"      // Grip strength
	BiohackCGM       = "BIOHACK:CGM"       // Continuous glucose monitoring

	// Interventions
	BiohackFasting   = "BIOHACK:FAST"      // Fasting protocols
	BiohackColdExposure = "BIOHACK:COLD"   // Cold exposure/cryotherapy
	BiohackHeatExposure = "BIOHACK:HEAT"   // Sauna/heat therapy
	BiohackSleep     = "BIOHACK:SLEEP"     // Sleep optimization
)

func ValidCodingSystems() []CodingSystem {
	return []CodingSystem{CodingICD10, CodingLOINC, CodingSNOMED, CodingRxNorm, CodingBIOHACK, CodingCustom}
}

func (cs CodingSystem) IsValid() bool {
	switch cs {
	case CodingICD10, CodingLOINC, CodingSNOMED, CodingRxNorm, CodingBIOHACK, CodingCustom:
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
	// ICD-10: Letter followed by 2 digits, optional decimal with 1-4 alphanumeric chars
	icd10Regex = regexp.MustCompile(`^[A-Z][0-9]{2}(\.[0-9A-Z]{1,4})?$`)
	// LOINC: 1-5 digits, hyphen, check digit
	loincRegex = regexp.MustCompile(`^[0-9]{1,5}-[0-9]$`)
	// SNOMED CT: 6-18 digits
	snomedRegex = regexp.MustCompile(`^[0-9]{6,18}$`)
	// RxNorm: Concept Unique Identifier (CUI) - typically 5-7 digits
	rxnormRegex = regexp.MustCompile(`^[0-9]{1,10}$`)
	// BIOHACK: BIOHACK:CODE format
	biohackRegex = regexp.MustCompile(`^BIOHACK:[A-Z0-9_]+$`)
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
	case CodingSNOMED:
		if !snomedRegex.MatchString(value) {
			return NewValidationError("code", fmt.Sprintf("invalid SNOMED CT format: %s", value))
		}
	case CodingRxNorm:
		if !rxnormRegex.MatchString(value) {
			return NewValidationError("code", fmt.Sprintf("invalid RxNorm format: %s", value))
		}
	case CodingBIOHACK:
		if !biohackRegex.MatchString(strings.ToUpper(value)) {
			return NewValidationError("code", fmt.Sprintf("invalid BIOHACK format: %s (expected BIOHACK:CODE)", value))
		}
	case CodingCustom:
		// Custom codes have no format restrictions
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
