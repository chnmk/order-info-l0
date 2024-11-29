package memory

import (
	"context"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

// Забирает все данные из БД, устанавливает значение maxId на максимальное id заказа из БД.
func (m *MemStore) RestoreData(ctx context.Context) {
	slog.Info("restoring data from database...")

	m.mu.Lock()
	defer m.mu.Unlock()

	m.orders = cfg.DB.RestoreData(ctx)

	if len(m.orders) == 0 {
		slog.Info("no data found, restoring cancelled")
		return
	}

	for _, order := range m.orders {
		if order.ID > m.maxId {
			m.maxId = order.ID
		}

	}

	slog.Info("data successfully restored")
}
