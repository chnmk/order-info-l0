package config

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/segmentio/kafka-go"
)

/*
	Устанавливает переменные для тестов.
*/

var (
	KafkaWriterAddr       net.Addr
	KafkaWriterTopic      string
	KafkaWriteExamples    bool
	KafkaWriterGoroutines int
)

// Получает глобальные переменные для тестов.
func getTestingVars() {
	KafkaWriterAddr = kafka.TCP(fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT")))
	KafkaWriterTopic = Env.Get("KAFKA_TOPIC")

	boolenv, err := envToBool("KAFKA_WRITE_EXAMPLES")
	if err != nil {
		slog.Error(err.Error())
	} else {
		KafkaWriteExamples = boolenv
	}

	env, err := envToInt("KAFKA_WRITER_GOROUTINES")
	if err != nil {
		slog.Error(err.Error())
	} else {
		KafkaWriterGoroutines = env
	}

}
