package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var EnvVariables map[string]string

func SetDefaultEnv() {
	EnvVariables = make(map[string]string)
	EnvVariables["PUBLISH_TEST_DATA"] = "0"
	EnvVariables["DB_NAME"] = "orders"
	EnvVariables["DB_USER"] = "user"
	EnvVariables["DB_PASSWORD"] = "12345"
}

func GetEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Warning: .env file not found")
	}

	log.Print("Reading environment variables...")

	for name, def := range EnvVariables {
		variable, exists := os.LookupEnv(name)
		if exists {
			log.Printf("%s: %s", name, variable)
			EnvVariables[name] = variable
		} else {
			log.Printf("%s not found, using default (%s)", name, def)
		}
	}
}
