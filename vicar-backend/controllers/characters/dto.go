package characters

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CreateCharacterDto struct {
	Name          string      `json:"name"`
	GenerationEra string      `json:"generationEra"`
	Generation    int         `json:"generation"`
	BookIDs       []uuid.UUID `json:"books"`
}

func (dto *CreateCharacterDto) Validate() error {
	if dto.Name == "" {
		return errors.New("name is required")
	}
	if len(dto.Name) > 255 {
		return errors.New("name must be less than 255 characters")
	}
	if dto.Generation < 2 || dto.Generation > 13 {
		return errors.New("generation must be between 2 and 13")
	}

	validEras := []string{"children", "newborn", "ancillae", "older", "elder", "cainesinheritance"}
	isValidEra := false
	for _, era := range validEras {
		if dto.GenerationEra == era {
			isValidEra = true
			break
		}
	}
	if !isValidEra {
		return errors.New("invalid generation era")
	}

	if len(dto.BookIDs) == 0 {
		return errors.New("at least one book must be selected")
	}

	return nil
}

type UpdateCharacterDto struct {
	Name      *string `json:"name"`
	Avatar    *string `json:"avatar"`
	Notes     *string `json:"notes"`
	Sex       *string `json:"sex"`
	Concept   *string `json:"concept"`
	Chronicle *string `json:"chronicle"`
	Exp       *int    `json:"exp"`
	UsedExp   *int    `json:"usedExp"`

	ChroniclePrinciples *string `json:"chroniclePrinciples"`
	AnchorsAndBeliefs   *string `json:"anchorsAndBeliefs"`
	Backstory           *string `json:"backstory"`
	Sire                *string `json:"sire"`
	Ambition            *string `json:"ambition"`
	Desire              *string `json:"desire"`

	ClanID         *uuid.UUID `json:"clanId"`
	PredatorTypeID *uuid.UUID `json:"predatorTypeId"`

	GenerationEra *string `json:"generationEra"`
	Generation    *int    `json:"generation"`

	Hunger       *int    `json:"hunger"`
	Humanity     *int    `json:"humanity"`
	Stains       *int    `json:"stains"`
	Resonance    *string `json:"resonance"`
	BloodPotency *int    `json:"bloodPotency"`

	Health          *int           `json:"health"`
	HealthDamage    pq.StringArray `json:"healthDamage"`
	Willpower       *int           `json:"willpower"`
	WillpowerDamage pq.StringArray `json:"willpowerDamage"`

	UseAdvancedDisciplines   *bool       `json:"useAdvancedDisciplines"`
	AllowLearningOfAllPowers *bool       `json:"allowLearningOfAllPowers"`
	FullCustomization        *bool       `json:"fullCustomization"`
	SkillSpreadType          *string     `json:"skillSpreadType"`
	BookIDs                  []uuid.UUID `json:"bookIds"`
}

func (dto *UpdateCharacterDto) Validate() error {
	if dto.Generation != nil && (*dto.Generation < 2 || *dto.Generation > 13) {
		return errors.New("generation must be between 2 and 13")
	}
	if dto.Hunger != nil && (*dto.Hunger < 0 || *dto.Hunger > 5) {
		return errors.New("hunger must be between 0 and 5")
	}
	if dto.Humanity != nil && (*dto.Humanity < 0 || *dto.Humanity > 10) {
		return errors.New("humanity must be between 0 and 10")
	}
	if dto.Stains != nil && (*dto.Stains < 0 || *dto.Stains > 10) {
		return errors.New("stains must be between 0 and 10")
	}
	if dto.BloodPotency != nil && (*dto.BloodPotency < 0 || *dto.BloodPotency > 10) {
		return errors.New("blood potency must be between 0 and 10")
	}
	if dto.Health != nil && (*dto.Health < 0 || *dto.Health > 20) {
		return errors.New("health must be between 0 and 20")
	}
	if dto.Willpower != nil && (*dto.Willpower < 0 || *dto.Willpower > 20) {
		return errors.New("willpower must be between 0 and 20")
	}

	return nil
}

type UpdateAttributeDto struct {
	Value int `json:"value"`
}

func (dto *UpdateAttributeDto) Validate() error {
	if dto.Value < 0 || dto.Value > 5 {
		return errors.New("attribute value must be between 0 and 5")
	}
	return nil
}

type UpdateSkillDto struct {
	Value          int            `json:"value"`
	Specialization pq.StringArray `json:"specialization"`
}

func (dto *UpdateSkillDto) Validate() error {
	if dto.Value < 0 || dto.Value > 5 {
		return errors.New("skill value must be between 0 and 5")
	}
	return nil
}

type CreateDisciplineSelectionDto struct {
	DisciplineID uuid.UUID `json:"disciplineId"`
}

type UpdateDisciplineSelectionDto struct {
	Points       *int `json:"points"`
	CurrentLevel *int `json:"currentLevel"`
}

func (dto *UpdateDisciplineSelectionDto) Validate() error {
	if dto.CurrentLevel != nil && (*dto.CurrentLevel < 0 || *dto.CurrentLevel > 5) {
		return errors.New("discipline level must be between 0 and 5")
	}
	if dto.Points != nil && *dto.Points < 0 {
		return errors.New("discipline points must be >= 0")
	}
	return nil
}

type CreateDisciplineAbilityDto struct {
	AbilityID uuid.UUID `json:"abilityId"`
	Level     int       `json:"level"`
}

type UpdateDisciplineAbilityDto struct {
	Level     *int `json:"level"`
	UsedLevel *int `json:"usedLevel"`
}

type CreateTraitPackUsageDto struct {
	Kind   string    `json:"kind"`
	PackID uuid.UUID `json:"packId"`
}

type AddTraitToPackDto struct {
	TraitID     uuid.UUID `json:"traitId"`
	IsLocked    bool      `json:"isLocked"`
	IsManual    bool      `json:"isManual"`
	CustomLevel *int      `json:"customLevel"`
	Suffix      *string   `json:"suffix"`
}

type UpdateTraitDto struct {
	IsLocked    *bool   `json:"isLocked"`
	IsManual    *bool   `json:"isManual"`
	CustomLevel *int    `json:"customLevel"`
	Suffix      *string `json:"suffix"`
}

type AddViewerDto struct {
	UserID uuid.UUID `json:"userId"`
}

type CreateDirectoryDto struct {
	Name string `json:"name"`
}

func (dto *CreateDirectoryDto) Validate() error {
	if dto.Name == "" {
		return errors.New("directory name is required")
	}
	if len(dto.Name) > 255 {
		return errors.New("directory name must be less than 255 characters")
	}
	return nil
}

type UpdateDirectoryDto struct {
	Name string `json:"name"`
}

func (dto *UpdateDirectoryDto) Validate() error {
	if dto.Name == "" {
		return errors.New("directory name is required")
	}
	if len(dto.Name) > 255 {
		return errors.New("directory name must be less than 255 characters")
	}
	return nil
}

type MoveToDirectoryDto struct {
	DirectoryID *uuid.UUID `json:"directoryId"`
}
