package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/apps/backend/internal/consent"
	"github.com/itspablomontes/fleming/apps/backend/internal/storage"
	"github.com/itspablomontes/fleming/apps/backend/internal/timeline"
	protocol "github.com/itspablomontes/fleming/pkg/protocol/audit"
	protoconsent "github.com/itspablomontes/fleming/pkg/protocol/consent"
	prototline "github.com/itspablomontes/fleming/pkg/protocol/timeline"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// noOpStorage is a no-op storage implementation for seeding (we don't upload files)
type noOpStorage struct{}

func (n *noOpStorage) Put(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	return "", fmt.Errorf("no-op storage: Put not implemented")
}

func (n *noOpStorage) Get(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("no-op storage: Get not implemented")
}

func (n *noOpStorage) Delete(ctx context.Context, bucketName, objectName string) error {
	return fmt.Errorf("no-op storage: Delete not implemented")
}

func (n *noOpStorage) GetURL(ctx context.Context, bucketName, objectName string) (string, error) {
	return "", fmt.Errorf("no-op storage: GetURL not implemented")
}

func (n *noOpStorage) CreateMultipartUpload(ctx context.Context, bucketName, objectName, contentType string) (string, error) {
	return "", fmt.Errorf("no-op storage: CreateMultipartUpload not implemented")
}

func (n *noOpStorage) UploadPart(ctx context.Context, bucketName, objectName, uploadID string, partNumber int, reader io.Reader, objectSize int64) (string, error) {
	return "", fmt.Errorf("no-op storage: UploadPart not implemented")
}

func (n *noOpStorage) CompleteMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string, parts []storage.Part) (string, error) {
	return "", fmt.Errorf("no-op storage: CompleteMultipartUpload not implemented")
}

func (n *noOpStorage) AbortMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string) error {
	return fmt.Errorf("no-op storage: AbortMultipartUpload not implemented")
}

type MockEvent struct {
	MockID string
	Event  timeline.TimelineEvent
}

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
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://fleming:fleming@localhost:5432/fleming?sslmode=disable"
	}

	overrideAddress := os.Getenv("DEV_OVERRIDE_WALLET_ADDRESS")

	// Configure GORM to suppress all logs (including "record not found" messages)
	// This is expected behavior when seeding a fresh database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Connected to database")

	// Clean up existing data (order matters due to foreign keys)
	log.Println("Cleaning up existing data...")

	// Use GORM models for safer deletion (handles missing tables gracefully)
	if err := db.Where("1 = 1").Delete(&timeline.EventFileAccess{}).Error; err != nil {
		// Ignore "relation does not exist" errors
		if !strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: failed to clean event_file_access: %v", err)
		}
	}
	if err := db.Where("1 = 1").Delete(&timeline.EventFile{}).Error; err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: failed to clean event_files: %v", err)
		}
	}
	if err := db.Where("1 = 1").Delete(&timeline.EventEdge{}).Error; err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: failed to clean event_edges: %v", err)
		}
	}
	if err := db.Where("1 = 1").Delete(&timeline.TimelineEvent{}).Error; err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: failed to clean timeline_events: %v", err)
		}
	}
	if err := db.Where("1 = 1").Delete(&consent.ConsentGrant{}).Error; err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: failed to clean consent_grants: %v", err)
		}
	}
	if err := db.Where("1 = 1").Delete(&audit.AuditEntry{}).Error; err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: failed to clean audit_entries: %v", err)
		}
	}
	if err := db.Where("1 = 1").Delete(&audit.AuditBatch{}).Error; err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: failed to clean audit_batches: %v", err)
		}
	}

	// Initialize services
	auditRepo := audit.NewRepository(db)
	consentRepo := consent.NewRepository(db)
	timelineRepo := timeline.NewRepository(db)
	storageService := &noOpStorage{}

	auditService := audit.NewService(auditRepo)
	consentService := consent.NewService(consentRepo, auditService)
	timelineService := timeline.NewService(timelineRepo, auditService, storageService, "fleming")

	// Mock Data Constants
	patientId := "0x742d35Cc6634C0532925a3b844Bc9e7595f"
	if overrideAddress != "" {
		patientId = overrideAddress
		log.Printf("Using override patient ID: %s", patientId)
	}

	// Doctor addresses for consent grants
	doctor1 := "0xDoctor1AddressForConsentTesting"
	doctor2 := "0xDoctor2AddressForConsentTesting"
	doctor3 := "0xDoctor3AddressForConsentTesting"

	// Helper to parse time
	parseTime := func(layout, value string) time.Time {
		t, err := time.Parse(layout, value)
		if err != nil {
			log.Fatalf("failed to parse time %s: %v", value, err)
		}
		return t
	}

	// ID Mapping: Mock ID -> UUID
	idMap := make(map[string]string)

	// Define raw events with logical chronological order
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
				Timestamp:   parseTime(time.RFC3339, "2024-01-16T14:00:00Z"), // After consultation
				IsEncrypted: false,
				Metadata: common.JSONMap{
					"glucoseLevel": 186,
					"hba1c":        7.8,
				},
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
				Timestamp:   parseTime(time.RFC3339, "2024-01-20T10:00:00Z"), // After lab results
				IsEncrypted: false,
				Metadata: common.JSONMap{
					"icd10":    "E11.9",
					"severity": "moderate",
				},
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
				Timestamp:   parseTime(time.RFC3339, "2024-01-20T10:15:00Z"), // Same day as diagnosis, slightly after
				IsEncrypted: false,
				Metadata: common.JSONMap{
					"medication": "Metformin",
					"dosage":     "500mg",
				},
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
				Timestamp:   parseTime(time.RFC3339, "2024-02-05T14:00:00Z"), // After prescription
				IsEncrypted: false,
				Metadata:    common.JSONMap{"specialty": "Endocrinology"},
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
				Timestamp:   parseTime(time.RFC3339, "2024-08-10T09:00:00Z"), // 6 months after treatment
				IsEncrypted: false,
				Metadata:    common.JSONMap{"hba1c": 6.8},
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
				Timestamp:   parseTime(time.RFC3339, "2025-01-22T11:00:00Z"), // Annual screening after diagnosis
				IsEncrypted: false,
				Metadata:    common.JSONMap{"results": "clear"},
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
				Timestamp:   parseTime(time.RFC3339, "2026-01-15T10:00:00Z"), // Follow-up after lab-2
				IsEncrypted: false,
				Metadata:    common.JSONMap{"status": "controlled"},
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
				Timestamp:   parseTime(time.RFC3339, "2026-06-20T14:00:00Z"), // Future planned event
				IsEncrypted: false,
				Metadata:    common.JSONMap{"planned": true},
			},
		},
	}

	// Validate chronological order
	for i := 1; i < len(rawEvents); i++ {
		if rawEvents[i].Event.Timestamp.Before(rawEvents[i-1].Event.Timestamp) {
			log.Fatalf("Event %s (timestamp: %v) is before previous event %s (timestamp: %v)",
				rawEvents[i].MockID, rawEvents[i].Event.Timestamp,
				rawEvents[i-1].MockID, rawEvents[i-1].Event.Timestamp)
		}
	}

	// 1. Seed Timeline Events using service (generates audit logs)
	log.Printf("Seeding %d timeline events...", len(rawEvents))
	for _, re := range rawEvents {
		uuid := generateUUID()
		idMap[re.MockID] = uuid
		re.Event.ID = uuid

		// Use service to add event (this generates audit logs)
		if err := timelineService.AddEvent(ctx, &re.Event); err != nil {
			log.Fatalf("Failed to seed event %s: %v", re.MockID, err)
		}
	}

	// 2. Create Edges with logical relationships
	type MockEdge struct {
		FromMockID string
		ToMockID   string
		Type       prototline.RelationshipType
		Timestamp  string
	}

	rawEdges := []MockEdge{
		// Baseline checkup leads to annual checkup (logical progression)
		{"evt-early-checkup", "evt-consultation-1", prototline.RelLeadTo, "2024-01-15T09:30:00Z"},
		// Consultation resulted in lab order
		{"evt-consultation-1", "evt-lab-1", prototline.RelResultedIn, "2024-01-15T09:30:00Z"},
		// Lab results support diagnosis
		{"evt-lab-1", "evt-diagnosis-1", prototline.RelSupports, "2024-01-20T10:30:00Z"},
		// Diagnosis leads to prescription
		{"evt-diagnosis-1", "evt-prescription-1", prototline.RelLeadTo, "2024-01-20T10:30:00Z"},
		// Prescription resulted in specialist consultation
		{"evt-prescription-1", "evt-consultation-2", prototline.RelResultedIn, "2024-02-05T15:00:00Z"},
		// Specialist consultation resulted in follow-up lab
		{"evt-consultation-2", "evt-lab-2", prototline.RelResultedIn, "2024-08-10T10:00:00Z"},
		// Lab results support annual review
		{"evt-lab-2", "evt-annual-2026", prototline.RelSupports, "2026-01-15T10:30:00Z"},
	}

	// Validate edge relationships make clinical sense
	for _, edge := range rawEdges {
		fromEvent := findEventByMockID(rawEvents, edge.FromMockID)
		toEvent := findEventByMockID(rawEvents, edge.ToMockID)

		if fromEvent == nil || toEvent == nil {
			log.Fatalf("Edge references non-existent event: from=%s, to=%s", edge.FromMockID, edge.ToMockID)
		}

		// Validate timestamp order: from event should be before or equal to to event
		if fromEvent.Event.Timestamp.After(toEvent.Event.Timestamp) {
			log.Fatalf("Edge from %s (timestamp: %v) to %s (timestamp: %v) violates chronological order",
				edge.FromMockID, fromEvent.Event.Timestamp,
				edge.ToMockID, toEvent.Event.Timestamp)
		}
	}

	log.Printf("Seeding %d edges...", len(rawEdges))
	for _, re := range rawEdges {
		fromUUID, ok1 := idMap[re.FromMockID]
		toUUID, ok2 := idMap[re.ToMockID]

		if !ok1 || !ok2 {
			log.Fatalf("Could not resolve edge from %s to %s", re.FromMockID, re.ToMockID)
		}

		// Use service to link events
		if _, err := timelineService.LinkEvents(ctx, fromUUID, toUUID, re.Type); err != nil {
			log.Fatalf("Failed to seed edge from %s to %s: %v", re.FromMockID, re.ToMockID, err)
		}
	}

	// 3. Seed Consent Grants with various states using service (generates audit logs)
	log.Println("Seeding consent grants...")

	now := time.Now()

	// Approved grant (active)
	grant1, err := consentService.RequestConsent(ctx, patientId, doctor1, "Primary care physician access", []string{"read", "write"}, now.Add(365*24*time.Hour))
	if err != nil {
		log.Fatalf("Failed to create consent grant 1: %v", err)
	}
	if err := consentService.ApproveConsent(ctx, grant1.ID); err != nil {
		log.Fatalf("Failed to approve consent grant 1: %v", err)
	}

	// Denied grant
	grant2, err := consentService.RequestConsent(ctx, patientId, doctor2, "Research study participation", []string{"read", "share"}, now.Add(180*24*time.Hour))
	if err != nil {
		log.Fatalf("Failed to create consent grant 2: %v", err)
	}
	if err := consentService.DenyConsent(ctx, grant2.ID); err != nil {
		log.Fatalf("Failed to deny consent grant 2: %v", err)
	}

	// Revoked grant (was approved, then revoked)
	grant3, err := consentService.RequestConsent(ctx, patientId, doctor3, "Specialist consultation", []string{"read"}, now.Add(90*24*time.Hour))
	if err != nil {
		log.Fatalf("Failed to create consent grant 3: %v", err)
	}
	if err := consentService.ApproveConsent(ctx, grant3.ID); err != nil {
		log.Fatalf("Failed to approve consent grant 3: %v", err)
	}
	if err := consentService.RevokeConsent(ctx, grant3.ID); err != nil {
		log.Fatalf("Failed to revoke consent grant 3: %v", err)
	}

	// Expired grant (approved but expired)
	grant4, err := consentService.RequestConsent(ctx, patientId, doctor2, "Temporary access for consultation", []string{"read"}, now.Add(-24*time.Hour)) // Expired yesterday
	if err != nil {
		log.Fatalf("Failed to create consent grant 4: %v", err)
	}
	if err := consentService.ApproveConsent(ctx, grant4.ID); err != nil {
		log.Fatalf("Failed to approve consent grant 4: %v", err)
	}
	// Manually set state to expired (since CheckPermission would auto-expire it)
	grant4Get, err := consentService.GetGrantByID(ctx, grant4.ID)
	if err != nil {
		log.Fatalf("Failed to get grant 4: %v", err)
	}
	grant4Get.State = protoconsent.StateExpired
	if err := consentRepo.Update(ctx, grant4Get); err != nil {
		log.Fatalf("Failed to expire grant 4: %v", err)
	}
	// Record expiration audit log
	_ = auditService.Record(ctx, grant4Get.Grantor, protocol.ActionConsentExpire, protocol.ResourceConsent, grant4Get.ID, nil)

	// Pending grant (requested but not yet approved/denied)
	_, err = consentService.RequestConsent(ctx, patientId, doctor1, "Extended access for ongoing treatment", []string{"read", "write", "share"}, now.Add(730*24*time.Hour))
	if err != nil {
		log.Fatalf("Failed to create consent grant 5: %v", err)
	}
	// Leave in requested state

	// Verify Data
	var eventCount int64
	db.Model(&timeline.TimelineEvent{}).Count(&eventCount)
	var edgeCount int64
	db.Model(&timeline.EventEdge{}).Count(&edgeCount)
	var consentCount int64
	db.Model(&consent.ConsentGrant{}).Count(&consentCount)
	var auditCount int64
	db.Model(&audit.AuditEntry{}).Count(&auditCount)

	log.Printf("Seeding completed successfully!")
	log.Printf("Total Events: %d", eventCount)
	log.Printf("Total Edges: %d", edgeCount)
	log.Printf("Total Consent Grants: %d", consentCount)
	log.Printf("Total Audit Entries: %d", auditCount)

	// Print summary JSON
	summary := map[string]interface{}{
		"events":    eventCount,
		"edges":     edgeCount,
		"consents":  consentCount,
		"auditLogs": auditCount,
		"patientId": patientId,
		"doctors":   []string{doctor1, doctor2, doctor3},
	}
	summaryJson, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(summaryJson))
}

func findEventByMockID(events []MockEvent, mockID string) *MockEvent {
	for i := range events {
		if events[i].MockID == mockID {
			return &events[i]
		}
	}
	return nil
}
