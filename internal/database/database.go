package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/chnmk/order-info-l0/internal/config"
	db_model "github.com/chnmk/order-info-l0/internal/database/model"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB   Database
	once sync.Once
)

type Database interface {
	Close()
	Ping()
	CreateTables()
	GetOrdersIDs() []int
	SelectOrderById(int) (int, models.Order)
	InsertOrder(models.Order, int) error
}

func NewDB(db *pgxpool.Pool, ctx context.Context, cfg string) Database {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.EnvVariables["DB_USER"],
		config.EnvVariables["DB_PASSWORD"],
		"postgres",
		"5432",
		config.EnvVariables["DB_NAME"],
	)

	once.Do(func() {
		var err error
		db, err = pgxpool.New(ctx, url)
		if err != nil {
			slog.Error("Unable to connect to database: " + err.Error())
		}

	})

	if cfg == "full" {
		return &db_model.PostgresDB{DB: db}
	} else {
		return &db_model.PostgresDB{DB: db}
	}
}
