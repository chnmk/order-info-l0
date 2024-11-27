package models

type Storage interface {
	HandleMessage([]byte)
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
