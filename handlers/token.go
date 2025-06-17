package handlers

//
//import (
//	"Backend/config"
//	"Backend/database"
//	"Backend/models"
//	"github.com/gofiber/fiber/v2"
//	"github.com/golang-jwt/jwt/v5"
//	"time"
//)
//
//func generateAccessToken(user *models.User) (string, error) {
//	token := jwt.New(jwt.SigningMethodHS256)
//	claims := token.Claims.(jwt.MapClaims)
//	claims["user_id"] = user.ID
//	claims["username"] = user.Username
//	claims["role"] = user.Role // добавь поле Role в модель
//	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
//	claims["iat"] = time.Now().Unix()
//	return token.SignedString([]byte(config.Config("SECRET")))
//}
//
//func generateRefreshToken(user *models.User) (string, error) {
//	token := jwt.New(jwt.SigningMethodHS256)
//	claims := token.Claims.(jwt.MapClaims)
//	claims["user_id"] = user.ID
//	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
//	return token.SignedString([]byte(config.Config("REFRESH_SECRET")))
//}
//
//func RefreshToken(c *fiber.Ctx) error {
//	user := c.Locals("user").(*jwt.Token)
//	claims := user.Claims.(jwt.MapClaims)
//	userID := uint(claims["user_id"].(float64))
//
//	var u models.User
//	if err := database.DB.First(&u, userID).Error; err != nil {
//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User not found"})
//	}
//
//	newToken := jwt.New(jwt.SigningMethodHS256)
//	newClaims := newToken.Claims.(jwt.MapClaims)
//	newClaims["user_id"] = u.ID
//	newClaims["username"] = u.Username
//	newClaims["exp"] = time.Now().Add(time.Hour * 72).Unix()
//
//	tokenString, err := newToken.SignedString([]byte(config.Config("SECRET")))
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create token"})
//	}
//
//	return c.JSON(fiber.Map{
//		"status": "success",
//		"token":  tokenString,
//	})
//}
