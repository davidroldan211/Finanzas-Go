package users

import (
	userhttp "finanzas-api/internal/users/adapter/in/http"
	"finanzas-api/internal/users/adapter/out/postgres"
	"finanzas-api/internal/users/application"
	"finanzas-api/internal/users/domain"

	"gorm.io/gorm"
)

type UsersModule struct {
	Handler    *userhttp.UserHandler
	UseCase    domain.UserUseCase
	repository domain.UserRepository
}

func NewUsersModule(db *gorm.DB) *UsersModule {
	var userRepo domain.UserRepository
	var userUseCase domain.UserUseCase
	var userHandler *userhttp.UserHandler

	userRepo = postgres.NewUserPostgresRepository(db)
	userUseCase = application.NewUserUseCase(userRepo)
	userHandler = userhttp.NewUserHandler(userUseCase)

	return &UsersModule{
		Handler:    userHandler,
		UseCase:    userUseCase,
		repository: userRepo,
	}
}
