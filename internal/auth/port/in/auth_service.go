package in

import (
	"context"

	"finanzas-api/internal/auth/domain"
)

// AuthService es el puerto de entrada del módulo auth.
type AuthService interface {
	Login(ctx context.Context, cmd LoginCommand) (domain.Token, error)
}

type LoginCommand struct {
	Email    string
	Password string
}
