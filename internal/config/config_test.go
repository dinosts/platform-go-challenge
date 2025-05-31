package config_test

import (
	"os"
	"platform-go-challenge/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func unsetEnvVars() {
	os.Unsetenv("ENV")
	os.Unsetenv("JWT_SECRET_KEY")
}

func TestGetRequiredEnvVariable(t *testing.T) {
	// Arrange
	unsetEnvVars()
	defer func() { _ = recover() }()

	// Act/Assert
	config.GetRequiredEnvVariable("JWT_SECRET_KEY")
	t.Errorf("did not panic")
}

func TestGetOptionalEnvVariableWithDefaultValue(t *testing.T) {
	t.Run("should return default value when env var is not set", func(t *testing.T) {
		// Arrange
		unsetEnvVars()
		expected_result := "dev"

		// Act
		actual_result := config.GetOptionalEnvVariableWithDefaultValue("ENV", expected_result)

		// Assert
		assert.Equal(t, actual_result, expected_result)
	})

	t.Run("should return environment variable when it is set", func(t *testing.T) {
		// Arrange
		unsetEnvVars()

		expected_result := "production"
		os.Setenv("ENV", expected_result)

		// Act
		actual_result := config.GetOptionalEnvVariableWithDefaultValue("ENV", "dev")

		// Assert
		assert.Equal(t, actual_result, expected_result)
	})
}
