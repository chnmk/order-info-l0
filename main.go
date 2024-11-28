package main

import (
	"log/slog"
	"os"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/consumer"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
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
	// Подключается к БД и создаёт отсутствующие таблицы.
	cfg.DB = database.NewDB(cfg.DB, cfg.ExitCtx)
	defer cfg.DB.Close()

	// Инициализирует хранилище в памяти и восстанавливает данные из БД.
	cfg.Data = memory.NewStorage(cfg.Data)
	defer cfg.Exit()

	// Проверяет подключение к Kafka.
	consumer.Connect()

	/*
		// Запускает сервер.
		http.HandleFunc("/orders", transport.GetOrder)
		http.HandleFunc("/", transport.DisplayPage)

		slog.Info(
			"starting server...",
			"port", cfg.ServerPort,
		)

		err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), nil)
		if err != nil {
			panic(err)
		}
	*/

	cfg.ExitWg.Wait()
	slog.Info("shutting down...")
}
