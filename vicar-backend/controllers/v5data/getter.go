package v5data

import (
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/serialize"

	"github.com/gofiber/fiber/v2"
)

func getBooks(c *fiber.Ctx) error {
	var books []entities.V5Book
	if err := db.DB.Order("name ASC").Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books",
		})
	}

	serializer := &serialize.V5BookSerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertBooksToPointers(books)))
}

func getClans(c *fiber.Ctx) error {
	var clans []entities.V5Clan
	if err := db.DB.Preload("Disciplines").Preload("Book").Order("name ASC").Find(&clans).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch clans",
		})
	}

	serializer := &serialize.V5ClanSerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertClansToPointers(clans)))
}

func getPredatorTypes(c *fiber.Ctx) error {
	var predatorTypes []entities.V5PredatorType
	if err := db.DB.Preload("Actions").Preload("Book").Order("name ASC").Find(&predatorTypes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch predator types",
		})
	}

	serializer := &serialize.V5PredatorTypeSerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertPredatorTypesToPointers(predatorTypes)))
}

func getBloodRituals(c *fiber.Ctx) error {
	var bloodRituals []entities.V5BloodRitual
	if err := db.DB.Preload("Book").Order("level ASC, name ASC").Find(&bloodRituals).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch blood rituals",
		})
	}

	serializer := &serialize.V5BloodRitualSerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertBloodRitualsToPointers(bloodRituals)))
}

func getOblivionCeremonies(c *fiber.Ctx) error {
	var oblivionCeremonies []entities.V5OblivionCeremony
	if err := db.DB.Preload("Book").Order("level ASC, name ASC").Find(&oblivionCeremonies).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch oblivion ceremonies",
		})
	}

	serializer := &serialize.V5OblivionCeremonySerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertOblivionCeremoniesToPointers(oblivionCeremonies)))
}

func getTraitPacks(c *fiber.Ctx) error {
	var traitPacks []entities.V5TraitPack
	if err := db.DB.Preload("Book").Preload("PackTraits.Trait").Order("type ASC, name ASC").Find(&traitPacks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch trait packs",
		})
	}

	serializer := &serialize.V5TraitPackSerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertTraitPacksToPointers(traitPacks)))
}

func getTraits(c *fiber.Ctx) error {
	var traits []entities.V5Trait
	if err := db.DB.Order("name ASC").Find(&traits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch traits",
		})
	}

	serializer := &serialize.V5TraitSerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertTraitsToPointers(traits)))
}

func getDisciplines(c *fiber.Ctx) error {
	var disciplines []entities.V5Discipline
	if err := db.DB.Preload("Abilities").Order("name ASC").Find(&disciplines).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch disciplines",
		})
	}

	serializer := &serialize.V5DisciplineSerializer{}
	return c.JSON(serialize.DoPtrArray(serializer, convertDisciplinesToPointers(disciplines)))
}

func convertBooksToPointers(books []entities.V5Book) []*entities.V5Book {
	result := make([]*entities.V5Book, len(books))
	for i := range books {
		result[i] = &books[i]
	}
	return result
}

func convertClansToPointers(clans []entities.V5Clan) []*entities.V5Clan {
	result := make([]*entities.V5Clan, len(clans))
	for i := range clans {
		result[i] = &clans[i]
	}
	return result
}

func convertPredatorTypesToPointers(predatorTypes []entities.V5PredatorType) []*entities.V5PredatorType {
	result := make([]*entities.V5PredatorType, len(predatorTypes))
	for i := range predatorTypes {
		result[i] = &predatorTypes[i]
	}
	return result
}

func convertBloodRitualsToPointers(bloodRituals []entities.V5BloodRitual) []*entities.V5BloodRitual {
	result := make([]*entities.V5BloodRitual, len(bloodRituals))
	for i := range bloodRituals {
		result[i] = &bloodRituals[i]
	}
	return result
}

func convertOblivionCeremoniesToPointers(oblivionCeremonies []entities.V5OblivionCeremony) []*entities.V5OblivionCeremony {
	result := make([]*entities.V5OblivionCeremony, len(oblivionCeremonies))
	for i := range oblivionCeremonies {
		result[i] = &oblivionCeremonies[i]
	}
	return result
}

func convertTraitPacksToPointers(traitPacks []entities.V5TraitPack) []*entities.V5TraitPack {
	result := make([]*entities.V5TraitPack, len(traitPacks))
	for i := range traitPacks {
		result[i] = &traitPacks[i]
	}
	return result
}

func convertTraitsToPointers(traits []entities.V5Trait) []*entities.V5Trait {
	result := make([]*entities.V5Trait, len(traits))
	for i := range traits {
		result[i] = &traits[i]
	}
	return result
}

func convertDisciplinesToPointers(disciplines []entities.V5Discipline) []*entities.V5Discipline {
	result := make([]*entities.V5Discipline, len(disciplines))
	for i := range disciplines {
		result[i] = &disciplines[i]
	}
	return result
}
