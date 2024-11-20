package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/database"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/transport"
	"github.com/chnmk/order-info-l0/test"
)

func init() {
	// TEMP
	content, err := os.ReadFile("test/model.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &test.E)
	if err != nil {
		log.Fatal(err)
	}

	// ===================

	config.SetDefaultEnv()
	config.GetEnv()
	/*
		if config.EnvVariables["PUBLISH_TEST_DATA"] == "1" {
			test.PublishTestData()
		}
	*/
}

func main() {
	database.DB = database.Connect()
	defer database.DB.Close(context.Background())

	database.Ping(database.DB)
	database.CreateTables(database.DB)

	memory.DATA = database.RestoreData(database.DB)

	// TEMP
	if len(memory.DATA) == 0 {
		database.InsertOrder(database.DB, test.E)
	}

	// ===================

	// broker.Consume()

	http.HandleFunc("/order", transport.GetOrder)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
