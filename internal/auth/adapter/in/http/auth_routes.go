package http

import (
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, h *AuthHandler) {
	router.POST("/api/v1/login", h.Login)
}
