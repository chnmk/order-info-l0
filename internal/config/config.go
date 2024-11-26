package config

import (
	"sync"

	"github.com/chnmk/order-info-l0/internal/models"
)

var (
	Env  models.Config
	once sync.Once
)

type EnvStorage struct {
	mu  sync.Mutex
	Env map[string]string
}

func NewConfig() models.Config {
	once.Do(func() {
		Env = &EnvStorage{}
		Env.InitEnv()
		Env.ReadEnv()
		getConsumerVars()
		getDatabaseVars()
		getTestVars()
	})

	return Env
}
