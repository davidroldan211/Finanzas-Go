package postgres

import (
	"testing"
	"time"

	"finanzas-api/internal/users/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestToDomain_Nil(t *testing.T) {
	if got := toDomain(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestToModel_Nil(t *testing.T) {
	if got := toModel(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestToDomain_ToModel_RoundTrip(t *testing.T) {
	id := uuid.New()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	m := &userModel{
		ID:        id,
		Email:     "a@b.com",
		Password:  "hashed-password",
		FirstName: "Ana",
		LastName:  "Pérez",
		Role:      "admin",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	u := toDomain(m)
	if u.ID != id {
		t.Fatalf("expected ID %v, got %v", id, u.ID)
	}
	if u.PasswordHash != "hashed-password" {
		t.Fatalf("expected PasswordHash to carry the model's Password, got %q", u.PasswordHash)
	}
	if u.Email != m.Email || u.FirstName != m.FirstName || u.LastName != m.LastName || u.Role != m.Role {
		t.Fatalf("field mismatch after toDomain: %+v", u)
	}
	if u.DeletedAt != nil {
		t.Fatalf("expected DeletedAt nil for a non-deleted model, got %v", u.DeletedAt)
	}

	back := toModel(u)
	if back.ID != id {
		t.Fatalf("expected round-tripped uuid.UUID ID %v, got %v", id, back.ID)
	}
	if back.Password != "hashed-password" {
		t.Fatalf("expected Password to carry PasswordHash back, got %q", back.Password)
	}
	if back.DeletedAt.Valid {
		t.Fatal("expected DeletedAt.Valid false when domain.DeletedAt is nil")
	}
}

func TestToModel_DeletedAt_MapsToGormDeletedAt(t *testing.T) {
	deletedAt := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)
	u := &domain.User{ID: uuid.New(), Email: "a@b.com", DeletedAt: &deletedAt}

	m := toModel(u)

	if !m.DeletedAt.Valid {
		t.Fatal("expected DeletedAt.Valid true when domain.DeletedAt is set")
	}
	if !m.DeletedAt.Time.Equal(deletedAt) {
		t.Fatalf("expected DeletedAt.Time %v, got %v", deletedAt, m.DeletedAt.Time)
	}
}

func TestToDomain_DeletedAt_MapsToPointer(t *testing.T) {
	deletedAt := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)
	m := &userModel{ID: uuid.New(), DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true}}

	u := toDomain(m)

	if u.DeletedAt == nil {
		t.Fatal("expected DeletedAt non-nil when model.DeletedAt.Valid is true")
	}
	if !u.DeletedAt.Equal(deletedAt) {
		t.Fatalf("expected DeletedAt %v, got %v", deletedAt, *u.DeletedAt)
	}
}
