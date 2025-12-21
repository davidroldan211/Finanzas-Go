package usecase

import (
	"finanzas-api/internal/verification/domain"
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
	// TODO: implementar generación de código, persistencia y envío.
	return "", nil
}

// // CleanupExpiredCodes implements domain.UseCase.
// func (uc *verificationUseCase) CleanupExpiredCodes() error {
// 	panic("unimplemented")
// }

// // ValidateVerificationCode implements domain.UseCase.
// func (uc *verificationUseCase) ValidateVerificationCode(ID uuid.UUID, email string, code string) (bool, error) {
// 	panic("unimplemented")
// }

// func (uc *verificationUseCase) SendVerificationCode(ID uuid.UUID, email string) error {
// 	panic("unimplemented")
// }
