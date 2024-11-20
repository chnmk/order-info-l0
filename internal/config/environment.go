package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var EnvVariables map[string]string

func SetDefaultEnv() {
	EnvVariables = make(map[string]string)
	EnvVariables["PUBLISH_TEST_DATA"] = "0"
	EnvVariables["DB_NAME"] = "orders"
	EnvVariables["DB_USER"] = "user"
	EnvVariables["DB_PASSWORD"] = "12345"
}

func GetEnv() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("Warning: .env file not found")
	}

	slog.Info("Reading environment variables...")

	for name, def := range EnvVariables {
		value, exists := os.LookupEnv(name)
		if exists {
			// TODO: переписать эти ошибки
			slog.Info(name + ": " + value)
			EnvVariables[name] = value
		} else {
			slog.Info(name + " not found, using default (" + def + ")")
		}
	}
}
