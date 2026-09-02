package postgres

import (
	"context"
	"errors"
	"fmt"

	"finanzas-api/internal/verification/domain"
	"finanzas-api/internal/verification/port/out"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type verificationPostgresRepository struct {
	db *gorm.DB
}

// NewVerificationPostgresRepository implementa out.VerificationRepository
// contra Postgres/GORM, usando verificationModel (no la entidad de
// dominio directamente).
func NewVerificationPostgresRepository(db *gorm.DB) out.VerificationRepository {
	return &verificationPostgresRepository{db: db}
}

func (r *verificationPostgresRepository) Create(ctx context.Context, v *domain.Verification) error {
	m := toModel(v)
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("verification: create: %w", err)
	}
	*v = *toDomain(m)
	return nil
}

func (r *verificationPostgresRepository) GetByEmail(ctx context.Context, email string) (*domain.Verification, error) {
	var m verificationModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("verification: get by email: %w", err)
	}
	return toDomain(&m), nil
}

func (r *verificationPostgresRepository) Update(ctx context.Context, v *domain.Verification) error {
	m := toModel(v)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("verification: update: %w", err)
	}
	*v = *toDomain(m)
	return nil
}

func (r *verificationPostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&verificationModel{}).Error; err != nil {
		return fmt.Errorf("verification: delete: %w", err)
	}
	return nil
}
