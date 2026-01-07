package users

import "github.com/gofiber/fiber/v2"

func Register(api fiber.Router) {
	users := api.Group("/users")

	users.Get("", getUsers)
	users.Get("/autocomplete", autocompleteUsers)
	users.Get("/:id", getUser)
	users.Get("/@me/sessions", getUserSessions)
	users.Post("/@me/totp", beginTotp)
	users.Post("/@me/totp/verify", verifyTotp)
	users.Get("/@me/totp/verify/url", requestUrlForVerification)
	users.Post("/@me/totp/disable", disableTotp)
	users.Delete("/:id/totp", removeTotp)
	users.Post("/@me/link/fido2/begin", beginFido2Link)
	users.Post("/@me/link/fido2/end", endFido2Link)
	users.Get("/@me/link/fido2/all", getFido2Links)
	users.Delete("/@me/link/fido2/:name", unlinkFido2)
	users.Patch("/:id/username", setUserUsername)
	users.Patch("/:id/password", setUserPassword)
}
