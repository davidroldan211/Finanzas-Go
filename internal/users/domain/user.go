package domain

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User es la entidad de dominio pura: sin tags de transporte (json) ni de
// persistencia (gorm). Nunca se serializa directamente ni se mapea 1:1 a
// una tabla; eso es responsabilidad de los adaptadores de entrada y salida.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time // soft delete
	// EmailVerifiedAt es nil mientras el correo no se haya confirmado.
	// NewUser siempre lo deja en nil; solo la confirmación (HU-02) lo llena.
	EmailVerifiedAt *time.Time
}

// NewUser construye un User validando las invariantes de negocio.
// passwordHash debe venir ya hasheado: el hashing es responsabilidad de la
// aplicación a través del puerto PasswordHasher, nunca del dominio.
func NewUser(email, firstName, lastName, passwordHash, role string) (*User, error) {
	email = normalizeEmail(email)
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if role == "" {
		role = RoleUser
	}

	if fields := validateUserFields(email, firstName, lastName); len(fields) > 0 {
		return nil, &ValidationError{Fields: fields}
	}

	return &User{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     true,
	}, nil
}

// ApplyUpdate aplica cambios parciales (nil = no cambiar) validando las
// invariantes del resultado antes de mutar el usuario.
func (u *User) ApplyUpdate(email, firstName, lastName, role *string, isActive *bool) error {
	newEmail := u.Email
	if email != nil {
		newEmail = normalizeEmail(*email)
	}
	newFirstName := u.FirstName
	if firstName != nil {
		newFirstName = strings.TrimSpace(*firstName)
	}
	newLastName := u.LastName
	if lastName != nil {
		newLastName = strings.TrimSpace(*lastName)
	}

	fields := validateUserFields(newEmail, newFirstName, newLastName)
	if role != nil && !IsValidRole(*role) {
		fields["role"] = "rol inválido"
	}
	if len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}

	u.Email = newEmail
	u.FirstName = newFirstName
	u.LastName = newLastName
	if role != nil {
		u.Role = *role
	}
	if isActive != nil {
		u.IsActive = *isActive
	}
	return nil
}

// IsValidRole indica si role es uno de los roles que existen.
func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleUser
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// emailMaxLen es el máximo práctico de una dirección (RFC 5321: 254).
const emailMaxLen = 254

// isValidEmail exige una dirección simple ya normalizada: que net/mail la
// acepte, que no traiga nombre visible ni <> (lo parseado debe ser igual a
// la entrada) y que el dominio tenga al menos un punto.
func isValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return false
	}
	at := strings.LastIndex(email, "@")
	return at > 0 && strings.Contains(email[at+1:], ".")
}

func validateUserFields(email, firstName, lastName string) map[string]string {
	fields := map[string]string{}

	switch {
	case email == "":
		fields["email"] = "el email es obligatorio"
	case len(email) > emailMaxLen:
		fields["email"] = "el email es demasiado largo"
	case !isValidEmail(email):
		fields["email"] = "formato de email inválido"
	}

	switch {
	case firstName == "":
		fields["first_name"] = "el nombre es obligatorio"
	case len(firstName) > 100:
		fields["first_name"] = "el nombre es demasiado largo"
	}

	switch {
	case lastName == "":
		fields["last_name"] = "el apellido es obligatorio"
	case len(lastName) > 100:
		fields["last_name"] = "el apellido es demasiado largo"
	}

	return fields
}

// GetFullName retorna el nombre completo del usuario.
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsValidForAuth verifica si el usuario puede autenticarse.
func (u *User) IsValidForAuth() bool {
	return u.IsActive && u.DeletedAt == nil
}
