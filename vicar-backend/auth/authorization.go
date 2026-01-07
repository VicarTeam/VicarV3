package auth

import (
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
)

func AuthorizeTeam(c *fiber.Ctx) (*entities.User, bool) {
	user := Extract(c)
	if user == nil {
		return nil, false
	}

	return user, user.IsTeam
}
