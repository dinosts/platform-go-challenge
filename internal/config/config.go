package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecretKey string
}

var notDefined = ""

func getRequiredEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == notDefined {
		log.Fatalf("%s env variable is required but not set", key)
	}

	return value
}

func getOptionalEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == notDefined {
		log.Printf("%s env variable is not set", key)
	}

	return value
}

func buildConfig() *Config {
	cfg := Config{
		JWTSecretKey: getRequiredEnvVariable("JWT_SECRET_KEY"),
	}

	return &cfg
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file could not be loaded")
	}

	cfg := buildConfig()

	return cfg
}
