package in

import (
	"context"

	"finanzas-api/internal/users/domain"

	"github.com/google/uuid"
)

// UserService es el puerto de entrada del módulo users: lo que la
// aplicación ofrece a sus adaptadores de entrada (p.ej. adapter/in/http).
type UserService interface {
	Create(ctx context.Context, cmd CreateUserCommand) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, cmd UpdateUserCommand) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, q ListUsersQuery) ([]*domain.User, error)
}

// CreateUserCommand transporta los datos para crear un usuario. Password va
// en texto plano hasta que la aplicación lo hashea vía out.PasswordHasher.
type CreateUserCommand struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
	Role      string
}

// UpdateUserCommand transporta una actualización parcial: nil significa "no
// cambiar este campo". La aplicación hace el merge, no el adaptador HTTP.
type UpdateUserCommand struct {
	ID        uuid.UUID
	Email     *string
	FirstName *string
	LastName  *string
	Role      *string
	IsActive  *bool
}

// ListUsersQuery transporta los parámetros de paginación.
type ListUsersQuery struct {
	Limit  int
	Offset int
}
