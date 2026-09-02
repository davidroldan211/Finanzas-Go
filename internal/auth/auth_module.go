package auth

import (
	"finanzas-api/config"
	authhttp "finanzas-api/internal/auth/adapter/in/http"
	"finanzas-api/internal/auth/application"
	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/shared/security"
	userRepo "finanzas-api/internal/users/adapter/out/postgres"

	"gorm.io/gorm"
)

type AuthModule struct {
	Handler    *authhttp.AuthHandler
	UseCase    domain.AuthUseCase
	Middleware *authhttp.Middleware
}

func NewAuthModule(db *gorm.DB, cfg *config.Config) *AuthModule {
	repo := userRepo.NewUserPostgresRepository(db)
	uc := application.NewAuthService(repo, cfg.JWT)
	h := authhttp.NewAuthHandler(uc)
	mw := authhttp.NewMiddleware(cfg.JWT.Secret, security.ParseToken)
	return &AuthModule{Handler: h, UseCase: uc, Middleware: mw}
}
