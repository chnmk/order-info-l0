package models

type Database interface {
	Ping()
	Close()
	CreateTables()
	RestoreData() []OrderStorage
	InsertOrder(OrderStorage)
}
