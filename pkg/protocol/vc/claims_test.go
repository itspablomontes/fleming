package vc

import (
	"testing"
)

func TestBloodworkRangeClaim_Validate(t *testing.T) {
	tests := []struct {
		name    string
		claim   BloodworkRangeClaim
		wantErr bool
	}{
		{
			name: "valid claim",
			claim: BloodworkRangeClaim{
				Marker:       "718-7",
				RangeMin:     13.5,
				RangeMax:     17.5,
				WindowMonths: 6,
				AllInRange:   true,
				SampleCount:  5,
			},
			wantErr: false,
		},
		{
			name: "missing marker",
			claim: BloodworkRangeClaim{
				RangeMin:     13.5,
				RangeMax:     17.5,
				WindowMonths: 6,
			},
			wantErr: true,
		},
		{
			name: "rangeMax < rangeMin",
			claim: BloodworkRangeClaim{
				Marker:       "718-7",
				RangeMin:     17.5,
				RangeMax:     13.5,
				WindowMonths: 6,
			},
			wantErr: true,
		},
		{
			name: "rangeMax == rangeMin",
			claim: BloodworkRangeClaim{
				Marker:       "718-7",
				RangeMin:     15.0,
				RangeMax:     15.0,
				WindowMonths: 6,
			},
			wantErr: false,
		},
		{
			name: "windowMonths <= 0",
			claim: BloodworkRangeClaim{
				Marker:       "718-7",
				RangeMin:     13.5,
				RangeMax:     17.5,
				WindowMonths: 0,
			},
			wantErr: true,
		},
		{
			name: "negative sampleCount",
			claim: BloodworkRangeClaim{
				Marker:       "718-7",
				RangeMin:     13.5,
				RangeMax:     17.5,
				WindowMonths: 6,
				SampleCount:  -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.claim.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBloodworkRangeClaim_ToMap(t *testing.T) {
	claim := BloodworkRangeClaim{
		Marker:       "718-7",
		RangeMin:     13.5,
		RangeMax:     17.5,
		WindowMonths: 6,
		AllInRange:   true,
		SampleCount:  5,
	}

	m := claim.ToMap()

	if m["marker"] != "718-7" {
		t.Errorf("ToMap() marker = %v, want 718-7", m["marker"])
	}
	if m["rangeMin"] != 13.5 {
		t.Errorf("ToMap() rangeMin = %v, want 13.5", m["rangeMin"])
	}
	if m["rangeMax"] != 17.5 {
		t.Errorf("ToMap() rangeMax = %v, want 17.5", m["rangeMax"])
	}
	if m["windowMonths"] != 6 {
		t.Errorf("ToMap() windowMonths = %v, want 6", m["windowMonths"])
	}
	if m["allInRange"] != true {
		t.Errorf("ToMap() allInRange = %v, want true", m["allInRange"])
	}
	if m["sampleCount"] != 5 {
		t.Errorf("ToMap() sampleCount = %v, want 5", m["sampleCount"])
	}
}

func TestProtocolAdherenceClaim_Validate(t *testing.T) {
	tests := []struct {
		name    string
		claim   ProtocolAdherenceClaim
		wantErr bool
	}{
		{
			name: "valid claim",
			claim: ProtocolAdherenceClaim{
				Intervention:        "BIOHACK:RAPA",
				MinDurationMonths:   6,
				ActualDurationMet:   true,
				ConfirmedByProvider: false,
			},
			wantErr: false,
		},
		{
			name: "missing intervention",
			claim: ProtocolAdherenceClaim{
				MinDurationMonths: 6,
			},
			wantErr: true,
		},
		{
			name: "minDurationMonths <= 0",
			claim: ProtocolAdherenceClaim{
				Intervention:      "BIOHACK:RAPA",
				MinDurationMonths: 0,
			},
			wantErr: true,
		},
		{
			name: "negative minDurationMonths",
			claim: ProtocolAdherenceClaim{
				Intervention:      "BIOHACK:RAPA",
				MinDurationMonths: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.claim.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProtocolAdherenceClaim_ToMap(t *testing.T) {
	claim := ProtocolAdherenceClaim{
		Intervention:        "BIOHACK:RAPA",
		MinDurationMonths:   6,
		ActualDurationMet:   true,
		ConfirmedByProvider: true,
	}

	m := claim.ToMap()

	if m["intervention"] != "BIOHACK:RAPA" {
		t.Errorf("ToMap() intervention = %v, want BIOHACK:RAPA", m["intervention"])
	}
	if m["minDurationMonths"] != 6 {
		t.Errorf("ToMap() minDurationMonths = %v, want 6", m["minDurationMonths"])
	}
	if m["actualDurationMet"] != true {
		t.Errorf("ToMap() actualDurationMet = %v, want true", m["actualDurationMet"])
	}
	if m["confirmedByProvider"] != true {
		t.Errorf("ToMap() confirmedByProvider = %v, want true", m["confirmedByProvider"])
	}
}

func TestBiometricPercentileClaim_Validate(t *testing.T) {
	tests := []struct {
		name    string
		claim   BiometricPercentileClaim
		wantErr bool
	}{
		{
			name: "valid claim",
			claim: BiometricPercentileClaim{
				Metric:          "BIOHACK:HRV",
				Percentile:      80,
				AboveThreshold:  true,
			},
			wantErr: false,
		},
		{
			name: "missing metric",
			claim: BiometricPercentileClaim{
				Percentile: 80,
			},
			wantErr: true,
		},
		{
			name: "percentile < 0",
			claim: BiometricPercentileClaim{
				Metric:    "BIOHACK:HRV",
				Percentile: -1,
			},
			wantErr: true,
		},
		{
			name: "percentile > 100",
			claim: BiometricPercentileClaim{
				Metric:    "BIOHACK:HRV",
				Percentile: 101,
			},
			wantErr: true,
		},
		{
			name: "percentile = 0",
			claim: BiometricPercentileClaim{
				Metric:    "BIOHACK:HRV",
				Percentile: 0,
			},
			wantErr: false,
		},
		{
			name: "percentile = 100",
			claim: BiometricPercentileClaim{
				Metric:    "BIOHACK:HRV",
				Percentile: 100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.claim.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgeOverClaim_Validate(t *testing.T) {
	tests := []struct {
		name    string
		claim   AgeOverClaim
		wantErr bool
	}{
		{
			name: "valid claim",
			claim: AgeOverClaim{
				AgeThreshold: 18,
				IsOver:       true,
			},
			wantErr: false,
		},
		{
			name: "ageThreshold <= 0",
			claim: AgeOverClaim{
				AgeThreshold: 0,
			},
			wantErr: true,
		},
		{
			name: "negative ageThreshold",
			claim: AgeOverClaim{
				AgeThreshold: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.claim.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseBloodworkRangeClaim(t *testing.T) {
	tests := []struct {
		name    string
		claims  map[string]any
		wantErr bool
	}{
		{
			name: "valid claim",
			claims: map[string]any{
				"marker":       "718-7",
				"rangeMin":     13.5,
				"rangeMax":     17.5,
				"windowMonths": 6,
				"allInRange":   true,
				"sampleCount":  5,
			},
			wantErr: false,
		},
		{
			name: "missing marker",
			claims: map[string]any{
				"rangeMin": 13.5,
			},
			wantErr: true,
		},
		{
			name: "invalid marker type",
			claims: map[string]any{
				"marker": 123,
			},
			wantErr: true,
		},
		{
			name: "windowMonths as int",
			claims: map[string]any{
				"marker":       "718-7",
				"rangeMin":     13.5,
				"rangeMax":     17.5,
				"windowMonths": int(6),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claim, err := ParseBloodworkRangeClaim(tt.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBloodworkRangeClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claim == nil {
				t.Error("ParseBloodworkRangeClaim() returned nil for valid claim")
			}
		})
	}
}

func TestParseProtocolAdherenceClaim(t *testing.T) {
	tests := []struct {
		name    string
		claims  map[string]any
		wantErr bool
	}{
		{
			name: "valid claim",
			claims: map[string]any{
				"intervention":        "BIOHACK:RAPA",
				"minDurationMonths":   6,
				"actualDurationMet":   true,
				"confirmedByProvider": false,
			},
			wantErr: false,
		},
		{
			name: "missing intervention",
			claims: map[string]any{
				"minDurationMonths": 6,
			},
			wantErr: true,
		},
		{
			name: "minDurationMonths as int",
			claims: map[string]any{
				"intervention":      "BIOHACK:RAPA",
				"minDurationMonths": int(6),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claim, err := ParseProtocolAdherenceClaim(tt.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProtocolAdherenceClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claim == nil {
				t.Error("ParseProtocolAdherenceClaim() returned nil for valid claim")
			}
		})
	}
}
