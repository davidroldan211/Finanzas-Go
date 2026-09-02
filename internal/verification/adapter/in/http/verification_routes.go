package http

import (
	"github.com/gin-gonic/gin"
)

func SetupVerificationRoutes(router *gin.Engine, verificationHandler *VerificationHandler) {
	verificationRoutes := router.Group("/api/v1/verification")
	{
		verificationRoutes.POST("/verify-email", verificationHandler.VerifyEmail)
	}
}
