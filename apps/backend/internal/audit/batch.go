package audit

import "time"

// AuditBatch tracks a batch of audit entries summarized by a Merkle root.
type AuditBatch struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	Actor string `json:"actor" gorm:"type:varchar(255);not null;index;uniqueIndex:idx_audit_batches_actor_root_hash,priority:1"`

	RootHash string `json:"rootHash" gorm:"type:varchar(64);not null;uniqueIndex:idx_audit_batches_actor_root_hash,priority:2"`

	StartTime  time.Time `json:"startTime" gorm:"index;not null"`
	EndTime    time.Time `json:"endTime" gorm:"index;not null"`
	EntryCount int       `json:"entryCount" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"index;not null"`

	AnchorTxHash      *string    `json:"anchorTxHash,omitempty" gorm:"type:varchar(66);index"`
	AnchorBlockNumber *uint64    `json:"anchorBlockNumber,omitempty" gorm:"index"`
	AnchoredAt        *time.Time `json:"anchoredAt,omitempty" gorm:"index"`
	AnchorStatus      string     `json:"anchorStatus" gorm:"type:varchar(20);not null;default:'pending';index"`
	AnchorError       *string    `json:"anchorError,omitempty" gorm:"type:text"`
}

// TableName returns the custom table name for audit batches.
func (AuditBatch) TableName() string {
	return "audit_batches"
}
