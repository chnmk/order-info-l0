package config

import (
	"context"
	"sync"

	"github.com/chnmk/order-info-l0/internal/models"
)

var (
	Env        models.Config // Глобальный конфиг.
	once       sync.Once     // Создать конфиг можно только один раз.
	ExitCtx    context.Context
	ExitCancel context.CancelFunc
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

		ExitCtx, ExitCancel = context.WithCancel(context.Background())
	})

	return Env
}
