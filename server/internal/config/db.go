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
	// Check if DATABASE_URL or POSTGRES_URL is provided (Render uses this)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("POSTGRES_URL")
	}
	
	var dsn string
	if databaseURL != "" {
		// Use DATABASE_URL directly (Render PostgreSQL format)
		dsn = databaseURL
	} else {
		// Build DSN from individual environment variables (local development)
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		
		if host == "" || port == "" || user == "" || password == "" || name == "" {
			log.Fatal("Database configuration incomplete. Please set DATABASE_URL or individual DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, and DB_NAME environment variables")
		}
		
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			host, port, user, password, name)
	}

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