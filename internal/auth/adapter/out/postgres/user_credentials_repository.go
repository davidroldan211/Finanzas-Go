package postgres

import (
	"context"
	"errors"
	"fmt"

	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/auth/port/out"

	"gorm.io/gorm"
)

type userCredentialsRepository struct {
	db *gorm.DB
}

// NewUserCredentialsRepository implementa out.UserFinder contra la misma
// tabla "users" que usa el módulo users, pero en solo lectura y sin
// importar su paquete de persistencia.
func NewUserCredentialsRepository(db *gorm.DB) out.UserFinder {
	return &userCredentialsRepository{db: db}
}

func (r *userCredentialsRepository) FindByEmail(ctx context.Context, email string) (*domain.Credentials, error) {
	var m userCredentialsModel
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("auth: find credentials by email: %w", err)
	}

	return &domain.Credentials{
		UserID:       m.ID.String(),
		Email:        m.Email,
		PasswordHash: m.Password,
		Role:         m.Role,
		Active:       m.IsActive,
	}, nil
}
