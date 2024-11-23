package database

import (
	"context"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/models"
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

func RestoreData(db *pgx.Conn) map[int][]models.Order {
	var result map[int][]models.Order

	ids := GetOrdersIDs(db)

	for _, id := range ids {
		key, order := SelectOrderById(db, id)
		memory.DATA.AddOrder(key, order)
	}

	return result
}
