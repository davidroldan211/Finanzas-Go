package main

import (
	"finanzas-api/config"
	"finanzas-api/internal/auth"
	authRoutes "finanzas-api/internal/auth/routes"
	"finanzas-api/internal/users"
	userRoutes "finanzas-api/internal/users/routes"
	DataBase "finanzas-api/shared/db"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {

	var r *gin.Engine
	var db *gorm.DB

	config, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %v", err))
	}

	switch config.App.Environment {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	r = gin.Default()

	db, err = DataBase.NewPostgresDB(config)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database: %v", err))

	}

	userModule := users.NewUsersModule(db)
	authModule := auth.NewAuthModule(db, config)

	authRoutes.SetupAuthRoutes(r, authModule.Handler)
	userRoutes.SetupUserRoutes(r, userModule.Handler, authModule.Middleware.Handler)

	log.Println("🚀 Servidor iniciado en " + config.Server.Host + ":" + strconv.Itoa(config.Server.Port))
	r.Run(config.Server.Host + ":" + strconv.Itoa(config.Server.Port))

}
