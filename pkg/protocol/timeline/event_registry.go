package timeline

import (
	"sync"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

var (
	// defaultEventTypeRegistry is the default registry for event types.
	defaultEventTypeRegistry types.TypeRegistry[EventType]

	// eventTypeRegistryOnce ensures the registry is initialized only once.
	eventTypeRegistryOnce sync.Once
)

func init() {
	// Initialize default registry on package load
	eventTypeRegistryOnce.Do(func() {
		defaultEventTypeRegistry = types.NewTypeRegistry[EventType]()
		RegisterDefaultEventTypes()
	})
}

// GetEventTypeRegistry returns the default event type registry.
func GetEventTypeRegistry() types.TypeRegistry[EventType] {
	return defaultEventTypeRegistry
}

// RegisterEventType registers a custom event type at runtime.
// This allows extensions without code changes.
func RegisterEventType(eventType EventType, metadata types.TypeMetadata) error {
	return defaultEventTypeRegistry.Register(eventType, metadata)
}

// ValidEventTypes returns all valid event types (backward compatibility).
func ValidEventTypes() []EventType {
	return defaultEventTypeRegistry.ValidTypes()
}

// RegisterDefaultEventTypes registers all built-in event types.
func RegisterDefaultEventTypes() {
	reg := defaultEventTypeRegistry
	types.RegisterBatch(reg, map[EventType]types.TypeMetadata{
		// Medical events
		EventConsultation: {
			Name:        "Consultation",
			Description: "Medical consultation or visit",
			Since:       "0.1.0",
		},
		EventDiagnosis: {
			Name:        "Diagnosis",
			Description: "Medical diagnosis",
			Since:       "0.1.0",
		},
		EventPrescription: {
			Name:        "Prescription",
			Description: "Medication prescription order",
			Since:       "0.1.0",
		},
		EventProcedure: {
			Name:        "Procedure",
			Description: "Medical procedure",
			Since:       "0.1.0",
		},
		EventLabResult: {
			Name:        "Lab Result",
			Description: "Laboratory test result",
			Since:       "0.1.0",
		},
		EventImaging: {
			Name:        "Imaging",
			Description: "Medical imaging study",
			Since:       "0.1.0",
		},
		EventNote: {
			Name:        "Note",
			Description: "Clinical note",
			Since:       "0.1.0",
		},
		EventVaccination: {
			Name:        "Vaccination",
			Description: "Vaccination record",
			Since:       "0.1.0",
		},
		EventAllergy: {
			Name:        "Allergy",
			Description: "Allergy record",
			Since:       "0.1.0",
		},
		EventVisitNote: {
			Name:        "Visit Note",
			Description: "Visit summary note",
			Since:       "0.1.0",
		},
		EventVitalSigns: {
			Name:        "Vital Signs",
			Description: "Vital signs measurement",
			Since:       "0.1.0",
		},
		EventReferral: {
			Name:        "Referral",
			Description: "Medical referral",
			Since:       "0.1.0",
		},
		EventInsuranceClaim: {
			Name:        "Insurance Claim",
			Description: "Insurance claim record",
			Since:       "0.1.0",
		},
		EventTombstone: {
			Name:        "Tombstone",
			Description: "Deleted event marker",
			Since:       "0.1.0",
		},
		EventOther: {
			Name:        "Other",
			Description: "Other event type",
			Since:       "0.1.0",
		},

		// Longevity/Biohacking specific
		EventMedication: {
			Name:        "Medication",
			Description: "Active medication record (current medications being taken)",
			Since:       "0.1.0",
		},
		EventSupplement: {
			Name:        "Supplement",
			Description: "Supplement intake (NAD+, NMN, peptides, etc.)",
			Since:       "0.1.0",
		},
		EventBiometric: {
			Name:        "Biometric",
			Description: "Wearable/biometric data (HRV, sleep, VO2max, DEXA, etc.)",
			Since:       "0.1.0",
		},
		EventIntervention: {
			Name:        "Intervention",
			Description: "Longevity intervention or protocol (rapamycin, fasting, etc.)",
			Since:       "0.1.0",
		},

		// Medical history
		EventFamilyHistory: {
			Name:        "Family History",
			Description: "Family health history record",
			Since:       "0.1.0",
		},
		EventSocialHistory: {
			Name:        "Social History",
			Description: "Social health factors (exercise, diet, lifestyle)",
			Since:       "0.1.0",
		},
		EventDocument: {
			Name:        "Document",
			Description: "General document or file attachment",
			Since:       "0.1.0",
		},

		// Aliases
		EventVital: {
			Name:        "Vital",
			Description: "Vital signs measurement (alias for vital_signs)",
			Since:       "0.1.0",
		},
	})
}
