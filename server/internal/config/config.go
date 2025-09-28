package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBHost  string
	DBPort  string
	DBUser  string
	DBPassword string
	DBName  string
	JWTSecret string
	FRONTEND_URL string
	PROD_URL string
}

func LoadConfig() *Config {
	// Only load .env file in development (when file exists)
	// In production (like Render), environment variables are set directly
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config := &Config{
		Port:     GetEnv("PORT", "8080"),
		DBHost:  os.Getenv("DB_HOST"),
		DBPort:  os.Getenv("DB_PORT"),
		DBUser:  os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:  os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWTSecret"),
		FRONTEND_URL: os.Getenv("FRONTEND_URL"),
		PROD_URL: os.Getenv("PROD_URL"),
	}

	// Log configuration (without sensitive data)
	log.Printf("Configuration loaded - Port: %s, DB Host: %s, DB Name: %s", 
		config.Port, config.DBHost, config.DBName)

	return config
}

func GetEnv(key string, fallback string) string {
	val := os.Getenv(key)

	if(val == "" ){
		return fallback
	}
	return val

}