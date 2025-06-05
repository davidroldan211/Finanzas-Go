package routes

import (
	"finanzas-api/internal/users/handler"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas para el módulo de usuarios
func SetupUserRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	// Grupo de rutas para usuarios
	userRoutes := router.Group("/api/v1/users")
	{
		// POST /api/v1/users - Crear usuario
		userRoutes.POST("", userHandler.CreateUser)

		// GET /api/v1/users - Listar usuarios
		userRoutes.GET("", userHandler.ListUsers)

		// GET /api/v1/users/:id - Obtener usuario por ID
		userRoutes.GET("/:id", userHandler.GetUser)

		// PUT /api/v1/users/:id - Actualizar usuario
		userRoutes.PUT("/:id", userHandler.UpdateUser)

		// DELETE /api/v1/users/:id - Eliminar usuario
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
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
