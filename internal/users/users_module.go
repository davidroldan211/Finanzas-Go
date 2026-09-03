package users

import (
	userhttp "finanzas-api/internal/users/adapter/in/http"
	"finanzas-api/internal/users/adapter/out/postgres"
	"finanzas-api/internal/users/application"
	"finanzas-api/internal/users/port/in"
	"finanzas-api/internal/users/port/out"

	"finanzas-api/internal/httpx"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module es el composition root del módulo users: arma el adaptador de
// salida, la aplicación y el adaptador de entrada.
type Module struct {
	Service in.UserService
	Handler *userhttp.UserHandler
}

// NewModule construye el módulo. hasher es un puerto de salida
// compartido (satisfecho hoy por *security.BcryptHasher desde main).
func NewModule(db *gorm.DB, hasher out.PasswordHasher) *Module {
	repo := postgres.NewUserPostgresRepository(db)
	svc := application.NewUserService(repo, hasher)
	handler := userhttp.NewUserHandler(svc)

	return &Module{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registra las rutas del módulo users.
func (m *Module) RegisterRoutes(router *gin.Engine, authGuard httpx.AuthGuard) {
	userhttp.SetupUserRoutes(router, m.Handler, authGuard)
}
