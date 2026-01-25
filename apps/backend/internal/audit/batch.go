package audit

import "time"

// AuditBatch tracks a batch of audit entries summarized by a Merkle root.
type AuditBatch struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RootHash   string    `json:"rootHash" gorm:"type:varchar(64);not null;uniqueIndex"`
	StartTime  time.Time `json:"startTime" gorm:"index;not null"`
	EndTime    time.Time `json:"endTime" gorm:"index;not null"`
	EntryCount int       `json:"entryCount" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"index;not null"`
}

// TableName returns the custom table name for audit batches.
func (AuditBatch) TableName() string {
	return "audit_batches"
}
