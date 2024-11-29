package models

type Consumer interface {
	Read() // Читает данные из Kafka.
}
