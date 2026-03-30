package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"`
	ContestID uint      `gorm:"not null;index"`
	SendAt    time.Time `gorm:"not null;index"`
	Sent      bool      `gorm:"default:false"`

	User    User    `gorm:"foreignKey:UserID"`
	Contest Contest `gorm:"foreignKey:ContestID"`
}
