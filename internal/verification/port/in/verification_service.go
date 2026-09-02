package in

import "context"

// VerificationService es el puerto de entrada del módulo verification.
//
// TODO: hoy solo valida el email y no persiste ni envía nada
// (GenerateVerificationCode devuelve "", nil). Falta CleanupExpiredCodes,
// ValidateVerificationCode y SendVerificationCode — quedan fuera de esta
// migración a hexagonal a propósito: son funcionalidad nueva (elegir
// proveedor de email, formato de código, rate limiting, política de
// cooldown), no una reestructuración arquitectónica. El hueco queda con
// la forma correcta (puertos out.CodeGenerator, out.EmailSender) para
// implementarse en un cambio aparte.
type VerificationService interface {
	GenerateVerificationCode(ctx context.Context, email string) (string, error)
}
