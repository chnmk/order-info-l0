package db_model

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Conn *pgxpool.Pool
}

func (db *PostgresDB) Close() {
	db.Conn.Close()
}

func (db *PostgresDB) Ping() {
	err := db.Conn.Ping(context.Background())
	if err != nil {
		slog.Error("Database ping failed: " + err.Error())
		os.Exit(1)
	}
}
