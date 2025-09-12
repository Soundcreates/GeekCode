package routes

import (
	"geekCode/internal/config"
	"geekCode/internal/handlers"
	"geekCode/internal/middleware"

	"geekCode/internal/ws"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB ,jwtSecret string, cfg *config.Config) {
	api := r.Group("/api")

	//inits handlers w db
	h := handlers.NewHandler(db, cfg)

	//for heallth check
	api.GET("/ping", handlers.Ping)
	api.GET("/ws/:roomId", ws.HandleWebSocket) //websocket route

	//for auth
	auth := api.Group("/auth")

	auth.POST("/login", h.Login)
	auth.POST("/register", h.Register)

	//protectedRoutes
	protected := api.Group("/")
	//profile route
	
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.GET("/profile", h.GetProfile)

	//room routes
	protected.POST("/rooms", h.CreateRoom)        // Create room
	protected.GET("/rooms", h.ListRooms)          // List user's rooms
	protected.GET("/rooms/:roomId", h.GetRoom)    // Get specific room



}