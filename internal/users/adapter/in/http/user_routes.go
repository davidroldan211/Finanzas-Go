package http

import (
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas para el módulo de usuarios
func SetupUserRoutes(router *gin.Engine, userHandler *UserHandler, authMiddleware func(...string) gin.HandlerFunc) {
	// Grupo de rutas para usuarios
	userRoutes := router.Group("/api/v1/users")
	{
		// POST /api/v1/users - Crear usuario (solo admin)
		userRoutes.POST("", userHandler.CreateUser)

		// GET /api/v1/users - Listar usuarios (solo admin)
		userRoutes.GET("", authMiddleware("admin"), userHandler.ListUsers)

		// GET /api/v1/users/:id - Obtener usuario por ID
		userRoutes.GET("/:id", authMiddleware("user", "admin"), userHandler.GetUser)

		// PUT /api/v1/users/:id - Actualizar usuario
		userRoutes.PUT("/:id", authMiddleware("user", "admin"), userHandler.UpdateUser)

		// DELETE /api/v1/users/:id - Eliminar usuario (solo admin)
		userRoutes.DELETE("/:id", authMiddleware("admin"), userHandler.DeleteUser)
	}
}
