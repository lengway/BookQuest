package database

import (
	"Backend/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectDB() {
	_ = godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}

	DB = db

	AMerr := db.AutoMigrate(&models.User{})
	if AMerr != nil {
		log.Fatal("User AutoMigrate error: ", AMerr.Error())
	}

	AMerr = db.AutoMigrate(&models.Book{})
	if AMerr != nil {
		log.Fatal("Book AutoMigrate error: ", AMerr.Error())
	}

}
