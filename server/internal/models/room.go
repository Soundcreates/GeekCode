package models

import "time"

type Room struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	CreatedBy uint
	Creator User `gorm:"foreignKey:CreatedBy"`
}

type RoomMember struct{
	ID uint `gorm:"primaryKey"`
	RoomID uint
	UserID uint
	JoinedAt time.Time

	Room Room `gorm:"foreignKey:RoomID"`
	User User `gorm:"foreignKey:UserID"`
}

