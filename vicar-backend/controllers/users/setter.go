package users

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func setUserUsername(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	userIdStr := c.Params("id")
	if userIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is missing"})
	}

	if userIdStr != "@me" && userIdStr != user.ID.String() {
		if !user.IsTeam {
			return fiber.ErrForbidden
		}
	}

	payload := struct {
		Username string `json:"username"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	targetUser := &entities.User{}
	if userIdStr == "@me" {
		targetUser = user
	} else if err := db.DB.Where("id = ?", userIdStr).First(targetUser).Error; err != nil {
		return err
	}

	targetUser.Username = payload.Username
	if res := db.DB.Save(targetUser); res.Error != nil {
		return res.Error
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func setUserPassword(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	userIdStr := c.Params("id")
	if userIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is missing"})
	}

	if userIdStr != "@me" && userIdStr != user.ID.String() {
		if !user.IsTeam {
			return fiber.ErrForbidden
		}
	}

	payload := struct {
		Password    string `json:"password"`
		OldPassword string `json:"oldPassword"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	targetUser := &entities.User{}
	if userIdStr == "@me" {
		targetUser = user
	} else if err := db.DB.Where("id = ?", userIdStr).First(targetUser).Error; err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(payload.OldPassword)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid old password"})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	targetUser.Password = string(passwordHash)
	if res := db.DB.Save(targetUser); res.Error != nil {
		return res.Error
	}

	return c.SendStatus(fiber.StatusNoContent)
}
