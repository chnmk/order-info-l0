package db_jsonb

import (
	"context"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

// Строка для добавления данных в таблицу orders.
const q_insert = `
	INSERT INTO orders(id, uid, created, order)
	VALUES (@id, @uid, @created, @order)
	RETURNING id
`

// Пробует добавить заказ в БД, выводит ошибку если заказ с таким id уже существует.
func (db *PostgresDB) InsertOrder(order models.OrderStorage, ctx context.Context) {
	slog.Info(
		"inserting order to database...",
		"id", order.ID,
	)

	args := pgx.NamedArgs{
		"id":      order.ID,
		"uid":     order.UID,
		"created": order.Date_created,
		"order":   order.Order,
	}
	row := db.Conn.QueryRow(context.TODO(), q_insert, args)

	var order_id int
	err := row.Scan(&order_id)
	if err != nil {
		slog.Error(
			"failed to insert data",
			"err", err.Error,
			"id", order.ID,
		)
		return
	}

	slog.Info(
		"finished inserting order to database",
		"id", order.ID,
	)
}
