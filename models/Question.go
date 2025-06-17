package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	QuizID uint   `gorm:"foreignkey:QuizID"`
	Text   string `gorm:"type:varchar(255)"`
	Answer string `gorm:"type:varchar(255)"`
}
