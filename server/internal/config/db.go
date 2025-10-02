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
	var dsn string
	
	// Check if DATABASE_URL is provided (Railway default)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Printf("Using DATABASE_URL for connection...")
		dsn = databaseURL
	} else {
		// Fallback to individual environment variables
		host := GetEnv("DB_HOST", "localhost")
		port := GetEnv("DB_PORT", "5432")
		user := GetEnv("DB_USER", "postgres")
		password := GetEnv("DB_PASSWORD", "")
		name := GetEnv("DB_NAME", "geekcode")
		sslMode := GetEnv("SSL_MODE", "disable")
		
		// Check if required environment variables are set
		if password == "" {
			log.Printf("ERROR: DB_PASSWORD environment variable is required")
			return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
		}
		
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, name, sslMode)
		log.Printf("Using individual DB variables for connection...")
	}
	
	log.Printf("Attempting to connect to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}
	
	log.Printf("Successfully connected to database!")

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Client{},
	); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}

	log.Printf("Database migration completed successfully!")
	
	return db, nil
}