 package main

import (
	
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"geekCode/internal/config"
	"geekCode/internal/routes"
	"github.com/gin-contrib/cors"
	"time"

)

func main () {
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

	db,err := config.ConnectDB(cfg)
	if err!=nil {
		log.Fatal("Failed to connect to database ", err)
	}



	routes.RegisterRoutes(r, db, cfg.JWTSecret, cfg)
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}