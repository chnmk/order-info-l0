package models

import "github.com/segmentio/kafka-go"

type Consumer interface {
	Read()
}

type MessageData struct {
	Reader  *kafka.Reader
	Message kafka.Message
}
