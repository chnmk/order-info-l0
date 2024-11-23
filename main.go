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
)

func init() {
	// Используется slog, поскольку он относится к стандартной библиотеке и обеспечивает простой вывод в формате JSON.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("initialization start...")

	// Переменные окружения.
	config.SetDefaultEnv()
	config.GetEnv()

	// Публикация пробных данных в Kafka (TODO: генерировать постоянно).
	/*
		if config.EnvVariables["PUBLISH_TEST_DATA"] == "1" {
			test.PublishTestData()
		}
	*/

	slog.Info("initialization complete")
}

func main() {
	// Подключение к БД (TODO: использовать connection pool), пингуем (?) и создаем таблицы.
	database.DB = database.Connect()
	defer database.DB.Close(context.Background())

	// database.Ping(database.DB)
	database.CreateTables(database.DB)

	// Инициализация хранилища и восстановление данных из БД (TODO: пошаманить с memory).
	memory.DATA.Init()
	database.RestoreData(database.DB)

	// Подключение к Kafka (TODO: многопоточность?).
	// go broker.Consume()
	go broker.Consume()

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
