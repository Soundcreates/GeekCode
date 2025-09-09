package config

import (
	"fmt"
	"log"

	"geekCode/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

func ConnectDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
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