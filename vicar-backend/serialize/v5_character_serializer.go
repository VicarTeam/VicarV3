package serialize

import (
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
)

type V5CharacterSummarySerializer struct{}

func (s *V5CharacterSummarySerializer) Serialize(char entities.V5Character, args ...any) any {
	m := fiber.Map{
		"id":            char.ID,
		"name":          char.Name,
		"avatar":        char.Avatar,
		"generation":    char.Generation,
		"generationEra": char.GenerationEra,
	}

	if char.Clan != nil {
		m["clan"] = fiber.Map{
			"id":   char.Clan.ID,
			"name": char.Clan.Name,
		}
	}

	if char.DirectoryID != nil {
		m["directoryId"] = char.DirectoryID
	}

	return m
}

type V5CharacterSerializer struct {
	IsOwner bool
}

func (s *V5CharacterSerializer) Serialize(char entities.V5Character, args ...any) any {
	m := fiber.Map{
		"id":        char.ID,
		"name":      char.Name,
		"avatar":    char.Avatar,
		"notes":     char.Notes,
		"sex":       char.Sex,
		"concept":   char.Concept,
		"chronicle": char.Chronicle,
		"exp":       char.Exp,
		"usedExp":   char.UsedExp,

		"chroniclePrinciples": char.ChroniclePrinciples,
		"anchorsAndBeliefs":   char.AnchorsAndBeliefs,
		"backstory":           char.Backstory,
		"sire":                char.Sire,
		"ambition":            char.Ambition,
		"desire":              char.Desire,

		"generationEra": char.GenerationEra,
		"generation":    char.Generation,

		"hunger":       char.Hunger,
		"humanity":     char.Humanity,
		"stains":       char.Stains,
		"resonance":    char.Resonance,
		"bloodPotency": char.BloodPotency,

		"health":          char.Health,
		"healthDamage":    char.HealthDamage,
		"willpower":       char.Willpower,
		"willpowerDamage": char.WillpowerDamage,

		"useAdvancedDisciplines":   char.UseAdavancedDisciplines,
		"allowLearningOfAllPowers": char.AllowLearningOfAllPowers,
		"fullCustomization":        char.FullCustomization,
		"version":                  char.Version,
		"skillSpreadType":          char.SkillSpreadType,
		"internalData":             char.InternalData,
	}

	if s.IsOwner {
		m["userId"] = char.UserID
	}

	if char.DirectoryID != nil {
		m["directoryId"] = char.DirectoryID
	}

	if char.ClanID != nil {
		m["clanID"] = char.ClanID
	}

	if char.PredatorTypeID != nil {
		m["predatorTypeID"] = char.PredatorTypeID
	}

	m["books"] = serializeBooks(char.Books)
	m["categories"] = serializeCategories(char.Categories)
	m["attributes"] = serializeAttributes(char.Attributes)
	m["skills"] = serializeSkills(char.Skills)
	m["requiredPointSpreads"] = serializeRequiredPointSpreads(char.RequiredPointSpreads)
	m["disciplineSelections"] = serializeDisciplineSelections(char.DisciplineSelections)
	m["traitPackUsages"] = serializeTraitPackUsages(char.TraitPackUsages)
	m["bloodRituals"] = serializeBloodRituals(char.BloodRituals)
	m["oblivionCeremonies"] = serializeOblivionCeremonies(char.OblivionCeremonies)
	m["levelHistory"] = serializeLevelHistory(char.LevelHistory)

	if s.IsOwner && len(char.Viewers) > 0 {
		m["viewers"] = DoPtrArray(&UserSerializer{ForMyself: false}, convertUsersToPointers(char.Viewers))
	}

	return m
}

func serializeClan(clan entities.V5Clan) fiber.Map {
	return fiber.Map{
		"id":          clan.ID,
		"name":        clan.Name,
		"slogan":      clan.Slogan,
		"description": clan.Description,
		"curse":       clan.Curse,
		"symbol":      clan.Symbol,
	}
}

func serializePredatorType(pt entities.V5PredatorType) fiber.Map {
	return fiber.Map{
		"id":          pt.ID,
		"name":        pt.Name,
		"description": pt.Description,
	}
}

func serializeBooks(books []entities.V5Book) []fiber.Map {
	result := make([]fiber.Map, len(books))
	for i, book := range books {
		result[i] = fiber.Map{
			"id":         book.ID,
			"oldVicarID": book.OldVicarID,
			"name":       book.Name,
		}
	}
	return result
}

func serializeCategories(categories []entities.V5CharacterCategory) []fiber.Map {
	result := make([]fiber.Map, len(categories))
	for i, cat := range categories {
		result[i] = fiber.Map{
			"id":   cat.ID,
			"name": cat.Name,
		}
	}
	return result
}

func serializeAttributes(attributes []entities.V5CharacterAttribute) []fiber.Map {
	result := make([]fiber.Map, len(attributes))
	for i, attr := range attributes {
		result[i] = fiber.Map{
			"id":       attr.ID,
			"category": attr.Category,
			"key":      attr.Key,
			"value":    attr.Value,
		}
	}
	return result
}

func serializeSkills(skills []entities.V5CharacterSkill) []fiber.Map {
	result := make([]fiber.Map, len(skills))
	for i, skill := range skills {
		result[i] = fiber.Map{
			"id":             skill.ID,
			"category":       skill.Category,
			"key":            skill.Key,
			"value":          skill.Value,
			"specialization": skill.Specialization,
		}
	}
	return result
}

func serializeRequiredPointSpreads(spreads []entities.V5RequiredPointSpread) []fiber.Map {
	result := make([]fiber.Map, len(spreads))
	for i, spread := range spreads {
		m := fiber.Map{
			"id":     spread.ID,
			"type":   spread.Type,
			"isFlaw": spread.IsFlaw,
			"points": spread.Points,
		}
		if spread.Pack != nil {
			m["pack"] = fiber.Map{
				"id":   spread.Pack.ID,
				"name": spread.Pack.Name,
			}
		}
		result[i] = m
	}
	return result
}

func serializeDisciplineSelections(selections []entities.V5CharacterDisciplineSelection) []fiber.Map {
	result := make([]fiber.Map, len(selections))
	for i, sel := range selections {
		result[i] = fiber.Map{
			"id":           sel.ID,
			"characterID":  sel.CharacterID,
			"disciplineID": sel.DisciplineID,
			"discipline":   serializeDiscipline(sel.Discipline),
			"points":       sel.Points,
			"currentLevel": sel.CurrentLevel,
			"abilities":    serializeDisciplineAbilities(sel.Abilities),
		}
	}
	return result
}

func serializeDiscipline(disc entities.V5Discipline) fiber.Map {
	return fiber.Map{
		"id":   disc.ID,
		"name": disc.Name,
	}
}

func serializeDisciplineAbilities(abilities []entities.V5CharacterDisciplineAbility) []fiber.Map {
	result := make([]fiber.Map, len(abilities))
	for i, ability := range abilities {
		result[i] = fiber.Map{
			"id":          ability.ID,
			"selectionID": ability.SelectionID,
			"abilityID":   ability.AbilityID,
			"level":       ability.Level,
			"usedLevel":   ability.UsedLevel,
		}
	}
	return result
}

func serializeTraitPackUsages(usages []entities.V5CharacterTraitPackUsage) []fiber.Map {
	result := make([]fiber.Map, len(usages))
	for i, usage := range usages {
		result[i] = fiber.Map{
			"id":         usage.ID,
			"kind":       usage.Kind,
			"pack":       serializeTraitPack(usage.Pack),
			"traits":     serializeTraits(usage.Traits),
			"flawTraits": serializeFlawTraits(usage.FlawTraits),
		}
	}
	return result
}

func serializeTraitPack(pack entities.V5TraitPack) fiber.Map {
	return fiber.Map{
		"id":   pack.ID,
		"name": pack.Name,
		"type": pack.Type,
	}
}

func serializeTraits(traits []entities.V5CharacterTrait) []fiber.Map {
	result := make([]fiber.Map, len(traits))
	for i, trait := range traits {
		result[i] = fiber.Map{
			"id":          trait.ID,
			"traitId":     trait.TraitID,
			"isLocked":    trait.IsLocked,
			"isManual":    trait.IsManual,
			"customLevel": trait.CustomLevel,
			"suffix":      trait.Suffix,
		}
	}
	return result
}

func serializeFlawTraits(traits []entities.V5CharacterFlawTrait) []fiber.Map {
	result := make([]fiber.Map, len(traits))
	for i, trait := range traits {
		result[i] = fiber.Map{
			"id":          trait.ID,
			"traitId":     trait.TraitID,
			"isLocked":    trait.IsLocked,
			"isManual":    trait.IsManual,
			"customLevel": trait.CustomLevel,
			"suffix":      trait.Suffix,
		}
	}
	return result
}

func serializeBloodRituals(rituals []entities.V5BloodRitual) []fiber.Map {
	result := make([]fiber.Map, len(rituals))
	for i, ritual := range rituals {
		result[i] = fiber.Map{
			"id":    ritual.ID,
			"name":  ritual.Name,
			"level": ritual.Level,
		}
	}
	return result
}

func serializeOblivionCeremonies(ceremonies []entities.V5OblivionCeremony) []fiber.Map {
	result := make([]fiber.Map, len(ceremonies))
	for i, ceremony := range ceremonies {
		result[i] = fiber.Map{
			"id":    ceremony.ID,
			"name":  ceremony.Name,
			"level": ceremony.Level,
		}
	}
	return result
}

func serializeLevelHistory(history []entities.V5LevelChange) []fiber.Map {
	result := make([]fiber.Map, len(history))
	for i, change := range history {
		result[i] = fiber.Map{
			"id":        change.ID,
			"type":      change.Type,
			"date":      change.Date,
			"text":      change.Text,
			"expUsed":   change.ExpUsed,
			"expBefore": change.ExpBefore,
			"expAfter":  change.ExpAfter,
		}
	}
	return result
}

func convertUsersToPointers(users []entities.User) []*entities.User {
	result := make([]*entities.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result
}

type CharacterDirectorySerializer struct {
	IncludeCharacterCount bool
}

func (s *CharacterDirectorySerializer) Serialize(dir entities.CharacterDirectory, args ...any) any {
	m := fiber.Map{
		"id":   dir.ID,
		"name": dir.Name,
		"open": dir.Open,
	}

	if s.IncludeCharacterCount {
		m["characterCount"] = len(dir.V5Characters)
	}

	return m
}
