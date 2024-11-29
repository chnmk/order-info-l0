package database

import (
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

// TODO: выбирать топ 100к (а лучше из переменной окружения) в обратном порядке по ID, где дата не раньше чем 14 дней назад (тоже по переменной)
const q_restore = `
	SELECT * FROM orders
	ORDER BY id DESC
	LIMIT @lim
`

// Пытается получить все данные из БД. В случае неудачи завершает работу сервиса.
func (db *PostgresDB) RestoreData() []models.OrderStorage {
	var result []models.OrderStorage

	args := pgx.NamedArgs{
		"lim": cfg.OrdersLimit,
	}

	rows, err := db.Conn.Query(cfg.ExitCtx, q_restore, args)
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
