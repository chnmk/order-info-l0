package config

import (
	"fmt"

	"github.com/chnmk/order-info-l0/internal/models"
)

/*
	Устанавливает переменные для пакета database.
*/

var (
	DB         models.Database // Глобальная база данных.
	PgxpoolUrl string          // Строка подключения к базе данных.
)

// Получает глобальные переменные для пакета database.
func getDatabaseVars() {
	PgxpoolUrl = fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		Env.Get("DB_PROTOCOL"),
		Env.Get("POSTGRES_USER"),
		Env.Get("POSTGRES_PASSWORD"),
		Env.Get("DB_HOST"),
		Env.Get("DB_PORT"),
		Env.Get("POSTGRES_DB"),
	)
}
