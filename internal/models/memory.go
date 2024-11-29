package models

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Storage interface {
	// Обрабатывает сообщение из Kafka.
	HandleMessage(context.Context)

	// Добавляет данные о заказе в память и возвращает сам заказ в том виде, в котором он хранится в памяти.
	AddOrder(order_uid string, date_created string, value []byte) OrderStorage

	ReadByID(int) OrderStorage
	ReadByUID(string) OrderStorage

	RestoreData(context.Context)
	ClearData(context.Context)
}

// Формат хранения данных в памяти.
type OrderStorage struct {
	ID           int
	UID          string
	Date_created string
	Order        []byte
}

// Формат передачи сообщений из Kafka обработчикам через канал.
type MessageData struct {
	Reader  *kafka.Reader
	Message kafka.Message
}
