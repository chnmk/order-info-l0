package database

import (
	"context"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

var q_orders_ids = "SELECT order_uid FROM orders"

func GetOrdersIDs(db *pgx.Conn) []string {
	rows, err := db.Query(context.Background(), q_orders_ids)
	if err != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}

	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
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

func RestoreData(db *pgx.Conn) []models.Order {
	var result []models.Order

	ids := GetOrdersIDs(db)

	// Для каждого заказа (один models.Order из result) читаются данные из всех таблиц в БД по его order_uid
	for _, id := range ids {
		result = append(result, SelectOrderById(db, id))
	}

	return result
}
