package handlers

import (
	"geekCode/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	cfg *config.Config
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message" : "pong"})
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler{
	return &Handler{DB : db, cfg: cfg}
}

