package audit

import (
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/pkg/protocol/audit"
)

// AuditEntry is the database model for cryptographic audit logs.
type AuditEntry struct {
	ID             string             `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Actor          string             `json:"actor" gorm:"index;index:idx_audit_actor_timestamp,priority:1;type:varchar(255);not null"`
	Action         audit.Action       `json:"action" gorm:"index:idx_audit_resource_type_action_timestamp,priority:2;type:varchar(50);not null"`
	ResourceType   audit.ResourceType `json:"resourceType" gorm:"index:idx_audit_resource_type_action_timestamp,priority:1;type:varchar(50);not null"`
	ResourceID     string             `json:"resourceId" gorm:"index;index:idx_audit_resource_timestamp,priority:1;type:varchar(255);not null"`
	Timestamp      time.Time          `json:"timestamp" gorm:"index;index:idx_audit_actor_timestamp,priority:2;index:idx_audit_resource_timestamp,priority:2;index:idx_audit_resource_type_action_timestamp,priority:3;not null"`
	Metadata       common.JSONMap     `json:"metadata,omitempty" gorm:"type:jsonb"`
	Hash           string             `json:"hash" gorm:"type:varchar(64);not null"`
	PreviousHash   string             `json:"previousHash" gorm:"type:varchar(64);not null"`
	SchemaVersion  string             `json:"schemaVersion,omitempty" gorm:"type:varchar(20)"`
}

// TableName returns the custom table name for audit entries.
func (AuditEntry) TableName() string {
	return "audit_entries"
}
