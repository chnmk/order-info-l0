package models

type Storage interface {
	HandleMessage([]byte)
	AddOrder([]byte)
	Read(int) Order
	RestoreData()
}

type OrderStorage struct {
	UID     string
	Expires int // Дата удаления заказа в формате epoch.
	Order   []byte
}
