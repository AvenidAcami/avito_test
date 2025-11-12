package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitENV() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("environment variable not found")
	}
	return value, nil
}
