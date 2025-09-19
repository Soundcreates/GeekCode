package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"` //corrected from Lastname to LastName
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCr]eateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
