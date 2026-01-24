package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/apps/backend/internal/timeline"
	prototline "github.com/itspablomontes/fleming/pkg/protocol/timeline"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://fleming:fleming@localhost:5432/fleming?sslmode=disable"
	}

	overrideAddress := os.Getenv("DEV_OVERRIDE_WALLET_ADDRESS")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Connected to database")

	// Clean up existing data
	log.Println("Cleaning up existing timeline data...")
	if err := db.Exec("DELETE FROM event_edges").Error; err != nil {
		log.Printf("Warning: failed to clean event_edges: %v", err)
	}
	if err := db.Exec("DELETE FROM event_files").Error; err != nil {
		log.Printf("Warning: failed to clean event_files: %v", err)
	}
	if err := db.Exec("DELETE FROM timeline_events").Error; err != nil {
		log.Printf("Warning: failed to clean timeline_events: %v", err)
	}

	// Mock Data Constants
	patientId := "0x742d35Cc6634C0532925a3b844Bc9e7595f"
	if overrideAddress != "" {
		patientId = overrideAddress
		log.Printf("Using override patient ID: %s", patientId)
	}

	// Helper to parse time
	parseTime := func(layout, value string) time.Time {
		t, err := time.Parse(layout, value)
		if err != nil {
			log.Fatalf("failed to parse time %s: %v", value, err)
		}
		return t
	}

	// ID Mapping: Old Mock ID -> New UUID
	idMap := make(map[string]string)

	// Define raw events first with their mock IDs
	type MockEvent struct {
		MockID string
		Event  timeline.TimelineEvent
	}

	rawEvents := []MockEvent{
		{
			MockID: "evt-early-checkup",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventConsultation,
				Title:       "Initial Baseline Checkup",
				Description: "Routine checkup before any chronic symptoms. All labs normal.",
				Provider:    "Dr. Sarah Chen",
				Timestamp:   parseTime(time.RFC3339, "2023-03-10T10:00:00Z"),
				IsEncrypted: false,
				Metadata:    common.JSONMap{"facility": "Main Street Medical Center"},
				CreatedAt:   parseTime(time.RFC3339, "2023-03-10T10:30:00Z"),
				UpdatedAt:   parseTime(time.RFC3339, "2023-03-10T10:30:00Z"),
			},
		},
		{
			MockID: "evt-consultation-1",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventConsultation,
				Title:       "Annual Checkup - Initial Symptoms",
				Description: "Routine physical examination. Patient reports fatigue and increased thirst. Ordered CMP.",
				Provider:    "Dr. Sarah Chen",
				Timestamp:   parseTime(time.RFC3339, "2024-01-15T09:00:00Z"),
				IsEncrypted: false,
				Metadata: common.JSONMap{
					"facility":  "Main Street Medical Center",
					"visitType": "preventive",
				},
				CreatedAt: parseTime(time.RFC3339, "2024-01-15T09:30:00Z"),
				UpdatedAt: parseTime(time.RFC3339, "2024-01-15T09:30:00Z"),
			},
		},
		{
			MockID: "evt-lab-1",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventLabResult,
				Title:       "Comprehensive Metabolic Panel",
				Description: "Fasting blood glucose: 186 mg/dL (High). HbA1c: 7.8% (High).",
				Provider:    "LabCorp",
				Timestamp:   parseTime(time.RFC3339, "2024-01-16T14:00:00Z"),
				IsEncrypted: false,
				Metadata: common.JSONMap{
					"glucoseLevel": 186,
					"hba1c":        7.8,
				},
				CreatedAt: parseTime(time.RFC3339, "2024-01-16T15:00:00Z"),
				UpdatedAt: parseTime(time.RFC3339, "2024-01-16T15:00:00Z"),
			},
		},
		{
			MockID: "evt-diagnosis-1",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventDiagnosis,
				Title:       "Type 2 Diabetes Mellitus",
				Description: "Diagnosis confirmed based on glucose and HbA1c levels. ICD-10: E11.9",
				Provider:    "Dr. Sarah Chen",
				Timestamp:   parseTime(time.RFC3339, "2024-01-20T10:00:00Z"),
				IsEncrypted: false,
				Metadata: common.JSONMap{
					"icd10":    "E11.9",
					"severity": "moderate",
				},
				CreatedAt: parseTime(time.RFC3339, "2024-01-20T10:30:00Z"),
				UpdatedAt: parseTime(time.RFC3339, "2024-01-20T10:30:00Z"),
			},
		},
		{
			MockID: "evt-prescription-1",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventPrescription,
				Title:       "Metformin 500mg",
				Description: "Start daily dosage with dinner.",
				Provider:    "Dr. Sarah Chen",
				Timestamp:   parseTime(time.RFC3339, "2024-01-20T10:15:00Z"),
				IsEncrypted: false,
				Metadata: common.JSONMap{
					"medication": "Metformin",
					"dosage":     "500mg",
				},
				CreatedAt: parseTime(time.RFC3339, "2024-01-20T10:30:00Z"),
				UpdatedAt: parseTime(time.RFC3339, "2024-01-20T10:30:00Z"),
			},
		},
		{
			MockID: "evt-consultation-2",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventConsultation,
				Title:       "Endocrinology Visit",
				Description: "Discussed blood glucose monitoring and dietary modifications.",
				Provider:    "Dr. Michael Park",
				Timestamp:   parseTime(time.RFC3339, "2024-02-05T14:00:00Z"),
				IsEncrypted: false,
				Metadata:    common.JSONMap{"specialty": "Endocrinology"},
				CreatedAt:   parseTime(time.RFC3339, "2024-02-05T15:00:00Z"),
				UpdatedAt:   parseTime(time.RFC3339, "2024-02-05T15:00:00Z"),
			},
		},
		{
			MockID: "evt-lab-2",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventLabResult,
				Title:       "6-Month HbA1c Review",
				Description: "HbA1c: 6.8% (Significant improvement).",
				Provider:    "LabCorp",
				Timestamp:   parseTime(time.RFC3339, "2024-08-10T09:00:00Z"),
				IsEncrypted: false,
				Metadata:    common.JSONMap{"hba1c": 6.8},
				CreatedAt:   parseTime(time.RFC3339, "2024-08-10T10:00:00Z"),
				UpdatedAt:   parseTime(time.RFC3339, "2024-08-10T10:00:00Z"),
			},
		},
		{
			MockID: "evt-procedure-1",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventProcedure,
				Title:       "Ophthalmology Screening",
				Description: "Diabetic retinopathy screening. No issues detected.",
				Provider:    "Dr. Elena Rossi",
				Timestamp:   parseTime(time.RFC3339, "2025-01-22T11:00:00Z"),
				IsEncrypted: false,
				Metadata:    common.JSONMap{"results": "clear"},
				CreatedAt:   parseTime(time.RFC3339, "2025-01-22T12:00:00Z"),
				UpdatedAt:   parseTime(time.RFC3339, "2025-01-22T12:00:00Z"),
			},
		},
		{
			MockID: "evt-annual-2026",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventConsultation,
				Title:       "2026 Health Maintenance",
				Description: "Long-term diabetes management review. HbA1c remains stable at 6.6%.",
				Provider:    "Dr. Sarah Chen",
				Timestamp:   parseTime(time.RFC3339, "2026-01-15T10:00:00Z"),
				IsEncrypted: false,
				Metadata:    common.JSONMap{"status": "controlled"},
				CreatedAt:   parseTime(time.RFC3339, "2026-01-15T10:30:00Z"),
				UpdatedAt:   parseTime(time.RFC3339, "2026-01-15T10:30:00Z"),
			},
		},
		{
			MockID: "evt-future-scan",
			Event: timeline.TimelineEvent{
				PatientID:   patientId,
				Type:        prototline.EventImaging,
				Title:       "Planned Cardiac Assessment",
				Description: "Scheduled preventative scan.",
				Provider:    "City Imaging",
				Timestamp:   parseTime(time.RFC3339, "2026-06-20T14:00:00Z"),
				IsEncrypted: false,
				Metadata:    common.JSONMap{"planned": true},
				CreatedAt:   parseTime(time.RFC3339, "2025-12-01T09:00:00Z"),
				UpdatedAt:   parseTime(time.RFC3339, "2025-12-01T09:00:00Z"),
			},
		},
	}

	// 1. Process and Insert Events
	log.Printf("Seeding %d events...", len(rawEvents))
	var eventsToInsert []timeline.TimelineEvent

	for _, re := range rawEvents {
		uuid := generateUUID()
		idMap[re.MockID] = uuid
		re.Event.ID = uuid
		eventsToInsert = append(eventsToInsert, re.Event)
	}

	if err := db.Create(&eventsToInsert).Error; err != nil {
		log.Fatalf("Failed to seed events: %v", err)
	}

	// 2. Create Edges
	type MockEdge struct {
		FromMockID string
		ToMockID   string
		Type       prototline.RelationshipType
		Timestamp  string
	}

	rawEdges := []MockEdge{
		{"evt-early-checkup", "evt-consultation-1", prototline.RelLeadTo, "2024-01-15T09:30:00Z"},
		{"evt-consultation-1", "evt-lab-1", prototline.RelResultedIn, "2024-01-15T09:30:00Z"},
		{"evt-lab-1", "evt-diagnosis-1", prototline.RelSupports, "2024-01-20T10:30:00Z"},
		{"evt-diagnosis-1", "evt-prescription-1", prototline.RelLeadTo, "2024-01-20T10:30:00Z"},
		{"evt-prescription-1", "evt-consultation-2", prototline.RelResultedIn, "2024-02-05T15:00:00Z"},
		{"evt-consultation-2", "evt-lab-2", prototline.RelResultedIn, "2024-08-10T10:00:00Z"},
		{"evt-lab-2", "evt-annual-2026", prototline.RelSupports, "2026-01-15T10:30:00Z"},
	}

	var edgesToInsert []timeline.EventEdge
	for _, re := range rawEdges {
		fromUUID, ok1 := idMap[re.FromMockID]
		toUUID, ok2 := idMap[re.ToMockID]

		if !ok1 || !ok2 {
			log.Fatalf("Could not resolve edge from %s to %s", re.FromMockID, re.ToMockID)
		}

		edge := timeline.EventEdge{
			ID:               generateUUID(),
			FromEventID:      fromUUID,
			ToEventID:        toUUID,
			RelationshipType: re.Type,
			CreatedAt:        parseTime(time.RFC3339, re.Timestamp),
		}
		edgesToInsert = append(edgesToInsert, edge)
	}

	log.Printf("Seeding %d edges...", len(edgesToInsert))
	if err := db.Create(&edgesToInsert).Error; err != nil {
		log.Fatalf("Failed to seed edges: %v", err)
	}

	// Verify Data
	var eventCount int64
	db.Model(&timeline.TimelineEvent{}).Count(&eventCount)
	var edgeCount int64
	db.Model(&timeline.EventEdge{}).Count(&edgeCount)

	log.Printf("Seeding completed successfully!")
	log.Printf("Total Events: %d", eventCount)
	log.Printf("Total Edges: %d", edgeCount)

	// Print JSON for verification
	eventsJson, _ := json.MarshalIndent(eventsToInsert, "", "  ")
	fmt.Println(string(eventsJson))
}
