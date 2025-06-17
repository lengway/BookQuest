package handlers

import (
	"Backend/database"
	"Backend/models"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
		Date        time.Time `json:"date" validate:"date"`
	}

	db := database.DB
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Пересмотри инпут",
			"errors":  err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"errors":  err.Error()})
	}

	if existingBook, err := IsBookAlreadyExists(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while checking book existence",
			"errors":  err.Error(),
		})
	} else if existingBook != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Book already exists",
			"errors":  "A book with the same title and author already exists",
		})
	}

	if err := db.Create(&book).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create book",
			"errors":  err.Error()})
	}

	newBook := NewBook{
		Author:      book.Author,
		Title:       book.Title,
		Description: book.Description,
		Date:        book.Date}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Created book",
		"data":    newBook})
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
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	if userRole != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only admins can delete books",
		})
	}

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
