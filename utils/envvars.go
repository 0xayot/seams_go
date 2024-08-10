package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	env := os.Getenv("APP_ENV")
	if env == "" || env == "development" {
		env = "development"
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
