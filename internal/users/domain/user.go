package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"` // El "-" oculta la contraseña en JSON
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"` // Soft delete
}

// UserRepository define la interfaz del repositorio de usuarios
type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*User, error)
	EmailExists(email string) (bool, error)
}

// UserUseCase define las operaciones de aplicación para usuarios.
type UserUseCase interface {
	CreateUser(user *User) error
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uuid.UUID) error
	ListUsers(limit, offset int) ([]*User, error)
	ValidateUserData(user *User) error
}

// TableName especifica el nombre de la tabla en la base de datos
func (User) TableName() string {
	return "users"
}

// GetFullName retorna el nombre completo del usuario
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsValidForAuth verifica si el usuario puede autenticarse
func (u *User) IsValidForAuth() bool {
	return u.IsActive && (u.DeletedAt == nil || u.DeletedAt.IsZero())
}

// ValidateEmail verifica si el email tiene un formato válido
func (u *User) ValidateEmail() bool {
	//TODO: Implementar una validación de email más robusta
	return len(u.Email) > 0 && len(u.Email) <= 255
}

// ValidateNames verifica si los nombres son válidos
func (u *User) ValidateNames() bool {
	//TODO: Implementar una validación de nombres más robusta
	return len(u.FirstName) > 0 && len(u.LastName) > 0
}
