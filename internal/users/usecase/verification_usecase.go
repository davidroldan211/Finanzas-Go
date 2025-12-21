package usecase

import (
	"finanzas-api/internal/users/domain"

	"github.com/google/uuid"
)

type verificationUseCase struct {
	verificationRepo domain.VerificationRepository
}

func NewVerificationUseCase(verificationRepo domain.VerificationRepository) domain.VerificationUseCase {
	return &verificationUseCase{
		verificationRepo: verificationRepo,
	}
}

// CleanupExpiredCodes implements domain.VerificationUseCase.
func (uc *verificationUseCase) CleanupExpiredCodes() error {
	panic("unimplemented")
}

// GenerateVerificationCode implements domain.VerificationUseCase.
func (uc *verificationUseCase) GenerateVerificationCode(email string) (string, error) {
	panic("unimplemented")
}

// ValidateVerificationCode implements domain.VerificationUseCase.
func (uc *verificationUseCase) ValidateVerificationCode(ID uuid.UUID, email string, code string) (bool, error) {
	panic("unimplemented")
}

func (uc *verificationUseCase) SendVerificationCode(ID uuid.UUID, email string) error {
	panic("unimplemented")
}
