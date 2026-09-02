package auth

import (
	authhttp "finanzas-api/internal/auth/adapter/in/http"
	"finanzas-api/internal/auth/adapter/out/postgres"
	"finanzas-api/internal/auth/application"
	"finanzas-api/internal/auth/port/in"
	"finanzas-api/internal/auth/port/out"
	"finanzas-api/internal/httpx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module es el composition root del módulo auth.
type Module struct {
	Service in.AuthService
	handler *authhttp.AuthHandler
	guard   *authhttp.Guard
}

// NewModule construye el módulo. verifier y tokens son puertos de salida
// compartidos (satisfechos hoy por *security.BcryptHasher y
// *security.HMACTokenProvider, construidos una sola vez en main). auth ya
// no importa config: el secreto y el TTL viven dentro de tokens.
func NewModule(db *gorm.DB, verifier out.PasswordVerifier, tokens out.TokenProvider) *Module {
	finder := postgres.NewUserCredentialsRepository(db)
	svc := application.NewAuthService(finder, verifier, tokens)

	return &Module{
		Service: svc,
		handler: authhttp.NewAuthHandler(svc),
		guard:   authhttp.NewGuard(tokens),
	}
}

// Guard expone el adaptador de entrada que otros módulos usan para
// proteger sus rutas. Esos módulos importan httpx.AuthGuard, no auth.
func (m *Module) Guard() httpx.AuthGuard {
	return m.guard.Handler
}

// RegisterRoutes registra las rutas del módulo auth.
func (m *Module) RegisterRoutes(router *gin.Engine) {
	authhttp.SetupAuthRoutes(router, m.handler)
}
