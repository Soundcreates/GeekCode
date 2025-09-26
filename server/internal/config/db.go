package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"geekCode/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	databaseURL := os.Getenv("POSTGRES_URL") 
	var dsn string

	if databaseURL != "" {
		// Fix Render's format: "postgresql://" → "postgres://"
		dsn = strings.Replace(databaseURL, "postgresql://", "postgres://", 1)
		// Ensure SSL mode is enforced
		if !strings.Contains(dsn, "sslmode=") {
			dsn += "?sslmode=require"
		}
	} else {
		// Local development: load from .env or config
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")

		if host == "" || port == "" || user == "" || password == "" || name == "" {
			log.Fatal("Database configuration incomplete. Please set DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, and DB_NAME")
		}

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, name,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Room{},
	); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	return db, nil
}
