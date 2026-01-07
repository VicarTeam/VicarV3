package controllers

import (
	"vicar-backend/controllers/auth"
	"vicar-backend/controllers/characters"
	"vicar-backend/controllers/users"
	"vicar-backend/controllers/v5data"

	"github.com/gofiber/fiber/v2"
)

func Register(api fiber.Router) {
	auth.Register(api)
	users.Register(api)
	characters.Register(api)
	v5data.Register(api)
}
