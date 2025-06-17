package handlers

import (
	"Backend/database"
	"Backend/models"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func IsBookAlreadyExists(book *models.Book) (*models.Book, error) {
	db := database.DB
	var existingBook models.Book
	if err := db.Where("title = ? AND author = ?", book.Title, book.Author).First(&existingBook).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &existingBook, nil
}

func CreateBook(c *fiber.Ctx) error {
	type NewBook struct {
		Title       string    `json:"title" validate:"required,max=255"`
		Author      string    `json:"author" validate:"required,max=255"`
		Description string    `json:"description" validate:"max=255"`
		Date        time.Time `json:"date"` // Removed validate:"date"
	}

	db := database.DB
	input := new(NewBook) // Parse into DTO

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error parsing request body", // Clearer message
			"errors":  err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed", // Clearer message
			"errors":  err.Error()})
	}

	// Create models.Book instance from DTO
	bookToCreate := &models.Book{
		Title:       input.Title,
		Author:      input.Author,
		Description: input.Description,
		Date:        input.Date,
		// UserID will be set if/when authentication is added and book creation is tied to a user
	}

	if existingBook, err := IsBookAlreadyExists(bookToCreate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error", // Added status
			"message": "Error while checking book existence",
			"errors":  err.Error(),
		})
	} else if existingBook != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ // Changed status to 409
			"status":  "error", // Added status
			"message": "Book already exists",
			"errors":  "A book with the same title and author already exists",
		})
	}

	if err := db.Create(&bookToCreate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // Changed status to 500
			"status":  "error",
			"message": "Couldn't create book",
			"errors":  err.Error()})
	}

	// Return the DTO in response
	responseBook := NewBook{
		Author:      bookToCreate.Author,
		Title:       bookToCreate.Title,
		Description: bookToCreate.Description,
		Date:        bookToCreate.Date,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{ // Changed status to 201
		"status":  "success",
		"message": "Book created successfully", // Clearer message
		"data":    responseBook})
}

func GetAllBooks(c *fiber.Ctx) error {
	db := database.DB
	var books []models.Book
	if err := db.Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve books",
			"errors":  err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved all books",
		"data":    books})
}

func GetBookById(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var book models.Book
	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve book",
			"errors":  err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved book",
		"data":    book})
}

func DeleteBook(c *fiber.Ctx) error {
	// The IsAdmin middleware now handles the role check.
	// user := c.Locals("user").(*jwt.Token) // No longer needed if only role was used
	// claims := user.Claims.(jwt.MapClaims) // No longer needed
	// userRole := claims["role"].(string) // No longer needed
	//
	// if userRole != "admin" { // This check is now handled by IsAdmin middleware
	// return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// "status": "error",
	// "message": "Only admins can delete books",
	// })
	// }

	id := c.Params("id")
	db := database.DB
	var book models.Book

	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "No book with this id found",
			"errors":  err.Error()})
	}
	if err := db.Delete(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete book",
			"errors":  err.Error()})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Book deleted successfully",
		})
	}
}
