package database

import (
	"context"

	db_model "github.com/chnmk/order-info-l0/internal/database/model"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5"
)

var DB Database

type Database interface {
	Connect(ctx context.Context)
	Close(ctx context.Context)
	Ping()
	CreateTables()
	GetOrdersIDs() []int
	SelectOrderById(int) (int, models.Order)
	InsertOrder(models.Order, int) error
}

func NewPostgresDB(db *pgx.Conn, cfg string) Database {
	if cfg == "full" {
		return &db_model.PostgresDB{DB: db}
	} else {
		return &db_model.PostgresDB{DB: db}
	}
}
