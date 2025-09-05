package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"log"
)

var globalEnv *Env

type Env struct {
	GinPort          string `envconfig:"GIN_PORT" default:"8000"`
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:""`
	PostgresUser     string `envconfig:"POSTGRES_USER" default:""`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:""`
	PostgresDB       string `envconfig:"POSTGRES_DB" default:""`
	PostgresPort     string `envconfig:"POSTGRES_PORT" default:""`
	JWTSecret        string `envconfig:"JWT_SECRET" default:"supersecret"`
}

func LoadEnv() *Env {
	// load .env file (optional, falls back to OS env if not found)
	if err := godotenv.Load(".env"); err != nil {
		GetLogger().Warn(".env file not loaded, using system env instead", zap.Error(err))
	}

	var env Env
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env: %v", err)
	}

	globalEnv = &env
	return globalEnv
}

func GetEnv() *Env {
	return globalEnv
}
