package main

import (
	"finanzas-api/config"
	"finanzas-api/internal/auth"
	authRoutes "finanzas-api/internal/auth/routes"
	"finanzas-api/internal/users"
	userRoutes "finanzas-api/internal/users/routes"
	"finanzas-api/internal/verification"
	verificationRoutes "finanzas-api/internal/verification/routes"
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

	r = gin.Default()

	switch config.App.Environment {
	case "development":
		gin.SetMode(gin.DebugMode)
		log.Println("⚙️ Running in development mode")
		// r.SetTrustedProxies(nil)
	case "production":
		gin.SetMode(gin.ReleaseMode)
		log.Println("💯 Running in production mode")
		r.SetTrustedProxies([]string{"192.168.1.100"}) // Ejemplo de IP confiable
	case "test":
		gin.SetMode(gin.TestMode)
		log.Println("🛠️ Running in test mode")
		r.SetTrustedProxies(nil)
	default:
		gin.SetMode(gin.DebugMode)
		log.Println("Running in default (development) mode")
		r.SetTrustedProxies(nil)
	}

	db, err = DataBase.NewPostgresDB(config)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database: %v", err))

	}

	userModule := users.NewUsersModule(db)
	authModule := auth.NewAuthModule(db, config)
	verifyModule := verification.NewVerificationModule(db)

	authRoutes.SetupAuthRoutes(r, authModule.Handler)
	userRoutes.SetupUserRoutes(r, userModule.Handler, authModule.Middleware.Handler)
	verificationRoutes.SetupVerificationRoutes(r, verifyModule.Handler)

	log.Println("🚀 Servidor iniciado en " + config.Server.Host + ":" + strconv.Itoa(config.Server.Port))
	r.Run(config.Server.Host + ":" + strconv.Itoa(config.Server.Port))

}
