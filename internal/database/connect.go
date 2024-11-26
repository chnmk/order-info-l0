package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Имплементация интерфейса models.Database.
type PostgresDB struct {
	Conn *pgxpool.Pool
}

// Проверяет подключение к БД. В случае ошибки завершает работу сервиса.
func (db *PostgresDB) Ping() {
	err := db.Conn.Ping(context.TODO())
	if err != nil {
		slog.Error("failed to ping database: " + err.Error())
		os.Exit(1)
	}
}

// Обёртка для Pool.Close(), чтобы вызывать её из main.go.
func (db *PostgresDB) Close() {
	db.Conn.Close()
}
