package memory

import (
	"context"
	"encoding/json"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Обрабатывает сообщение в горутине. TODO: переписать этот комментарий.
func (m *MemStore) HandleMessage(b []byte) {
	slog.Error("handling message...")

	var orderData models.Order

	err := json.Unmarshal(b, &orderData)
	if err != nil {
		slog.Error(
			"failed to unmarshal message",
			"err", err,
		)
		return
	}

	err = ValidateMsg(orderData)
	if err != nil {
		slog.Error(
			"failed to validate message",
			"err", err,
		)
		return
	}

	orderStruct := m.AddOrder(orderData.Order_uid, orderData.Date_created, b)

	// TODO: хендлер БД будет отдельно жить?
	cfg.DB.InsertOrder(orderStruct, context.TODO())

	slog.Info("message handling finished")
}
