package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/chnmk/order-info-l0/internal/broker"
	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/transport"
	"github.com/chnmk/order-info-l0/internal/web"
	"github.com/chnmk/order-info-l0/test"
)

func init() {
	// Используется slog, поскольку он относится к стандартной библиотеке и обеспечивает простой вывод в формате JSON.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("initialization start...")

	// TODO: Удали меня.
	slog.Info(gofakeit.Emoji())

	// TODO: Удали меня, и из пакета test тоже.
	test.ReadModelFile()

	// Переменные окружения.
	config.SetDefaultEnv()
	config.GetEnv()

	// Публикация пробных данных в Kafka (TODO: сгенерировать нормальные данные).
	if config.EnvVariables["PUBLISH_TEST_DATA"] == "1" {
		test.PublishTestData()
	}

	slog.Info("initialization complete")
}

func main() {
	// Подключение к БД (TODO: использовать connection pool), пингуем (?) и создаем таблицы.
	database.DB = database.Connect()
	defer database.DB.Close(context.Background())

	// database.Ping(database.DB)
	database.CreateTables(database.DB)

	// Восстановление данных из БД в память (TODO: пошаманить с memory).
	memory.DATA = database.RestoreData(database.DB)

	// TODO: удали меня.
	if len(memory.DATA) == 0 {
		database.InsertOrder(database.DB, test.E)
	}

	// Подключение к Kafka (TODO: многопоточность?).
	broker.Consume()

	// Запуск сервера (TODO: обновить хендлеры).
	http.HandleFunc("/order", transport.GetOrder)
	http.HandleFunc("/", web.DisplayTemplate)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
