package repository

import (
	"finanzas-api/internal/verification/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type verificationPostgresRepository struct {
	db *gorm.DB
}

func NewVerificationPostgresRepository(db *gorm.DB) domain.VerificationRepository {
	return &verificationPostgresRepository{db: db}
}

func (r *verificationPostgresRepository) Create(v *domain.Verification) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return r.db.Create(v).Error
}

func (r *verificationPostgresRepository) GetByEmail(email string) (*domain.Verification, error) {
	var v domain.Verification
	if err := r.db.Where("email = ?", email).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *verificationPostgresRepository) Update(v *domain.Verification) error {
	return r.db.Save(v).Error
}

func (r *verificationPostgresRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Verification{}, id).Error
}
