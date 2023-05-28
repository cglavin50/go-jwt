package initializers

import (
	"fmt"
	"log"

	"github.com/cglavin50/go-jwt/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(dsn string) *gorm.DB {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the DB")
	}

	configureTable() // configure table to match our user model

	fmt.Println("DB connection established")
	return db
}

func configureTable() {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error formatting table to model", err)
	}
}

// paul was here :D
