package serialize

import (
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
)

type UserSerializer struct {
	ForMyself      bool
	IncludeBlocked bool
}

func (s *UserSerializer) Serialize(user entities.User, args ...any) any {
	m := fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
		"isTeam":   user.IsTeam,
	}

	if s.IncludeBlocked {
		m["isBlocked"] = user.IsBlocked
	}

	if s.ForMyself {
		m["otpActive"] = user.OtpActive
		m["otpVerified"] = user.OtpVerified
	}

	return m
}
