package users

import (
	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/handler"
	"finanzas-api/internal/users/repository"
	"finanzas-api/internal/users/usecase"

	"gorm.io/gorm"
)

type UsersModule struct {
	Handler    *handler.UserHandler
	UseCase    domain.UserUseCase
	repository domain.UserRepository
}

func NewUsersModule(db *gorm.DB) *UsersModule {
	var userRepo domain.UserRepository
	var userUseCase domain.UserUseCase
	var userHandler *handler.UserHandler

	userRepo = repository.NewUserPostgresRepository(db)
	userUseCase = usecase.NewUserUseCase(userRepo)
	userHandler = handler.NewUserHandler(userUseCase)

	return &UsersModule{
		Handler:    userHandler,
		UseCase:    userUseCase,
		repository: userRepo,
	}
}
