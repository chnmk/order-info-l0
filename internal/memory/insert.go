package memory

import (
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
)

// Добавляет заказ value в память.
func (d *MemStore) AddOrder(value models.Order) {
	slog.Info("adding order to memory storage...")

	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.orders[d.currentkey]
	if ok {
		slog.Error("Failed to add order: id already exists")
		return
	}

	d.orders[d.currentkey] = value
	slog.Info("finished adding order to memory storage")
	d.currentkey++
}
