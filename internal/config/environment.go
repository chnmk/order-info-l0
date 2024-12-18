package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Записывает стандартные значения переменных окружения, после чего вызывает функции для чтения самого окружение.
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
	e.Env["KAFKA_TOPIC"] = "go-orders"
	e.Env["KAFKA_GROUP_ID"] = "go-orders-1"
	e.Env["KAFKA_MAX_BYTES"] = "20000" // 20kb
	e.Env["KAFKA_READER_GOROUTINES"] = "1"
	e.Env["KAFKA_WRITE_EXAMPLES"] = "0"
	e.Env["KAFKA_WRITER_GOROUTINES"] = "1"
	e.Env["MEMORY_RESTORE_DATA"] = "1"
	e.Env["MEMORY_CLEANUP_MINUTES_INTERVAL"] = "10"
	e.Env["MEMORY_ORDERS_LIMIT"] = "100000" // исходя из максимального размера сообщения 20kb, 100 тысяч сообщений это не более 2гб памяти.
	e.Env["MEMORY_HANDLER_GOROUTINES"] = "1"
	e.Env["SERVER_PORT"] = "3000"
	e.Env["SERVER_GET_ORDER_BY_ID"] = "1"

	Env.ReadEnv()

	getConsumerVars()
	getDatabaseVars()
	getMemoryVars()
	getTestingVars()
	getServerVars()
}

// Читает переменные окружения и записывает их в мапу.
func (e *EnvStorage) ReadEnv() {
	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not found")
	}

	slog.Info("reading environment variables...")

	for name, def := range e.Env {
		value, exists := os.LookupEnv(name)
		if exists {
			slog.Info(
				"env variable found",
				"name", name,
				"value", value,
			)
			e.Env[name] = value
		} else {
			slog.Info(
				"env variable not found, using default",
				"name", name,
				"value", def,
			)
		}
	}
}

// Получает значение из мапы переменных.
func (e *EnvStorage) Get(key string) string {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Env[key]
}

// Конвертирует строковую переменную окружения s из мапы в тип int.
func envToInt(s string) int {
	att := Env.Get(s)

	result, err := strconv.Atoi(att)
	if err != nil {
		slog.Error(
			"invalid environment variable",
			"name", s,
			"err", err,
		)
		Exit()
	}

	return result
}

// Конвертирует строковую переменную окружения s из мапы в тип bool.
func envToBool(s string) bool {
	att := Env.Get(s)

	result, err := strconv.ParseBool(att)
	if err != nil {
		slog.Error(
			"invalid environment variable",
			"name", s,
			"err", err,
		)
		Exit()
	}

	return result
}
