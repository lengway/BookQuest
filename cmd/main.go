package main

import (
	"Backend/config"
	"Backend/database"
	router "Backend/routes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	config.LoadConfig()

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "BookQuest",
	})

	database.ConnectDB()

	router.SetupRoutes(app)
	// Запуск сервера
	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port)))
}
