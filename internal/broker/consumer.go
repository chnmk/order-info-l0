package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

func Consume() {
	slog.Info("connecting to kafka...")

	// make a new reader that consumes from orders
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "orders",
		Partition: 0,
		MaxBytes:  100e3, // 100kb
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			slog.Info(err.Error())
			break
		}

		storeMsg(m.Value)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

// Проверяет что нужные поля не пустые и соответствуют нашим требованиям.
//
// Пока что нам точно нужны те данные, которые выводятся в веб-интерфейсе.
func validateMsg(order models.Order) bool {
	if order.Order_uid == "" ||
		order.Delivery.Name == "" ||
		order.Delivery.City == "" ||
		order.Delivery.Address == "" ||
		order.Delivery.Phone == "" ||
		len(order.Items) < 1 {

		return false
	}

	for _, i := range order.Items {
		if i.Chrt_id == 0 ||
			i.Name == "" ||
			i.Total_price == 0 {
			return false
		}
	}

	return true
}

func storeMsg(m []byte) {
	var order models.Order
	err := json.Unmarshal(m, &order)
	if err != nil {
		slog.Info("failed to unmarshal, skipping")
	} else {
		if ok := validateMsg(order); !ok {
			slog.Info("failed to validate, skipping")
		} else {
			fmt.Println(order)
			// database.InsertOrder(database.DB, order)
		}
	}
}
