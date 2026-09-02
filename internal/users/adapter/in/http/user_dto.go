package http

import (
	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/port/in"

	"github.com/google/uuid"
)

// CreateUserRequest representa la estructura de la petición para crear usuario.
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"omitempty,oneof=admin user"`
}

func (r CreateUserRequest) toCommand() in.CreateUserCommand {
	return in.CreateUserCommand{
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Password:  r.Password,
		Role:      r.Role,
	}
}

// UpdateUserRequest representa la estructura de la petición para actualizar
// usuario. Un campo ausente en el JSON queda nil y no se toca.
type UpdateUserRequest struct {
	Email     *string `json:"email" binding:"omitempty,email"`
	FirstName *string `json:"first_name" binding:"omitempty"`
	LastName  *string `json:"last_name" binding:"omitempty"`
	IsActive  *bool   `json:"is_active" binding:"omitempty"`
	Role      *string `json:"role" binding:"omitempty,oneof=admin user"`
}

func (r UpdateUserRequest) toCommand(id uuid.UUID) in.UpdateUserCommand {
	return in.UpdateUserCommand{
		ID:        id,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Role:      r.Role,
		IsActive:  r.IsActive,
	}
}

// UserResponse representa la respuesta de usuario (sin contraseña).
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		FullName:  user.GetFullName(),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
