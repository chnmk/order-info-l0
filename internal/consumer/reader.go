package consumer

import (
	"context"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

// Читает сообщения из Kafka.
//
// Использован самый базовый функционал, логика обработки сообщений, включая коммиты, находится в пакете memory.
func (c *KafkaReader) Read(ctx context.Context) {
	slog.Info("reading messages from kafka...")

	for {
		// Чтобы не терять данные отдельно вызывает FetchMessage и, после записи в БД, CommitMessages.
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			slog.Error(err.Error())
			break
		}

		slog.Info("=== new message fetched ===")

		// Запускаем обработчик в горутине, чтобы мы не ждали, пока он доработает. TODO: подумать.
		// Он же его и закоммитит?
		// По-хорошему надо как-то через каналы это делать...
		go cfg.Data.HandleMessage(m.Value)

		// TODO: вынести это отсюда.
		if err := c.Reader.CommitMessages(ctx, m); err != nil {
			slog.Error(err.Error())
		}

	}

	if err := c.Reader.Close(); err != nil {
		slog.Error("failed to close reader: " + err.Error())
	}

	slog.Info("closing consumer connection...")
}
