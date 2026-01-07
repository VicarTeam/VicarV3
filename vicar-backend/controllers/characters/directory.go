package characters

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/serialize"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func getDirectories(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	directories := []entities.CharacterDirectory{}
	if err := db.DB.
		Preload("V5Characters", "user_id = ?", user.ID).
		Find(&directories).Error; err != nil {
		return err
	}

	return serialize.JSONArray(c, &serialize.CharacterDirectorySerializer{IncludeCharacterCount: true}, directories)
}

func createDirectory(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	var dto CreateDirectoryDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	directory := &entities.CharacterDirectory{
		Name: dto.Name,
		Open: false,
	}

	if err := db.DB.Create(directory).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": directory.ID})
}

func updateDirectory(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	dirIDStr := c.Params("id")
	dirID, err := uuid.Parse(dirIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid directory ID"})
	}

	var dto UpdateDirectoryDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	directory := &entities.CharacterDirectory{}
	if err := db.DB.First(directory, dirID).Error; err != nil {
		return err
	}

	directory.Name = dto.Name
	if err := db.DB.Save(directory).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func deleteDirectory(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	dirIDStr := c.Params("id")
	dirID, err := uuid.Parse(dirIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid directory ID"})
	}

	directory := &entities.CharacterDirectory{}
	if err := db.DB.First(directory, dirID).Error; err != nil {
		return err
	}

	if err := db.DB.Delete(directory).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
