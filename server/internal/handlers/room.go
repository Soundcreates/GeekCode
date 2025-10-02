package handlers

import (
	"geekCode/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"

	// "gorm.io/gorm"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name   string `json:"name" binding:"required"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest

	if err := c.ShouldBindJSON(&req) ; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		log.Println(err)
		return
	}

	// Get user ID from JWT token via middleware
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	//finding the user for their name
	var user models.User
	if err  := h.DB.First(&user, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "User not found"})
		return
	}

	room  := models.Room{
		Name: req.Name,
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
		Creator: user,
		RoomID: uuid.New().String(),
	}

	if err := h.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to create room"})
		return
	}

	// member := models.Client{
	// 	RoomID: room.RoomID,
	// 	UserID: req.UserID,
	// 	JoinedAt: time.Now(),
	// }

	// if err := h.DB.Create(&member).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to add creator as member"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"roomID" : room.ID, "link" : "/code/" +fmt.Sprint(room.RoomID)})
}

func (h *Handler) ListRooms(c *gin.Context) {
	const limit = 10
	var rooms []models.Room
	userId, exists := c.Get("userId")
	if !exists {
		log.Println("User doesnt exist!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get all rooms created by user, ordered by creation date
	if err := h.DB.Where("created_by = ?", userId).Order("created_at DESC").Limit(limit).Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

// Get active rooms (rooms with status Active)
func (h *Handler) GetActiveRooms(c *gin.Context) {
	var rooms []models.Room
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	if err := h.DB.Where("created_by = ? AND status = ?", userId, models.Active).Order("created_at DESC").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active rooms!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activeRooms": rooms,
	})
}

// Get ended rooms (rooms with status Ended)
func (h *Handler) GetEndedRooms(c *gin.Context) {
	var rooms []models.Room
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	if err := h.DB.Where("created_by = ? AND status = ?", userId, models.Ended).Order("created_at DESC").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ended rooms!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endedRooms": rooms,
	})
}

// End a room (change status to Ended)
func (h *Handler) EndRoom(c *gin.Context) {
	roomId := c.Param("roomId")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	var room models.Room
	if err := h.DB.Where("room_id = ? AND created_by = ?", roomId, userId).First(&room).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found!"})
		return
	}
	
	room.Status = models.Ended
	if err := h.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end room!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room ended successfully",
		"room": room,
	})
}

func (h *Handler) GetRoom(c *gin.Context){
	roomId := c.Param("roomId")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "User not authenticated!"})
		return
	}
	var room models.Room

	if err := h.DB.Where("room_id = ? AND created_by = ?", roomId, userId).First(&room).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Room not found!"})
		return
	}
	
	//Now returning the specific rooms members
	// var roomMembers []models.Client

	// if err := h.DB.Where("RoomID = ?", roomId).Find(&roomMembers).Error ; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"Error": "The room has no members!"})
	// 	return
	// }


	c.JSON(http.StatusOK , gin.H{"room" : room})

}