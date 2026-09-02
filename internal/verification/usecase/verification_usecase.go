package usecase

import (
	"errors"
	"finanzas-api/internal/verification/domain"
	"strings"
)

type verificationUseCase struct {
	repo domain.VerificationRepository
}

func NewVerificationUseCase(repo domain.VerificationRepository) domain.VerificationUseCase {
	return &verificationUseCase{
		repo: repo,
	}
}

// GenerateVerificationCode implements domain.UseCase.
func (uc *verificationUseCase) GenerateVerificationCode(email string) (string, error) {
	if strings.TrimSpace(email) == "" {
		return "", errors.New("email is required")
	}
	return "", nil
}

//TODO: CleanupExpiredCodes,ValidateVerificationCode,SendVerificationCode
