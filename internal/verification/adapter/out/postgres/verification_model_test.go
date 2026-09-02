package postgres

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestVerificationModel_TableName(t *testing.T) {
	// Regresión: el dominio antes de esta migración declaraba
	// TableName() = "verifications", pero db/migrations/verifications.sql
	// crea "email_verifications". Cualquier query explotaba con
	// "relation does not exist".
	if got := (verificationModel{}).TableName(); got != "email_verifications" {
		t.Fatalf("expected table name 'email_verifications', got %q", got)
	}
}

func TestVerificationToDomain_Nil(t *testing.T) {
	if got := toDomain(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestVerificationToModel_Nil(t *testing.T) {
	if got := toModel(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestVerificationToDomain_ToModel_RoundTrip(t *testing.T) {
	id := uuid.New()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cooldown := now.Add(time.Minute)
	m := &verificationModel{
		ID:            id,
		Email:         "a@b.com",
		CodeHash:      "hash",
		ExpiresAt:     now.Add(time.Hour),
		Attempts:      1,
		MaxAttempts:   5,
		CooldownUntil: &cooldown,
		UsedAt:        nil,
		CreatedAt:     now,
	}

	v := toDomain(m)
	if v.ID != id || v.Email != m.Email || v.CodeHash != m.CodeHash {
		t.Fatalf("field mismatch after toDomain: %+v", v)
	}
	if v.CooldownUntil == nil || !v.CooldownUntil.Equal(cooldown) {
		t.Fatalf("expected CooldownUntil %v, got %v", cooldown, v.CooldownUntil)
	}
	if v.UsedAt != nil {
		t.Fatalf("expected UsedAt nil, got %v", v.UsedAt)
	}

	back := toModel(v)
	if back.ID != id {
		t.Fatalf("expected round-tripped ID %v, got %v", id, back.ID)
	}
	if back.CooldownUntil == nil || !back.CooldownUntil.Equal(cooldown) {
		t.Fatalf("expected round-tripped CooldownUntil %v, got %v", cooldown, back.CooldownUntil)
	}
}
