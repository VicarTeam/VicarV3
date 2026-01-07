package serialize

import (
	"vicar-backend/db/entities"

	"github.com/gofiber/fiber/v2"
)

type V5BookSerializer struct{}

func (s *V5BookSerializer) Serialize(book entities.V5Book, args ...any) any {
	return fiber.Map{
		"id":         book.ID,
		"name":       book.Name,
		"isOfficial": book.IsOfficial,
	}
}

type V5ClanSerializer struct{}

func (s *V5ClanSerializer) Serialize(clan entities.V5Clan, args ...any) any {
	m := fiber.Map{
		"id":          clan.ID,
		"name":        clan.Name,
		"slogan":      clan.Slogan,
		"description": clan.Description,
		"curse":       clan.Curse,
		"symbol":      clan.Symbol,
		"isHomebrew":  clan.IsHomebrew,
		"creator":     clan.Creator,
	}

	if clan.Book.ID.String() != "00000000-0000-0000-0000-000000000000" {
		m["book"] = fiber.Map{
			"id":   clan.Book.ID,
			"name": clan.Book.Name,
		}
	}

	if len(clan.Disciplines) > 0 {
		disciplines := make([]fiber.Map, len(clan.Disciplines))
		for i, disc := range clan.Disciplines {
			disciplines[i] = fiber.Map{
				"id":   disc.ID,
				"name": disc.Name,
			}
		}
		m["disciplines"] = disciplines
	}

	return m
}

type V5PredatorTypeSerializer struct{}

func (s *V5PredatorTypeSerializer) Serialize(pt entities.V5PredatorType, args ...any) any {
	m := fiber.Map{
		"id":          pt.ID,
		"name":        pt.Name,
		"description": pt.Description,
	}

	if pt.Book.ID.String() != "00000000-0000-0000-0000-000000000000" {
		m["book"] = fiber.Map{
			"id":   pt.Book.ID,
			"name": pt.Book.Name,
		}
	}

	if pt.Restriction != nil {
		m["restriction"] = fiber.Map{
			"type": pt.Restriction.Type,
			"data": pt.Restriction.Data,
		}
	}

	if len(pt.Actions) > 0 {
		actions := make([]fiber.Map, len(pt.Actions))
		for i, action := range pt.Actions {
			actionMap := fiber.Map{
				"id":          action.ID,
				"description": action.Description,
				"type":        action.Type,
				"data":        action.Data,
			}
			if action.Restriction != nil {
				actionMap["restriction"] = fiber.Map{
					"type": action.Restriction.Type,
					"data": action.Restriction.Data,
				}
			}
			actions[i] = actionMap
		}
		m["actions"] = actions
	}

	return m
}

type V5BloodRitualSerializer struct{}

func (s *V5BloodRitualSerializer) Serialize(ritual entities.V5BloodRitual, args ...any) any {
	m := fiber.Map{
		"id":          ritual.ID,
		"level":       ritual.Level,
		"name":        ritual.Name,
		"description": ritual.Description,
		"ingredients": ritual.Ingredients,
		"execution":   ritual.Execution,
		"system":      ritual.System,
	}

	if ritual.Book.ID.String() != "00000000-0000-0000-0000-000000000000" {
		m["book"] = fiber.Map{
			"id":   ritual.Book.ID,
			"name": ritual.Book.Name,
		}
	}

	return m
}

type V5OblivionCeremonySerializer struct{}

func (s *V5OblivionCeremonySerializer) Serialize(ceremony entities.V5OblivionCeremony, args ...any) any {
	m := fiber.Map{
		"id":          ceremony.ID,
		"level":       ceremony.Level,
		"name":        ceremony.Name,
		"cost":        ceremony.Cost,
		"roll":        ceremony.Roll,
		"summary":     ceremony.Summary,
		"requires":    ceremony.Requires,
		"cult":        ceremony.Cult,
		"ingredients": ceremony.Ingredients,
		"execution":   ceremony.Execution,
		"system":      ceremony.System,
		"duration":    ceremony.Duration,
	}

	if ceremony.Book.ID.String() != "00000000-0000-0000-0000-000000000000" {
		m["book"] = fiber.Map{
			"id":   ceremony.Book.ID,
			"name": ceremony.Book.Name,
		}
	}

	return m
}

type V5TraitPackSerializer struct{}

func (s *V5TraitPackSerializer) Serialize(pack entities.V5TraitPack, args ...any) any {
	m := fiber.Map{
		"id":           pack.ID,
		"type":         pack.Type,
		"name":         pack.Name,
		"description":  pack.Description,
		"specialRules": pack.SpecialRules,
	}

	if pack.Book.ID.String() != "00000000-0000-0000-0000-000000000000" {
		m["book"] = fiber.Map{
			"id":   pack.Book.ID,
			"name": pack.Book.Name,
		}
	}

	if pack.Restriction != nil {
		m["restriction"] = fiber.Map{
			"type": pack.Restriction.Type,
			"data": pack.Restriction.Data,
		}
	}

	if len(pack.PackTraits) > 0 {
		packTraits := make([]fiber.Map, len(pack.PackTraits))
		for i, pt := range pack.PackTraits {
			packTraits[i] = fiber.Map{
				"traitId": pt.TraitID,
				"side":    pt.Side,
				"trait": fiber.Map{
					"id":    pt.Trait.ID,
					"level": pt.Trait.Level,
					"name":  pt.Trait.Name,
				},
			}
		}
		m["packTraits"] = packTraits
	}

	return m
}

type V5TraitSerializer struct{}

func (s *V5TraitSerializer) Serialize(trait entities.V5Trait, args ...any) any {
	return fiber.Map{
		"id":           trait.ID,
		"level":        trait.Level,
		"name":         trait.Name,
		"description":  trait.Description,
		"isRepeatable": trait.IsRepeatable,
		"repeatAmount": trait.RepeatAmount,
		"repeatSize":   trait.RepeatSize,
		"requirement":  trait.Requirement,
		"actions":      trait.Actions,
	}
}

type V5DisciplineSerializer struct{}

func (s *V5DisciplineSerializer) Serialize(disc entities.V5Discipline, args ...any) any {
	m := fiber.Map{
		"id":         disc.ID,
		"name":       disc.Name,
		"summary":    disc.Summary,
		"note":       disc.Note,
		"isHomebrew": disc.IsHomebrew,
		"creator":    disc.Creator,
	}

	if len(disc.Abilities) > 0 {
		abilities := make([]fiber.Map, len(disc.Abilities))
		for i, ability := range disc.Abilities {
			abilities[i] = fiber.Map{
				"id":               ability.ID,
				"level":            ability.Level,
				"name":             ability.Name,
				"combinationRefId": ability.CombinationRefID,
				"combinationLevel": ability.CombinationLevel,
				"requirement":      ability.Requirement,
				"minBloodPotency":  ability.MinBloodPotency,
				"summary":          ability.Summary,
				"costs":            ability.Costs,
				"diceSupplies":     ability.DiceSupplies,
				"system":           ability.System,
				"alternatives":     ability.Alternatives,
				"duration":         ability.Duration,
			}
		}
		m["abilities"] = abilities
	}

	return m
}
