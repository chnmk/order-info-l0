package config

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

var (
	KafkaInitNework        string
	KafkaInitAddress       string
	KafkaInitTopic         string
	KafkaInitPartition     int
	KafkaReconnectAttempts int
	KafkaReaderGoroutines  int
	KafkaWriterGoroutines  int
	KafkaReaderConfig      kafka.ReaderConfig
)

func getConsumerVars() {
	KafkaInitNework = Env.Get("KAFKA_NETWORK")
	KafkaInitAddress = fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT"))
	KafkaInitTopic = Env.Get("KAFKA_TOPIC")
	KafkaInitPartition = 0
	KafkaReconnectAttempts = envToInt("KAFKA_RECONNECT_ATTEMPTS")
	KafkaReaderGoroutines = envToInt("KAFKA_READER_GOROUTINES")
	KafkaReaderConfig = kafka.ReaderConfig{
		Brokers:     []string{fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT"))},
		GroupID:     Env.Get("KAFKA_GROUP_ID"),
		Topic:       Env.Get("KAFKA_TOPIC"),
		MaxBytes:    envToInt("KAFKA_MAX_BYTES"),
		MaxAttempts: envToInt("KAFKA_RECONNECT_ATTEMPTS"),
	}
}
