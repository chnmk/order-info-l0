package memory

import (
	"log/slog"
	"slices"

	"github.com/chnmk/order-info-l0/internal/database"
)

// Забирает все данные из БД, устанавливает значение currentkey на максимальное id заказа из БД.
func (d *MemStore) RestoreData() {
	d.mu.Lock()
	defer d.mu.Unlock()

	slog.Info("restoring data from DB...")

	ids := database.DB.GetOrdersIDs()
	if len(ids) == 0 {
		slog.Info("no data found in DB, restoring canceled")
		return
	}

	// TODO: у нас будут интерфейсы, так что скорее всего че-т поумнее придумаем.
	for _, id := range ids {
		key, order := database.DB.SelectOrderById(id)
		d.orders[key] = order
	}

	d.currentkey = slices.Max(ids) + 1

	slog.Info("data successfully restored")
}
