package verification

import (
	verificationhttp "finanzas-api/internal/verification/adapter/in/http"
	"finanzas-api/internal/verification/adapter/out/postgres"
	"finanzas-api/internal/verification/application"
	"finanzas-api/internal/verification/port/in"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module es el composition root del módulo verification.
type Module struct {
	Service in.VerificationService
	Handler *verificationhttp.VerificationHandler
}

func NewVerificationModule(db *gorm.DB) *Module {
	repo := postgres.NewVerificationPostgresRepository(db)
	svc := application.NewVerificationService(repo)
	handler := verificationhttp.NewVerificationHandler(svc)

	return &Module{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registra las rutas del módulo verification.
func (m *Module) RegisterRoutes(router *gin.Engine) {
	verificationhttp.SetupVerificationRoutes(router, m.Handler)
}
