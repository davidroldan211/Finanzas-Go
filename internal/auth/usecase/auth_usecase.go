package usecase

import (
	"context"
	"errors"

	"finanzas-api/config"
	userOut "finanzas-api/internal/users/port/out"
	"finanzas-api/shared/security"
)

// NOTA: auth reutiliza aquí el puerto de salida de users (userOut.UserRepository)
// como solución interina. La fase 2 lo sustituye por un puerto propio de auth
// (out.UserFinder), read-only y sin acoplar los dos módulos.
type AuthService struct {
	userRepo  userOut.UserRepository
	jwtConfig config.JWTConfig
}

func NewAuthService(repo userOut.UserRepository, cfg config.JWTConfig) *AuthService {
	return &AuthService{userRepo: repo, jwtConfig: cfg}
}

func (uc *AuthService) Login(email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(context.Background(), email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	if !user.IsValidForAuth() {
		return "", errors.New("user inactive")
	}
	return security.GenerateToken(user.ID, user.Role, uc.jwtConfig.Secret, uc.jwtConfig.Expires)
}
