package routes

import (
	"github.com/gin-gonic/gin"

	"finanzas-api/internal/verification/handler"
)

func SetupVerificationRoutes(router *gin.Engine, verificationHandler *handler.VerificationHandler) {
	verificationRoutes := router.Group("/api/v1/verification")
	{
		verificationRoutes.POST("/verify-email", verificationHandler.VerifyEmail)
	}
}
