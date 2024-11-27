package models

type Storage interface {
	HandleMessage([]byte)
	AddOrder([]byte)
	Read(int) Order
	RestoreData()
}

type OrderStorage struct {
	ID      int
	UID     string
	Expires string // Дата удаления заказа в формате epoch.
	Order   []byte
}
