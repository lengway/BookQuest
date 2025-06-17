package models

import "gorm.io/gorm"

type Chapter struct {
	gorm.Model
	Title  string
	BookID uint
}
