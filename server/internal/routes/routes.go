package routes


import (

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"geekCode/internal/handlers"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	api.GET("/ping", handlers.Ping)

}