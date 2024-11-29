package memory

import (
	"encoding/json"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Обрабатывает сообщение в горутине. TODO: переписать этот комментарий.
func (m *MemStore) HandleMessage() {
	defer cfg.ExitWg.Done()

	for {
		select {
		case msg := <-cfg.MessagesChan:

			slog.Info("handling message...")

			var orderData models.Order

			err := json.Unmarshal(msg.Message.Value, &orderData)
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

			orderStruct := m.AddOrder(orderData.Order_uid, orderData.Date_created, msg.Message.Value)

			cfg.DB.InsertOrder(orderStruct)

			if err := msg.Reader.CommitMessages(cfg.ExitCtx, msg.Message); err != nil {
				slog.Error(err.Error())
			}

			slog.Info("message handling finished")

		case <-cfg.ExitCtx.Done():
			slog.Info("message handling canceled")

			return
		}
	}
}
