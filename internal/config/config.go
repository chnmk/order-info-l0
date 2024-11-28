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
	Env     models.Config // Глобальный конфиг.
	once    sync.Once     // Создать конфиг можно только один раз.
	Exit    context.CancelFunc
	ExitWg  sync.WaitGroup
	ExitCtx context.Context
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
