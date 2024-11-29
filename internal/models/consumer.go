package models

import "context"

type Consumer interface {
	Read(context.Context) // Читает данные из Kafka.
}
