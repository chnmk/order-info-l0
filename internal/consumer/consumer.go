package consumer

import (
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

var Cons models.Consumer

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func newConsumer() models.Consumer {
	return &KafkaConsumer{}
}
