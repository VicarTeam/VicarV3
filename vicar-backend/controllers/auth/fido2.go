package auth

import (
	"vicar-backend/auth"
	"vicar-backend/log"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
)

func beginLoginFido2(c *fiber.Ctx) error {
	state, options, err := auth.BeginFido2Login()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"state": state, "options": options})
}

func endLoginFido2(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "State is missing"})
	}

	pcc, err := protocol.ParseCredentialRequestResponseBytes(c.Body())
	if err != nil {
		log.Error(log.Server, "❌", "Failed to parse credential request response body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse credential request response body"})
	}

	user, err := auth.FinishFido2Login(state, pcc)
	if err != nil {
		return err
	}

	token, err := auth.CreateTokenPairForUser(user, auth.GetDeviceName(c))
	if err != nil {
		return err
	}

	setRefrehTokenCookie(c, token.RefreshToken)

	return c.JSON(token)
}
