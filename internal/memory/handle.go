package memory

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Обрабатывает сообщение в горутине.
func (m *MemStore) HandleMessage(b []byte) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var order models.Order

		err := json.Unmarshal(b, &order)
		if err != nil {
			slog.Info("failed to unmarshal, skipping")
			return
		}

		if ok := ValidateMsg(order); !ok {
			slog.Info("failed to validate, skipping")
			return
		}

		err = database.DB.InsertOrder(m.currentkey, b, context.Background())
		if err != nil {
			slog.Error("failed to add order: order already exists")
			return
		}

		m.AddOrder(order)
	}()

	wg.Wait()
	slog.Info("order handling finished")
}
