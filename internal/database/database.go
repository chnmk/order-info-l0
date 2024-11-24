package database

import (
	"context"
	"log/slog"
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	db_jsonb "github.com/chnmk/order-info-l0/internal/database/jsonb"
	db_model "github.com/chnmk/order-info-l0/internal/database/model"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Не обязательно, но пусть будет.
var once sync.Once

// Создаёт подключение к PostgreSQL.
//
// От переменной окружения DB_INTERFACE_MODE зависит схема базы данных заказов.
// При "model" заказ делится на таблицы с отдельным полем под каждое значение.
// При любом другом значении (значение по умолчанию) заказ записывается в поле типа JSONB.
func NewDB(db *pgxpool.Pool, ctx context.Context) models.Database {
	once.Do(func() {
		var err error
		db, err = pgxpool.New(ctx, cfg.PgxpoolUrl)
		if err != nil {
			slog.Error("Unable to connect to database: " + err.Error())
		}

	})

	if cfg.Env.Get("DB_INTERFACE_MODE") == "model" {
		return &db_model.PostgresDB{DB: db}
	} else {
		return &db_jsonb.PostgresDB{DB: db}
	}
}
