package auth

import (
	"errors"
	"os"
	"time"
	"vicar-backend/auth"
	"vicar-backend/cache"
	"vicar-backend/log"

	"github.com/gofiber/fiber/v2"
)

var cookieDomain = os.Getenv("COOKIE_DOMAIN")
var cookieSecure = os.Getenv("COOKIE_SECURE") == "true"

func login(c *fiber.Ctx) error {
	dto := &loginDto{}
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Given data is invalid"})
	}

	token, _, err := auth.Login(dto.Username, dto.Password, auth.GetDeviceName(c), nil)
	if errors.Is(err, auth.ErrUserNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if errors.Is(err, auth.ErrInvalidCredentials) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	if errors.Is(err, auth.ErrUserBlocked) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is blocked"})
	}

	if errors.Is(err, auth.ErrUserOtpMissing) {
		otpState := cache.BeginState("login_otp", map[string]string{
			"username": dto.Username,
			"password": dto.Password,
		}, time.Minute*3)

		return c.Status(fiber.StatusPreconditionRequired).JSON(fiber.Map{
			"error": "OTP is wrong",
			"state": otpState,
		})
	}

	if err != nil {
		return err
	}

	setRefrehTokenCookie(c, token.RefreshToken)

	return c.JSON(token)
}

func loginTotp(c *fiber.Ctx) error {
	dto := &loginTotpDto{}
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Given data is invalid"})
	}

	var payload map[string]string
	if ok := cache.EndState("login_otp", dto.State, &payload); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid state"})
	}

	token, _, err := auth.Login(payload["username"], payload["password"], auth.GetDeviceName(c), &dto.Code)
	if errors.Is(err, auth.ErrUserNotFound) {
		c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if errors.Is(err, auth.ErrInvalidCredentials) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password"})
	}

	if errors.Is(err, auth.ErrUserBlocked) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is blocked"})
	}

	if errors.Is(err, auth.ErrUserOtpWrong) {
		otpState := cache.BeginState("login_otp", map[string]string{
			"username": payload["username"],
			"password": payload["password"],
		}, time.Minute*3)

		return c.Status(fiber.StatusPreconditionRequired).JSON(fiber.Map{
			"error": "OTP is wrong",
			"state": otpState,
		})
	}

	if err != nil {
		return err
	}

	setRefrehTokenCookie(c, token.RefreshToken)

	return c.JSON(token)
}

func refreshToken(c *fiber.Ctx) error {
	refreshToken, err := extractRefreshToken(c)
	if err != nil {
		log.Warning(log.Auth, "🔒", "Refresh token is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Refresh token is missing"})
	}

	token, err := auth.Refresh(auth.GetDeviceName(c), refreshToken)
	if errors.Is(err, auth.ErrUsedRefreshToken) {
		log.Warning(log.Auth, "🔒", "Refresh token is used")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token invalid"})
	}

	if err != nil {
		log.Error(log.Auth, "🔒", "Refresh token error: %s", err.Error())
		return err
	}

	setRefrehTokenCookie(c, token.RefreshToken)

	return c.JSON(token)
}

func logout(c *fiber.Ctx) error {
	refreshToken, err := extractRefreshToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Refresh token is missing"})
	}

	err = auth.Logout(auth.GetDeviceName(c), refreshToken, false)
	if err != nil {
		return err
	}

	unsetRefrehTokenCookie(c)

	return c.JSON(fiber.Map{"message": "User logged out successfully"})
}

func logoutAll(c *fiber.Ctx) error {
	refreshToken, err := extractRefreshToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Refresh token is missing"})
	}

	err = auth.Logout(auth.GetDeviceName(c), refreshToken, true)
	if err != nil {
		return err
	}

	unsetRefrehTokenCookie(c)

	return c.JSON(fiber.Map{"message": "User logged out successfully"})
}

func extractRefreshToken(c *fiber.Ctx) (string, error) {
	refreshToken := c.Cookies("rf_token")
	if refreshToken == "" {
		return "", errors.New("refresh token is missing")
	}

	return refreshToken, nil
}

func setRefrehTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "rf_token",
		Value:    token,
		Expires:  time.Now().Add(60 * 60 * 24 * 30 * time.Second),
		Path:     "/",
		Domain:   cookieDomain,
		Secure:   cookieSecure,
		SameSite: "Lax",
		HTTPOnly: true,
	})
}

func unsetRefrehTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "rf_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Second),
		Path:     "/",
		Domain:   cookieDomain,
		Secure:   cookieSecure,
		SameSite: "Lax",
		HTTPOnly: true,
	})
}
