package auth

import "errors"

var (
	ErrAbsoluteExpReached   = errors.New("absolute expiration reached")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserBlocked          = errors.New("user blocked")
	ErrUserOtpMissing       = errors.New("user otp missing")
	ErrUserOtpWrong         = errors.New("user otp wrong")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUsedRefreshToken     = errors.New("used refresh token")
)
