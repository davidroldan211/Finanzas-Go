package domain

import "errors"

var (
	// ErrInvalidCredentials cubre tanto email desconocido como password
	// incorrecta: deliberadamente indistinguibles, defensa contra
	// enumeración de usuarios.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserInactive indica que las credenciales son correctas pero el
	// usuario no puede autenticarse (inactivo o borrado).
	ErrUserInactive = errors.New("user inactive")
	// ErrInvalidToken cubre token malformado, firma inválida o expirado.
	ErrInvalidToken = errors.New("invalid token")
)
