package verification

import (
	"finanzas-api/internal/verification/domain"
	"finanzas-api/internal/verification/handler"
	"finanzas-api/internal/verification/repository"
	"finanzas-api/internal/verification/usecase"

	"gorm.io/gorm"
)

type verificationModule struct {
	Handler    *handler.VerificationHandler
	UseCase    domain.VerificationUseCase
	repository domain.VerificationRepository
}

func NewVerificationModule(db *gorm.DB) *verificationModule {
	repo := repository.NewVerificationPostgresRepository(db)
	uc := usecase.NewVerificationUseCase(repo)
	h := handler.NewVerificationHandler(uc)
	return &verificationModule{
		Handler:    h,
		UseCase:    uc,
		repository: repo,
	}
}
