package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/chnmk/order-info-l0/internal/broker"
	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/transport"
	"github.com/chnmk/order-info-l0/internal/web"
	"github.com/chnmk/order-info-l0/test"
	"github.com/jackc/pgx/v5"
)

func init() {
	// Используется slog, поскольку он относится к стандартной библиотеке и обеспечивает простой вывод в формате JSON.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("initialization start...")

	// Переменные окружения.
	config.SetDefaultEnv()
	config.GetEnv()

	slog.Info("initialization complete")
}

func main() {
	database.DB = database.NewPostgresDB(&pgx.Conn{}, "full")

	// Подключение к БД (TODO: использовать connection pool), пингуем (?) и создаем таблицы.
	ctx_database := context.Background()
	database.DB.Connect(ctx_database)
	defer database.DB.Close(ctx_database)

	// database.DB.Ping()
	database.DB.CreateTables()

	// Инициализация хранилища и восстановление данных из БД (TODO: пошаманить с memory?).
	memory.DATA.Init()
	memory.DATA.RestoreData()

	broker.Init()

	// Создание горутин для Кафки.
	ctx_consumers := context.Background()
	for i := 0; i < 1; i++ {
		go broker.Consume(ctx_consumers)
	}

	// Генерация данных для Kafka
	if config.EnvVariables["PUBLISH_TEST_DATA"] == "1" {
		test.GofakeInit()
		go test.PublishTestData()
	}

	// Запуск сервера (TODO: обновить хендлеры).
	http.HandleFunc("/orders", transport.GetOrder)
	http.HandleFunc("/", web.DisplayTemplate)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
