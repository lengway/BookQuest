package database

import (
	"Backend/config"
	"Backend/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() {
	dsn := config.AppConfig.DBSource
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}

	DB = db

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("User AutoMigrate error: ", err.Error())
	}

	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal("Book AutoMigrate error: ", err.Error())
	}

	err = db.AutoMigrate(&models.Chapter{})
	if err != nil {
		log.Fatal("Chapter AutoMigrate error: ", err.Error())
	}

	err = db.AutoMigrate(&models.Question{})
	if err != nil {
		log.Fatal("Question AutoMigrate error: ", err.Error())
	}

	err = db.AutoMigrate(&models.Quiz{})
	if err != nil {
		log.Fatal("Quiz AutoMigrate error: ", err.Error())
	}

	err = db.AutoMigrate(&models.UserProgress{})
	if err != nil {
		log.Fatal("UserProgress AutoMigrate error: ", err.Error())
	}

	err = db.AutoMigrate(&models.RefreshToken{})
	if err != nil {
		log.Fatal("RefreshToken AutoMigrate error: ", err.Error())
	}
}
