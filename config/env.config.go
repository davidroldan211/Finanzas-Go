package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	SSLMode  string `validate:"required"`
}

type JWTConfig struct {
	Secret  string        `validate:"required"`
	Expires time.Duration `validate:"required"`
}

type ServerConfig struct {
	Host           string `validate:"required"`
	Port           int    `validate:"required"`
	TrustedProxies []string
}

type AppConfig struct {
	Environment string `validate:"required"`
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtExpiresIn, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))

	if err != nil {
		jwtExpiresIn = 24 * time.Hour // 24h por defecto si no se puede parsear
	}

	Config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "mydb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expires: jwtExpiresIn,
		},
		Server: ServerConfig{
			Host:           getEnv("SERVER_HOST", "localhost"),
			Port:           getEnvAsInt("SERVER_PORT", 8080),
			TrustedProxies: getEnvAsList("TRUSTED_PROXIES"),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
		},
	}
	return Config, nil
}
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsList(key string) []string {
	value := os.Getenv(key)
	if value == "" {
		return nil
	}

	var items []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			items = append(items, item)
		}
	}
	return items
}
