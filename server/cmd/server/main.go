 package main

import (
	
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"geekCode/internal/config"
	"geekCode/internal/routes"
	"github.com/gin-contrib/cors"
)

func main () {
	cfg := config.LoadConfig()

	//handling cors
	corsConfig := cors.DefaultConfig() //the default config function only returns one variable 



	corsConfig.AllowOrigins  = []string{cfg.FRONTEND_URL}

	db,err := config.ConnectDB(cfg)
	if err!=nil {
		log.Fatal("Failed to connect to database ", err)
	}

	r := gin.Default()

	routes.RegisterRoutes(r, db, cfg.JWTSecret)
	r.Use(cors.New(config))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}