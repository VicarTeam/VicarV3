package v5data

import "github.com/gofiber/fiber/v2"

func Register(api fiber.Router) {
	v5data := api.Group("/v5data")

	v5data.Get("/books", getBooks)
	v5data.Get("/clans", getClans)
	v5data.Get("/predator-types", getPredatorTypes)
	v5data.Get("/blood-rituals", getBloodRituals)
	v5data.Get("/oblivion-ceremonies", getOblivionCeremonies)
	v5data.Get("/trait-packs", getTraitPacks)
	v5data.Get("/traits", getTraits)
	v5data.Get("/disciplines", getDisciplines)
}
