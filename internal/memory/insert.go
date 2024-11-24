package memory

import (
	"context"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Добавляет заказ value в память.
func (d *MemStore) AddOrder(value models.Order) {
	d.mu.Lock()
	defer d.mu.Unlock()

	slog.Info("adding order to memory storage...")

	_, ok := d.orders[d.currentkey]
	if ok {
		slog.Error("Failed to add order: id already exists")
		return
	}

	err := database.DB.InsertOrder(d.currentkey, value, context.Background())
	if err != nil {
		slog.Error("Failed to add order: order already exists")
		return
	}

	d.orders[d.currentkey] = value
	slog.Info("added order to memory storage")
	d.currentkey++
}
