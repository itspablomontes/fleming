// Package timeline provides timeline event primitives for the Protocol layer.
// These types represent medical events and their relationships without persistence concerns.
package timeline

import (
	"slices"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type EventType string

const (
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
)

func ValidEventTypes() []EventType {
	return []EventType{
		EventConsultation, EventDiagnosis, EventPrescription,
		EventProcedure, EventLabResult, EventImaging,
		EventNote, EventVaccination, EventAllergy,
		EventVisitNote, EventVitalSigns, EventReferral,
		EventInsuranceClaim, EventTombstone, EventOther,
	}
}

func (et EventType) IsValid() bool {
	return slices.Contains(ValidEventTypes(), et)
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

	CreatedAt time.Time `json:"createdAt,omitempty"`

	UpdatedAt time.Time `json:"updatedAt,omitempty"`
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
