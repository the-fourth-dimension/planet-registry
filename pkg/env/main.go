package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	PORT = iota
	DB_HOST
	DB_DRIVER
	DB_USER
	DB_PASSWORD
	DB_NAME
	DB_PORT
	JWT_SECRET
)

var keys = []string{
	"PORT",
	"DB_HOST",
	"DB_DRIVER",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
	"DB_PORT",
	"JWT_SECRET",
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	for idx := range keys {
		key := keys[idx]
		_, isPresent := os.LookupEnv(key)
		if !isPresent {
			panic(fmt.Sprintf("missing environment variable: %s", key))
		}

	}
}

func GetEnv(key int) string {
	if key < 0 || key >= len(keys) {
		panic("Invalid env key")
	}

	return os.Getenv(keys[key])
}
