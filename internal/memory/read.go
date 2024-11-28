package memory

import (
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
)

// Получает заказ со всеми дополнительными данными по его ID. Ожидается, что это более эффективно, чем ReadByUID.
func (m *MemStore) ReadByID(id int) (order models.OrderStorage) {
	slog.Info(
		"reading from memory storage...",
		"id", id,
	)

	m.mu.Lock()
	defer m.mu.Unlock()

	if id > len(m.orders)-1 {
		slog.Error(
			"order not found",
			"id", id,
		)
		return
	}

	order = m.orders[id]

	slog.Info(
		"finished reading from memory storage",
		"id", id,
	)

	return
}

// Получает заказ со всеми дополнительными данными по его UID. Ожидается, что это менее эффективно, чем ReadByID.
func (m *MemStore) ReadByUID(uid string) (order models.OrderStorage) {
	slog.Info(
		"reading from memory storage...",
		"uid", uid,
	)

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, o := range m.orders {
		if o.UID == uid {
			slog.Info(
				"finished reading from memory storage",
				"uid", uid,
			)
			return o
		}
	}

	slog.Error(
		"order not found",
		"uid", uid,
	)

	return
}
