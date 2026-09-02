package postgres

import (
	"context"
	"errors"
	"fmt"

	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/port/out"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPostgresRepository struct {
	db *gorm.DB
}

// NewUserPostgresRepository implementa out.UserRepository contra Postgres/GORM.
func NewUserPostgresRepository(db *gorm.DB) out.UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) Save(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("users: save: %w", err)
	}
	*user = *toDomain(m)
	return nil
}

func (r *userPostgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m userModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("users: find by id: %w", err)
	}
	return toDomain(&m), nil
}

func (r *userPostgresRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m userModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("users: find by email: %w", err)
	}
	return toDomain(&m), nil
}

func (r *userPostgresRepository) Update(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("users: update: %w", err)
	}
	*user = *toDomain(m)
	return nil
}

func (r *userPostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Desactiva y soft-elimina en una única transacción: hoy estas dos
	// escrituras se hacían por separado y el error de la primera se
	// descartaba en silencio.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&userModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
			return fmt.Errorf("users: deactivate: %w", err)
		}
		if err := tx.Where("id = ?", id).Delete(&userModel{}).Error; err != nil {
			return fmt.Errorf("users: delete: %w", err)
		}
		return nil
	})
}

func (r *userPostgresRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var models []*userModel
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("users: list: %w", err)
	}
	users := make([]*domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, toDomain(m))
	}
	return users, nil
}

func (r *userPostgresRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&userModel{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("users: exists by email: %w", err)
	}
	return count > 0, nil
}
