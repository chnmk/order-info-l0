package config

import (
	"fmt"
	"net"

	"github.com/segmentio/kafka-go"
)

var (
	KafkaWriterAddr  net.Addr
	KafkaWriterTopic string
)

func getTestVars() {
	KafkaWriterAddr = kafka.TCP(fmt.Sprintf("%s:%s", Env.Get("KAFKA_PROTOCOL"), Env.Get("KAFKA_PORT")))
	KafkaWriterTopic = Env.Get("KAFKA_TOPIC")
}
