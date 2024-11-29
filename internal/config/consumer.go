package config

import (
	"fmt"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

/*
	Устанавливает переменные для пакета consumer.
*/

var (
	KafkaInitNework        string
	KafkaInitAddress       string
	KafkaInitTopic         string
	KafkaInitPartition     int
	KafkaReconnectAttempts int
	KafkaReaderGoroutines  int
	KafkaReaderConfig      kafka.ReaderConfig

	MessagesChan chan models.MessageData
)

// Получает глобальные переменные для пакета consumer.
func getConsumerVars() {
	KafkaInitNework = Env.Get("KAFKA_NETWORK")
	KafkaInitAddress = fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT"))
	KafkaInitTopic = Env.Get("KAFKA_TOPIC")
	KafkaInitPartition = 0

	env, err := envToInt("KAFKA_RECONNECT_ATTEMPTS")
	if err != nil {
		slog.Error(err.Error())
	} else {
		KafkaReconnectAttempts = env
	}

	env, err = envToInt("KAFKA_READER_GOROUTINES")
	if err != nil {
		slog.Error(err.Error())
	} else {
		KafkaReaderGoroutines = env
	}

	var max_bytes int
	var reconnect_attempts int

	env, err = envToInt("KAFKA_MAX_BYTES")
	if err != nil {
		slog.Error(err.Error())
	} else {
		max_bytes = env
	}

	env, err = envToInt("KAFKA_RECONNECT_ATTEMPTS")
	if err != nil {
		slog.Error(err.Error())
	} else {
		reconnect_attempts = env
	}

	KafkaReaderConfig = kafka.ReaderConfig{
		Brokers:     []string{fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT"))},
		GroupID:     Env.Get("KAFKA_GROUP_ID"),
		Topic:       Env.Get("KAFKA_TOPIC"),
		MaxBytes:    max_bytes,
		MaxAttempts: reconnect_attempts,
	}

	MessagesChan = make(chan models.MessageData, 10) // TODO: заменить 10 на ожидаемое число хендлеров
}
