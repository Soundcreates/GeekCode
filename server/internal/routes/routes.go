package routes

import (
	"geekCode/internal/handlers"
	"geekCode/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB ,jwtSecret string) {
	api := r.Group("/api")

	//inits handlers w db
	h := handlers.NewHandler(db)

	//for heallth check
	api.GET("/ping", handlers.Ping)

	//for auth
	auth := api.Group("/auth")

	auth.POST("/login", h.Login)
	auth.POST("/register", h.Register)

	//protectedRoutes
	protected := api.Group("/")
	//profile route
	
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.GET("/profile", h.GetProfile)
	protected.POST("/rooms", h.CreateRoom)
	protected.GET("/rooms", h.ListRooms)
	protected.GET("/rooms/:id", h.GetRoom)


}