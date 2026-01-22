package auth

import "time"

type Challenge struct {
	Address   string    `gorm:"primaryKey;type:varchar(255)"`
	Message   string    `gorm:"type:text;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
}

func (Challenge) TableName() string {
	return "auth_challenges"
}
