package repository

import "finanzas-api/internal/users/domain"

type UserRepository interface {
	GetAll() ([]domain.User, error)
	GetByID(id int) (*domain.User, error)
	Create(user domain.User) (*domain.User, error)
	Update(user domain.User) (*domain.User, error)
	Delete(id int) error
}
