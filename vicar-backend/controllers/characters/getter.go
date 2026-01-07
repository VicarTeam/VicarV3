package characters

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/serialize"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func getCharacters(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	qb := db.NewSearchBuilder(&entities.V5Character{}).
		OrderByDesc("updated_at").
		SearchFilter("name", false, true).
		PreFilter("directoryId", "directory_id = ?").
		PreFilter("clanId", "clan_id = ?").
		Preload("Clan").
		Joins("LEFT JOIN v5_character_viewers ON v5_characters.id = v5_character_viewers.v5_character_id").
		Where("v5_characters.user_id = ? OR v5_character_viewers.user_id = ?", user.ID, user.ID)

	qb.Extract(c)

	sort := c.Query("sort", "-updated_at")
	if sort == "name" {
		qb.OrderByAsc("name")
	} else if sort == "-name" {
		qb.OrderByDesc("name")
	}

	items := []*entities.V5Character{}
	total, err := qb.Execute(&items)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"total": total,
		"items": serialize.DoPtrArray(&serialize.V5CharacterSummarySerializer{}, items),
	})
}

func getCharacter(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	character, perm, err := RequireViewOrOwner(c, user, charID)
	if err != nil {
		return err
	}

	if err := db.DB.
		Preload("Directory").
		Preload("Clan").
		Preload("PredatorType").
		Preload("Books").
		Preload("Categories").
		Preload("Attributes").
		Preload("Skills").
		Preload("RequiredPointSpreads.Pack").
		Preload("DisciplineSelections.Discipline").
		Preload("DisciplineSelections.Abilities.Ability").
		Preload("TraitPackUsages.Pack").
		Preload("TraitPackUsages.Traits.Trait").
		Preload("TraitPackUsages.FlawTraits.Trait").
		Preload("BloodRituals").
		Preload("OblivionCeremonies").
		Preload("LevelHistory").
		Preload("Viewers").
		First(character, charID).Error; err != nil {
		return err
	}

	isOwner := perm == PermissionOwner
	return serialize.JSON(c, &serialize.V5CharacterSerializer{IsOwner: isOwner}, *character)
}

func getCharacterViewers(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	character, err := RequireOwner(c, user, charID)
	if err != nil {
		return err
	}

	if err := db.DB.Preload("Viewers").First(character, charID).Error; err != nil {
		return err
	}

	viewers := make([]*entities.User, len(character.Viewers))
	for i := range character.Viewers {
		viewers[i] = &character.Viewers[i]
	}

	return serialize.JSONPtrArray(c, &serialize.UserSerializer{ForMyself: false}, viewers)
}
