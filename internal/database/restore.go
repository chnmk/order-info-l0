package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

var q_orders_ids = "SELECT id FROM orders"

func GetOrdersIDs(db *pgx.Conn) []int {
	rows, err := db.Query(context.Background(), q_orders_ids)
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
