package memory

import (
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
)

// Добавляет заказ value в память.
func (m *MemStore) AddOrder(value models.Order) {
	slog.Info("adding order to memory storage...")

	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.orders[m.currentkey]
	if ok {
		slog.Error("Failed to add order: id already exists")
		return
	}

	m.orders[m.currentkey] = value
	slog.Info("finished adding order to memory storage")
	m.currentkey++
}
