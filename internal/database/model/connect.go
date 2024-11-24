package db_model

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	DB *pgxpool.Pool
}

func (db *PostgresDB) Close() {
	db.DB.Close()
}

func (db *PostgresDB) Ping() {
	err := db.DB.Ping(context.Background())
	if err != nil {
		slog.Error("Database ping failed: " + err.Error())
		os.Exit(1)
	}
}
