package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	KafkaInitNework        string
	KafkaInitAddress       string
	KafkaInitTopic         string
	KafkaInitPartition     int
	KafkaReconnectAttempts int
	KafkaReaderConfig      kafka.ReaderConfig
)

func getConsumerVars() {
	KafkaInitNework = Env.Get("KAFKA_NETWORK")
	KafkaInitAddress = fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT"))
	KafkaInitTopic = Env.Get("KAFKA_TOPIC")
	KafkaInitPartition = 0
	KafkaReconnectAttempts = getKafkaReconnectAttempts()
	KafkaReaderConfig = getKafkaReaderConfig()
}

func getKafkaReconnectAttempts() int {
	att := Env.Get("KAFKA_RECONNECT_ATTEMPTS")

	result, err := strconv.Atoi(att)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return result
}

func getKafkaReaderConfig() kafka.ReaderConfig {
	maxBytes, err := strconv.Atoi(Env.Get("KAFKA_MAX_BYTES"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	att, err := strconv.Atoi(Env.Get("KAFKA_RECONNECT_ATTEMPTS"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	interv, err := strconv.Atoi(Env.Get("KAFKA_COMMIT_INVERVAL_SECONDS"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return kafka.ReaderConfig{
		Brokers:        []string{fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT"))},
		GroupID:        Env.Get("KAFKA_GROUP_ID"),
		Topic:          Env.Get("KAFKA_TOPIC"),
		MaxBytes:       maxBytes,
		MaxAttempts:    att,
		CommitInterval: time.Duration(interv) * time.Second,
	}
}
