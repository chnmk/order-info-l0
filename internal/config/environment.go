package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func (e *EnvStorage) InitEnv() {
	e.Env = make(map[string]string)
	e.Env["POSTGRES_DB"] = "orders"
	e.Env["POSTGRES_USER"] = "user"
	e.Env["POSTGRES_PASSWORD"] = "12345"
	e.Env["DB_PROTOCOL"] = "postgres"
	e.Env["DB_PORT"] = "5432"
	e.Env["DB_HOST"] = "postgres"
	e.Env["KAFKA_RECONNECT_ATTEMPTS"] = "20"
	e.Env["KAFKA_NETWORK"] = "tcp"
	e.Env["KAFKA_PROTOCOL"] = "kafka"
	e.Env["KAFKA_PORT"] = "9092"
	e.Env["KAFKA_TOPIC"] = "9092"
	e.Env["KAFKA_GROUP_ID"] = "go-orders-1"
	e.Env["KAFKA_MAX_BYTES"] = "100000" // 100kb
	e.Env["KAFKA_COMMIT_INVERVAL_SECONDS"] = "1"
	e.Env["CONSUMER_GOROUTINES"] = "1"
	e.Env["SERVER_PORT"] = "3000"
	e.Env["TEST_PUBLISH_DATA"] = "0"
}

func (e *EnvStorage) GetEnv() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("Warning: .env file not found")
	}

	slog.Info("Reading environment variables...")

	for name, def := range e.Env {
		value, exists := os.LookupEnv(name)
		if exists {
			// TODO: переписать эти ошибки
			slog.Info(name + ": " + value)
			e.Env[name] = value
		} else {
			slog.Info(name + " not found, using default (" + def + ")")
		}
	}
}
