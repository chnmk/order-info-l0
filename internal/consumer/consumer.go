package consumer

import (
	"log/slog"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

// Имплементация интерфейса models.KafkaReader.
type KafkaReader struct {
	Reader *kafka.Reader
}

// Возвращает новый reader в соответствии с настройками из переменных окружения.
func newReader() models.Consumer {
	slog.Info("creating new kafka reader...")
	return &KafkaReader{Reader: kafka.NewReader(cfg.KafkaReaderConfig)}
}
