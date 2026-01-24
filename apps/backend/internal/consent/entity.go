package consent

import (
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/pkg/protocol/consent"
)

// ConsentGrant is the database model for patient-controlled access.
type ConsentGrant struct {
	ID          string             `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Grantor     string             `json:"grantor" gorm:"index;type:varchar(255);not null"` // Patient
	Grantee     string             `json:"grantee" gorm:"index;type:varchar(255);not null"` // Doctor/Researcher
	Scope       common.JSONStrings `json:"scope,omitempty" gorm:"type:jsonb"`               // List of event IDs or categories
	Permissions common.JSONStrings `json:"permissions" gorm:"type:jsonb"`                   // Read, Write, Share
	State       consent.State      `json:"state" gorm:"type:varchar(50);not null"`
	Reason      string             `json:"reason,omitempty" gorm:"type:text"`
	ExpiresAt   time.Time          `json:"expiresAt,omitempty" gorm:"index"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

// TableName returns the custom table name for consent grants.
func (ConsentGrant) TableName() string {
	return "consent_grants"
}
