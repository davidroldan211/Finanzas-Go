package repository

import (
	"time"

	"finanzas-api/internal/users/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userModel representa la estructura de persistencia para Gorm.
type userModel struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Role      string    `gorm:"default:'user'"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (userModel) TableName() string { return "users" }

func toDomain(m *userModel) *domain.User {
	if m == nil {
		return nil
	}
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	return &domain.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Role:      m.Role,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func toModel(u *domain.User) *userModel {
	if u == nil {
		return nil
	}
	m := &userModel{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *u.DeletedAt, Valid: true}
	}
	return m
}
