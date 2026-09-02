package out

import (
	"context"

	"finanzas-api/internal/auth/domain"
)

// UserFinder es el puerto de salida que auth usa para leer credenciales.
// NO es el repositorio de users: auth solo necesita esta proyección
// read-only, implementada por su propio adaptador.
type UserFinder interface {
	FindByEmail(ctx context.Context, email string) (*domain.Credentials, error)
}
