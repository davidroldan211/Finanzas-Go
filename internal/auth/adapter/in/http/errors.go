package http

import (
	"errors"

	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/httpx"
)

// toAppError traduce errores de dominio a códigos HTTP. Un error no
// reconocido cae en httpx.Wrap: 500 genérico, sin filtrar el mensaje
// interno al cliente.
func toAppError(err error) *httpx.AppError {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return httpx.Unauthorized("Credenciales inválidas.").WithErr(err)
	case errors.Is(err, domain.ErrUserInactive):
		return httpx.Unauthorized("El usuario está inactivo.").WithErr(err)
	default:
		return httpx.Wrap(err)
	}
}
