package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

var Cons Consumer

type Consumer interface {
	Read(context.Context)
}

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func NewConsumer() Consumer {
	return &KafkaConsumer{}
}
