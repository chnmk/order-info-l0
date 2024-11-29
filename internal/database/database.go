package database

import (
	"context"
	"log/slog"
	"sync"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Ожидается, что БД будет создана только один раз.
// Подстраховка от неожиданного поведения.
var once sync.Once

// Имплементация интерфейса models.Database.
type PostgresDB struct {
	Conn *pgxpool.Pool
}

// Создаёт подключение к PostgreSQL, пингует, создает таблицы если их нет.
func NewDB(db models.Database, ctx context.Context) models.Database {
	slog.Info("initializing new database connection pool...")

	once.Do(func() {
		var err error
		conn, err := pgxpool.New(ctx, cfg.PgxpoolUrl)
		if err != nil {
			slog.Error(
				"unable to connect to database",
				"err", err.Error(),
			)
		}

		db = &PostgresDB{Conn: conn}

		db.Ping()
		db.CreateTables()
	})

	slog.Info("database connection pool successfully initialized")

	return db

}

// Проверяет подключение к БД. В случае ошибки завершает работу сервиса.
func (db *PostgresDB) Ping() {
	err := db.Conn.Ping(cfg.ExitCtx)
	if err != nil {
		slog.Error(
			"failed to ping database",
			"err", err.Error(),
		)
		cfg.Exit()
	}
}

// Обёртка для Pool.Close(), чтобы вызывать её из main.go.
func (db *PostgresDB) Close() {
	db.Conn.Close()
}
