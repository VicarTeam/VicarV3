package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Register(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/register", createAccount)
	auth.Post("/login", login)
	auth.Post("/login/totp", loginTotp)
	auth.Post("/login/fido2/begin", beginLoginFido2)
	auth.Post("/login/fido2/end", endLoginFido2)
	auth.Post("/refresh", refreshToken)
	auth.Post("/logout", logout)
	auth.Post("/logout/all", logoutAll)
}
