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

type User struct {
	Address        string    `gorm:"primaryKey;type:varchar(255)"`
	EncryptionSalt string    `gorm:"type:varchar(64);not null"` // Hex string
	CreatedAt      time.Time `gorm:"index;not null;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"index;not null;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
