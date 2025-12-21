package handler

import (
	"finanzas-api/internal/users/domain"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	verificationUseCase domain.VerificationUseCase
}

func NewVerificationHandler(verificationUseCase domain.VerificationUseCase) *VerificationHandler {
	return &VerificationHandler{
		verificationUseCase: verificationUseCase,
	}
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Lógica para manejar la verificación del correo electrónico
	c.JSON(200, gin.H{"message": "Verification email sent"})
}
