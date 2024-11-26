package models

import "time"

type Storage interface {
	HandleMessage([]byte)
	AddOrder([]byte)
	Read(int) Order
	RestoreData()
}

type OrderStorage struct {
	UID     string
	Order   []byte
	Expires time.Time
}
