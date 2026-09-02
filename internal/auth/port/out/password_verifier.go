package out

// PasswordVerifier es el puerto de salida que usa la aplicación para
// verificar una contraseña contra su hash. Auth nunca hashea, solo
// verifica: un único método (users declara su propio PasswordHasher con
// Hash; ambos puertos los satisface el mismo *security.BcryptHasher).
type PasswordVerifier interface {
	Matches(plain, hash string) bool
}
