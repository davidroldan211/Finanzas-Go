package domain

import (
	"time"

	"github.com/google/uuid"
)

// Verification es la entidad de dominio pura: sin tags de persistencia ni
// hooks de GORM. La identidad y el mapeo a tabla son responsabilidad del
// adaptador de salida (adapter/out/postgres).
type Verification struct {
	ID            uuid.UUID
	Email         string
	CodeHash      string
	ExpiresAt     time.Time
	Attempts      int
	MaxAttempts   int
	CooldownUntil *time.Time
	UsedAt        *time.Time
	CreatedAt     time.Time
}
