package application

import (
	"context"
	"errors"
	"testing"

	"finanzas-api/internal/verification/domain"
)

func TestGenerateVerificationCode_EmptyEmail(t *testing.T) {
	svc := NewVerificationService(nil)

	_, err := svc.GenerateVerificationCode(context.Background(), "   ")
	if !errors.Is(err, domain.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestGenerateVerificationCode_ValidEmail(t *testing.T) {
	svc := NewVerificationService(nil)

	code, err := svc.GenerateVerificationCode(context.Background(), "a@b.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "" {
		t.Fatalf("expected empty code (feature not implemented yet), got %q", code)
	}
}
