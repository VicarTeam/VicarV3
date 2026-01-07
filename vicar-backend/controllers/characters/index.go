package characters

import "github.com/gofiber/fiber/v2"

func Register(api fiber.Router) {
	characters := api.Group("/characters")

	characters.Get("", getCharacters)
	characters.Post("", createCharacter)
	characters.Post("/migrate", migrateCharacter)

	characters.Get("/:id", getCharacter)
	characters.Patch("/:id", updateCharacter)
	characters.Delete("/:id", deleteCharacter)

	characters.Patch("/:id/attributes/:attrId", updateAttribute)
	characters.Patch("/:id/skills/:skillId", updateSkill)

	characters.Post("/:id/disciplines", addDisciplineSelection)
	characters.Patch("/:id/disciplines/:disciplineId", updateDisciplineSelection)
	characters.Delete("/:id/disciplines/:disciplineId", deleteDisciplineSelection)

	characters.Patch("/:id/directory", moveCharacterToDirectory)

	characters.Get("/:id/viewers", getCharacterViewers)
	characters.Post("/:id/viewers", addViewer)
	characters.Delete("/:id/viewers/:userId", removeViewer)

	directories := api.Group("/directories")
	directories.Get("", getDirectories)
	directories.Post("", createDirectory)
	directories.Patch("/:id", updateDirectory)
	directories.Delete("/:id", deleteDirectory)
}
