package characters

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func addViewer(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	character, err := RequireOwner(c, user, charID)
	if err != nil {
		return err
	}

	var dto AddViewerDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	targetUser := &entities.User{}
	if err := db.DB.First(targetUser, dto.UserID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	if targetUser.ID == user.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot add owner as viewer"})
	}

	var count int64
	if err := db.DB.Table("v5_character_viewers").
		Where("v5_character_id = ? AND user_id = ?", charID, dto.UserID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already a viewer"})
	}

	if err := db.DB.Model(character).Association("Viewers").Append(targetUser); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func removeViewer(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	userIDStr := c.Params("userId")
	targetUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	character, err := RequireOwner(c, user, charID)
	if err != nil {
		return err
	}

	targetUser := &entities.User{}
	if err := db.DB.First(targetUser, targetUserID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	if err := db.DB.Model(character).Association("Viewers").Delete(targetUser); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
