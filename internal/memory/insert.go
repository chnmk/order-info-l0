package memory

import (
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
)

// Добавляет данные о заказе в память и возвращает сам заказ в том виде, в котором он хранится в памяти.
func (m *MemStore) AddOrder(order_uid string, date_created string, value []byte) models.OrderStorage {
	slog.Info("adding order to memory storage...")

	var order models.OrderStorage

	order.UID = order_uid
	order.Date_created = date_created
	order.Order = value

	m.mu.Lock()
	defer m.mu.Unlock()

	m.maxId++
	order.ID = m.maxId
	m.orders = append(m.orders, order)

	slog.Info(
		"finished adding order to memory storage",
		"id", order.ID,
	)

	return order
}
