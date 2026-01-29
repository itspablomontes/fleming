package timeline

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestEventType_IsValid(t *testing.T) {
	tests := []struct {
		et   EventType
		want bool
	}{
		// Medical events
		{EventConsultation, true},
		{EventDiagnosis, true},
		{EventPrescription, true},
		{EventProcedure, true},
		{EventLabResult, true},
		{EventImaging, true},
		{EventNote, true},
		{EventVaccination, true},
		{EventAllergy, true},
		{EventVisitNote, true},
		{EventVitalSigns, true},
		{EventReferral, true},
		{EventInsuranceClaim, true},
		{EventTombstone, true},
		{EventOther, true},
		// Longevity/Biohacking
		{EventMedication, true},
		{EventSupplement, true},
		{EventBiometric, true},
		{EventIntervention, true},
		// Medical history
		{EventFamilyHistory, true},
		{EventSocialHistory, true},
		{EventDocument, true},
		// Alias
		{EventVital, true},
		// Invalid
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.et), func(t *testing.T) {
			if got := tt.et.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_Validate(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")

	tests := []struct {
		name    string
		event   Event
		wantErr bool
	}{
		{
			name: "valid event",
			event: Event{
				PatientID: validAddr,
				Type:      EventLabResult,
				Title:     "Blood Test Results",
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid with codes",
			event: Event{
				PatientID: validAddr,
				Type:      EventDiagnosis,
				Title:     "Type 2 Diabetes",
				Timestamp: time.Now(),
				Codes: types.Codes{
					{System: types.CodingICD10, Value: "E11.9"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing patient ID",
			event: Event{
				Type:      EventLabResult,
				Title:     "Blood Test",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid event type",
			event: Event{
				PatientID: validAddr,
				Type:      "invalid",
				Title:     "Test",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing title",
			event: Event{
				PatientID: validAddr,
				Type:      EventLabResult,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing timestamp",
			event: Event{
				PatientID: validAddr,
				Type:      EventLabResult,
				Title:     "Blood Test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEvent_AddCode(t *testing.T) {
	event := Event{}

	code, _ := types.NewCodeWithDisplay(types.CodingICD10, "E11.9", "Type 2 diabetes")
	if err := event.AddCode(code); err != nil {
		t.Errorf("AddCode() error = %v", err)
	}

	if len(event.Codes) != 1 {
		t.Errorf("Expected 1 code, got %d", len(event.Codes))
	}

	if !event.HasCode(types.CodingICD10) {
		t.Error("Expected HasCode(ICD10) to be true")
	}

	retrieved, found := event.GetCode(types.CodingICD10)
	if !found {
		t.Error("Expected to find ICD10 code")
	}
	if retrieved.Value != "E11.9" {
		t.Errorf("Expected code E11.9, got %s", retrieved.Value)
	}
}
