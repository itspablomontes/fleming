package audit

import (
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/pkg/protocol/audit"
)

// AuditEntry is the database model for cryptographic audit logs.
type AuditEntry struct {
	ID           string             `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Actor        string             `json:"actor" gorm:"index;type:varchar(255);not null"`
	Action       audit.Action       `json:"action" gorm:"type:varchar(50);not null"`
	ResourceType audit.ResourceType `json:"resourceType" gorm:"type:varchar(50);not null"`
	ResourceID   string             `json:"resourceId" gorm:"index;type:varchar(255);not null"`
	Timestamp    time.Time          `json:"timestamp" gorm:"index;not null"`
	Metadata     common.JSONMap     `json:"metadata,omitempty" gorm:"type:jsonb"`
	Hash         string             `json:"hash" gorm:"type:varchar(64);not null"`
	PreviousHash string             `json:"previousHash" gorm:"type:varchar(64);not null"`
}

// TableName returns the custom table name for audit entries.
func (AuditEntry) TableName() string {
	return "audit_entries"
}
