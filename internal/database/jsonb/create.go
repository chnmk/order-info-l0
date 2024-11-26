package db_jsonb

import (
	"context"
	"log/slog"
)

// Строка для создания таблицы.
var createJSONOrders = `CREATE TABLE IF NOT EXISTS jsonorders (
	id serial PRIMARY KEY, 
	jsonorder JSONB
	)`

// Создаёт отсутствующие таблицы в базе данных.
//
// Не использует индексы из-за потенциально значительно большего количества операций записи чем чтения.
func (db *PostgresDB) CreateTables() {
	_, err := db.Conn.Exec(context.Background(), createJSONOrders)
	if err != nil {
		slog.Error("Failed to create table " + err.Error())
	}
}
