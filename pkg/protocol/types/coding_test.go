package types

import "testing"

func TestCodingSystem_IsValid(t *testing.T) {
	tests := []struct {
		system CodingSystem
		want   bool
	}{
		{CodingICD10, true},
		{CodingLOINC, true},
		{CodingSNOMED, true},
		{CodingRxNorm, true},
		{CodingBIOHACK, true},
		{CodingCustom, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.system), func(t *testing.T) {
			if got := tt.system.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCode_Validate_ICD10(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"E11 - diabetes", "E11", false},
		{"E11.9 - with decimal", "E11.9", false},
		{"J06.9 - respiratory", "J06.9", false},
		{"Z23 - vaccination", "Z23", false},
		{"A00.0 - cholera", "A00.0", false},
		{"lowercase valid", "e11.9", false}, // Should normalize
		{"invalid - no letter", "123.4", true},
		{"invalid - wrong format", "ABC", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCode(CodingICD10, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCode(ICD10, %q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestCode_Validate_LOINC(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"8480-6 blood pressure", "8480-6", false},
		{"2339-0 glucose", "2339-0", false},
		{"55284-4 blood pressure panel", "55284-4", false},
		{"1-1 short code", "1-1", false},
		{"invalid - no hyphen", "84806", true},
		{"invalid - letters", "8480-A", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCode(CodingLOINC, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCode(LOINC, %q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestCode_Custom(t *testing.T) {
	// Custom codes should accept any non-empty value
	code, err := NewCode(CodingCustom, "my-custom-code-123")
	if err != nil {
		t.Errorf("Custom code should accept any value: %v", err)
	}
	if code.Value != "my-custom-code-123" {
		t.Errorf("Code value mismatch")
	}
}

func TestCode_Equals(t *testing.T) {
	code1, _ := NewCode(CodingICD10, "E11.9")
	code2, _ := NewCode(CodingICD10, "e11.9") // lowercase
	code3, _ := NewCode(CodingLOINC, "8480-6")

	if !code1.Equals(code2) {
		t.Error("Codes with same system and value should be equal (case-insensitive)")
	}
	if code1.Equals(code3) {
		t.Error("Codes with different systems should not be equal")
	}
}

func TestCode_Validate_SNOMED(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"6 digits", "123456", false},
		{"8 digits", "12345678", false},
		{"18 digits", "123456789012345678", false},
		{"7 digits", "1234567", false},
		{"invalid - too short", "12345", true},
		{"invalid - contains letters", "12345A", true},
		{"invalid - contains dash", "12345-6", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCode(CodingSNOMED, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCode(SNOMED, %q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestCode_Validate_RxNorm(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"5 digits", "12345", false},
		{"7 digits", "1234567", false},
		{"10 digits", "1234567890", false},
		{"1 digit", "1", false},
		{"invalid - contains letters", "12345A", true},
		{"invalid - contains dash", "123-45", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCode(CodingRxNorm, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCode(RxNorm, %q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestCode_Validate_BIOHACK(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"Rapamycin", "BIOHACK:RAPA", false},
		{"NAD", "BIOHACK:NAD", false},
		{"NMN", "BIOHACK:NMN", false},
		{"Peptides", "BIOHACK:PEPT", false},
		{"HRV", "BIOHACK:HRV", false},
		{"VO2Max", "BIOHACK:VO2MAX", false},
		{"lowercase", "biohack:rapa", false}, // Should normalize
		{"invalid - missing prefix", "RAPA", true},
		{"invalid - wrong prefix", "CUSTOM:RAPA", true},
		{"invalid - no code", "BIOHACK:", true},
		{"invalid - contains spaces", "BIOHACK:RA PA", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCode(CodingBIOHACK, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCode(BIOHACK, %q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestCodes_BySystem(t *testing.T) {
	codes := Codes{
		{System: CodingICD10, Value: "E11.9", Display: "Type 2 diabetes"},
		{System: CodingLOINC, Value: "8480-6", Display: "Systolic BP"},
		{System: CodingSNOMED, Value: "123456", Display: "SNOMED code"},
		{System: CodingBIOHACK, Value: "BIOHACK:RAPA", Display: "Rapamycin"},
	}

	icd, found := codes.BySystem(CodingICD10)
	if !found {
		t.Error("Should find ICD-10 code")
	}
	if icd.Value != "E11.9" {
		t.Errorf("Wrong code value: %s", icd.Value)
	}

	snomed, found := codes.BySystem(CodingSNOMED)
	if !found {
		t.Error("Should find SNOMED code")
	}
	if snomed.Value != "123456" {
		t.Errorf("Wrong SNOMED value: %s", snomed.Value)
	}

	biohack, found := codes.BySystem(CodingBIOHACK)
	if !found {
		t.Error("Should find BIOHACK code")
	}
	if biohack.Value != "BIOHACK:RAPA" {
		t.Errorf("Wrong BIOHACK value: %s", biohack.Value)
	}

	_, found = codes.BySystem(CodingCustom)
	if found {
		t.Error("Should not find custom code")
	}
}
