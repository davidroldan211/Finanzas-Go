package auth

import (
	"finanzas-api/config"
	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/auth/handler"
	"finanzas-api/internal/auth/usecase"
	"finanzas-api/internal/middleware"
	userRepo "finanzas-api/internal/users/repository"
	"finanzas-api/shared/security"

	"gorm.io/gorm"
)

type AuthModule struct {
	Handler    *handler.AuthHandler
	UseCase    domain.AuthUseCase
	Middleware *middleware.Middleware
}

func NewAuthModule(db *gorm.DB, cfg *config.Config) *AuthModule {
	repo := userRepo.NewUserPostgresRepository(db)
	uc := usecase.NewAuthUseCase(repo, cfg.JWT)
	h := handler.NewAuthHandler(uc)
	mw := middleware.NewMiddleware(cfg.JWT.Secret, security.ParseToken)
	return &AuthModule{Handler: h, UseCase: uc, Middleware: mw}
}
