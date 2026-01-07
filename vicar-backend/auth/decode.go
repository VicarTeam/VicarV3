package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type decodedToken struct {
	Subject   uuid.UUID
	Audience  string
	ExpiresAt int64
}

func decodeToken(secret string, tokenString string) (*decodedToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		subClaim, ok := claims["sub"]
		if !ok {
			return nil, ErrInvalidToken
		}

		expClaim, ok := claims["exp"]
		if !ok {
			return nil, ErrInvalidToken
		}

		audClaim, ok := claims["aud"]
		if !ok {
			return nil, ErrInvalidToken
		}

		subject, err := uuid.Parse(subClaim.(string))
		if err != nil {
			return nil, err
		}

		audience, ok := audClaim.(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		return &decodedToken{
			Subject:   subject,
			Audience:  audience,
			ExpiresAt: int64(expClaim.(float64)),
		}, nil
	} else {
		return nil, ErrInvalidToken
	}
}
