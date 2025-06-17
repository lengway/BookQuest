package main

import (
	"Backend/database"
	router "Backend/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App Name",
	})

	database.ConnectDB()

	router.SetupRoutes(app)
	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}
