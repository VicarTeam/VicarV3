package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generatePair(userID uuid.UUID, deviceName string, absoluteExpiration *int64) (TokenPair, error) {
	rexp := int64(0)
	if absoluteExpiration != nil {
		rexp = *absoluteExpiration
	} else {
		rexp = time.Now().Add(time.Hour * 24 * 30).Unix()
	}

	if rexp < time.Now().Unix() {
		return TokenPair{}, ErrAbsoluteExpReached
	}

	userIdStr := userID.String()

	exp := time.Now().Add(time.Minute * 30).Unix()
	identityToken, err := generateToken(getIdentityTokenSecret(), userIdStr, exp)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := generateToken(getRefreshTokenSecret(), userIdStr, rexp)
	if err != nil {
		return TokenPair{}, err
	}

	if err := registerIdentityTokenInCache(userIdStr, deviceName, identityToken, exp); err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		RefreshToken: refreshToken,
		IdenityToken: identityToken,
		ExpiresIn:    exp * 1000,
	}, nil
}

func generateToken(secret string, sub string, exp int64) (string, error) {
	currTime := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "Nauri Auth",
		"aud": "nauri.io",
		"sub": sub,
		"exp": exp,
		"iat": currTime,
		"nbf": currTime,
		"jti": uuid.New().String(),
	})

	return token.SignedString([]byte(secret))
}

func getIdentityTokenSecret() string {
	secret, ok := os.LookupEnv("JWT_IDENTITY_SECRET")
	if !ok {
		log.Println("WARNING!!! JWT_IDENTITY_SECRET not set, using default value")
	}

	return secret
}

func getRefreshTokenSecret() string {
	secret, ok := os.LookupEnv("JWT_REFRESH_SECRET")
	if !ok {
		log.Println("WARNING!!! JWT_REFRESH_SECRET not set, using default value")
		return "secret"
	}

	return secret
}
