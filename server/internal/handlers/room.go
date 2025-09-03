package handlers


import (
	"github.com/gin-gonic/gin"
	"net/http"
	"geekCode/internal/models"
// "gorm.io/gorm"
"time"
"fmt"
"log"
)

type CreateRoomRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID uint   `json:"userId" binding:"required"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest

	if err := c.ShouldBindJSON(&req) ; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	room  := models.Room{
		Name: req.Name,
		CreatedBy: req.UserID,
		CreatedAt: time.Now(),
	}

	if err := h.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to create room"})
		return
	}

	member := models.RoomMember {
		RoomID: room.ID,
		UserID: req.UserID,
		JoinedAt: time.Now(),
	}

	if err := h.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to add creator as member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roomID" : room.ID, "link" : "/rooms/" +fmt.Sprint(room.ID)})
}

func (h *Handler) ListRooms(c *gin.Context) {
	const limit = 5
	var rooms []models.Room
	userId, exists := c.Get("userId")
	if !exists {
		log.Println("User doesnt exist!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
if err := h.DB.Where("CreatedBy = ?", userId).Limit(limit).Find(&rooms).Error; err !=nil {
	c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch rooms!"})
	return
}

	c.JSON(http.StatusOK, rooms)


}

func (h *Handler) GetRoom(c *gin.Context){
	roomId := c.Param("roomId")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "User not authenticated!"})
		return
	}
	var room models.Room

	if err := h.DB.Where("ID = ? AND CreatedBy = ?", roomId,userId).First(&room).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Room not found!"})
		return
	}
	
	//Now returning the specific rooms members
	var roomMembers []models.RoomMember

	if err := h.DB.Where("RoomID = ?", roomId).Find(&roomMembers).Error ; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "The room has no members!"})
		return
	}


	c.JSON(http.StatusOK , gin.H{"room" : room, "members" : roomMembers})

}