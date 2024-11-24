package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/consumer"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/transport"
	"github.com/chnmk/order-info-l0/internal/web"
	"github.com/chnmk/order-info-l0/test"
	"github.com/jackc/pgx/v5/pgxpool"
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
	// Подключение к БД, пингуем и создаем таблицы.
	ctx_data := context.Background()
	database.DB = database.NewDB(&pgxpool.Pool{}, ctx_data, "")
	defer database.DB.Close()

	database.DB.Ping()
	database.DB.CreateTables()

	// Инициализация хранилища и восстановление данных из БД.
	memory.DATA = memory.NewStorage()
	memory.DATA.RestoreData()

	// Инициализация подключения к брокеру сообщений и создание горутин.
	ctx_consumers := context.Background()
	consumer.Init()
	for i := 0; i < 1; i++ {
		c := consumer.NewConsumer()
		go c.Read(ctx_consumers)
	}

	// Генерация данных для брокера.
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
