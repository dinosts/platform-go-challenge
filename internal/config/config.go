package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment       string
	JWTSecretKey      string
	HashingSalt       string
	HashingIterations int
}

const notDefined = ""

func GetRequiredEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == notDefined {
		log.Panicf("%s env variable is required but not set", key)
	}

	return value
}

func GetOptionalEnvVariableWithDefaultValue(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == notDefined {
		log.Printf("%s env variable is not set", key)

		if defaultValue != "" {
			return defaultValue
		}
	}

	return value
}

func buildConfig() *Config {
	cfg := Config{
		Environment:  GetOptionalEnvVariableWithDefaultValue("APP_ENV", "dev"),
		JWTSecretKey: GetRequiredEnvVariable("JWT_SECRET_KEY"),
		HashingSalt:  GetRequiredEnvVariable("HASHING_SALT"),
	}

	return &cfg
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	return buildConfig()
}
