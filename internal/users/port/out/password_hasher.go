package out

// PasswordHasher es el puerto de salida que usa la aplicación para hashear
// contraseñas antes de persistirlas. Users nunca verifica contraseñas
// (eso es responsabilidad de auth), solo las hashea: un único método.
type PasswordHasher interface {
	Hash(plain string) (string, error)
}
