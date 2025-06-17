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

	//app.Get("/", func(c *fiber.Ctx) error {
	//	return c.SendString("BookBattle API работает!")
	//})
	//
	//app.Post("/register", handlers.CreateUser)
	//app.Post("/login", handlers.Login)
	//app.Get("/me", middleware.Protected(), handlers.GetMe)

	router.SetupRoutes(app)
	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}
