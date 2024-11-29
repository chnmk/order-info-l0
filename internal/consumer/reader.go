package consumer

import (
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
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
				slog.Error("failed to close reader: " + err.Error())
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

			// Запускаем обработчик в горутине, чтобы мы не ждали, пока он доработает. TODO: подумать.
			// Он же его и закоммитит?
			// По-хорошему надо как-то через каналы это делать...
			cfg.ExitWg.Add(1)
			go cfg.Data.HandleMessage(c.Reader, m)
		}
	}
}
