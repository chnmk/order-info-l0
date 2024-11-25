package memory

import (
	"log/slog"
	"slices"
	"strconv"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

// Забирает все данные из БД, устанавливает значение currentkey на максимальное id заказа из БД.
func (m *MemStore) RestoreData() {
	slog.Info("restoring data from DB...")

	m.mu.Lock()
	defer m.mu.Unlock()

	ids := cfg.DB.GetOrdersIDs()
	if len(ids) == 0 {
		slog.Info("no data found in DB, restoring canceled")
		return
	}

	// TODO: у нас будут интерфейсы, так что скорее всего че-т поумнее придумаем.
	for _, id := range ids {
		key, order := cfg.DB.SelectOrderById(id)
		slog.Info(strconv.Itoa(key))
		m.orders[key] = order
	}

	m.currentkey = slices.Max(ids) + 1

	slog.Info("data successfully restored")
}
