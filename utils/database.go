package utils

import (
	"log"
	"os"

	models "seams_go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitialiseDB() {
	LoadEnvironmentVariables()
	log.Print("Initialising Database...")

	dsn := os.Getenv("DATABASE_URL")
	log.Printf("db is running on %s", dsn)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err.Error())
	}
	log.Print("Successfully connected!")

	err = DB.AutoMigrate(models.User{})
	if err != nil {
		log.Print("Failed Migration")
		log.Panic(err.Error())
	}
}
