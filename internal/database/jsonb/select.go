package db_jsonb

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

// Строка для получения данных из таблицы orders.
//
// Не получает данные из таблиц delivery, payments и items.
var q_select_by_id = `
	SELECT *
	FROM orders 
	WHERE id = @id
`

var q_select_by_uid = `
	SELECT *
	FROM orders 
	WHERE uid = @uid
`

// Возвращает один заказ из базы данных по его id (положение в массиве в кэше).
func (db *PostgresDB) SelectOrderById(id int) (int, models.Order) {
	var key int
	var order models.Order
	var orderjson []byte

	args := pgx.NamedArgs{"id": id}
	err := db.Conn.QueryRow(context.Background(), q_select_by_id, args).Scan(&key, &orderjson)
	if err != nil {
		slog.Error("queryRow failed: " + err.Error())
	}

	err = json.Unmarshal(orderjson, &order)
	if err != nil {
		slog.Error("unmarshalling failed: " + err.Error())
	}

	slog.Info("order successfully selected from database")
	return key, order
}

// Возвращает один заказ из базы данных по его order_uid.
func (db *PostgresDB) SelectOrderByUID(uid string) (int, models.Order) {
	var key int
	var order models.Order
	var orderjson []byte

	args := pgx.NamedArgs{"id": uid}
	err := db.Conn.QueryRow(context.Background(), q_select_by_uid, args).Scan(&key, &orderjson)
	if err != nil {
		slog.Error("queryRow failed: " + err.Error())
	}

	err = json.Unmarshal(orderjson, &order)
	if err != nil {
		slog.Error("unmarshalling failed: " + err.Error())
	}

	slog.Info("order successfully selected from database")
	return key, order
}
