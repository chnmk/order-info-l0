package memory

import (
	"encoding/json"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

// Обрабатывает сообщение в горутине. TODO: переписать этот комментарий.
func (m *MemStore) HandleMessage(reader *kafka.Reader, msg kafka.Message) {
	defer cfg.ExitWg.Done()

	slog.Error("handling message...")

	var orderData models.Order

	err := json.Unmarshal(msg.Value, &orderData)
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

	orderStruct := m.AddOrder(orderData.Order_uid, orderData.Date_created, msg.Value)

	cfg.DB.InsertOrder(orderStruct)

	// TODO: вынести это отсюда.
	if err := reader.CommitMessages(cfg.ExitCtx, msg); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("message handling finished")
}
