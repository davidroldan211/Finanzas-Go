package http

import (
	"errors"

	"finanzas-api/internal/httpx"
	"finanzas-api/internal/users/domain"
)

// toAppError traduce errores de dominio a códigos HTTP concretos. Un error
// no reconocido cae en httpx.Wrap: 500 genérico, sin filtrar el mensaje
// interno al cliente.
func toAppError(err error) *httpx.AppError {
	var ve *domain.ValidationError

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return httpx.NotFound("Usuario no encontrado.").WithErr(err)
	case errors.Is(err, domain.ErrEmailTaken):
		return httpx.Conflict("El email ya está registrado.").WithErr(err)
	case errors.Is(err, domain.ErrInvalidUserID):
		return httpx.BadRequest("ID de usuario inválido.").WithErr(err)
	case errors.As(err, &ve):
		return httpx.Validation(ve.Fields).WithErr(err)
	default:
		return httpx.Wrap(err)
	}
}
