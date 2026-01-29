package vc

import (
	"fmt"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// ClaimValidator is the interface for validating claims against timeline events.
// Implementations should verify that the claim criteria are met by the source events.
type ClaimValidator interface {
	// Validate checks if the claim criteria are satisfied by the given event IDs.
	// The implementation should query the timeline to verify the claim.
	Validate(eventIDs []types.ID) error

	// ClaimType returns the type of claim this validator handles.
	ClaimType() ClaimType
}

// BloodworkRangeClaim proves that biomarkers are within optimal ranges
// over a specified time window.
type BloodworkRangeClaim struct {
	// Marker is the LOINC code for the biomarker
	Marker string `json:"marker"`

	// RangeMin is the minimum acceptable value
	RangeMin float64 `json:"rangeMin"`

	// RangeMax is the maximum acceptable value
	RangeMax float64 `json:"rangeMax"`

	// WindowMonths is the time window in months to check
	WindowMonths int `json:"windowMonths"`

	// AllInRange indicates if ALL values in the window are within range
	AllInRange bool `json:"allInRange"`

	// SampleCount is the number of samples checked
	SampleCount int `json:"sampleCount"`
}

// Validate validates the BloodworkRangeClaim structure.
func (c *BloodworkRangeClaim) Validate() error {
	var errs types.ValidationErrors

	if c.Marker == "" {
		errs.Add("marker", "marker (LOINC code) is required")
	}

	if c.RangeMax < c.RangeMin {
		errs.Add("rangeMax", "rangeMax must be >= rangeMin")
	}

	if c.WindowMonths <= 0 {
		errs.Add("windowMonths", "windowMonths must be positive")
	}

	if c.SampleCount < 0 {
		errs.Add("sampleCount", "sampleCount cannot be negative")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ToMap converts the claim to a map for inclusion in credentials.
func (c *BloodworkRangeClaim) ToMap() map[string]any {
	return map[string]any{
		"marker":       c.Marker,
		"rangeMin":     c.RangeMin,
		"rangeMax":     c.RangeMax,
		"windowMonths": c.WindowMonths,
		"allInRange":   c.AllInRange,
		"sampleCount":  c.SampleCount,
	}
}

// ProtocolAdherenceClaim proves adherence to an intervention protocol
// for a minimum duration.
type ProtocolAdherenceClaim struct {
	// Intervention is the intervention code (e.g., BIOHACK:RAPA)
	Intervention string `json:"intervention"`

	// MinDurationMonths is the minimum required duration
	MinDurationMonths int `json:"minDurationMonths"`

	// ActualDurationMet indicates if the actual duration meets/exceeds minimum
	ActualDurationMet bool `json:"actualDurationMet"`

	// ConfirmedByProvider indicates if a provider confirmed the protocol
	ConfirmedByProvider bool `json:"confirmedByProvider,omitempty"`
}

// Validate validates the ProtocolAdherenceClaim structure.
func (c *ProtocolAdherenceClaim) Validate() error {
	var errs types.ValidationErrors

	if c.Intervention == "" {
		errs.Add("intervention", "intervention code is required")
	}

	if c.MinDurationMonths <= 0 {
		errs.Add("minDurationMonths", "minDurationMonths must be positive")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ToMap converts the claim to a map for inclusion in credentials.
func (c *ProtocolAdherenceClaim) ToMap() map[string]any {
	return map[string]any{
		"intervention":        c.Intervention,
		"minDurationMonths":   c.MinDurationMonths,
		"actualDurationMet":   c.ActualDurationMet,
		"confirmedByProvider": c.ConfirmedByProvider,
	}
}

// BiometricPercentileClaim proves biometric values rank above a specified percentile.
type BiometricPercentileClaim struct {
	// Metric is the biometric metric (e.g., BIOHACK:HRV, BIOHACK:VO2MAX)
	Metric string `json:"metric"`

	// Percentile is the minimum percentile threshold
	Percentile int `json:"percentile"`

	// AboveThreshold indicates if the subject is above the percentile
	AboveThreshold bool `json:"aboveThreshold"`

	// ReferencePopulation describes the comparison population
	ReferencePopulation string `json:"referencePopulation,omitempty"`
}

// Validate validates the BiometricPercentileClaim structure.
func (c *BiometricPercentileClaim) Validate() error {
	var errs types.ValidationErrors

	if c.Metric == "" {
		errs.Add("metric", "metric is required")
	}

	if c.Percentile < 0 || c.Percentile > 100 {
		errs.Add("percentile", "percentile must be between 0 and 100")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ToMap converts the claim to a map for inclusion in credentials.
func (c *BiometricPercentileClaim) ToMap() map[string]any {
	return map[string]any{
		"metric":              c.Metric,
		"percentile":          c.Percentile,
		"aboveThreshold":      c.AboveThreshold,
		"referencePopulation": c.ReferencePopulation,
	}
}

// AgeOverClaim proves the subject is over a specified age without revealing exact age.
type AgeOverClaim struct {
	// AgeThreshold is the minimum age threshold
	AgeThreshold int `json:"ageThreshold"`

	// IsOver indicates if the subject is over the threshold
	IsOver bool `json:"isOver"`
}

// Validate validates the AgeOverClaim structure.
func (c *AgeOverClaim) Validate() error {
	if c.AgeThreshold <= 0 {
		return types.NewValidationError("ageThreshold", "ageThreshold must be positive")
	}
	return nil
}

// ToMap converts the claim to a map for inclusion in credentials.
func (c *AgeOverClaim) ToMap() map[string]any {
	return map[string]any{
		"ageThreshold": c.AgeThreshold,
		"isOver":       c.IsOver,
	}
}

// ParseBloodworkRangeClaim parses a BloodworkRangeClaim from a claims map.
func ParseBloodworkRangeClaim(claims map[string]any) (*BloodworkRangeClaim, error) {
	c := &BloodworkRangeClaim{}

	if marker, ok := claims["marker"].(string); ok {
		c.Marker = marker
	} else {
		return nil, fmt.Errorf("missing or invalid marker")
	}

	if rangeMin, ok := claims["rangeMin"].(float64); ok {
		c.RangeMin = rangeMin
	}

	if rangeMax, ok := claims["rangeMax"].(float64); ok {
		c.RangeMax = rangeMax
	}

	if windowMonths, ok := claims["windowMonths"].(float64); ok {
		c.WindowMonths = int(windowMonths)
	} else if windowMonths, ok := claims["windowMonths"].(int); ok {
		c.WindowMonths = windowMonths
	}

	if allInRange, ok := claims["allInRange"].(bool); ok {
		c.AllInRange = allInRange
	}

	if sampleCount, ok := claims["sampleCount"].(float64); ok {
		c.SampleCount = int(sampleCount)
	} else if sampleCount, ok := claims["sampleCount"].(int); ok {
		c.SampleCount = sampleCount
	}

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid claim: %w", err)
	}

	return c, nil
}

// ParseProtocolAdherenceClaim parses a ProtocolAdherenceClaim from a claims map.
func ParseProtocolAdherenceClaim(claims map[string]any) (*ProtocolAdherenceClaim, error) {
	c := &ProtocolAdherenceClaim{}

	if intervention, ok := claims["intervention"].(string); ok {
		c.Intervention = intervention
	} else {
		return nil, fmt.Errorf("missing or invalid intervention")
	}

	if minDuration, ok := claims["minDurationMonths"].(float64); ok {
		c.MinDurationMonths = int(minDuration)
	} else if minDuration, ok := claims["minDurationMonths"].(int); ok {
		c.MinDurationMonths = minDuration
	}

	if actualMet, ok := claims["actualDurationMet"].(bool); ok {
		c.ActualDurationMet = actualMet
	}

	if confirmed, ok := claims["confirmedByProvider"].(bool); ok {
		c.ConfirmedByProvider = confirmed
	}

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid claim: %w", err)
	}

	return c, nil
}
