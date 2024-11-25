package memory

import (
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
)

func (d *MemStore) Read(id int) models.Order {
	slog.Info("reading from memory storage...")

	d.mu.Lock()
	defer d.mu.Unlock()

	order, ok := d.orders[id]
	if ok {
		slog.Error("Failed to read order from memory: id doesn't exist")
		return order
	}

	slog.Info("finished reading from memory storage")

	return order
}
