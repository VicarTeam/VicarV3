package characters

import (
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/serialize"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OldCharacterData struct {
	ID                          string                   `json:"id"`
	Name                        string                   `json:"name"`
	Avatar                      string                   `json:"avatar"`
	Notes                       string                   `json:"notes"`
	Sex                         string                   `json:"sex"`
	Concept                     string                   `json:"concept"`
	Chronicle                   string                   `json:"chronicle"`
	Exp                         int                      `json:"exp"`
	UsedExp                     int                      `json:"usedExp"`
	ChroniclePrinciples         string                   `json:"chroniclePrinciples"`
	AnchorsAndBeliefs           string                   `json:"anchorsAndBeliefs"`
	Backstory                   string                   `json:"backstory"`
	Sire                        string                   `json:"sire"`
	Ambition                    string                   `json:"ambition"`
	Desire                      string                   `json:"desire"`
	GenerationEra               string                   `json:"generationEra"`
	Generation                  int                      `json:"generation"`
	Hunger                      int                      `json:"hunger"`
	Humanity                    int                      `json:"humanity"`
	Stains                      int                      `json:"stains"`
	Resonance                   string                   `json:"resonance"`
	BloodPotency                int                      `json:"bloodPotency"`
	Health                      int                      `json:"health"`
	HealthDamage                []string                 `json:"healthDamage"`
	Willpower                   int                      `json:"willpower"`
	WillpowerDamage             []string                 `json:"willpowerDamage"`
	UseAdvancedDisciplines      bool                     `json:"useAdavancedDisciplines"`
	AllowLearningOfAllPowers    bool                     `json:"allowLearningOfAllPowers"`
	FullCustomization           bool                     `json:"fullCustomization"`
	Version                     int                      `json:"version"`
	Clan                        *OldClanData             `json:"clan"`
	PredatorType                *OldPredatorTypeData     `json:"predatorType"`
	Books                       []OldBookData            `json:"books"`
	Categories                  []OldCategoryData        `json:"categories"`
	Disciplines                 []OldDisciplineData      `json:"disciplines"`
	Merits                      map[string]OldMeritData  `json:"merits"`
	Backgrounds                 map[string]OldMeritData  `json:"backgrounds"`
	BloodRituals                []OldBloodRitualData     `json:"bloodRituals"`
	OblivionCeremonies          []OldOblivionCeremonyData `json:"oblivionCeremonies"`
	RequiredPointSpreads        []OldPointSpreadData     `json:"requiredPointSpreads"`
	LevelHistory                []OldLevelHistoryData    `json:"levelHistory"`
	Skillspread                 map[string]interface{}   `json:"skillspread"`
}

type OldClanData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OldPredatorTypeData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OldBookData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OldCategoryData struct {
	Name       string                  `json:"name"`
	Attributes []OldAttributeData      `json:"attributes"`
	Skills     []OldSkillData          `json:"skills"`
}

type OldAttributeData struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type OldSkillData struct {
	Key            string   `json:"key"`
	Value          int      `json:"value"`
	Specialization []string `json:"specialization"`
}

type OldDisciplineData struct {
	DisciplineID int                    `json:"disciplineId"`
	Name         string                 `json:"name"`
	Points       int                    `json:"points"`
	CurrentLevel int                    `json:"currentLevel"`
	Abilities    []OldDisciplineAbility `json:"abilities"`
}

type OldDisciplineAbility struct {
	AbilityID  int `json:"abilityId"`
	Level      int `json:"level"`
	UsedLevel  int `json:"usedLevel"`
}

type OldMeritData struct {
	TraitPackID int           `json:"traitPackId"`
	Name        string        `json:"name"`
	Traits      []OldTraitData `json:"traits"`
	FlawTraits  []OldTraitData `json:"flawTraits"`
}

type OldTraitData struct {
	TraitID     int     `json:"traitId"`
	IsLocked    bool    `json:"isLocked"`
	IsManual    bool    `json:"isManual"`
	CustomLevel *int    `json:"customLevel"`
	Suffix      *string `json:"suffix"`
}

type OldBloodRitualData struct {
	ID int `json:"id"`
}

type OldOblivionCeremonyData struct {
	ID int `json:"id"`
}

type OldPointSpreadData struct {
	Type        string `json:"type"`
	IsFlaw      bool   `json:"isFlaw"`
	Points      int    `json:"points"`
	TraitPackID *int   `json:"traitPackId"`
}

type OldLevelHistoryData struct {
	Type      string `json:"type"`
	Date      string `json:"date"`
	Text      string `json:"text"`
	ExpUsed   int    `json:"expUsed"`
	ExpBefore int    `json:"expBefore"`
	ExpAfter  int    `json:"expAfter"`
}

func migrateCharacter(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	var oldData OldCharacterData
	if err := c.BodyParser(&oldData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	char := entities.V5Character{
		BaseSheet: entities.BaseSheet{
			UserID:    user.ID,
			Name:      oldData.Name,
			Avatar:    oldData.Avatar,
			Notes:     oldData.Notes,
			Sex:       entities.Sex(oldData.Sex),
			Concept:   oldData.Concept,
			Chronicle: oldData.Chronicle,
			Exp:       oldData.Exp,
			UsedExp:   oldData.UsedExp,
		},
		ChroniclePrinciples:      oldData.ChroniclePrinciples,
		AnchorsAndBeliefs:        oldData.AnchorsAndBeliefs,
		Backstory:                oldData.Backstory,
		Sire:                     oldData.Sire,
		Ambition:                 oldData.Ambition,
		Desire:                   oldData.Desire,
		GenerationEra:            entities.GenerationEra(oldData.GenerationEra),
		Generation:               oldData.Generation,
		Hunger:                   oldData.Hunger,
		Humanity:                 oldData.Humanity,
		Stains:                   oldData.Stains,
		Resonance:                entities.V5Resonance(oldData.Resonance),
		BloodPotency:             oldData.BloodPotency,
		Health:                   oldData.Health,
		HealthDamage:             pq.StringArray(oldData.HealthDamage),
		Willpower:                oldData.Willpower,
		WillpowerDamage:          pq.StringArray(oldData.WillpowerDamage),
		UseAdavancedDisciplines:  oldData.UseAdvancedDisciplines,
		AllowLearningOfAllPowers: oldData.AllowLearningOfAllPowers,
		FullCustomization:        oldData.FullCustomization,
		Version:                  oldData.Version,
	}

	if oldData.Clan != nil {
		var clan entities.V5Clan
		if err := db.DB.Where("old_vicar_id = ?", oldData.Clan.ID).First(&clan).Error; err == nil {
			char.ClanID = &clan.ID
		}
	}

	if oldData.PredatorType != nil {
		var pt entities.V5PredatorType
		if err := db.DB.Where("old_vicar_id = ?", oldData.PredatorType.ID).First(&pt).Error; err == nil {
			char.PredatorTypeID = &pt.ID
		}
	}

	if err := db.DB.Create(&char).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create character",
		})
	}

	for _, oldBook := range oldData.Books {
		var book entities.V5Book
		if err := db.DB.Where("old_vicar_id = ?", oldBook.ID).First(&book).Error; err == nil {
			db.DB.Exec("INSERT INTO v5_character_books (v5_character_id, v5_book_id) VALUES (?, ?)", char.ID, book.ID)
		}
	}

	for _, oldCat := range oldData.Categories {
		category := entities.V5CharacterCategory{
			CharacterID: char.ID,
			Name:        entities.CategoryKey(oldCat.Name),
		}
		db.DB.Create(&category)

		for _, oldAttr := range oldCat.Attributes {
			attr := entities.V5CharacterAttribute{
				CharacterID: char.ID,
				Category:    entities.CategoryKey(oldCat.Name),
				Key:         entities.AttributeKey(oldAttr.Key),
				Value:       oldAttr.Value,
			}
			db.DB.Create(&attr)
		}

		for _, oldSkill := range oldCat.Skills {
			skill := entities.V5CharacterSkill{
				CharacterID:    char.ID,
				Category:       entities.CategoryKey(oldCat.Name),
				Key:            entities.SkillKey(oldSkill.Key),
				Value:          oldSkill.Value,
				Specialization: pq.StringArray(oldSkill.Specialization),
			}
			db.DB.Create(&skill)
		}
	}

	for _, oldDisc := range oldData.Disciplines {
		var discipline entities.V5Discipline
		if err := db.DB.Where("old_vicar_id = ?", oldDisc.DisciplineID).First(&discipline).Error; err != nil {
			continue
		}

		selection := entities.V5CharacterDisciplineSelection{
			CharacterID:  char.ID,
			DisciplineID: discipline.ID,
			Points:       oldDisc.Points,
			CurrentLevel: oldDisc.CurrentLevel,
		}
		db.DB.Create(&selection)

		for _, oldAbility := range oldDisc.Abilities {
			var ability entities.V5DisciplineAbility
			if err := db.DB.Where("old_vicar_id = ?", oldAbility.AbilityID).First(&ability).Error; err != nil {
				continue
			}

			charAbility := entities.V5CharacterDisciplineAbility{
				SelectionID: selection.ID,
				AbilityID:   ability.ID,
				Level:       oldAbility.Level,
				UsedLevel:   oldAbility.UsedLevel,
			}
			db.DB.Create(&charAbility)
		}
	}

	migrateTraitPacks(char.ID, oldData.Merits, entities.V5TraitPackKindMerits)
	migrateTraitPacks(char.ID, oldData.Backgrounds, entities.V5TraitPackKindBackgrounds)

	for _, oldRitual := range oldData.BloodRituals {
		var ritual entities.V5BloodRitual
		if err := db.DB.Where("old_vicar_id = ?", oldRitual.ID).First(&ritual).Error; err == nil {
			db.DB.Exec("INSERT INTO v5_character_blood_rituals (v5_character_id, v5_blood_ritual_id) VALUES (?, ?)", char.ID, ritual.ID)
		}
	}

	for _, oldCeremony := range oldData.OblivionCeremonies {
		var ceremony entities.V5OblivionCeremony
		if err := db.DB.Where("old_vicar_id = ?", oldCeremony.ID).First(&ceremony).Error; err == nil {
			db.DB.Exec("INSERT INTO v5_character_oblivion_ceremonies (v5_character_id, v5_oblivion_ceremony_id) VALUES (?, ?)", char.ID, ceremony.ID)
		}
	}

	for _, oldSpread := range oldData.RequiredPointSpreads {
		spread := entities.V5RequiredPointSpread{
			CharacterID: char.ID,
			Type:        entities.V5CharacterTraitPackKind(oldSpread.Type),
			IsFlaw:      oldSpread.IsFlaw,
			Points:      oldSpread.Points,
		}

		if oldSpread.TraitPackID != nil {
			var pack entities.V5TraitPack
			if err := db.DB.Where("old_vicar_id = ?", *oldSpread.TraitPackID).First(&pack).Error; err == nil {
				spread.PackID = &pack.ID
			}
		}

		db.DB.Create(&spread)
	}

	for _, oldHistory := range oldData.LevelHistory {
		history := entities.V5LevelChange{
			CharacterID: char.ID,
			Type:        entities.LevelChangeType(oldHistory.Type),
			Date:        oldHistory.Date,
			Text:        oldHistory.Text,
			ExpUsed:     oldHistory.ExpUsed,
			ExpBefore:   oldHistory.ExpBefore,
			ExpAfter:    oldHistory.ExpAfter,
		}
		db.DB.Create(&history)
	}

	if oldData.Skillspread != nil {
		spreadType, ok := oldData.Skillspread["type"].(string)
		if ok {
			char.SkillSpreadType = entities.SkillSpreadType(spreadType)
			db.DB.Save(&char)
		}
	}

	var fullChar entities.V5Character
	db.DB.Preload("Clan.Disciplines").
		Preload("Clan.Book").
		Preload("PredatorType.Actions").
		Preload("PredatorType.Book").
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
		Preload("BloodRituals.Book").
		Preload("OblivionCeremonies.Book").
		Preload("LevelHistory").
		Preload("Viewers").
		First(&fullChar, char.ID)

	serializer := serialize.V5CharacterSerializer{IsOwner: true}
	return c.JSON(serializer.Serialize(fullChar))
}

func migrateTraitPacks(characterID uuid.UUID, oldPacks map[string]OldMeritData, kind entities.V5CharacterTraitPackKind) {
	for _, oldPack := range oldPacks {
		var pack entities.V5TraitPack
		if err := db.DB.Where("old_vicar_id = ?", oldPack.TraitPackID).First(&pack).Error; err != nil {
			continue
		}

		usage := entities.V5CharacterTraitPackUsage{
			CharacterID: characterID,
			Kind:        kind,
			PackID:      pack.ID,
		}
		db.DB.Create(&usage)

		for _, oldTrait := range oldPack.Traits {
			var trait entities.V5Trait
			if err := db.DB.Where("old_vicar_id = ?", oldTrait.TraitID).First(&trait).Error; err != nil {
				continue
			}

			charTrait := entities.V5CharacterTrait{
				UsageID:     usage.ID,
				TraitID:     trait.ID,
				IsLocked:    oldTrait.IsLocked,
				IsManual:    oldTrait.IsManual,
				CustomLevel: oldTrait.CustomLevel,
				Suffix:      oldTrait.Suffix,
			}
			db.DB.Create(&charTrait)
		}

		for _, oldTrait := range oldPack.FlawTraits {
			var trait entities.V5Trait
			if err := db.DB.Where("old_vicar_id = ?", oldTrait.TraitID).First(&trait).Error; err != nil {
				continue
			}

			charTrait := entities.V5CharacterFlawTrait{
				UsageID:     usage.ID,
				TraitID:     trait.ID,
				IsLocked:    oldTrait.IsLocked,
				IsManual:    oldTrait.IsManual,
				CustomLevel: oldTrait.CustomLevel,
				Suffix:      oldTrait.Suffix,
			}
			db.DB.Create(&charTrait)
		}
	}
}
