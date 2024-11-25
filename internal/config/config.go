package config

import "sync"

// TODO:  добавить разные структуры конфигов, например cfg.Kafka
var (
	Env  Config
	once sync.Once
)

type Config interface {
	InitEnv()
	GetEnv()
}

type EnvStorage struct {
	Env map[string]string
}

func (e *EnvStorage) init() {
}

func NewConfig() Config {
	once.Do(func() {
		Env = &EnvStorage{}
		Env.InitEnv()
	})

	return Env
}
