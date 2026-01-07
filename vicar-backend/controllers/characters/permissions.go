package characters

import (
	"vicar-backend/db"
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PermissionLevel int

const (
	PermissionNone PermissionLevel = iota
	PermissionView
	PermissionOwner
)

func CheckCharacterPermission(c *fiber.Ctx, user *entities.User, charID uuid.UUID) (*entities.V5Character, PermissionLevel, error) {
	character := &entities.V5Character{}

	if err := db.DB.Where("id = ?", charID).First(character).Error; err != nil {
		return nil, PermissionNone, err
	}

	if character.UserID == user.ID {
		return character, PermissionOwner, nil
	}

	var count int64
	if err := db.DB.Table("v5_character_viewers").
		Where("v5_character_id = ? AND user_id = ?", charID, user.ID).
		Count(&count).Error; err != nil {
		return nil, PermissionNone, err
	}

	if count > 0 {
		return character, PermissionView, nil
	}

	return character, PermissionNone, nil
}

func RequireOwner(c *fiber.Ctx, user *entities.User, charID uuid.UUID) (*entities.V5Character, error) {
	character, perm, err := CheckCharacterPermission(c, user, charID)
	if err != nil {
		return nil, err
	}

	if perm != PermissionOwner {
		return nil, fiber.ErrForbidden
	}

	return character, nil
}

func RequireViewOrOwner(c *fiber.Ctx, user *entities.User, charID uuid.UUID) (*entities.V5Character, PermissionLevel, error) {
	character, perm, err := CheckCharacterPermission(c, user, charID)
	if err != nil {
		return nil, PermissionNone, err
	}

	if perm == PermissionNone {
		return nil, PermissionNone, fiber.ErrForbidden
	}

	return character, perm, nil
}
