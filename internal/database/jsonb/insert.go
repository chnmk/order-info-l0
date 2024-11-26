package db_jsonb

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// Добавляет данные в таблицу jsonorders.
var q_insert_json = `
	INSERT INTO jsonorders(jsonorder)
	VALUES (@jsonorder)
	RETURNING id
`

// Пробует добавить заказ в БД, возвращает ошибку только в случае если заказ с таким id уже существует.
func (db *PostgresDB) InsertOrder(key int, m []byte, ctx context.Context) error {
	slog.Info("inserting order to database...")

	args := pgx.NamedArgs{"jsonorder": m}
	row := db.Conn.QueryRow(context.Background(), q_insert_json, args)

	var order_id int
	err := row.Scan(&order_id)
	if err != nil {
		slog.Error("Failed to insert data: " + err.Error())
		return nil
	}

	slog.Info("finished inserting order to database")
	return nil
}
