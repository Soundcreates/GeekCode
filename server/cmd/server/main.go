package main

import (
	"geekCode/internal/config"
	"geekCode/internal/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main () {
	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "" && os.Getenv("PORT") != "" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	cfg := config.LoadConfig()
	r := gin.Default()

	//handling cors
	// Build allowed origins list
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:5174"}
	
	// Add PROD_URL if it's not empty
	if cfg.PROD_URL != "" {
		allowedOrigins = append(allowedOrigins, cfg.PROD_URL)
	}
	
	// Add FRONTEND_URL if it's not empty and different from PROD_URL
	if cfg.FRONTEND_URL != "" && cfg.FRONTEND_URL != cfg.PROD_URL {
		allowedOrigins = append(allowedOrigins, cfg.FRONTEND_URL)
	}
	
	// If no production URLs are set, use a more permissive approach for development
	if len(allowedOrigins) <= 2 { // Only localhost URLs
		// For development, we'll disable credentials and allow all origins
		corsConfig := cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false, // Must be false when using "*"
			MaxAge:           12 * time.Hour,
		}
		r.Use(cors.New(corsConfig))
	} else {
		// For production with specific origins, we can use credentials
		corsConfig := cors.Config{
			AllowOrigins:     allowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
		r.Use(cors.New(corsConfig))
	}

	db,err := config.ConnectDB()
	if err!=nil {
		log.Fatal("Failed to connect to database ", err)
	}

	if cfg.JWTSecret == "" {
		log.Fatal("Jwtsecret is needed but hasn't been set up yet")
	}

	routes.RegisterRoutes(r, db, cfg.JWTSecret)
	
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}