package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chnmk/order-info-l0/internal/broker"
	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/test"
)

func init() {
	fmt.Fprintf(os.Stderr, "Fail")
	os.Exit(1)

	config.SetDefaultEnv()
	config.GetEnv()
	if config.EnvVariables["PUBLISH_TEST_DATA"] == "1" {
		test.PublishTestData()
	}
}

func main() {
	broker.Consume()

	db_conn := database.Connect()
	defer db_conn.Close(context.Background())

	database.Ping(db_conn)
}
