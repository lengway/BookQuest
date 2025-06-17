package handlers

import (
	"Backend/database"
	"Backend/models"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

// CreateChapterDTO defines the structure for creating a new chapter.
type CreateChapterDTO struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required"` // Assuming chapters have content
}

// UpdateChapterDTO defines the structure for updating an existing chapter.
type UpdateChapterDTO struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content *string `json:"content,omitempty"`
}

// ChapterResponseDTO defines the structure for chapter responses.
// We might want to exclude/include certain fields from models.Chapter in the future.
type ChapterResponseDTO struct {
	ID        uint   `json:"id"`
	BookID    uint   `json:"book_id"`
	Title     string `json:"title"`
	Content   string `json:"content"` // Assuming content is part of the response
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Helper to convert models.Chapter to ChapterResponseDTO
func toChapterResponseDTO(chapter models.Chapter) ChapterResponseDTO {
	return ChapterResponseDTO{
		ID:        chapter.ID,
		BookID:    chapter.BookID,
		Title:     chapter.Title,
		Content:   chapter.Content,
		CreatedAt: chapter.CreatedAt.String(),
		UpdatedAt: chapter.UpdatedAt.String(),
	}
}

// CreateChapter creates a new chapter for a book.
// @Summary Create a new chapter
// @Description Adds a new chapter to a specified book. Requires admin privileges.
// @Tags Chapters
// @Accept json
// @Produce json
// @Param bookId path int true "Book ID"
// @Param chapter body CreateChapterDTO true "Chapter data"
// @Success 201 {object} ChapterResponseDTO
// @Failure 400 {object} fiber.Map "Invalid input or validation error"
// @Failure 404 {object} fiber.Map "Book not found"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/books/{bookId}/chapters [post]
func CreateChapter(c *fiber.Ctx) error {
	bookIDStr := c.Params("bookId")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Book ID provided",
			"errors":  err.Error(),
		})
	}

	// Check if book exists
	var book models.Book
	if err := database.DB.First(&book, uint(bookID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Book not found with ID: " + bookIDStr,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error checking for book existence",
			"errors":  err.Error(),
		})
	}

	input := new(CreateChapterDTO)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error parsing request body",
			"errors":  err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	chapter := models.Chapter{
		BookID:  uint(bookID),
		Title:   input.Title,
		Content: input.Content, // Assuming models.Chapter has a Content field
	}

	if err := database.DB.Create(&chapter).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create chapter",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(toChapterResponseDTO(chapter))
}

// GetBookChapters retrieves all chapters for a specific book.
// @Summary Get all chapters for a book
// @Description Retrieves a list of all chapters associated with a given book ID.
// @Tags Chapters
// @Produce json
// @Param bookId path int true "Book ID"
// @Success 200 {array} ChapterResponseDTO
// @Failure 404 {object} fiber.Map "Book not found (optional check)"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/books/{bookId}/chapters [get]
func GetBookChapters(c *fiber.Ctx) error {
	bookIDStr := c.Params("bookId")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Book ID format",
			"errors":  err.Error(),
		})
	}

	// Optional: Check if book exists to return a more specific 404 for the book.
	var book models.Book
	if err := database.DB.First(&book, uint(bookID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Book not found with ID: " + bookIDStr,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error checking book existence",
			"errors":  err.Error(),
		})
	}

	var chapters []models.Chapter
	if err := database.DB.Where("book_id = ?", uint(bookID)).Order("created_at asc").Find(&chapters).Error; err != nil {
		// It's unlikely to get gorm.ErrRecordNotFound here with Find, it would return an empty slice.
		// But good to handle general DB errors.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve chapters for book",
			"errors":  err.Error(),
		})
	}

	responseChapters := make([]ChapterResponseDTO, len(chapters))
	for i, chapter := range chapters {
		responseChapters[i] = toChapterResponseDTO(chapter)
	}

	return c.Status(fiber.StatusOK).JSON(responseChapters)
}

// GetChapter retrieves a specific chapter by its ID.
// @Summary Get a specific chapter
// @Description Retrieves details of a chapter by its ID.
// @Tags Chapters
// @Produce json
// @Param chapterId path int true "Chapter ID"
// @Success 200 {object} ChapterResponseDTO
// @Failure 404 {object} fiber.Map "Chapter not found"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/chapters/{chapterId} [get]
func GetChapter(c *fiber.Ctx) error {
	chapterIDStr := c.Params("chapterId")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Chapter ID format",
			"errors":  err.Error(),
		})
	}

	var chapter models.Chapter
	if err := database.DB.First(&chapter, uint(chapterID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Chapter not found with ID: " + chapterIDStr,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error retrieving chapter",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(toChapterResponseDTO(chapter))
}

// UpdateChapter updates an existing chapter.
// @Summary Update a chapter
// @Description Updates details of an existing chapter by its ID. Requires admin privileges.
// @Tags Chapters
// @Accept json
// @Produce json
// @Param chapterId path int true "Chapter ID"
// @Param chapter body UpdateChapterDTO true "Chapter data to update"
// @Success 200 {object} ChapterResponseDTO
// @Failure 400 {object} fiber.Map "Invalid input or validation error"
// @Failure 404 {object} fiber.Map "Chapter not found"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/chapters/{chapterId} [put]
func UpdateChapter(c *fiber.Ctx) error {
	chapterIDStr := c.Params("chapterId")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Chapter ID format",
			"errors":  err.Error(),
		})
	}

	var chapter models.Chapter
	if err := database.DB.First(&chapter, uint(chapterID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Chapter not found with ID: " + chapterIDStr,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error retrieving chapter for update",
			"errors":  err.Error(),
		})
	}

	input := new(UpdateChapterDTO)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error parsing request body",
			"errors":  err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	if input.Title != nil {
		chapter.Title = *input.Title
	}
	if input.Content != nil {
		chapter.Content = *input.Content
	}
	// Note: BookID is not updatable through this endpoint by design.

	if err := database.DB.Save(&chapter).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update chapter",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(toChapterResponseDTO(chapter))
}

// DeleteChapter deletes a chapter by its ID.
// @Summary Delete a chapter
// @Description Deletes a chapter by its ID. Requires admin privileges.
// @Tags Chapters
// @Produce json
// @Param chapterId path int true "Chapter ID"
// @Success 200 {object} fiber.Map "Chapter deleted successfully"
// @Failure 204 {object} fiber.Map "Chapter deleted successfully (alternative)"
// @Failure 404 {object} fiber.Map "Chapter not found"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/chapters/{chapterId} [delete]
func DeleteChapter(c *fiber.Ctx) error {
	chapterIDStr := c.Params("chapterId")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Chapter ID format",
			"errors":  err.Error(),
		})
	}

	var chapter models.Chapter
	// First, check if the chapter exists
	if err := database.DB.First(&chapter, uint(chapterID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Chapter not found with ID: " + chapterIDStr,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error checking chapter existence for deletion",
			"errors":  err.Error(),
		})
	}

	// If found, proceed with deletion
	if err := database.DB.Delete(&chapter).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete chapter",
			"errors":  err.Error(),
		})
	}

	// Return 200 with a message or 204 No Content
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Chapter deleted successfully",
	})
	// Alternatively, for 204 No Content:
	// return c.SendStatus(fiber.StatusNoContent)
}
