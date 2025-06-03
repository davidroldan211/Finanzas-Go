package main

import (
	"finanzas-api/config"
	"finanzas-api/users"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.GetEnv()

	userModule := users.NewUsersModule(r)
	_ = userModule.Handler
	_ = userModule.UseCase
	fmt.Println("Iniciando el servidor en " + config.Envs.AppHost + ":" + strconv.Itoa(config.Envs.AppPort))
	r.Run(config.Envs.AppHost + ":" + strconv.Itoa(config.Envs.AppPort)) // por defecto en :8080
}
