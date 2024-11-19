package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/chnmk/order-info-l0/internal/broker"
	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/chnmk/order-info-l0/test"
)

func init() {
	// TEMP
	var E models.Order

	content, err := os.ReadFile("test/model.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &E)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

	// ===================

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
