package repository

import (
	"errors"
	"finanzas-api/internal/users/domain"
)

type userRepositoryMemory struct {
	data []domain.User
	id   int
}

func NewUserRepositoryMemory() UserRepository {
	return &userRepositoryMemory{data: []domain.User{}, id: 1}
}

func (r *userRepositoryMemory) GetAll() ([]domain.User, error) {
	return r.data, nil
}

func (r *userRepositoryMemory) GetByID(id int) (*domain.User, error) {
	for _, u := range r.data {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepositoryMemory) Create(user domain.User) (*domain.User, error) {
	user.ID = r.id
	r.id++
	r.data = append(r.data, user)
	return &user, nil
}

func (r *userRepositoryMemory) Update(user domain.User) (*domain.User, error) {
	for i, u := range r.data {
		if u.ID == user.ID {
			r.data[i] = user
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepositoryMemory) Delete(id int) error {
	for i, u := range r.data {
		if u.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}
