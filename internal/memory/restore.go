package memory

import (
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

// Забирает все данные из БД, устанавливает значение currentkey на максимальное id заказа из БД.
func (m *MemStore) RestoreData() {
	slog.Info("restoring data from DB...")

	m.mu.Lock()
	defer m.mu.Unlock()

	m.orders = cfg.DB.RestoreData()

	for _, order := range m.orders {
		if order.ID > m.maxId {
			m.maxId = order.ID
		}

	}

	slog.Info("data successfully restored")
}
