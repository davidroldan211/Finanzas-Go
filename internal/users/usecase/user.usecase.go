package usecase

import (
	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/repository"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) GetAll() ([]domain.User, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) GetByID(id int) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) Create(user domain.User) (*domain.User, error) {
	return uc.repo.Create(user)
}

func (uc *UserUseCase) Update(user domain.User) (*domain.User, error) {
	return uc.repo.Update(user)
}

func (uc *UserUseCase) Delete(id int) error {
	return uc.repo.Delete(id)
}
