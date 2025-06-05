package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // El "-" oculta la contraseña en JSON
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete
}

// UserRepository define la interfaz del repositorio de usuarios
type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(limit, offset int) ([]*User, error)
	EmailExists(email string) (bool, error)
}

type UserUseCase interface {
	CreateUser(user *User) error
	GetUserByID(id uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
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
	return u.IsActive && u.DeletedAt.Time.IsZero()
}

// ValidateEmail verifica si el email tiene un formato válido
func (u *User) ValidateEmail() bool {
	// Implementación básica de validación de email
	return len(u.Email) > 0 && len(u.Email) <= 255
}

// ValidateNames verifica si los nombres son válidos
func (u *User) ValidateNames() bool {
	return len(u.FirstName) > 0 && len(u.LastName) > 0
}
