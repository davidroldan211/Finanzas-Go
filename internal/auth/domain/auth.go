package domain

// AuthUseCase defines authentication methods
type AuthUseCase interface {
	Login(email, password string) (string, error)
}
