package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/models"
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

	// time.Sleep(25 * time.Second)

	// ===================

	config.SetDefaultEnv()
	config.GetEnv()
	if config.EnvVariables["PUBLISH_TEST_DATA"] == "1" {
		// test.PublishTestData()
	}
}

func main() {
	// broker.Consume()

	database.DB = database.Connect()
	defer database.DB.Close(context.Background())

	database.Ping(database.DB)
	database.CreateTables(database.DB)

}
