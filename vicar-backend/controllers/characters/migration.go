package characters

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"vicar-backend/auth"
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/log"
	"vicar-backend/serialize"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OldCharacterData struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Avatar              string   `json:"avatar"`
	Notes               string   `json:"notes"`
	Sex                 string   `json:"sex"`
	Concept             string   `json:"concept"`
	Chronicle           string   `json:"chronicle"`
	Exp                 int      `json:"exp"`
	UsedExp             int      `json:"usedExp"`
	ChroniclePrinciples string   `json:"chroniclePrinciples"`
	AnchorsAndBeliefs   string   `json:"anchorsAndBeliefs"`
	Backstory           string   `json:"backstory"`
	Sire                string   `json:"sire"`
	Ambition            string   `json:"ambition"`
	Desire              string   `json:"desire"`
	GenerationEra       string   `json:"generationEra"`
	Generation          int      `json:"generation"`
	Hunger              int      `json:"hunger"`
	Humanity            int      `json:"humanity"`
	Stains              int      `json:"stains"`
	Resonance           string   `json:"resonance"`
	BloodPotency        int      `json:"bloodPotency"`
	Health              int      `json:"health"`
	HealthDamage        []string `json:"healthDamage"`
	Willpower           int      `json:"willpower"`
	WillpowerDamage     []string `json:"willpowerDamage"`

	UseAdvancedDisciplines   bool `json:"useAdavancedDisciplines"`
	AllowLearningOfAllPowers bool `json:"allowLearningOfAllPowers"`
	FullCustomization        bool `json:"fullCustomization"`
	Version                  int  `json:"version"`

	Clan         *OldRef `json:"clan"`
	PredatorType *OldRef `json:"predatorType"`

	Books                []OldIDRef            `json:"books"`
	Categories           []OldCategoryData     `json:"categories"`
	Disciplines          []OldDisciplineData   `json:"disciplines"`
	Merits               OldPacksWrapper       `json:"merits"`
	Backgrounds          OldPacksWrapper       `json:"backgrounds"`
	BloodRituals         []OldIDRef            `json:"bloodRituals"`
	OblivionCeremonies   []OldIDRef            `json:"oblivionCeremonies"`
	RequiredPointSpreads []OldPointSpreadData  `json:"requiredPointSpreads"`
	LevelHistory         []OldLevelHistoryData `json:"levelHistory"`
	Skillspread          *OldSkillspreadData   `json:"skillspread"`

	Inventory *OldInventoryData `json:"inventory"`
}

type OldRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OldIDRef struct {
	ID int `json:"id"`
}

func (r *OldIDRef) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return nil
	}

	if b[0] >= '0' && b[0] <= '9' || b[0] == '-' {
		var id int
		if err := json.Unmarshal(b, &id); err != nil {
			return err
		}
		r.ID = id
		return nil
	}

	var tmp struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	r.ID = tmp.ID
	return nil
}

type OldCategoryData struct {
	Name       string             `json:"name"`
	Attributes []OldAttributeData `json:"attributes"`
	Skills     []OldSkillData     `json:"skills"`
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
	Discipline struct {
		ID int `json:"id"`
	} `json:"discipline"`
	Points       int                    `json:"points"`
	CurrentLevel int                    `json:"currentLevel"`
	Abilities    []OldDisciplineAbility `json:"abilities"`
}

type OldDisciplineAbility struct {
	ID        int `json:"id"`
	Level     int `json:"level"`
	UsedLevel int `json:"usedLevel"`
}

type OldPacksWrapper struct {
	Packs []OldMeritPack `json:"packs"`
}

type OldMeritPack struct {
	Pack struct {
		ID int `json:"id"`
	} `json:"pack"`
	Traits     []OldTraitData `json:"traits"`
	FlawTraits []OldTraitData `json:"flawTraits"`
}

type OldTraitData struct {
	ID          int     `json:"id"`
	IsLocked    bool    `json:"isLocked"`
	IsManual    bool    `json:"isManual"`
	CustomLevel *int    `json:"customLevel"`
	Suffix      *string `json:"suffix"`
}

type OldPointSpreadData struct {
	Type        string `json:"type"`
	IsFlaw      bool   `json:"isFlaw"`
	Points      int    `json:"points"`
	TraitPackID *int   `json:"traitPackId"`
}

type OldLevelHistoryData struct {
	Type string `json:"type"`
	Date string `json:"date"`
	Text string `json:"text"`
	Exp  struct {
		Used   int `json:"used"`
		Before int `json:"before"`
		After  int `json:"after"`
	} `json:"exp"`
}

type OldSkillspreadData struct {
	ID      int `json:"id"`
	Spreads []struct {
		Skills []string `json:"skills"`
		Value  int      `json:"value"`
	} `json:"spreads"`
}

type OldInventoryData struct {
	Bank         int `json:"bank"`
	Cash         int `json:"cash"`
	CarriedItems []struct {
		Amount int `json:"amount"`
		Item   struct {
			IsCustom    bool   `json:"isCustom"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Category    string `json:"category"`
		} `json:"item"`
	} `json:"carriedItems"`
	OwnedItems []struct {
		Amount int `json:"amount"`
		Item   struct {
			IsCustom    bool   `json:"isCustom"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Category    string `json:"category"`
		} `json:"item"`
	} `json:"ownedItems"`
}

type V5InventoryJSON struct {
	Cash         int                   `json:"cash"`
	Bank         int                   `json:"bank"`
	CarriedItems []V5InventoryItemJSON `json:"carriedItems"`
	OwnedItems   []V5InventoryItemJSON `json:"ownedItems"`
}

type V5InventoryItemJSON struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Category    string `json:"category"`
}

func migrateCharacter(c *fiber.Ctx) error {
	user := auth.Extract(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	var oldData OldCharacterData
	rawData := c.Body()
	if err := json.Unmarshal(rawData, &oldData); err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			problemPart := rawData[jsonErr.Offset-10 : jsonErr.Offset+10]
			err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
		}

		log.Error(log.Server, "❌", "Failed to parse migration data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var createdID uuid.UUID

	err := db.DB.Transaction(func(tx *gorm.DB) error {
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
			SkillSpreadType:          entities.SkillSpreadType("balanced"),
		}

		if invJSON, ok := buildInventoryJSON(oldData.Inventory); ok {
			raw, _ := json.Marshal(invJSON)
			char.Inventory = datatypes.JSON(raw)
		}

		internalRaw := buildInternalDataJSON(oldData)
		char.InternalData = internalRaw

		if oldData.Clan != nil {
			var clan entities.V5Clan
			if err := tx.Where("old_vicar_id = ?", oldData.Clan.ID).First(&clan).Error; err == nil {
				char.ClanID = &clan.ID
			}
		}

		if oldData.PredatorType != nil {
			var pt entities.V5PredatorType
			if err := tx.Where("old_vicar_id = ?", oldData.PredatorType.ID).First(&pt).Error; err == nil {
				char.PredatorTypeID = &pt.ID
			}
		}

		if oldData.Skillspread != nil {
			char.SkillSpreadType = guessSkillSpreadType(*oldData.Skillspread)
		}

		if err := tx.Create(&char).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create character")
		}

		createdID = char.ID

		for _, oldBook := range oldData.Books {
			var book entities.V5Book
			if err := tx.Where("old_vicar_id = ?", oldBook.ID).First(&book).Error; err == nil {
				tx.Exec("INSERT INTO v5_character_books (v5_character_id, v5_book_id) VALUES (?, ?)", char.ID, book.ID)
			}
		}

		for _, oldCat := range oldData.Categories {
			category := entities.V5CharacterCategory{
				CharacterID: char.ID,
				Name:        entities.CategoryKey(oldCat.Name),
			}
			tx.Create(&category)

			for _, oldAttr := range oldCat.Attributes {
				attr := entities.V5CharacterAttribute{
					CharacterID: char.ID,
					Category:    entities.CategoryKey(oldCat.Name),
					Key:         entities.AttributeKey(oldAttr.Key),
					Value:       oldAttr.Value,
				}
				tx.Create(&attr)
			}

			for _, oldSkill := range oldCat.Skills {
				skill := entities.V5CharacterSkill{
					CharacterID:    char.ID,
					Category:       entities.CategoryKey(oldCat.Name),
					Key:            entities.SkillKey(oldSkill.Key),
					Value:          oldSkill.Value,
					Specialization: pq.StringArray(oldSkill.Specialization),
				}
				tx.Create(&skill)
			}
		}

		for _, oldDisc := range oldData.Disciplines {
			var discipline entities.V5Discipline
			if err := tx.Where("old_vicar_id = ?", oldDisc.Discipline.ID).First(&discipline).Error; err != nil {
				continue
			}

			selection := entities.V5CharacterDisciplineSelection{
				CharacterID:  char.ID,
				DisciplineID: discipline.ID,
				Points:       oldDisc.Points,
				CurrentLevel: oldDisc.CurrentLevel,
			}
			tx.Create(&selection)

			for _, oldAbility := range oldDisc.Abilities {
				var ability entities.V5DisciplineAbility
				if err := tx.Where("old_vicar_id = ?", oldAbility.ID).First(&ability).Error; err != nil {
					continue
				}

				charAbility := entities.V5CharacterDisciplineAbility{
					SelectionID: selection.ID,
					AbilityID:   ability.ID,
					Level:       oldAbility.Level,
					UsedLevel:   oldAbility.UsedLevel,
				}
				tx.Create(&charAbility)
			}
		}

		migrateTraitPacks(tx, char.ID, oldData.Merits, entities.V5TraitPackKindMerits)
		migrateTraitPacks(tx, char.ID, oldData.Backgrounds, entities.V5TraitPackKindBackgrounds)

		for _, oldRitual := range oldData.BloodRituals {
			var ritual entities.V5BloodRitual
			if err := tx.Where("old_vicar_id = ?", oldRitual.ID).First(&ritual).Error; err == nil {
				tx.Exec("INSERT INTO v5_character_blood_rituals (v5_character_id, v5_blood_ritual_id) VALUES (?, ?)", char.ID, ritual.ID)
			}
		}

		for _, oldCeremony := range oldData.OblivionCeremonies {
			var ceremony entities.V5OblivionCeremony
			if err := tx.Where("old_vicar_id = ?", oldCeremony.ID).First(&ceremony).Error; err == nil {
				tx.Exec("INSERT INTO v5_character_oblivion_ceremonies (v5_character_id, v5_oblivion_ceremony_id) VALUES (?, ?)", char.ID, ceremony.ID)
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
				if err := tx.Where("old_vicar_id = ?", *oldSpread.TraitPackID).First(&pack).Error; err == nil {
					spread.PackID = &pack.ID
				}
			}

			tx.Create(&spread)
		}

		for _, oldHistory := range oldData.LevelHistory {
			history := entities.V5LevelChange{
				CharacterID: char.ID,
				Type:        entities.LevelChangeType(oldHistory.Type),
				Date:        oldHistory.Date,
				Text:        oldHistory.Text,
				ExpUsed:     oldHistory.Exp.Used,
				ExpBefore:   oldHistory.Exp.Before,
				ExpAfter:    oldHistory.Exp.After,
			}
			tx.Create(&history)
		}

		return nil
	})

	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return fe
		}
		return err
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
		First(&fullChar, createdID)

	serializer := serialize.V5CharacterSerializer{IsOwner: true}
	return c.JSON(serializer.Serialize(fullChar))
}

func migrateTraitPacks(tx *gorm.DB, characterID uuid.UUID, wrapper OldPacksWrapper, kind entities.V5CharacterTraitPackKind) {
	for _, oldPack := range wrapper.Packs {
		var pack entities.V5TraitPack
		if err := tx.Where("old_vicar_id = ?", oldPack.Pack.ID).First(&pack).Error; err != nil {
			continue
		}

		usage := entities.V5CharacterTraitPackUsage{
			CharacterID: characterID,
			Kind:        kind,
			PackID:      pack.ID,
		}
		tx.Create(&usage)

		for _, oldTrait := range oldPack.Traits {
			var trait entities.V5Trait
			if err := tx.Where("old_vicar_id = ?", oldTrait.ID).First(&trait).Error; err != nil {
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
			tx.Create(&charTrait)
		}

		for _, oldTrait := range oldPack.FlawTraits {
			var trait entities.V5Trait
			if err := tx.Where("old_vicar_id = ?", oldTrait.ID).First(&trait).Error; err != nil {
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
			tx.Create(&charTrait)
		}
	}
}

func buildInventoryJSON(oldInv *OldInventoryData) (V5InventoryJSON, bool) {
	if oldInv == nil {
		return V5InventoryJSON{}, false
	}

	next := V5InventoryJSON{
		Cash:         max0(oldInv.Cash),
		Bank:         max0(oldInv.Bank),
		CarriedItems: []V5InventoryItemJSON{},
		OwnedItems:   []V5InventoryItemJSON{},
	}

	for _, s := range oldInv.CarriedItems {
		name := strings.TrimSpace(s.Item.Name)
		desc := strings.TrimSpace(s.Item.Description)
		if name == "" && desc == "" {
			continue
		}
		next.CarriedItems = append(next.CarriedItems, V5InventoryItemJSON{
			ID:          "inv-" + uuid.NewString(),
			Name:        name,
			Description: desc,
			Amount:      max1(s.Amount),
			Category:    strings.TrimSpace(s.Item.Category),
		})
	}

	for _, s := range oldInv.OwnedItems {
		name := strings.TrimSpace(s.Item.Name)
		desc := strings.TrimSpace(s.Item.Description)
		if name == "" && desc == "" {
			continue
		}
		next.OwnedItems = append(next.OwnedItems, V5InventoryItemJSON{
			ID:          "inv-" + uuid.NewString(),
			Name:        name,
			Description: desc,
			Amount:      max1(s.Amount),
			Category:    strings.TrimSpace(s.Item.Category),
		})
	}

	return next, true
}

func buildInternalDataJSON(old OldCharacterData) datatypes.JSON {
	out := map[string]any{}

	legacy := map[string]any{}
	if old.Skillspread != nil {
		legacy["skillspread"] = old.Skillspread
	}
	if old.ID != "" {
		legacy["oldCharacterId"] = old.ID
	}
	if old.Clan != nil {
		legacy["oldClanId"] = old.Clan.ID
	}
	if old.PredatorType != nil {
		legacy["oldPredatorTypeId"] = old.PredatorType.ID
	}
	if len(legacy) > 0 {
		out["legacy"] = legacy
	}

	raw, _ := json.Marshal(out)
	return datatypes.JSON(raw)
}

func guessSkillSpreadType(s OldSkillspreadData) entities.SkillSpreadType {
	switch s.ID {
	case 1:
		return entities.SkillSpreadType("balanced")
	case 2:
		return entities.SkillSpreadType("specialist")
	case 3:
		return entities.SkillSpreadType("balanced")
	default:
		return entities.SkillSpreadType("balanced")
	}
}

func max0(v int) int {
	if v < 0 {
		return 0
	}
	return v
}

func max1(v int) int {
	if v < 1 {
		return 1
	}
	return v
}
