package timeline

import (
	"time"

	"github.com/itspablomontes/fleming/api/internal/common"
)

type TimelineEventType string

const (
	EventTypeConsultation   TimelineEventType = "consultation"
	EventTypeDiagnosis      TimelineEventType = "diagnosis"
	EventTypePrescription   TimelineEventType = "prescription"
	EventTypeProcedure      TimelineEventType = "procedure"
	EventTypeLabResult      TimelineEventType = "lab_result"
	EventTypeImaging        TimelineEventType = "imaging"
	EventTypeNote           TimelineEventType = "note"
	EventTypeVaccination    TimelineEventType = "vaccination"
	EventTypeAllergy        TimelineEventType = "allergy"
	EventTypeVisitNote      TimelineEventType = "visit_note"
	EventTypeVitalSigns     TimelineEventType = "vital_signs"
	EventTypeReferral       TimelineEventType = "referral"
	EventTypeInsuranceClaim TimelineEventType = "insurance_claim"
	EventTypeOther          TimelineEventType = "other"
)

type RelationshipType string

const (
	RelationshipResultedIn  RelationshipType = "resulted_in"  // consultation → diagnosis
	RelationshipLeadTo      RelationshipType = "lead_to"      // diagnosis → treatment
	RelationshipRequestedBy RelationshipType = "requested_by" // lab test ← consultation
	RelationshipSupports    RelationshipType = "supports"     // lab report → diagnosis
	RelationshipFollowsUp   RelationshipType = "follows_up"   // visit → previous visit
	RelationshipContradicts RelationshipType = "contradicts"  // finding → previous finding
	RelationshipAttachedTo  RelationshipType = "attached_to"  // file → event
)

type TimelineEvent struct {
	ID           string            `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PatientID    string            `json:"patientId" gorm:"index;type:varchar(255);not null"`
	Type         TimelineEventType `json:"type" gorm:"type:varchar(50);not null"`
	Title        string            `json:"title" gorm:"type:varchar(255);not null"`
	Description  string            `json:"description,omitempty" gorm:"type:text"`
	Provider     string            `json:"provider,omitempty" gorm:"type:varchar(255)"`
	Code         string            `json:"code,omitempty" gorm:"type:varchar(50)"`         // e.g. "E11.9" or "8480-6"
	CodingSystem string            `json:"codingSystem,omitempty" gorm:"type:varchar(50)"` // e.g. "ICD-10", "LOINC"
	Timestamp    time.Time         `json:"timestamp" gorm:"index;not null"`
	BlobRef      string            `json:"blobRef,omitempty" gorm:"type:varchar(255)"`
	IsEncrypted  bool              `json:"isEncrypted" gorm:"not null;default:false"`
	Metadata     common.JSONMap    `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`

	OutgoingEdges []EventEdge `json:"outgoingEdges,omitempty" gorm:"foreignKey:FromEventID"`
	IncomingEdges []EventEdge `json:"incomingEdges,omitempty" gorm:"foreignKey:ToEventID"`
	Files         []EventFile `json:"files,omitempty" gorm:"foreignKey:EventID"`
}

type EventEdge struct {
	ID               string           `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FromEventID      string           `json:"fromEventId" gorm:"type:uuid;not null;index"`
	ToEventID        string           `json:"toEventId" gorm:"type:uuid;not null;index"`
	RelationshipType RelationshipType `json:"relationshipType" gorm:"type:varchar(50);not null"`
	Metadata         common.JSONMap   `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt        time.Time        `json:"createdAt"`

	FromEvent *TimelineEvent `json:"fromEvent,omitempty" gorm:"foreignKey:FromEventID"`
	ToEvent   *TimelineEvent `json:"toEvent,omitempty" gorm:"foreignKey:ToEventID"`
}

func (EventEdge) TableName() string {
	return "event_edges"
}

type EventFile struct {
	ID         string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EventID    string         `json:"eventId" gorm:"type:uuid;not null;index"`
	BlobRef    string         `json:"blobRef" gorm:"type:varchar(255);not null"`
	FileName   string         `json:"fileName" gorm:"type:varchar(255);not null"`
	MimeType   string         `json:"mimeType" gorm:"type:varchar(100);not null"`
	FileSize   int64          `json:"fileSize" gorm:"not null"`
	WrappedDEK []byte         `json:"-" gorm:"type:bytea;not null"` // Never expose to client
	Metadata   common.JSONMap `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt  time.Time      `json:"createdAt"`

	Event *TimelineEvent `json:"event,omitempty" gorm:"foreignKey:EventID"`
}

func (EventFile) TableName() string {
	return "event_files"
}
