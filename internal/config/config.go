package config

import (
	"sync"

	"github.com/chnmk/order-info-l0/internal/models"
)

var (
	Env  models.Config // Глобальный конфиг.
	once sync.Once     // Создать конфиг можно только один раз.
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
	})

	return Env
}
