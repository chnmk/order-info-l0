package database

import (
	"context"
	"log/slog"
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Не обязательно, но пусть будет.
var once sync.Once

// Создаёт подключение к PostgreSQL, пингует, создает таблицы если их нет.
func NewDB(db models.Database, ctx context.Context) models.Database {
	slog.Info("initializing new database connection pool...")

	once.Do(func() {
		var err error
		conn, err := pgxpool.New(ctx, cfg.PgxpoolUrl)
		if err != nil {
			slog.Error("unable to connect to database: " + err.Error())
		}

		db = &PostgresDB{Conn: conn}

		db.Ping()
		db.CreateTables()
	})

	slog.Info("database successfully initialized")

	return db

}
