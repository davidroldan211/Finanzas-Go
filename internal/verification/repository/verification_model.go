package repository

import (
	"time"

	"github.com/google/uuid"
)

type VerificationModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email         string    `gorm:"type:varchar(255);not null;index"`
	CodeHash      string    `gorm:"type:varchar(255);not null"`
	ExpiresAt     time.Time `gorm:"type:timestamp;not null"`
	Attempts      int       `gorm:"type:int;default:0"`
	MaxAttempts   int       `gorm:"type:int;default:5"`
	CooldownUntil time.Time `gorm:"type:timestamp"`
	UsedAt        time.Time `gorm:"type:timestamp"`
	CreatedAt     time.Time `gorm:"type:timestamp;autoCreateTime"`
}
