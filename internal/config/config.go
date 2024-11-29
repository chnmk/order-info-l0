package config

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chnmk/order-info-l0/internal/models"
)

var (
	Env     models.Config      // Глобальный конфиг.
	once    sync.Once          // Создать конфиг можно только один раз.
	Exit    context.CancelFunc // Функция завершения работы (graceful shutdown).
	ExitWg  sync.WaitGroup     // WaitGroup для ожидания выхода из всех горутин функциями main и server.
	ExitCtx context.Context    // Контекст для отмены всех процессов при завершении работы.
)

// Имплементация интерфейса models.Config.
type EnvStorage struct {
	mu  sync.Mutex
	Env map[string]string
}

// Записывает стандартные значения окружения, читает окружение, возвращает новый конфиг.
func NewConfig() models.Config {
	once.Do(func() {
		Env = &EnvStorage{}
		Env.InitEnv()

		ExitCtx, Exit = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	})

	return Env
}
