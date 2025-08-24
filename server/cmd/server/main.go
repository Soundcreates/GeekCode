 package main

import (
	
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"geekCode/internal/config"
	"geekCode/internal/routes"
)

func main () {
	cfg := config.LoadConfig()

	db,err := config.ConnectDB(cfg)
	if err!=nil {
		log.Fatal("Failed to connect to database ", err)
	}

	r:=gin.Default()

	routes.RegisterRoutes(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}