package models

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	BookID    uint       `gorm:"foreignkey:BookID"`
	ChapterID uint       `gorm:"foreignkey:ChapterID"`
	Title     string     `gorm:"type:varchar(255)"`
	Question  []Question `gorm:"foreignkey:QuizID"`
}
