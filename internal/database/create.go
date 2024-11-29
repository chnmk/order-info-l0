package database

import (
	"context"
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
func (db *PostgresDB) CreateTables(ctx context.Context) {
	_, err := db.Conn.Exec(ctx, q_create)
	if err != nil {
		slog.Error(
			"failed to create tables",
			"err", err.Error(),
		)
		cfg.Exit()
	}
}
