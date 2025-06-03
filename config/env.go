package config

import (
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
)

type Env struct {
	AppPort    int    `validate:"required"`
	AppHost    string `validate:"required"`
	DbPassword string `validate:"required"`
	DbName     string `validate:"required"`
	DbHost     string `validate:"required"`
	DbPort     int    `validate:"required"`
	DbUser     string `validate:"required"`
	JWTSecret  string `validate:"required"`
}

var Envs Env

func GetEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Envs = Env{
		AppPort:    mustParseInt(os.Getenv("APP_PORT")),
		AppHost:    os.Getenv("APP_HOST"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     mustParseInt(os.Getenv("DB_PORT")),
		DbUser:     os.Getenv("DB_USER"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	validate := validator.New()
	if err := validate.Struct(Envs); err != nil {
		log.Fatalf("Error validando las variables de entorno: %v", err)

	}
}

func mustParseInt(val string) int {
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Error convirtiendo a int: %s", val)
	}
	return parsed
}
