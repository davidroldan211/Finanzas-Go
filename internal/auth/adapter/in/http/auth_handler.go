package http

import (
	nethttp "net/http"

	"finanzas-api/internal/auth/port/in"
	"finanzas-api/internal/httpx"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service in.AuthService
}

func NewAuthHandler(service in.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	token, err := h.service.Login(c.Request.Context(), in.LoginCommand{Email: req.Email, Password: req.Password})
	if err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusOK, LoginResponse{Token: token.Value})
}
