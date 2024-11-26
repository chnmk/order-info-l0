package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/consumer"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/transport"
)

func init() {
	// Используется slog, поскольку он относится к стандартной библиотеке и обеспечивает простой вывод в формате JSON.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("initialization start...")

	// Получает переменные окружения и их значения по умолчанию, проставляет переменные для других пакетов.
	cfg.Env = cfg.NewConfig()

	slog.Info("initialization complete")
}

func main() {
	// Подключается к БД и создаёт таблицы.
	cfg.DB = database.NewDB(cfg.DB, context.TODO())
	defer cfg.DB.Close()

	// Инициализирует хранилище в памяти и восстанавливает данные из БД.
	cfg.Data = memory.NewStorage(cfg.Data)

	// Пытается подключиться к Kafka в соответствии с переменной окружения KAFKA_RECONNECT_ATTEMPTS.
	consumer.Connect()

	// Запуск сервера.
	http.HandleFunc("/orders", transport.GetOrder)
	http.HandleFunc("/", transport.DisplayPage)

	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Env.Get("SERVER_PORT")), nil)
	if err != nil {
		panic(err)
	}
}
