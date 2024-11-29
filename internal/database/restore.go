package database

import (
	"log/slog"
	"os"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

const q_restore = `
	SELECT * FROM orders
`

// Пытается получить все данные из БД. В случае неудачи завершает работу сервиса.
func (db *PostgresDB) RestoreData() []models.OrderStorage {
	var result []models.OrderStorage
	rows, err := db.Conn.Query(cfg.ExitCtx, q_restore)
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
		err = rows.Scan(&order.ID, &order.UID, &order.Date_created, &order.Order)
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
