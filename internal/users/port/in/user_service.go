package in

import (
	"context"

	"finanzas-api/internal/users/domain"

	"github.com/google/uuid"
)

// UserService es el puerto de entrada del módulo users: lo que la
// aplicación ofrece a sus adaptadores de entrada (p.ej. adapter/in/http).
type UserService interface {
	// Create es el registro público: siempre crea un usuario con rol user.
	Create(ctx context.Context, cmd CreateUserCommand) (*domain.User, error)
	// GetByID solo devuelve el usuario si actor es ese usuario o un admin.
	GetByID(ctx context.Context, actor Actor, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	// Update solo permite actualizar al propio usuario, o a cualquiera si
	// actor es admin. Role e IsActive solo los puede cambiar un admin.
	Update(ctx context.Context, actor Actor, cmd UpdateUserCommand) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, q ListUsersQuery) ([]*domain.User, error)
}

// Actor es quien ejecuta la operación, tomado de un token ya verificado
// por el adaptador de entrada.
type Actor struct {
	ID   uuid.UUID
	Role string
}

// IsAdmin indica si el actor tiene rol de administrador.
func (a Actor) IsAdmin() bool { return a.Role == domain.RoleAdmin }

// CanAccess indica si el actor puede ver o modificar al usuario userID.
func (a Actor) CanAccess(userID uuid.UUID) bool {
	return a.IsAdmin() || (a.ID != uuid.Nil && a.ID == userID)
}

// CreateUserCommand transporta los datos del registro. Password va en texto
// plano hasta que la aplicación lo hashea vía out.PasswordHasher. No lleva
// rol a propósito: el registro público nunca elige su rol.
type CreateUserCommand struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
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
