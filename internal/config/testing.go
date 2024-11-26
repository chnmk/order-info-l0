package config

import (
	"fmt"
	"net"

	"github.com/segmentio/kafka-go"
)

/*
	Устанавливает переменные для тестов.
*/

var (
	KafkaWriterAddr       net.Addr
	KafkaWriterTopic      string
	KafkaWriterGoroutines int
)

// Получает глобальные переменные для тестов.
func getTestingVars() {
	KafkaWriterAddr = kafka.TCP(fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT")))
	KafkaWriterTopic = Env.Get("KAFKA_TOPIC")
	KafkaReaderGoroutines = envToInt("KAFKA_WRITER_GOROUTINES")

}
