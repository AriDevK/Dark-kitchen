package models

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}
