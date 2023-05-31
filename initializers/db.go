package initializers

import (
	"fmt"
	"log"

	"github.com/cglavin50/go-jwt/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	//db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// url := `host=localhost user=jwtuser password=jwt db=jwt_db port=5432 sslmode=disable TimeZone='Pacific Standard Time'"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the DB")
	}

	configureTable() // configure table to match our user model

	fmt.Println("DB connection established")
}

func configureTable() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error formatting table to model", err)
	}
}

// paul was here :D
