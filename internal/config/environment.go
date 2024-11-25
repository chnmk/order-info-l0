package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var Env map[string]string

func SetDefaultEnv() {
	Env = make(map[string]string)
	Env["POSTGRES_DB"] = "orders"
	Env["POSTGRES_USER"] = "user"
	Env["POSTGRES_PASSWORD"] = "12345"
	Env["DB_PROTOCOL"] = "postgres"
	Env["DB_PORT"] = "5432"
	Env["DB_HOST"] = "postgres"
	Env["KAFKA_RECONNECT_ATTEMPTS"] = "20"
	Env["KAFKA_NETWORK"] = "tcp"
	Env["KAFKA_PROTOCOL"] = "kafka"
	Env["KAFKA_PORT"] = "9092"
	Env["KAFKA_TOPIC"] = "9092"
	Env["KAFKA_GROUP_ID"] = "go-orders-1"
	Env["KAFKA_MAX_BYTES"] = "100000" // 100kb
	Env["KAFKA_COMMIT_INVERVAL_SECONDS"] = "1"
	Env["CONSUMER_GOROUTINES"] = "1"
	Env["SERVER_PORT"] = "3000"
	Env["TEST_PUBLISH_DATA"] = "0"
}

func GetEnv() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("Warning: .env file not found")
	}

	slog.Info("Reading environment variables...")

	for name, def := range Env {
		value, exists := os.LookupEnv(name)
		if exists {
			// TODO: переписать эти ошибки
			slog.Info(name + ": " + value)
			Env[name] = value
		} else {
			slog.Info(name + " not found, using default (" + def + ")")
		}
	}
}
