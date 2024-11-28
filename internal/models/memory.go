package models

import "github.com/segmentio/kafka-go"

type Storage interface {
	HandleMessage(*kafka.Reader, kafka.Message)
	AddOrder(string, string, []byte) OrderStorage
	ReadByID(int) OrderStorage
	ReadByUID(string) OrderStorage
	RestoreData()
}

type OrderStorage struct {
	ID           int
	UID          string
	Date_created string
	Order        []byte
}
