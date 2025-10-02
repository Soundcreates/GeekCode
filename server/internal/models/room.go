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
    CreatedBy uint      `gorm:"not null"` // the user ID who created the room
    Creator   User      `gorm:"foreignKey:CreatedBy;references:ID"`
    Status    Status    `gorm:"default:0"` // Default to Active
}

type Client struct {
    ID       uint      `gorm:"primaryKey"`
    RoomID   string    `gorm:"not null"` // This references Room.RoomID
    UserID   uint      `gorm:"not null"` // This references User.ID
    JoinedAt time.Time

    // Only reference User to avoid circular dependency
    User User `gorm:"foreignKey:UserID;references:ID"`
}