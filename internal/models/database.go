package models

import (
	"context"
)

type Database interface {
	Ping()
	CreateTables()
	RestoreData() []OrderStorage
	SelectOrderById(int) (int, Order)
	SelectOrderByUID(string) (int, Order)
	InsertOrder(OrderStorage, context.Context)
}
