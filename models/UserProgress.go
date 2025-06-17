package models

import "gorm.io/gorm"

type UserProgress struct {
	gorm.Model
	UserID    uint `gorm:"foreignkey:UserID"`
	BookID    uint `gorm:"foreignkey:BookID"`
	ChapterID uint `gorm:"foreignkey:ChapterID"`
	Completed bool `gorm:"default:false"`
}
