package main

import (
	"finanzas-api/config"
	"finanzas-api/internal/auth"
	DataBase "finanzas-api/internal/shared/db"
	"finanzas-api/internal/shared/security"
	"finanzas-api/internal/users"
	"finanzas-api/internal/verification"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {

	var r *gin.Engine
	var db *gorm.DB

	config, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %v", err))
	}

	// El modo de Gin se fija ANTES de gin.Default(): si se fija después, el
	// logger/recovery ya se construyeron en modo debug.
	switch config.App.Environment {
	case "development":
		gin.SetMode(gin.DebugMode)
		log.Println("⚙️ Running in development mode")
	case "production":
		gin.SetMode(gin.ReleaseMode)
		log.Println("💯 Running in production mode")
	case "test":
		gin.SetMode(gin.TestMode)
		log.Println("🛠️ Running in test mode")
	default:
		gin.SetMode(gin.DebugMode)
		log.Println("Running in default (development) mode")
	}

	r = gin.Default()

	if config.App.Environment == "production" {
		r.SetTrustedProxies(config.Server.TrustedProxies)
	} else {
		r.SetTrustedProxies(nil)
	}

	db, err = DataBase.NewPostgresDB(config)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database: %v", err))

	}

	hasher := security.NewBcryptHasher(bcrypt.DefaultCost)
	tokens := security.NewHMACTokenProvider(config.JWT.Secret, config.JWT.Expires)

	userModule := users.NewModule(db, hasher)
	authModule := auth.NewModule(db, hasher, tokens)
	verifyModule := verification.NewVerificationModule(db)

	authModule.RegisterRoutes(r)
	userModule.RegisterRoutes(r, authModule.Guard())
	verifyModule.RegisterRoutes(r)

	log.Println("🚀 Servidor iniciado en " + config.Server.Host + ":" + strconv.Itoa(config.Server.Port))
	r.Run(config.Server.Host + ":" + strconv.Itoa(config.Server.Port))

}
