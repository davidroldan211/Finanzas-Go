package usecase

import (
	"errors"
	"finanzas-api/internal/users/domain"
	"finanzas-api/shared/security"
	"strings"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(UserRepo domain.UserRepository) domain.UserUseCase {
	return &UserUseCase{
		userRepo: UserRepo,
	}
}

// CreateUser implements domain.UserUseCase.
func (uc *UserUseCase) CreateUser(user *domain.User) error {
	// Validar datos del usuario
	if err := uc.ValidateUserData(user); err != nil {
		return err
	}

	// Verificar si el email ya existe
	exists, err := uc.userRepo.EmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Crear usuario
	return uc.userRepo.Create(user)
}

// DeleteUser implements domain.UserUseCase.
func (uc *UserUseCase) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}

	// Verificar que el usuario existe
	_, err := uc.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return uc.userRepo.Delete(id)
}

// GetUserByEmail implements domain.UserUseCase.
func (uc *UserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	return uc.userRepo.GetByEmail(email)
}

// GetUserByID implements domain.UserUseCase.
func (uc *UserUseCase) GetUserByID(id uint) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	return uc.userRepo.GetByID(id)
}

// ListUsers implements domain.UserUseCase.
func (uc *UserUseCase) ListUsers(limit int, offset int) ([]*domain.User, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset must be non-negative")
	}

	// Valor por defecto para limit
	if limit == 0 {
		limit = 10
	}

	// Máximo 100 usuarios por página
	if limit > 100 {
		limit = 100
	}
	return uc.userRepo.List(limit, offset)
}

// UpdateUser implements domain.UserUseCase.
func (uc *UserUseCase) UpdateUser(user *domain.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}

	// Validar datos del usuario
	if err := uc.ValidateUserData(user); err != nil {
		return err
	}

	// Verificar que el usuario existe
	existingUser, err := uc.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}

	// Si cambió el email, verificar que no exista
	if existingUser.Email != user.Email {
		exists, err := uc.userRepo.EmailExists(user.Email)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("email already exists")
		}
	}

	return uc.userRepo.Update(user)
}

// ValidateUserData implements domain.UserUseCase.
func (uc *UserUseCase) ValidateUserData(user *domain.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	// Validar email
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	// Validación básica de email
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}

	// Validar nombres
	if strings.TrimSpace(user.FirstName) == "" {
		return errors.New("first name is required")
	}

	if strings.TrimSpace(user.LastName) == "" {
		return errors.New("last name is required")
	}

	// Limpiar espacios en nombres
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	// Validar longitud de campos
	if len(user.Email) > 255 {
		return errors.New("email too long")
	}

	if len(user.FirstName) > 100 {
		return errors.New("first name too long")
	}

	if len(user.LastName) > 100 {
		return errors.New("last name too long")
	}

	return nil
}
