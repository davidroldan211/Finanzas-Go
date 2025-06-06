package usecase

import (
	"errors"

	"finanzas-api/config"
	userDomain "finanzas-api/internal/users/domain"
	"finanzas-api/shared/security"
)

type AuthUseCase struct {
	userRepo  userDomain.UserRepository
	jwtConfig config.JWTConfig
}

func NewAuthUseCase(repo userDomain.UserRepository, cfg config.JWTConfig) *AuthUseCase {
	return &AuthUseCase{userRepo: repo, jwtConfig: cfg}
}

func (uc *AuthUseCase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !security.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}
	if !user.IsValidForAuth() {
		return "", errors.New("user inactive")
	}
	return security.GenerateToken(user.ID, user.Role, uc.jwtConfig.Secret, uc.jwtConfig.Expires)
}
