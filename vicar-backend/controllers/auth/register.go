package auth

import (
	"vicar-backend/db"
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func createAccount(c *fiber.Ctx) error {
	dto := &registerDto{}
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request: " + err.Error(),
		})
	}

	if dto.Username == "" || dto.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request: missing required fields",
		})
	}

	count := int64(0)
	if err := db.DB.Model(&entities.User{}).Where("LOWER(username) = LOWER(?)", dto.Username).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Conflict: username already exists",
		})
	}

	user := &entities.User{}
	user.Username = dto.Username

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(passwordHash)

	if err := db.DB.Save(user).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Created: User created and invite code consumed",
	})
}
