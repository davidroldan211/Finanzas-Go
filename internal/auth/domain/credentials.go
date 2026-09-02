package domain

// Credentials es la proyección de solo lectura que auth necesita para
// autenticar a un usuario. No es la entidad users/domain.User: es una
// vista propia de auth sobre la misma tabla, de 5 columnas.
type Credentials struct {
	UserID       string
	Email        string
	PasswordHash string
	Role         string
	Active       bool
}
