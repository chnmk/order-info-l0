package consumer

import (
	"context"
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/segmentio/kafka-go"
)

/*
TODO: написать объяснительную.

Consumer group используются по нескольким причинам:
	- Возможность напрямую использовать коммиты и не перечитывать старые сообщения.
	- Потенциальная поддержка многопоточности.
*/

func (c *KafkaConsumer) Read(ctx context.Context) {
	slog.Info("creating new kafka reader...")

	c.Reader = kafka.NewReader(cfg.KafkaReaderConfig)

	slog.Info("reader created, reading messages...")

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

		// TODO: мы в основной горутине делаем то что стоило бы делать асинхронно?
		if err := c.Reader.CommitMessages(ctx, m); err != nil {
			slog.Error(err.Error())
		}

	}

	if err := c.Reader.Close(); err != nil {
		slog.Error("failed to close reader: " + err.Error())
	}

	slog.Info("closing consumer connection...")
}