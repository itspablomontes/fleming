package types

import "testing"

func TestCodingSystem_IsValid(t *testing.T) {
	tests := []struct {
		system CodingSystem
		want   bool
	}{
		{CodingICD10, true},
		{CodingLOINC, true},
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

func TestCodes_BySystem(t *testing.T) {
	codes := Codes{
		{System: CodingICD10, Value: "E11.9", Display: "Type 2 diabetes"},
		{System: CodingLOINC, Value: "8480-6", Display: "Systolic BP"},
	}

	icd, found := codes.BySystem(CodingICD10)
	if !found {
		t.Error("Should find ICD-10 code")
	}
	if icd.Value != "E11.9" {
		t.Errorf("Wrong code value: %s", icd.Value)
	}

	_, found = codes.BySystem(CodingCustom)
	if found {
		t.Error("Should not find custom code")
	}
}
