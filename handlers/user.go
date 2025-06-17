package handlers

import (
	"Backend/database"
	"Backend/models"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB
	var user models.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Теперь можно найти юзера в БД
	var u models.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":   "success",
		"username": u.Username,
		"email":    u.Email,
		"user_id":  u.ID,
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found with ID: " + id,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error while fetching user",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User found",
		"data":    user})
}

func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username" validate:"required,min=3,max=32"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	db := database.DB
	input := new(NewUser) // Parse into DTO

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // Changed status to 400
			"status":  "error",
			"message": "Error parsing request body", // Clearer message
			"errors":  err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",             // Added status
			"message": "Validation failed", // Clearer message
			"errors":  err.Error()})
	}

	// Password length check is implicitly handled by validate:"min=8" and bcrypt's limit.
	// Explicit check for > 72 might still be useful if bcrypt had a lower internal limit than validator's max.
	// For now, relying on bcrypt's error handling during hashPassword if it's too long.

	if existingUser, _ := getUserByEmail(input.Email); existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ // Changed status to 409
			"status":  "error", // Added status
			"message": "Email already exists",
			"errors":  "A user with this email already exists",
		})
	}

	// getUserByUsername was checking input.Email, should be input.Username
	if existingUser, _ := getUserByUsername(input.Username); existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ // Changed status to 409
			"status":  "error", // Added status
			"message": "Username already exists",
			"errors":  "A user with this username already exists",
		})
	}

	hash, err := hashPassword(input.Password)
	if err != nil {
		// This could happen if password is too long for bcrypt (e.g., > 72 bytes)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't hash password",
			"errors":  err.Error()})
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hash,
		// Role will be its zero value (e.g., empty string or specific default if set in model)
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // Changed status to 500
			"status":  "error",
			"message": "Couldn't create user",
			"errors":  err.Error()})
	}

	// Return DTO without password
	responseUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{ // Changed status to 201
		"status":  "success",
		"message": "User created successfully", // Clearer message
		"data":    responseUser})
}

func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error()})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token id",
			"data":    nil})
	}

	db := database.DB
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found with ID: " + id,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error while fetching user for update",
			"errors":  err.Error(),
		})
	}

	user.Names = uui.Names
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update user",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully updated",
		"data":    user})
}

func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error()})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token id",
			"data":    nil})
	}

	if !validUser(id, pi.Password) {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Not valid user",
			"data":    nil})
	}

	db := database.DB
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found with ID: " + id,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error while fetching user for deletion",
			"errors":  err.Error(),
		})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete user",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully deleted",
		"data":    nil})
}
