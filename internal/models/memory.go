package models

type Storage interface {
	Init()
	HandleMessage([]byte)
	AddOrder(Order)
	Read(int) Order
	RestoreData()
}
