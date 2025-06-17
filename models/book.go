package models

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	gorm.Model
	Title       string    `gorm:"not null"`
	Author      string    `gorm:"not null"`
	Description string    `gorm:"size:255"`
	Date        time.Time `gorm:"default:null"`
}
