package repository

import (
	"finanzas-api/internal/users/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) domain.UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) Create(user *domain.User) error {
	m := toModel(user)
	// Genera ID si viene vacío
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	*user = *toDomain(m)
	return nil
}

func (r *userPostgresRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var m userModel
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *userPostgresRepository) GetByEmail(email string) (*domain.User, error) {
	var m userModel
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *userPostgresRepository) Update(user *domain.User) error {
	m := toModel(user)
	if err := r.db.Save(m).Error; err != nil {
		return err
	}
	*user = *toDomain(m)
	return nil
}

func (r *userPostgresRepository) Delete(id uuid.UUID) error {
	r.db.Model(&userModel{}).Where("id = ?", id).Update("is_active", false)
	return r.db.Delete(&userModel{}, id).Error // soft delete
}

func (r *userPostgresRepository) List(limit, offset int) ([]*domain.User, error) {
	var models []*userModel
	if err := r.db.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, err
	}
	users := make([]*domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, toDomain(m))
	}
	return users, nil
}

func (r *userPostgresRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&userModel{}).Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error
	return count > 0, err
}
