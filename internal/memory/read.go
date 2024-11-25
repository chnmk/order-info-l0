package memory

import (
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
)

func (m *MemStore) Read(id int) models.Order {
	slog.Info("reading from memory storage...")

	m.mu.Lock()
	defer m.mu.Unlock()

	order, ok := m.orders[id]
	if ok {
		slog.Error("Failed to read order from memory: id doesn't exist")
		return order
	}

	slog.Info("finished reading from memory storage")

	return order
}
