package router

import (
	"Backend/handlers"
	"Backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handlers.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)

	// User
	user := api.Group("/user")
	user.Get("/me", middleware.Protected(), handlers.GetMe)
	user.Get("/:id", middleware.IsOwnerOrAdmin(), handlers.GetUser)
	user.Post("/register", handlers.CreateUser)
	user.Patch("/:id", middleware.IsOwnerOrAdmin(), handlers.UpdateUser)
	user.Delete("/:id", middleware.IsOwnerOrAdmin(), handlers.DeleteUser)

	//book
	book := api.Group("/book")

	book.Use(middleware.Protected()) // только авторизованные пользователи могут получить доступ к книгам

	book.Get("/", handlers.GetAllBooks)
	book.Get("/:id", handlers.GetBookById)
	book.Post("/", middleware.IsOwnerOrAdmin(), handlers.CreateBook)      // теперь только админы могут получить
	book.Delete("/:id", middleware.IsOwnerOrAdmin(), handlers.DeleteBook) // доступ к эти рутам

	token := api.Group("/token")
	token.Get("/refresh", middleware.Protected(), handlers.Refresh)
}
