package out

import (
	"context"

	"finanzas-api/internal/verification/domain"

	"github.com/google/uuid"
)

// VerificationRepository es el puerto de salida para persistir códigos de
// verificación de email. Aún sin implementación de negocio (ver
// application/verification_service.go): existe para que el adaptador de
// salida tenga la forma correcta cuando se implemente la funcionalidad.
type VerificationRepository interface {
	Create(ctx context.Context, v *domain.Verification) error
	GetByEmail(ctx context.Context, email string) (*domain.Verification, error)
	Update(ctx context.Context, v *domain.Verification) error
	Delete(ctx context.Context, id uuid.UUID) error
}
