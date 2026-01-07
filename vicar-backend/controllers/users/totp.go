package users

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/util"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/pquerna/otp/totp"
)

func beginTotp(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if user.OtpActive {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is already active"})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Nauri",
		AccountName: user.Username,
	})
	if err != nil {
		return err
	}

	secret := key.Secret()
	user.OtpSecret = &secret
	user.OtpActive = true

	if res := db.DB.Save(user); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{"url": key.String()})
}

func verifyTotp(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if !user.OtpActive {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is not active"})
	}

	if user.OtpVerified {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is already verified"})
	}

	dto := &struct {
		Code string `json:"code"`
	}{}
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Given data is invalid"})
	}

	if !totp.Validate(dto.Code, *user.OtpSecret) {
		return c.Status(403).JSON(fiber.Map{"error": "Invalid code"})
	}

	backupCodes, err := generateRecoveryCodes(10)
	if err != nil {
		return err
	}

	user.OtpVerified = true
	user.OtpBackupCodes = pq.StringArray(backupCodes)

	if res := db.DB.Save(user); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{
		"message": "TOTP verified successfully",
		"codes":   user.OtpBackupCodes,
	})
}

func disableTotp(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if !user.OtpActive {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is not active"})
	}

	if !user.OtpVerified {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is not verified"})
	}

	dto := &struct {
		Code string `json:"code"`
	}{}
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Given data is invalid"})
	}

	if !totp.Validate(dto.Code, *user.OtpSecret) {
		if err := auth.TryBackupCodes(user, &dto.Code); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Invalid code"})
		}
	}

	user.OtpActive = false
	user.OtpVerified = false
	user.OtpSecret = nil
	user.OtpBackupCodes = pq.StringArray{}

	if res := db.DB.Save(user); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{"message": "TOTP disabled successfully"})
}

func requestUrlForVerification(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if !user.OtpActive {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is not active"})
	}

	if user.OtpVerified {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is already verified"})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Nauri",
		AccountName: user.Username,
	})
	if err != nil {
		return err
	}

	secret := key.Secret()
	user.OtpSecret = &secret
	user.OtpActive = true

	if res := db.DB.Save(user); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{"url": key.String()})
}

func removeTotp(c *fiber.Ctx) error {
	if _, ok := auth.AuthorizeTeam(c); !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	user := &entities.User{}
	if res := db.DB.Where("id = ?", c.Params("id")).First(user); res.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if !user.OtpActive {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP is not active"})
	}

	user.OtpActive = false
	user.OtpVerified = false
	user.OtpSecret = nil
	user.OtpBackupCodes = nil

	if res := db.DB.Save(user); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{"message": "TOTP removed successfully"})
}

func generateRecoveryCodes(count int) ([]string, error) {
	codes := make([]string, count)

	for i := 0; i < count; i++ {
		code, err := util.GenerateRandomString(8)
		if err != nil {
			return nil, err
		}

		codes[i] = code
	}

	return codes, nil
}
