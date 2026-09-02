// Package security contiene los adaptadores de salida compartidos que
// implementan los puertos de seguridad de los distintos módulos (hashing de
// contraseñas, emisión/verificación de tokens). No es dominio ni aplicación:
// es infraestructura, y por eso vive bajo internal/shared.
package security

import "golang.org/x/crypto/bcrypt"

// BcryptHasher implementa, con el mismo tipo, tanto users/port/out.PasswordHasher
// (un solo método: Hash) como auth/port/out.PasswordVerifier (Matches) — el
// puerto lo define quien lo consume, y ambos puertos son satisfechos sin que
// los módulos users y auth se conozcan entre sí.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher construye un BcryptHasher con el costo indicado.
func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost: cost}
}

// Hash satisface users/port/out.PasswordHasher.
func (h *BcryptHasher) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Matches satisface auth/port/out.PasswordVerifier.
func (h *BcryptHasher) Matches(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
