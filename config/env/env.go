package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// this is public method
func LoadEnv() {
	// Find/load .env file
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}


func GetString(key string, fallBack string) string {
	// Getting and using a value from .env
	// value := os.Getenv(key)

	value, ok := os.LookupEnv(key)

	if !ok {
		return fallBack
	}

	return value
}

func GetInt(key string, fallBack int) int {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallBack
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Error converting %s to int: %v\n", key, err)
		return fallBack
	}

	return intValue
}

func GetBool(key string, fallBack bool) bool {

	value, ok := os.LookupEnv(key)

	if !ok {
		return fallBack
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		fmt.Printf("Error converting %s to int: %v\n", key, err)
		return fallBack
	}

	return boolValue
}
