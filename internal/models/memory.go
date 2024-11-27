package models

type Storage interface {
	HandleMessage([]byte)
	AddOrder([]byte)
	Read(int) Order
	RestoreData()
}

type OrderStorage struct {
	ID           int
	UID          string
	Date_created string
	Order        []byte
}
