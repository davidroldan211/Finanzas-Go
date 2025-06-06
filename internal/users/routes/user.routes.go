package routes

import (
	"finanzas-api/internal/users/handler"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas para el módulo de usuarios
func SetupUserRoutes(router *gin.Engine, userHandler *handler.UserHandler, authMiddleware func(...string) gin.HandlerFunc) {
	// Grupo de rutas para usuarios
	userRoutes := router.Group("/api/v1/users")
	{
		// POST /api/v1/users - Crear usuario (solo admin)
		userRoutes.POST("", authMiddleware("admin"), userHandler.CreateUser)

		// GET /api/v1/users - Listar usuarios (solo admin)
		userRoutes.GET("", authMiddleware("admin"), userHandler.ListUsers)

		// GET /api/v1/users/:id - Obtener usuario por ID
		userRoutes.GET("/:id", authMiddleware("admin", "user"), userHandler.GetUser)

		// PUT /api/v1/users/:id - Actualizar usuario
		userRoutes.PUT("/:id", authMiddleware("admin", "user"), userHandler.UpdateUser)

		// DELETE /api/v1/users/:id - Eliminar usuario (solo admin)
		userRoutes.DELETE("/:id", authMiddleware("admin"), userHandler.DeleteUser)
	}
}

// O en el main.go o donde inicialices la aplicación:
func setupRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	// Rutas individuales
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.ListUsers)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
}
