package users

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
)

func beginFido2Link(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	displayName := c.Query("display_name")
	if displayName == "" {
		displayName = user.Username + " Passkey"
	}

	state, options, err := auth.BeginFido2Register(user, displayName)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"state": state, "options": options})
}

func endFido2Link(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	state := c.Query("state")
	if state == "" {
		return c.Status(400).JSON(fiber.Map{"error": "State is missing"})
	}

	pcc, err := protocol.ParseCredentialCreationResponseBytes(c.Body())
	if err != nil {
		return err
	}

	if err = auth.FinishFido2Register(state, user, pcc); err != nil {
		return err
	}

	return c.SendStatus(204)
}

func unlinkFido2(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	name := c.Params("name")
	if name == "" {
		return c.JSON(fiber.Map{"error": "Name is missing"})
	}

	if err := auth.DeleteFido2Login(user, user.ID, name); err != nil {
		return err
	}

	return c.SendStatus(204)
}

func getFido2Links(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})

	}

	logins := []entities.Fido2Login{}
	if res := db.DB.Where("user_id = ?", user.ID).Find(&logins); res.Error != nil {
		return res.Error
	}

	formatted := []string{}
	for _, login := range logins {
		formatted = append(formatted, login.DisplayName)
	}

	return c.JSON(formatted)
}
