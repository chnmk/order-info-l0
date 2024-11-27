package models

import (
	"context"
)

type Database interface {
	Ping()
	Close()
	CreateTables()
	RestoreData() []OrderStorage
	SelectOrderById(int) (int, Order)
	SelectOrderByUID(string) (int, Order)
	InsertOrder(OrderStorage, context.Context)
}
