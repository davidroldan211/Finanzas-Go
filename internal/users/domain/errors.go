package domain

import "errors"

var (
	// ErrUserNotFound indica que no existe un usuario con el ID o email dados.
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailTaken indica que el email ya está en uso por otro usuario.
	ErrEmailTaken = errors.New("email already in use")
	// ErrInvalidUserID indica que el ID de usuario es inválido (nulo/vacío).
	ErrInvalidUserID = errors.New("invalid user id")
)

// ValidationError agrupa errores de campo. Los adaptadores de entrada la
// traducen a una respuesta 422 con el mapa de campos completo.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string { return "invalid user data" }
