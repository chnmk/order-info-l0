package memory

import (
	"context"
	"encoding/json"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Обрабатывает сообщение в горутине.
func (m *MemStore) HandleMessage(b []byte) {
	// TODO: подумать над горутинами, плюс может унмаршалер отдельно будет жить
	var order models.Order

	// TODO: валидировать будем отдельно в горутине
	err := json.Unmarshal(b, &order)
	if err != nil {
		slog.Info("failed to unmarshal, skipping")
		return
	}

	// TODO: валидировать будем отдельно в горутине
	if ok := ValidateMsg(order); !ok {
		slog.Info("failed to validate, skipping")
		return
	}

	// TODO: хендлер БД будет отдельно жить?
	err = cfg.DB.InsertOrder(len(m.orders), b, context.Background())
	if err != nil {
		slog.Error("failed to add order: order already exists")
		return
	}

	// Это очень странно если у нас в память добавляется ПОСЛЕ бд
	// Если ты хочешь проверять повторный заказ... Нафига?
	m.AddOrder(b)

	slog.Info("order handling finished")
}
