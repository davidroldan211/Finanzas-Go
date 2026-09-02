package application

import (
	"context"
	"strings"

	"finanzas-api/internal/verification/domain"
	"finanzas-api/internal/verification/port/in"
	"finanzas-api/internal/verification/port/out"
)

type verificationService struct {
	repo out.VerificationRepository
}

func NewVerificationService(repo out.VerificationRepository) in.VerificationService {
	return &verificationService{repo: repo}
}

// GenerateVerificationCode implements in.VerificationService.
//
// TODO: implementar generación de código, persistencia y envío (ver
// port/in/verification_service.go). repo no se usa todavía a propósito.
func (s *verificationService) GenerateVerificationCode(ctx context.Context, email string) (string, error) {
	if strings.TrimSpace(email) == "" {
		return "", domain.ErrEmailRequired
	}
	return "", nil
}
