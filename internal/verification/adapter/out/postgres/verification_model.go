package postgres

import (
	"time"

	"finanzas-api/internal/verification/domain"

	"github.com/google/uuid"
)

// verificationModel representa la estructura de persistencia para GORM.
// No es la entidad de dominio: solo los mappers de este archivo cruzan la
// frontera.
type verificationModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string     `gorm:"type:varchar(255);not null;index"`
	CodeHash      string     `gorm:"type:varchar(255);not null"`
	ExpiresAt     time.Time  `gorm:"type:timestamptz;not null"`
	Attempts      int        `gorm:"type:int;default:0"`
	MaxAttempts   int        `gorm:"type:int;default:5"`
	CooldownUntil *time.Time `gorm:"type:timestamptz"`
	UsedAt        *time.Time `gorm:"type:timestamptz"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;autoCreateTime"`
}

// TableName corrige el nombre real de la tabla: la migración
// (db/migrations/verifications.sql) crea "email_verifications", no
// "verifications" como declaraba la entidad de dominio antes de esta
// migración a hexagonal. Con el modelo de dominio persistido directamente
// (sin mapper) ese desajuste nunca se ejercitaba, así que quedó latente.
func (verificationModel) TableName() string { return "email_verifications" }

func toDomain(m *verificationModel) *domain.Verification {
	if m == nil {
		return nil
	}
	return &domain.Verification{
		ID:            m.ID,
		Email:         m.Email,
		CodeHash:      m.CodeHash,
		ExpiresAt:     m.ExpiresAt,
		Attempts:      m.Attempts,
		MaxAttempts:   m.MaxAttempts,
		CooldownUntil: m.CooldownUntil,
		UsedAt:        m.UsedAt,
		CreatedAt:     m.CreatedAt,
	}
}

func toModel(v *domain.Verification) *verificationModel {
	if v == nil {
		return nil
	}
	return &verificationModel{
		ID:            v.ID,
		Email:         v.Email,
		CodeHash:      v.CodeHash,
		ExpiresAt:     v.ExpiresAt,
		Attempts:      v.Attempts,
		MaxAttempts:   v.MaxAttempts,
		CooldownUntil: v.CooldownUntil,
		UsedAt:        v.UsedAt,
		CreatedAt:     v.CreatedAt,
	}
}
