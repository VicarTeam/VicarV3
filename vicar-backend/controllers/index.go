package controllers

import (
	"vicar-backend/controllers/auth"
	"vicar-backend/controllers/users"

	"github.com/gofiber/fiber/v2"
)

func Register(api fiber.Router) {
	auth.Register(api)
	users.Register(api)
}
