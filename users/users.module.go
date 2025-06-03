package users

import (
	"finanzas-api/users/handler"
	"finanzas-api/users/repository"
	"finanzas-api/users/usecase"

	"github.com/gin-gonic/gin"
)

type UsersModule struct {
	Handler *handler.UserHandler
	UseCase *usecase.UserUseCase
	Repo    repository.UserRepository
}

func NewUsersModule(r *gin.Engine) *UsersModule {
	userRepo := repository.NewUserRepositoryMemory()
	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(r, userUC)

	return &UsersModule{
		Handler: userHandler,
		UseCase: userUC,
		Repo:    userRepo,
	}
}
