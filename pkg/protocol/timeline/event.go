// Package timeline provides timeline event primitives for the Protocol layer.
// These types represent medical events and their relationships without persistence concerns.
package timeline

import (
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type EventType string

const (
	// Medical events
	EventConsultation   EventType = "consultation"
	EventDiagnosis      EventType = "diagnosis"
	EventPrescription   EventType = "prescription"
	EventProcedure      EventType = "procedure"
	EventLabResult      EventType = "lab_result"
	EventImaging        EventType = "imaging"
	EventNote           EventType = "note"
	EventVaccination    EventType = "vaccination"
	EventAllergy        EventType = "allergy"
	EventVisitNote      EventType = "visit_note"
	EventVitalSigns     EventType = "vital_signs"
	EventReferral       EventType = "referral"
	EventInsuranceClaim EventType = "insurance_claim"
	EventTombstone      EventType = "tombstone"
	EventOther          EventType = "other"

	// Longevity/Biohacking specific
	EventMedication   EventType = "medication"   // Active medication (vs prescription order)
	EventSupplement   EventType = "supplement"   // Supplements (NAD+, NMN, etc.)
	EventBiometric    EventType = "biometric"    // Wearable data (HRV, sleep, etc.)
	EventIntervention EventType = "intervention" // Longevity interventions (rapamycin protocol, etc.)

	// Medical history
	EventFamilyHistory EventType = "family_history" // Family health history
	EventSocialHistory EventType = "social_history" // Social health factors
	EventDocument      EventType = "document"       // General documents

	// Alias for backward compatibility
	EventVital EventType = "vital" // Alias for vital_signs
)

func (et EventType) IsValid() bool {
	return GetEventTypeRegistry().IsValid(et)
}

type Event struct {
	ID types.ID `json:"id"`

	PatientID types.WalletAddress `json:"patientId"`

	Type EventType `json:"type"`

	Title string `json:"title"`

	Description string `json:"description,omitempty"`

	Provider string `json:"provider,omitempty"`

	Codes types.Codes `json:"codes,omitempty"`

	Timestamp time.Time `json:"timestamp"`

	Metadata types.Metadata `json:"metadata,omitempty"`

	SchemaVersion string `json:"schemaVersion,omitempty"` // Protocol schema version (e.g., "timeline.v1")

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *Event) Validate() error {
	var errs types.ValidationErrors

	if e.PatientID.IsEmpty() {
		errs.Add("patientId", "patient ID is required")
	}

	if !e.Type.IsValid() {
		errs.Add("type", "invalid event type")
	}

	if e.Title == "" {
		errs.Add("title", "title is required")
	}

	if e.Timestamp.IsZero() {
		errs.Add("timestamp", "timestamp is required")
	}

	for i, code := range e.Codes {
		if err := code.Validate(); err != nil {
			errs.Add("codes", err.Error()+" (index: "+string(rune('0'+i))+")")
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func (e *Event) HasCode(system types.CodingSystem) bool {
	return e.Codes.HasSystem(system)
}

func (e *Event) GetCode(system types.CodingSystem) (types.Code, bool) {
	return e.Codes.BySystem(system)
}

func (e *Event) AddCode(code types.Code) error {
	if err := code.Validate(); err != nil {
		return err
	}
	e.Codes = append(e.Codes, code)
	return nil
}
