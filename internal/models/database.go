package models

import (
	"context"
)

type Database interface {
	Close()
	Ping()
	CreateTables()
	GetOrdersIDs() []int
	SelectOrderById(int) (int, Order)
	InsertOrder(int, []byte, context.Context) error
}
