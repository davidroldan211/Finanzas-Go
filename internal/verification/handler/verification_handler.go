package handler

import (
	"net/http"

	"finanzas-api/internal/verification/domain"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	uc domain.VerificationUseCase
}

func NewVerificationHandler(uc domain.VerificationUseCase) *VerificationHandler {
	return &VerificationHandler{uc: uc}
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if _, err := h.uc.GenerateVerificationCode(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}
