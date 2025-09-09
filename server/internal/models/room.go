package models

import "time"

type Room struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"uniqueIndex;not null"`
    RoomID    string    `gorm:"uniqueIndex;not null"`
    CreatedAt time.Time
    CreatedBy uint      // Match User.ID type
    Creator   User      `gorm:"foreignKey:CreatedBy;references:ID"`
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