package models

import "context"

type Database interface {
	Close()
	Ping(context.Context)
	CreateTables(context.Context)
	InsertOrder(context.Context, OrderStorage)
	RestoreData(context.Context) []OrderStorage
}
