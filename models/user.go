package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
	Names    string `json:"names"`
	Role     string `gorm:"default:user" json:"role"` // либо "user", либо "admin"
}
