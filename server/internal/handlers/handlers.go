package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message" : "pong"})
}

func NewHandler(db *gorm.DB) *Handler{
	return &Handler{DB : db}
}

