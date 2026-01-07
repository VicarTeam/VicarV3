package characters

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createCharacter(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	var dto CreateCharacterDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var bookCount int64
	if err := db.DB.Model(&entities.V5Book{}).Where("id IN ?", dto.BookIDs).Count(&bookCount).Error; err != nil {
		return err
	}
	if bookCount != int64(len(dto.BookIDs)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "One or more books not found"})
	}

	character := &entities.V5Character{
		BaseSheet: entities.BaseSheet{
			UserID: user.ID,
			Name:   dto.Name,
		},
		GenerationEra: entities.GenerationEra(dto.GenerationEra),
		Generation:    dto.Generation,
		Humanity:      7,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(character).Error; err != nil {
			return err
		}

		var books []entities.V5Book
		if err := tx.Where("id IN ?", dto.BookIDs).Find(&books).Error; err != nil {
			return err
		}
		if err := tx.Model(character).Association("Books").Append(books); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": character.ID})
}

func updateCharacter(c *fiber.Ctx) error {
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

	var dto UpdateCharacterDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	changes := fiber.Map{}

	if dto.Name != nil {
		character.Name = *dto.Name
		changes["name"] = *dto.Name
	}
	if dto.Avatar != nil {
		character.Avatar = *dto.Avatar
		changes["avatar"] = *dto.Avatar
	}
	if dto.Notes != nil {
		character.Notes = *dto.Notes
		changes["notes"] = *dto.Notes
	}
	if dto.Sex != nil {
		character.Sex = entities.Sex(*dto.Sex)
		changes["sex"] = *dto.Sex
	}
	if dto.Concept != nil {
		character.Concept = *dto.Concept
		changes["concept"] = *dto.Concept
	}
	if dto.Chronicle != nil {
		character.Chronicle = *dto.Chronicle
		changes["chronicle"] = *dto.Chronicle
	}
	if dto.Exp != nil {
		character.Exp = *dto.Exp
		changes["exp"] = *dto.Exp
	}
	if dto.UsedExp != nil {
		character.UsedExp = *dto.UsedExp
		changes["usedExp"] = *dto.UsedExp
	}

	if dto.ChroniclePrinciples != nil {
		character.ChroniclePrinciples = *dto.ChroniclePrinciples
		changes["chroniclePrinciples"] = *dto.ChroniclePrinciples
	}
	if dto.AnchorsAndBeliefs != nil {
		character.AnchorsAndBeliefs = *dto.AnchorsAndBeliefs
		changes["anchorsAndBeliefs"] = *dto.AnchorsAndBeliefs
	}
	if dto.Backstory != nil {
		character.Backstory = *dto.Backstory
		changes["backstory"] = *dto.Backstory
	}
	if dto.Sire != nil {
		character.Sire = *dto.Sire
		changes["sire"] = *dto.Sire
	}
	if dto.Ambition != nil {
		character.Ambition = *dto.Ambition
		changes["ambition"] = *dto.Ambition
	}
	if dto.Desire != nil {
		character.Desire = *dto.Desire
		changes["desire"] = *dto.Desire
	}

	if dto.ClanID != nil {
		character.ClanID = dto.ClanID
		changes["clanId"] = dto.ClanID
	}
	if dto.PredatorTypeID != nil {
		character.PredatorTypeID = dto.PredatorTypeID
		changes["predatorTypeId"] = dto.PredatorTypeID
	}

	if dto.GenerationEra != nil {
		character.GenerationEra = entities.GenerationEra(*dto.GenerationEra)
		changes["generationEra"] = *dto.GenerationEra
	}
	if dto.Generation != nil {
		character.Generation = *dto.Generation
		changes["generation"] = *dto.Generation
	}

	if dto.Hunger != nil {
		character.Hunger = *dto.Hunger
		changes["hunger"] = *dto.Hunger
	}
	if dto.Humanity != nil {
		character.Humanity = *dto.Humanity
		changes["humanity"] = *dto.Humanity
	}
	if dto.Stains != nil {
		character.Stains = *dto.Stains
		changes["stains"] = *dto.Stains
	}
	if dto.Resonance != nil {
		character.Resonance = entities.V5Resonance(*dto.Resonance)
		changes["resonance"] = *dto.Resonance
	}
	if dto.BloodPotency != nil {
		character.BloodPotency = *dto.BloodPotency
		changes["bloodPotency"] = *dto.BloodPotency
	}

	if dto.Health != nil {
		character.Health = *dto.Health
		changes["health"] = *dto.Health
	}
	if dto.HealthDamage != nil {
		character.HealthDamage = dto.HealthDamage
		changes["healthDamage"] = dto.HealthDamage
	}
	if dto.Willpower != nil {
		character.Willpower = *dto.Willpower
		changes["willpower"] = *dto.Willpower
	}
	if dto.WillpowerDamage != nil {
		character.WillpowerDamage = dto.WillpowerDamage
		changes["willpowerDamage"] = dto.WillpowerDamage
	}

	if dto.UseAdvancedDisciplines != nil {
		character.UseAdavancedDisciplines = *dto.UseAdvancedDisciplines
		changes["useAdvancedDisciplines"] = *dto.UseAdvancedDisciplines
	}
	if dto.AllowLearningOfAllPowers != nil {
		character.AllowLearningOfAllPowers = *dto.AllowLearningOfAllPowers
		changes["allowLearningOfAllPowers"] = *dto.AllowLearningOfAllPowers
	}
	if dto.FullCustomization != nil {
		character.FullCustomization = *dto.FullCustomization
		changes["fullCustomization"] = *dto.FullCustomization
	}
	if dto.SkillSpreadType != nil {
		character.SkillSpreadType = entities.SkillSpreadType(*dto.SkillSpreadType)
		changes["skillSpreadType"] = *dto.SkillSpreadType
	}

	if dto.BookIDs != nil && len(dto.BookIDs) > 0 {
		var books []entities.V5Book
		if err := db.DB.Where("id IN ?", dto.BookIDs).Find(&books).Error; err != nil {
			return err
		}
		if len(books) != len(dto.BookIDs) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "One or more books not found"})
		}
		if err := db.DB.Model(character).Association("Books").Replace(books); err != nil {
			return err
		}
		changes["bookIds"] = dto.BookIDs
	}

	if err := db.DB.Save(character).Error; err != nil {
		return err
	}

	if len(changes) > 0 {
		for key, value := range changes {
			sync.SyncCharacterChanges(charID.String(), fiber.Map{
				"path":  key,
				"value": value,
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func deleteCharacter(c *fiber.Ctx) error {
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

	if err := db.DB.Delete(character).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func updateAttribute(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	attrIDStr := c.Params("attrId")
	attrID, err := uuid.Parse(attrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid attribute ID"})
	}

	if _, err := RequireOwner(c, user, charID); err != nil {
		return err
	}

	var dto UpdateAttributeDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	attribute := &entities.V5CharacterAttribute{}
	if err := db.DB.Where("id = ? AND character_id = ?", attrID, charID).First(attribute).Error; err != nil {
		return err
	}

	attribute.Value = dto.Value
	if err := db.DB.Save(attribute).Error; err != nil {
		return err
	}

	sync.SyncCharacterChanges(charID.String(), fiber.Map{
		"path": "attributes." + string(attribute.Key),
		"value": fiber.Map{
			"id":       attribute.ID,
			"category": attribute.Category,
			"key":      attribute.Key,
			"value":    attribute.Value,
		},
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func updateSkill(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	skillIDStr := c.Params("skillId")
	skillID, err := uuid.Parse(skillIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid skill ID"})
	}

	if _, err := RequireOwner(c, user, charID); err != nil {
		return err
	}

	var dto UpdateSkillDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	skill := &entities.V5CharacterSkill{}
	if err := db.DB.Where("id = ? AND character_id = ?", skillID, charID).First(skill).Error; err != nil {
		return err
	}

	skill.Value = dto.Value
	skill.Specialization = dto.Specialization
	if err := db.DB.Save(skill).Error; err != nil {
		return err
	}

	sync.SyncCharacterChanges(charID.String(), fiber.Map{
		"path": "skills." + string(skill.Key),
		"value": fiber.Map{
			"id":             skill.ID,
			"category":       skill.Category,
			"key":            skill.Key,
			"value":          skill.Value,
			"specialization": skill.Specialization,
		},
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func addDisciplineSelection(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	if _, err := RequireOwner(c, user, charID); err != nil {
		return err
	}

	var dto CreateDisciplineSelectionDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	var discipline entities.V5Discipline
	if err := db.DB.First(&discipline, dto.DisciplineID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Discipline not found"})
	}

	selection := &entities.V5CharacterDisciplineSelection{
		CharacterID:  charID,
		DisciplineID: dto.DisciplineID,
		Points:       0,
		CurrentLevel: 0,
	}

	if err := db.DB.Create(selection).Error; err != nil {
		return err
	}

	db.DB.Preload("Discipline").First(selection, selection.ID)

	sync.SyncCharacterChanges(charID.String(), fiber.Map{
		"path": "disciplines.add",
		"value": fiber.Map{
			"id": selection.ID,
			"discipline": fiber.Map{
				"id":   discipline.ID,
				"name": discipline.Name,
			},
			"points":       0,
			"currentLevel": 0,
			"abilities":    []fiber.Map{},
		},
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": selection.ID})
}

func updateDisciplineSelection(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	discIDStr := c.Params("disciplineId")
	discID, err := uuid.Parse(discIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid discipline ID"})
	}

	if _, err := RequireOwner(c, user, charID); err != nil {
		return err
	}

	var dto UpdateDisciplineSelectionDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	selection := &entities.V5CharacterDisciplineSelection{}
	if err := db.DB.Where("id = ? AND character_id = ?", discID, charID).First(selection).Error; err != nil {
		return err
	}

	if dto.Points != nil {
		selection.Points = *dto.Points
	}
	if dto.CurrentLevel != nil {
		selection.CurrentLevel = *dto.CurrentLevel
	}

	if err := db.DB.Save(selection).Error; err != nil {
		return err
	}

	sync.SyncCharacterChanges(charID.String(), fiber.Map{
		"path": "disciplines." + discID.String(),
		"value": fiber.Map{
			"id":           selection.ID,
			"points":       selection.Points,
			"currentLevel": selection.CurrentLevel,
		},
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func deleteDisciplineSelection(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	charIDStr := c.Params("id")
	charID, err := uuid.Parse(charIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid character ID"})
	}

	discIDStr := c.Params("disciplineId")
	discID, err := uuid.Parse(discIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid discipline ID"})
	}

	if _, err := RequireOwner(c, user, charID); err != nil {
		return err
	}

	selection := &entities.V5CharacterDisciplineSelection{}
	if err := db.DB.Where("id = ? AND character_id = ?", discID, charID).First(selection).Error; err != nil {
		return err
	}

	if err := db.DB.Delete(selection).Error; err != nil {
		return err
	}

	sync.SyncCharacterChanges(charID.String(), fiber.Map{
		"path":  "disciplines.remove",
		"value": discID.String(),
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func moveCharacterToDirectory(c *fiber.Ctx) error {
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

	var dto MoveToDirectoryDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if dto.DirectoryID != nil {
		var directory entities.CharacterDirectory
		if err := db.DB.First(&directory, dto.DirectoryID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Directory not found"})
		}
	}

	character.DirectoryID = dto.DirectoryID
	if err := db.DB.Save(character).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
