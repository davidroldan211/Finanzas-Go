package routes

import (
	"finanzas-api/internal/auth/handler"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, h *handler.AuthHandler) {
	router.POST("/api/v1/login", h.Login)
}
