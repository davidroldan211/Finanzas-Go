package http

import (
	"errors"
	nethttp "net/http"

	"finanzas-api/internal/httpx"
	"finanzas-api/internal/verification/domain"
	"finanzas-api/internal/verification/port/in"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	service in.VerificationService
}

func NewVerificationHandler(service in.VerificationService) *VerificationHandler {
	return &VerificationHandler{service: service}
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	if _, err := h.service.GenerateVerificationCode(c.Request.Context(), req.Email); err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"message": "Verification email sent"})
}

func toAppError(err error) *httpx.AppError {
	if errors.Is(err, domain.ErrEmailRequired) {
		return httpx.BadRequest("El email es obligatorio.").WithErr(err)
	}
	return httpx.Wrap(err)
}
