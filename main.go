package main

import (
	"github.com/chnmk/order-info-l0/internal/broker"
	"github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/test"
)

func init() {
	config.SetDefaultEnv()
	config.GetEnv()
	if config.EnvVariables["PUBLISH_TEST_DATA"] == "1" {
		test.PublishTestData()
	}
}

func main() {
	broker.Consume()
}
