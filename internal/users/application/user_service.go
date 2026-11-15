package application

import (
	"context"
	"errors"
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

// Create es el registro público. El rol siempre es user; un admin solo se
// crea promoviendo a un usuario existente vía Update.
//
// Orden (SPEC 01): validar datos y contraseña, hashear y solo después
// consultar si el correo existe. Hashear antes iguala el costo de bcrypt
// entre un correo nuevo y uno ya registrado (A5).
//
// Un correo ya registrado (también de un usuario eliminado) devuelve nil sin
// guardar nada, igual que un registro nuevo. Una carrera entre dos registros
// con el mismo correo llega como ErrEmailTaken desde Save y también termina
// en nil. Ninguno de los dos caminos se registra en logs: el correo es
// Confidencial y su presencia en un log delataría que existe.
//
// TODO: exigir correo verificado antes de activar la cuenta cuando se
// implemente internal/verification.
func (s *userService) Create(ctx context.Context, cmd in.CreateUserCommand) error {
	user, err := validateRegistration(cmd)
	if err != nil {
		return err
	}

	hash, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return err
	}

	exists, err := s.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	user.PasswordHash = hash
	if err := s.repo.Save(ctx, user); err != nil {
		if errors.Is(err, domain.ErrEmailTaken) {
			return nil
		}
		return err
	}
	return nil
}

// validateRegistration construye el usuario y aplica la política de
// contraseñas, reuniendo todos los campos inválidos en un solo
// ValidationError. No toca el repositorio ni el hasher.
func validateRegistration(cmd in.CreateUserCommand) (*domain.User, error) {
	fields := map[string]string{}

	user, err := domain.NewUser(cmd.Email, cmd.FirstName, cmd.LastName, "", domain.RoleUser)
	if err != nil {
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			return nil, err
		}
		for k, v := range ve.Fields {
			fields[k] = v
		}
	}

	if msg := domain.ValidatePassword(cmd.Password, cmd.Email, cmd.FirstName, cmd.LastName); msg != "" {
		fields["password"] = msg
	}

	if len(fields) > 0 {
		return nil, &domain.ValidationError{Fields: fields}
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, actor in.Actor, id uuid.UUID) (*domain.User, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidUserID
	}
	// Un usuario ajeno responde igual que uno inexistente: no se revela
	// qué IDs existen.
	if !actor.CanAccess(id) {
		return nil, domain.ErrUserNotFound
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

func (s *userService) Update(ctx context.Context, actor in.Actor, cmd in.UpdateUserCommand) (*domain.User, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidUserID
	}
	if !actor.CanAccess(cmd.ID) {
		return nil, domain.ErrUserNotFound
	}
	// Rol y estado de la cuenta son decisiones de administración.
	if !actor.IsAdmin() && (cmd.Role != nil || cmd.IsActive != nil) {
		return nil, domain.ErrForbidden
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
