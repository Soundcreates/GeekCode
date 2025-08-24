package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message" : "pong",
		})
	})

	router.Run() //listens and serves on the localhost
}