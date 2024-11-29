package consumer

import (
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Читает сообщения из Kafka.
//
// Использован самый базовый функционал, логика обработки сообщений, включая коммиты, находится в пакете memory.
func (c *KafkaReader) Read() {
	defer cfg.ExitWg.Done()

	slog.Info("reading messages from kafka...")

	for {
		select {

		case <-cfg.ExitCtx.Done():
			if err := c.Reader.Close(); err != nil {
				slog.Error(
					"failed to close reader",
					"err", err.Error(),
				)
			}

			slog.Info("closing consumer connection...")
			return

		default:
			// Чтобы не терять данные отдельно вызывает FetchMessage и, после записи в БД, CommitMessages.
			m, err := c.Reader.FetchMessage(cfg.ExitCtx)
			if err != nil {
				slog.Error(err.Error())
				break
			}

			slog.Info("new message fetched")

			// Отправляет сообщение в канал для обработки. Поле reader необходимо, чтобы после обработки закоммитить сообщение.
			cfg.MessagesChan <- models.MessageData{Reader: c.Reader, Message: m}
		}
	}
}
