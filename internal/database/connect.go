package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect() *pgx.Conn {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.EnvVariables["DB_USER"],
		config.EnvVariables["DB_PASSWORD"],
		"postgres",
		"5432",
		config.EnvVariables["DB_NAME"],
	)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		slog.Error("Unable to connect to database: " + err.Error())
	}

	return conn
}

func Ping(conn *pgx.Conn) {
	err := conn.Ping(context.Background())
	if err != nil {
		slog.Error("QueryRow failed: " + err.Error())
	}
}
