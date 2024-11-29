package main

import (
	"log/slog"
	"os"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/consumer"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/server"
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
	// Подключается к БД, создаёт отсутствующие таблицы.
	cfg.DB = database.NewDB(cfg.ExitCtx, cfg.DB)
	defer cfg.DB.Close()

	// Инициализирует хранилище в памяти, восстанавливает данные из БД.
	cfg.Data = memory.NewStorage(cfg.ExitCtx, cfg.Data)
	defer cfg.Exit()

	// Создаёт пул обработчиков сообщений.
	for i := 0; i < cfg.MemoryHandlerGoroutines; i++ {
		cfg.ExitWg.Add(1)
		go cfg.Data.HandleMessage(cfg.ExitCtx)
	}

	// Создаёт обработчик устаревших сообщений.
	cfg.ExitWg.Add(1)
	go cfg.Data.ClearData(cfg.ExitCtx)

	// Проверяет подключение к Kafka, читает сообщения.
	consumer.Connect(cfg.ExitCtx)

	// Запускает сервер.
	server.StartServer(cfg.ExitCtx)

	// Ожидает завершения всех процессов.
	cfg.ExitWg.Wait()

	// Ожидает завершения работы сервера.
	server.ServWg.Wait()

	slog.Info("shutdown complete")
}
