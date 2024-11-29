package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

var ServWg sync.WaitGroup

// Запускает сервер, обрабатывает завершение его работы.
func StartServer(ctx context.Context) {
	http.HandleFunc("/orders", GetOrder)
	http.HandleFunc("/", DisplayPage)

	server := &http.Server{Addr: ":" + cfg.ServerPort, Handler: nil}

	// Ожидает завершения работы всех остальных горутин, затем завершает работу сервера.
	ServWg.Add(1)
	go func() {
		defer ServWg.Done()

		// Ожидает завершения работы всех остальных горутин.
		cfg.ExitWg.Wait()

		// Завершает работу сервера.
		// https://pkg.go.dev/net/http#Server.Shutdown
		if err := server.Shutdown(ctx); err != nil {
			slog.Info(
				"server shutdown",
				"err", err,
			)
		}

		slog.Info("server closed")
	}()

	slog.Info(
		"starting server...",
		"port", cfg.ServerPort,
	)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"unexpected server shutdown error",
			"err", err,
		)
	}
}
