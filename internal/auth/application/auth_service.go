package application

import (
	"context"

	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/auth/port/in"
	"finanzas-api/internal/auth/port/out"
)

type authService struct {
	users    out.UserFinder
	verifier out.PasswordVerifier
	tokens   out.TokenProvider
}

// NewAuthService construye el puerto de entrada in.AuthService. Depende
// solo de puertos de salida: ni config ni bcrypt ni el HMAC aparecen aquí.
func NewAuthService(users out.UserFinder, verifier out.PasswordVerifier, tokens out.TokenProvider) in.AuthService {
	return &authService{users: users, verifier: verifier, tokens: tokens}
}

func (s *authService) Login(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
	creds, err := s.users.FindByEmail(ctx, cmd.Email)
	if err != nil {
		// Email desconocido y password incorrecta devuelven el mismo
		// error, indistinguible: defensa contra enumeración de usuarios.
		return domain.Token{}, domain.ErrInvalidCredentials
	}
	if !s.verifier.Matches(cmd.Password, creds.PasswordHash) {
		return domain.Token{}, domain.ErrInvalidCredentials
	}
	if !creds.Active {
		return domain.Token{}, domain.ErrUserInactive
	}

	return s.tokens.Issue(domain.Claims{Subject: creds.UserID, Role: creds.Role})
}
