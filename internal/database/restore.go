package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/chnmk/order-info-l0/internal/models"
)

const q_restore = `
	SELECT * FROM orders
`

// Пытается получить все данные из БД. В случае неудачи завершает работу сервиса.
func (db *PostgresDB) RestoreData() []models.OrderStorage {
	slog.Info("restoring data from database...")

	var result []models.OrderStorage
	rows, err := db.Conn.Query(context.TODO(), q_restore)
	if err != nil {
		slog.Error(
			"failed to restore data",
			"err", err,
		)
		os.Exit(1)
	}

	defer rows.Close()

	for rows.Next() {
		var order models.OrderStorage
		err = rows.Scan(&order)
		if err != nil {
			slog.Error(
				"failed to restore data",
				"err", err,
			)
			os.Exit(1)
		}
		result = append(result, order)
	}

	if rows.Err() != nil {
		slog.Error(
			"failed to restore data",
			"err", err,
		)
		os.Exit(1)
	}

	return result
}
