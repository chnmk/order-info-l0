package database

import (
	"context"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

// Строка для восстановления данных из БД.
const q_restore = `
	SELECT * FROM orders
	ORDER BY id DESC
	LIMIT @lim
`

// Пытается получить все данные из БД. В случае неудачи завершает работу сервиса.
func (db *PostgresDB) RestoreData(ctx context.Context) []models.OrderStorage {
	var result []models.OrderStorage

	args := pgx.NamedArgs{
		"lim": cfg.OrdersLimit,
	}

	rows, err := db.Conn.Query(ctx, q_restore, args)
	if err != nil {
		slog.Error(
			"failed to restore data",
			"err", err,
		)
		cfg.Exit()
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
			cfg.Exit()
		}
		result = append(result, order)
	}

	if rows.Err() != nil {
		slog.Error(
			"failed to restore data",
			"err", err,
		)
		cfg.Exit()
	}

	return result
}
