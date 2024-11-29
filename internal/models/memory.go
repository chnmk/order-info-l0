package models

type Storage interface {
	HandleMessage()
	AddOrder(string, string, []byte) OrderStorage
	ReadByID(int) OrderStorage
	ReadByUID(string) OrderStorage
	RestoreData()
	ClearData()
}

type OrderStorage struct {
	ID           int
	UID          string
	Date_created string
	Order        []byte
}
