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
var q_jsonorders = `
	SELECT *
	FROM jsonorders 
	WHERE id = @id
`

// Возвращает один заказ из базы данных по его order_uid.
func (db *PostgresDB) SelectOrderById(id int) (int, models.Order) {
	var key int
	var order models.Order
	var orderjson []byte

	args := pgx.NamedArgs{"id": id}
	err := db.Conn.QueryRow(context.Background(), q_jsonorders, args).Scan(&key, &orderjson)
	if err != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}

	err = json.Unmarshal(orderjson, &order)
	if err != nil {
		slog.Error("Unmarshalling failed: " + err.Error())
	}

	slog.Info("SelectOrderById: success")
	return key, order
}
