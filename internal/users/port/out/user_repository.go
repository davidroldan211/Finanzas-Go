package out

import (
	"context"

	"finanzas-api/internal/users/domain"

	"github.com/google/uuid"
)

// UserRepository es el puerto de salida que la aplicación usa para
// persistir y consultar usuarios. Lo implementa un adaptador de salida
// (p.ej. adapter/out/postgres).
type UserRepository interface {
	Save(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
