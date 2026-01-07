package auth

import (
	"errors"
	"strings"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/log"
	"vicar-backend/util"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckUserLoginAllowance(user *entities.User) error {
	if user.Password == "" {
		return errors.New("login with this account not allowed")
	}

	return util.Ternary(user.IsBlocked, ErrUserBlocked, nil)
}

// Login tries to login a user with the given username and password.
// Returns a token pair if successful, otherwise an error.
func Login(username, password, deviceName string, otp *string) (*TokenPair, *entities.User, error) {
	user := &entities.User{}

	if res := db.DB.Where("username = ?", username).First(user); res.Error != nil {
		return nil, nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := CheckUserLoginAllowance(user); err != nil {
		return nil, user, err
	}

	if user.HasTwoFactor() {
		if otp == nil {
			return nil, user, ErrUserOtpMissing
		}

		if !totp.Validate(*otp, *user.OtpSecret) {
			if err := TryBackupCodes(user, otp); err != nil {
				return nil, user, err
			}
		}
	}

	if res := db.DB.Model(&entities.RefreshToken{}).Where("user_id = ? AND device_name = ?", user.ID, deviceName).Update("is_revoked", true); res.Error != nil {
		return nil, user, res.Error
	}

	tokenPair, err := generatePair(user.ID, deviceName, nil)
	if err != nil {
		return nil, user, err
	}

	refreshToken := &entities.RefreshToken{
		UserID:     user.ID,
		Token:      tokenPair.RefreshToken,
		DeviceName: deviceName,
	}

	if res := db.DB.Create(refreshToken); res.Error != nil {
		return nil, user, res.Error
	}

	return &tokenPair, user, nil
}

func CreateTokenPairForUser(user *entities.User, deviceName string) (*TokenPair, error) {
	tokenPair, err := generatePair(user.ID, deviceName, nil)
	if err != nil {
		return nil, err
	}

	refreshToken := &entities.RefreshToken{
		UserID:     user.ID,
		Token:      tokenPair.RefreshToken,
		DeviceName: deviceName,
	}

	if res := db.DB.Create(refreshToken); res.Error != nil {
		return nil, res.Error
	}

	return &tokenPair, nil
}

// TryBackupCodes tries to use a backup code to do a two factor authentication.
func TryBackupCodes(user *entities.User, otp *string) error {
	if len(user.OtpBackupCodes) == 0 {
		return ErrUserOtpWrong
	}

	found := false
	for i, code := range user.OtpBackupCodes {
		if code == *otp {
			found = true
			user.OtpBackupCodes = append(user.OtpBackupCodes[:i], user.OtpBackupCodes[i+1:]...)
			break
		}
	}

	if !found {
		return ErrUserOtpWrong
	}

	if res := db.DB.Save(user); res.Error != nil {
		return res.Error
	}

	return nil
}

// Logout logs out the user with the given refresh token.
// If all is true, all refresh tokens for the user will be deleted.
func Logout(deviceName, refreshToken string, all bool) error {
	decoded, err := decodeToken(getRefreshTokenSecret(), refreshToken)
	if err != nil {
		return err
	}

	if all {
		if err := unregisterAllIdentityTokensInCache(decoded.Subject.String()); err != nil {
			log.Error(log.Auth, "❌", "Failed to unregister all identity tokens in cache: %v", err)
		}

		if res := db.DB.Where("user_id = ?", decoded.Subject).Delete(&entities.RefreshToken{}); res.Error != nil {
			return res.Error
		}
	} else {
		if err := unregisterAllIdentityTokensForDeviceInCache(decoded.Subject.String(), deviceName); err != nil {
			log.Error(log.Auth, "❌", "Failed to unregister all identity tokens for device in cache: %v", err)
		}

		if res := db.DB.Where("token = ? OR (user_id = ? AND (device_name = ? OR ?) AND is_revoked = ?)", refreshToken, decoded.Subject, deviceName, deviceName == "", true).Delete(&entities.RefreshToken{}); res.Error != nil {
			return res.Error
		}
	}

	return nil
}

// Refresh refreshes the access token with the given refresh token. If the refresh token is invalid, an error is returned.
// The refresh token will be rotated, so that the old one is no longer valid.
func Refresh(deviceName, refreshToken string) (*TokenPair, error) {
	decoded, err := decodeToken(getRefreshTokenSecret(), refreshToken)
	if err != nil {
		return nil, ErrUsedRefreshToken
	}

	existing := &entities.RefreshToken{}
	if res := db.DB.Where("token = ? AND user_id = ?", refreshToken, decoded.Subject).First(existing); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, res.Error
		}
	} else if existing.IsRevoked {
		_ = Logout("", refreshToken, true)
		return nil, ErrUsedRefreshToken
	}

	tokenPair, err := generatePair(decoded.Subject, deviceName, nil)
	if err != nil {
		return nil, err
	}

	if res := db.DB.Model(&entities.RefreshToken{}).Where("user_id = ? AND token = ?", decoded.Subject, refreshToken).Update("is_revoked", true); res.Error != nil {
		return nil, res.Error
	}

	refreshTokenEntity := &entities.RefreshToken{
		UserID:     decoded.Subject,
		Token:      tokenPair.RefreshToken,
		DeviceName: deviceName,
	}

	if res := db.DB.Create(refreshTokenEntity); res.Error != nil {
		return nil, res.Error
	}

	return &tokenPair, nil
}

// Extract extracts the user from the given context by extracting the token from the Authorization header.
func Extract(ctx *fiber.Ctx) *entities.User {
	idenityToken := ctx.Get("Authorization")
	if idenityToken == "" {
		idenityToken = ctx.Query("identity")
	} else {
		_, cut, ok := strings.Cut(idenityToken, "Bearer")
		if !ok {
			log.Error(log.Auth, "❌", "Failed to extract identity token from Authorization header")
			return nil
		}

		idenityToken = cut
	}

	idenityToken = strings.TrimSpace(idenityToken)
	decoded, err := decodeToken(getIdentityTokenSecret(), idenityToken)
	if err != nil {
		log.Error(log.Auth, "❌", "Failed to decode identity token: %v", err)
		return nil
	}

	ignoreCache := ctx.Locals("ignoreCache")
	if ignoreCache == nil {
		if valid, err := isIdentityTokenValidInCache(decoded.Subject.String(), GetDeviceName(ctx), idenityToken); err != nil || !valid {
			if err != nil {
				log.Error(log.Auth, "❌", "Failed to check if identity token is valid in cache: %v", err)
			}
			log.Error(log.Auth, "❌", "Identity token is not valid in cache")
			return nil
		}
	}

	log.Debug(log.Auth, "✅", "Decoded: %+v", decoded)

	user := &entities.User{}
	if res := db.DB.Where("id = ?", decoded.Subject).First(user); res.Error != nil {
		log.Error(log.Auth, "❌", "Failed to find user with id %s: %v", decoded.Subject, res.Error)
		return nil
	}

	if CheckUserLoginAllowance(user) != nil {
		log.Error(log.Auth, "❌", "User %s is blocked", user.Username)
		return nil
	}

	return user
}

func ExtractForWebSocket(c *websocket.Conn) *entities.User {
	dn := c.Headers("X-Device-Name")
	if dn == "" {
		dn = c.Headers("User-Agent")
	}

	dn = transformDeviceName(dn)

	idenityToken := c.Query("identity")
	if idenityToken == "" {
		log.Error(log.Auth, "❌", "No identity token provided in WebSocket connection")
		return nil
	}

	decoded, err := decodeToken(getIdentityTokenSecret(), idenityToken)
	if err != nil {
		log.Error(log.Auth, "❌", "Failed to decode identity token: %v", err)
		return nil
	}

	if valid, err := isIdentityTokenValidInCache(decoded.Subject.String(), dn, idenityToken); err != nil || !valid {
		if err != nil {
			log.Error(log.Auth, "❌", "Failed to check if identity token is valid in cache: %v", err)
		}
		log.Error(log.Auth, "❌", "Identity token is not valid in cache")
		return nil
	}

	log.Debug(log.Auth, "✅", "Decoded: %+v", decoded)

	user := &entities.User{}
	if res := db.DB.Where("id = ?", decoded.Subject).First(user); res.Error != nil {
		log.Error(log.Auth, "❌", "Failed to find user with id %s: %v", decoded.Subject, res.Error)
		return nil
	}

	if CheckUserLoginAllowance(user) != nil {
		log.Error(log.Auth, "❌", "User %s is blocked", user.Username)
		return nil
	}

	return user
}

func GetDeviceName(c *fiber.Ctx) string {
	dn := c.Get("X-Device-Name")
	if dn == "" {
		dn = c.Get("User-Agent")
	}

	return transformDeviceName(dn)
}

func transformDeviceName(dn string) string {
	dn = strings.ReplaceAll(dn, "ä", "ae")
	dn = strings.ReplaceAll(dn, "ö", "oe")
	dn = strings.ReplaceAll(dn, "ü", "ue")
	dn = strings.ReplaceAll(dn, "Ä", "Ae")
	dn = strings.ReplaceAll(dn, "Ö", "Oe")
	dn = strings.ReplaceAll(dn, "Ü", "Ue")
	dn = strings.ReplaceAll(dn, "ß", "ss")

	dn = strings.Map(func(r rune) rune {
		if r > 127 {
			return -1
		}
		return r
	}, dn)

	if dn == "" {
		dn = "Unknown Device"
	}

	return dn
}
