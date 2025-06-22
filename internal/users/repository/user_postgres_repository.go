package repository

import (
	"finanzas-api/internal/users/domain"

	"gorm.io/gorm"
)

type userPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) domain.UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userPostgresRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error // soft delete
}

func (r *userPostgresRepository) List(limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userPostgresRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error
	return count > 0, err
}
