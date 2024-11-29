package database

import (
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

// Строка для создания таблицы.
const q_create = `
	CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY, 
	uid VARCHAR(32),
	created VARCHAR(32),
	orderdata JSONB
	)`

// Создаёт отсутствующие таблицы в базе данных.
//
// Не использует индексы из-за потенциально значительно большего количества операций записи чем чтения.
func (db *PostgresDB) CreateTables() {
	_, err := db.Conn.Exec(cfg.ExitCtx, q_create)
	if err != nil {
		slog.Error(
			"failed to create tables",
			"err", err.Error(),
		)
		cfg.Exit()
	}
}
