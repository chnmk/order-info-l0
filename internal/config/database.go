package config

import "fmt"

var PgxpoolUrl string

func getDatabaseVars() {
	PgxpoolUrl = getPgxpoolUrl()
}

func getPgxpoolUrl() string {
	pgxpoolUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		Env.Get("DB_PROTOCOL"),
		Env.Get("POSTGRES_USER"),
		Env.Get("POSTGRES_PASSWORD"),
		Env.Get("DB_HOST"),
		Env.Get("DBS_PORT"),
		Env.Get("POSTGRES_DB"),
	)

	return pgxpoolUrl
}
