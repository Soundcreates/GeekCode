package models

import "time"

type Status int

const (
    Active Status = iota
    Ended
    Archived
)

type Room struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    RoomID    string    `gorm:"uniqueIndex;not null"`
    CreatedAt time.Time
    CreatedBy string    // the username of the user
    Creator   User      `gorm:"foreignKey:CreatedBy;references:ID"`
    Status  Status    `gorm:"default:0"` // Default to Active
}

type Client struct {
    ID       uint      `gorm:"primaryKey"`
    RoomID   string    `gorm:"not null"` // This references Room.RoomID
    UserID   uint      `gorm:"not null"` // This references User.ID
    JoinedAt time.Time

    // Correct foreign key relationships
    Room Room `gorm:"foreignKey:RoomID;references:RoomID"`
    User User `gorm:"foreignKey:UserID;references:ID"`
}