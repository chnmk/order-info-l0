package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
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

	/*
		// Подключается к БД и создаёт отсутствующие таблицы.
		cfg.DB = database.NewDB(cfg.DB, cfg.ExitCtx)
		defer cfg.DB.Close()

		// Инициализирует хранилище в памяти и восстанавливает данные из БД.
		cfg.Data = memory.NewStorage(cfg.Data)
		defer cfg.Exit()

		// Проверяет подключение к Kafka.
		consumer.Connect()
	*/

	// Запускает сервер.
	http.HandleFunc("/orders", transport.GetOrder)
	http.HandleFunc("/", transport.DisplayPage)

	server := &http.Server{Addr: fmt.Sprintf(":%s", cfg.ServerPort), Handler: nil}

	slog.Info(
		"starting server...",
		"port", cfg.ServerPort,
	)

	var ServWg sync.WaitGroup
	ServWg.Add(1)

	go func() {
		defer ServWg.Done()

		cfg.ExitWg.Wait()
		if err := server.Shutdown(cfg.ExitCtx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		} else {
			slog.Info("test")
		}

		slog.Info("server closed")
	}()

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"shutdown error",
			"err", err,
		)
	}

	cfg.ExitWg.Wait()
	ServWg.Wait()
	slog.Info("shutting down...")
}
