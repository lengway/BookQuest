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
	book.Post("/", middleware.IsAdmin(), handlers.CreateBook)      // Changed to IsAdmin
	book.Delete("/:id", middleware.IsAdmin(), handlers.DeleteBook) // Changed to IsAdmin

	// Chapter routes
	// POST /api/books/:bookId/chapters - Create a new chapter for a book (Admin only)
	// GET /api/books/:bookId/chapters - Get all chapters for a book (Protected)
	// GET /api/chapters/:chapterId - Get a specific chapter (Protected)
	// PUT /api/chapters/:chapterId - Update a chapter (Admin only)
	// DELETE /api/chapters/:chapterId - Delete a chapter (Admin only)

	// Routes for chapters related to a specific book
	bookChapters := api.Group("/books/:bookId/chapters")
	bookChapters.Post("/", middleware.Protected(), middleware.IsAdmin(), handlers.CreateChapter) // Added Protected for consistency, IsAdmin will be primary
	bookChapters.Get("/", middleware.Protected(), handlers.GetBookChapters)

	// Routes for individual chapters (not nested under book for GET, PUT, DELETE by specific chapter ID)
	chapters := api.Group("/chapters")
	chapters.Get("/:chapterId", middleware.Protected(), handlers.GetChapter)
	chapters.Put("/:chapterId", middleware.Protected(), middleware.IsAdmin(), handlers.UpdateChapter)    // Added Protected
	chapters.Delete("/:chapterId", middleware.Protected(), middleware.IsAdmin(), handlers.DeleteChapter) // Added Protected

	token := api.Group("/token")
	token.Get("/refresh", middleware.Protected(), handlers.Refresh)
}
