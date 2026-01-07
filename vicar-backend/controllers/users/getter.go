package users

import (
	"strings"
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/serialize"

	"github.com/gofiber/fiber/v2"
)

func getUser(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	userIdStr := c.Params("id")
	if userIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is missing"})
	}

	targetUser := &entities.User{}
	if userIdStr == "@me" {
		targetUser = user
	} else if err := db.DB.Where("id = ?", userIdStr).First(targetUser).Error; err != nil {
		return err
	}

	forMyself := true
	if userIdStr != "@me" {
		if c.Query("view") == "admin" {
			if !user.IsTeam {
				return fiber.ErrUnauthorized
			}

			forMyself = false
		}
	}

	return serialize.JSON(c, &serialize.UserSerializer{ForMyself: forMyself}, *targetUser)
}

func getUserSessions(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	refreshTokenSessions := []string{}
	if err := db.DB.Raw("SELECT DISTINCT device_name FROM refresh_tokens WHERE user_id = ? AND device_name <> ?", user.ID, "").Scan(&refreshTokenSessions).Error; err != nil {
		return err
	}

	return c.JSON(refreshTokenSessions)
}

func getUsers(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	if !user.IsTeam {
		return fiber.ErrForbidden
	}

	qb := db.NewSearchBuilder(&entities.User{}).OrderByAsc("username").SearchFilter("username", false, true)

	items := []*entities.User{}
	total, err := qb.Extract(c).Execute(&items)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"total": total,
		"items": serialize.DoPtrArray(&serialize.UserSerializer{ForMyself: false}, items),
	})
}

func autocompleteUsers(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	name := strings.TrimSpace(c.Query("search"))
	if len(name) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is too short"})
	}

	users := []*entities.User{}
	if err := db.DB.Where("username ILIKE ?", "%"+name+"%").Limit(10).Order("username ASC").Find(&users).Error; err != nil {
		return err
	}

	return c.JSON(serialize.DoPtrArray(&serialize.UserSerializer{ForMyself: false}, users))
}
