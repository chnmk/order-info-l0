package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

func Consume() {
	topic := "orders"
	partition := 0

	var conn *kafka.Conn
	var err error

	// Попробуем подключиться несколько раз, почему бы и нет?
	for retry := 0; retry < 10; retry++ {
		conn, err = kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
		if err == nil {
			break
		}

		slog.Info("connection failed for kafka consumer, retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		slog.Info("failed to dial leader:" + err.Error())
		os.Exit(1)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	slog.Info("reading orders...")

	batch := conn.ReadBatch(10, 1e6) // fetch 10B min, 1MB max

	b := make([]byte, 100e3) // 100KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}

		// Проверяет данные и добавляет в БД если всё нормально
		var order models.Order
		err = json.Unmarshal(b[:n], &order)
		if err != nil {
			slog.Info("failed to unmarshal, skipping")
		} else {
			if ok := ValidateMsg(order); !ok {
				slog.Info("failed to validate, skipping")
			} else {
				fmt.Println(order)
				// database.InsertOrder(database.DB, order)
			}
		}
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}

	// TODO: чтение никогда не заканчивается. Либо через 10 минут.
	slog.Info("finished reading orders...")
}

/*
Проверяет что нужные поля не пустые и соответствуют нашим требованиям.

Пока что нам точно нужны те данные, которые выводятся в веб-интерфейсе.
*/
func ValidateMsg(order models.Order) bool {
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
