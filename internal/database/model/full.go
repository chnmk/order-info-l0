package db_model

import "github.com/jackc/pgx/v5"

type PostgresDB struct {
	DB *pgx.Conn
}
