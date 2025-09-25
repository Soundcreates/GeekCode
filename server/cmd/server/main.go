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
	corsConfig := cors.Config{
    AllowOrigins:     []string{"http://localhost:5173", cfg.PROD_URL},
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))

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