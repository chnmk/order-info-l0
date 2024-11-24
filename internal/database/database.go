package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	db_jsonb "github.com/chnmk/order-info-l0/internal/database/jsonb"
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
	InsertOrder(int, []byte, context.Context) error
}

// Создаёт подключение к PostgreSQL.
//
// От переменной окружения DB_INTERFACE_MODE зависит схема базы данных заказов.
// При "model" заказ делится на таблицы с отдельным полем под каждое значение.
// При любом другом значении (значение по умолчанию) заказ записывается в поле типа JSONB.
func NewDB(db *pgxpool.Pool, ctx context.Context) Database {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		cfg.Env["POSTGRES_PROTOCOL"],
		cfg.Env["POSTGRES_USER"],
		cfg.Env["POSTGRES_PASSWORD"],
		cfg.Env["POSTGRES_HOST"],
		cfg.Env["POSTGRES_PORT"],
		cfg.Env["POSTGRES_NAME"],
	)

	once.Do(func() {
		var err error
		db, err = pgxpool.New(ctx, url)
		if err != nil {
			slog.Error("Unable to connect to database: " + err.Error())
		}

	})

	if cfg.Env["DB_INTERFACE_MODE"] == "model" {
		return &db_model.PostgresDB{DB: db}
	} else {
		return &db_jsonb.PostgresDB{DB: db}
	}
}
