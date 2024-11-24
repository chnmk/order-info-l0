package db_model

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/jackc/pgx/v5"
)

func (db *PostgresDB) Connect(ctx context.Context) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.EnvVariables["DB_USER"],
		config.EnvVariables["DB_PASSWORD"],
		"postgres",
		"5432",
		config.EnvVariables["DB_NAME"],
	)

	var err error

	db.DB, err = pgx.Connect(ctx, url)
	if err != nil {
		slog.Error("Unable to connect to database: " + err.Error())
	}
}

func (db *PostgresDB) Close(ctx context.Context) {
	db.DB.Close(ctx)
}

func (db *PostgresDB) Ping() {
	err := db.DB.Ping(context.Background())
	if err != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}
}
