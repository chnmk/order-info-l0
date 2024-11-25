package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/consumer"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/transport"
	"github.com/jackc/pgx/v5/pgxpool"
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
	// Подключается к БД, пингует и создает таблицы.
	ctx_db := context.Background()
	cfg.DB = database.NewDB(&pgxpool.Pool{}, ctx_db)
	defer cfg.DB.Close()

	cfg.DB.Ping()
	cfg.DB.CreateTables()

	// Инициализация хранилища и восстановление данных из БД.
	cfg.Data = memory.NewStorage(cfg.Data)
	cfg.Data.RestoreData()

	// Инициализация подключения к брокеру сообщений и создание горутин.
	ctx_consumers := context.Background()
	consumer.Init()

	routines, err := strconv.Atoi(cfg.Env.Get("CONSUMER_GOROUTINES"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	for i := 0; i < routines; i++ {
		c := consumer.NewConsumer()
		go c.Read(ctx_consumers)
	}

	// Генерация данных для брокера.
	if cfg.Env.Get("CONSUMER_PUBLISH_EXAMPLES") == "1" {
		consumer.GofakeInit()
		go consumer.PublishExampleData()
	}

	// Запуск сервера (TODO: обновить хендлеры).
	http.HandleFunc("/orders", transport.GetOrder)
	http.HandleFunc("/", transport.DisplayPage)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Env.Get("SERVER_PORT")), nil)
	if err != nil {
		panic(err)
	}
}
