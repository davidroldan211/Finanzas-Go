package routes

import (
	"finanzas-api/internal/users/handler"

	"github.com/gin-gonic/gin"
)

func setupVerificationRoutes(router *gin.Engine, verificationHandler *handler.VerificationHandler) {
	verificationRoutes := router.Group("/api/v1/users")
	{
		verificationRoutes.POST("/verify-email", verificationHandler.VerifyEmail)
	}

}
