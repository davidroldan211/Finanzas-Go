package application

import (
	"context"
	"strings"

	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/port/in"
	"finanzas-api/internal/users/port/out"

	"github.com/google/uuid"
)

type userService struct {
	repo   out.UserRepository
	hasher out.PasswordHasher
}

// NewUserService construye el puerto de entrada in.UserService. Depende
// solo de puertos de salida: ni gin ni gorm ni bcrypt aparecen aquí.
func NewUserService(repo out.UserRepository, hasher out.PasswordHasher) in.UserService {
	return &userService{repo: repo, hasher: hasher}
}

// TODO: se debe ajustar el crear usuario para que valide bearer token una
// vez que se implemente la validación de correos (internal/verification).
func (s *userService) Create(ctx context.Context, cmd in.CreateUserCommand) (*domain.User, error) {
	// Validar y normalizar antes de tocar el repositorio o el hasher.
	user, err := domain.NewUser(cmd.Email, cmd.FirstName, cmd.LastName, "", cmd.Role)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailTaken
	}

	hash, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hash

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidUserID
	}
	return s.repo.FindByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, &domain.ValidationError{Fields: map[string]string{"email": "el email es obligatorio"}}
	}
	return s.repo.FindByEmail(ctx, email)
}

func (s *userService) Update(ctx context.Context, cmd in.UpdateUserCommand) (*domain.User, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidUserID
	}

	user, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	previousEmail := user.Email
	if err := user.ApplyUpdate(cmd.Email, cmd.FirstName, cmd.LastName, cmd.Role, cmd.IsActive); err != nil {
		return nil, err
	}

	// Solo se consulta unicidad de email si realmente cambió.
	if user.Email != previousEmail {
		exists, err := s.repo.ExistsByEmail(ctx, user.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrEmailTaken
		}
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ErrInvalidUserID
	}
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, q in.ListUsersQuery) ([]*domain.User, error) {
	limit, offset := q.Limit, q.Offset
	if limit < 0 || offset < 0 {
		return nil, &domain.ValidationError{Fields: map[string]string{"pagination": "limit y offset deben ser no negativos"}}
	}

	// Valor por defecto para limit.
	if limit == 0 {
		limit = 10
	}
	// Máximo 100 usuarios por página.
	if limit > 100 {
		limit = 100
	}

	return s.repo.List(ctx, limit, offset)
}
