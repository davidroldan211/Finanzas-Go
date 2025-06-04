package main

import (
	"finanzas-api/config"
	"finanzas-api/internal/users"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	config, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %v", err))
	}

	userModule := users.NewUsersModule(r)
	_ = userModule.Handler
	_ = userModule.UseCase

	log.Println("Iniciando el servidor en " + config.Server.Host + ":" + strconv.Itoa(config.Server.Port))
	r.Run(config.Server.Host + ":" + strconv.Itoa(config.Server.Port)) // por defecto en :8080

}
