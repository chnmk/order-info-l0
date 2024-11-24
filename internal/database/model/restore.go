package db_model

import (
	"context"
	"log/slog"
)

var q_orders_ids = "SELECT id FROM orders"

// Возвращает все id (не order_uid) заказов из таблицы orders.
func (db *PostgresDB) GetOrdersIDs() []int {
	rows, err := db.DB.Query(context.Background(), q_orders_ids)
	if err != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}

	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			slog.Error("QueryRow failed: " + err.Error())
		}
		ids = append(ids, id)
	}

	if rows.Err() != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}

	return ids
}
