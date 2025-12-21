package domain

import (
	"time"

	"github.com/google/uuid"
)

type Verification struct {
	ID            uuid.UUID
	Email         string
	CodeHash      string
	ExpiresAt     time.Time
	Attempts      int
	MaxAttempts   int
	CooldownUntil time.Time
	UsedAt        time.Time
	CreatedAt     time.Time
}

type VerificationRepository interface {
	Create(v *Verification) error
	GetByEmail(email string) (*Verification, error)
	Update(v *Verification) error
	Delete(id uuid.UUID) error
}

type VerificationUseCase interface {
	GenerateVerificationCode(email string) (string, error)
	// ValidateVerificationCode(ID uuid.UUID, email, code string) (bool, error)
	// SendVerificationCode(ID uuid.UUID, email string) error
	// CleanupExpiredCodes() error
}

func (Verification) TableName() string {
	return "verifications"
}

func (v *Verification) BeforeCreate(tx interface{}) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}
