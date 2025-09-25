package config

import (
	"fmt"
	"log"
	"os"

	"geekCode/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
		host, port, user, password, name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatal("Failed to connect to database ", err)
	}

	//automigrating models
	err = db.AutoMigrate(
				&models.User{},
				&models.Room{},
				// &models.Client{},
				
			)

	if err != nil {
		log.Fatal("Failed to migrate database ", err)
	}

	return db, nil

}