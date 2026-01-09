package characters

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/log"
	"vicar-backend/sync"
	"vicar-backend/util"

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
		log.Error(log.Server, "❌", "Failed to parse update character dto: %v", err)
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
		changes["clanID"] = dto.ClanID
	}
	if dto.PredatorTypeID != nil {
		character.PredatorTypeID = dto.PredatorTypeID
		changes["predatorTypeID"] = dto.PredatorTypeID
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

	if dto.Books != nil && len(dto.Books) > 0 {
		bookIDs := util.MapArray(dto.Books, func(b struct {
			ID uuid.UUID `json:"id"`
		}) uuid.UUID {
			return b.ID
		})

		var books []entities.V5Book
		if err := db.DB.Where("id IN ?", bookIDs).Find(&books).Error; err != nil {
			return err
		}
		if len(books) != len(bookIDs) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "One or more books not found"})
		}
		if err := db.DB.Model(character).Association("Books").Replace(books); err != nil {
			return err
		}
		changes["books"] = bookIDs
	}

	if dto.Inventory != nil {
		character.Inventory = *dto.Inventory
		changes["inventory"] = *dto.Inventory
	}

	if dto.InternalData != nil {
		character.InternalData = *dto.InternalData
		changes["internalData"] = *dto.InternalData
	}

	tx := db.DB.Begin()
	defer tx.Rollback()

	if err := tx.Save(character).Error; err != nil {
		return err
	}

	if dto.Attributes != nil && len(dto.Attributes) > 0 {
		if err := tx.Where("character_id = ?", character.ID).Delete(&entities.V5CharacterAttribute{}).Error; err != nil {
			return err
		}

		for _, attr := range dto.Attributes {
			newAttr := &entities.V5CharacterAttribute{
				CharacterID: character.ID,
				Category:    attr.Category,
				Key:         attr.Key,
				Value:       attr.Value,
			}

			if err := tx.Create(newAttr).Error; err != nil {
				return err
			}
		}

		changes["attributes"] = dto.Attributes
	}

	if dto.Skills != nil && len(dto.Skills) > 0 {
		if err := tx.Where("character_id = ?", character.ID).Delete(&entities.V5CharacterSkill{}).Error; err != nil {
			return err
		}

		for _, skill := range dto.Skills {
			newSkill := &entities.V5CharacterSkill{
				CharacterID:    character.ID,
				Category:       skill.Category,
				Key:            skill.Key,
				Value:          skill.Value,
				Specialization: skill.Specialization,
			}

			if err := tx.Create(newSkill).Error; err != nil {
				return err
			}
		}
	}

	if dto.DisciplineSelections != nil && len(dto.DisciplineSelections) > 0 {
		if err := tx.Where("character_id = ?", character.ID).Delete(&entities.V5CharacterDisciplineSelection{}).Error; err != nil {
			return err
		}

		for _, selection := range dto.DisciplineSelections {
			newSelection := &entities.V5CharacterDisciplineSelection{
				CharacterID:  character.ID,
				DisciplineID: selection.DisciplineID,
				Points:       selection.Points,
				CurrentLevel: selection.CurrentLevel,
			}

			if err := tx.Create(newSelection).Error; err != nil {
				return err
			}

			for _, ability := range selection.Abilities {
				newAbility := &entities.V5CharacterDisciplineAbility{
					SelectionID: newSelection.ID,
					AbilityID:   ability.AbilityID,
					Level:       ability.Level,
					UsedLevel:   ability.UsedLevel,
				}

				if err := tx.Create(newAbility).Error; err != nil {
					return err
				}
			}
		}

		changes["disciplineSelections"] = dto.DisciplineSelections
	}

	if dto.TraitPackUsages != nil && len(dto.TraitPackUsages) > 0 {
		if err := tx.Where("character_id = ?", character.ID).Delete(&entities.V5CharacterTraitPackUsage{}).Error; err != nil {
			return err
		}

		for _, usage := range dto.TraitPackUsages {
			newUsage := &entities.V5CharacterTraitPackUsage{
				CharacterID: character.ID,
				PackID:      usage.PackID,
				Kind:        usage.Kind,
			}

			if err := tx.Create(newUsage).Error; err != nil {
				return err
			}

			for _, trait := range usage.Traits {
				newTrait := &entities.V5CharacterTrait{
					UsageID:     newUsage.ID,
					TraitID:     trait.TraitID,
					CustomLevel: &trait.CustomLevel,
					IsLocked:    trait.IsLocked,
					IsManual:    trait.IsManual,
					Suffix:      trait.Suffix,
				}

				if err := tx.Create(newTrait).Error; err != nil {
					return err
				}
			}

			for _, trait := range usage.FlawTraits {
				newTrait := &entities.V5CharacterFlawTrait{
					UsageID:     newUsage.ID,
					TraitID:     trait.TraitID,
					CustomLevel: &trait.CustomLevel,
					IsLocked:    trait.IsLocked,
					IsManual:    trait.IsManual,
					Suffix:      trait.Suffix,
				}

				if err := tx.Create(newTrait).Error; err != nil {
					return err
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
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
