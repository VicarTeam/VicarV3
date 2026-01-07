package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type RestrictionType string

const (
	RestrictionSpecificClans         RestrictionType = "only_clans"
	RestrictionExcludeClans          RestrictionType = "exclude_clans"
	RestrictionMinimumCharacterValue RestrictionType = "minimum_character_value"
	RestrictionBookActivated         RestrictionType = "book_activated"
	RestrictionMaxGeneration         RestrictionType = "max_generation"
	RestrictionMaxBloodPotency       RestrictionType = "max_blood_potency"
)

type V5Restriction struct {
	Type *RestrictionType `gorm:"type:text"`
	Data datatypes.JSON   `gorm:"type:jsonb"`
}

type PTActionType string

const (
	PTActionAdditionalSpecialization     PTActionType = "additional_specialization"
	PTActionDisciplinePoint              PTActionType = "discipline_point"
	PTActionHumanityChange               PTActionType = "humanity_change"
	PTActionAddMerit                     PTActionType = "add_merit"
	PTActionAddBackground                PTActionType = "add_background"
	PTActionAddBackgroundPoints          PTActionType = "add_background_points"
	PTActionAddFlaw                      PTActionType = "add_flaw"
	PTActionBloodPotencyChange           PTActionType = "blood_potency_change"
	PTActionSpendBackgroundPointsBetween PTActionType = "spend_background_points_between"
	PTActionSpendFlawPointsBetween       PTActionType = "spend_flaw_points_between"
)

type TraitSpecialRules string

const (
	TraitSpecialRulesNone   TraitSpecialRules = "none"
	TraitSpecialRulesAllies TraitSpecialRules = "allies"
	TraitSpecialRulesHaven  TraitSpecialRules = "haven"
	TraitSpecialRulesMask   TraitSpecialRules = "mask"
)

type TraitPackType string

const (
	TraitPackTypeMerits      TraitPackType = "merits"
	TraitPackTypeBackgrounds TraitPackType = "backgrounds"
)

type TraitActionType string

const (
	TraitActionCapSkill TraitActionType = "cap_skill"
)

type TraitPackSide string

const (
	TraitPackSideAdvantage    TraitPackSide = "advantage"
	TraitPackSideDisadvantage TraitPackSide = "disadvantage"
)

type V5Book struct {
	OldVicarID         uint                 `gorm:"type:bigint"`
	ID                 uuid.UUID            `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name               string               `gorm:"type:varchar(255);not null;unique"`
	IsOfficial         bool                 `gorm:"type:boolean;default:false"`
	BloodRituals       []V5BloodRitual      `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OblivionCeremonies []V5OblivionCeremony `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Clans              []V5Clan             `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TraitPacks         []V5TraitPack        `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PredatorTypes      []V5PredatorType     `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5BloodRitual struct {
	OldVicarID  uint      `gorm:"type:bigint"`
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Book        V5Book    `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Level       int       `gorm:"type:int;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	Ingredients string    `gorm:"type:text;not null"`
	Execution   string    `gorm:"type:text;not null"`
	System      string    `gorm:"type:text;not null"`
}

type V5OblivionCeremony struct {
	OldVicarID  uint      `gorm:"type:bigint"`
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Book        V5Book    `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Level       int       `gorm:"type:int;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Cost        string    `gorm:"type:text;not null"`
	Roll        string    `gorm:"type:text;not null"`
	Summary     string    `gorm:"type:text;not null"`
	Requires    *int      `gorm:"type:int"`
	Cult        *string   `gorm:"type:text"`
	Ingredients string    `gorm:"type:text;not null"`
	Execution   string    `gorm:"type:text;not null"`
	System      string    `gorm:"type:text;not null"`
	Duration    *string   `gorm:"type:text"`
}

type V5Clan struct {
	OldVicarID  uint           `gorm:"type:bigint"`
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	Book        V5Book         `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string         `gorm:"type:varchar(255);not null;unique"`
	Slogan      string         `gorm:"type:text;not null"`
	Description string         `gorm:"type:text;not null"`
	Curse       string         `gorm:"type:text;not null"`
	Symbol      *string        `gorm:"type:text"`
	IsHomebrew  bool           `gorm:"type:boolean;default:false"`
	Creator     *string        `gorm:"type:text"`
	Disciplines []V5Discipline `gorm:"many2many:v5_clan_disciplines;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5Discipline struct {
	OldVicarID uint                  `gorm:"type:bigint"`
	ID         uuid.UUID             `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name       string                `gorm:"type:varchar(255);not null;unique"`
	Summary    *string               `gorm:"type:text"`
	Note       *string               `gorm:"type:text"`
	IsHomebrew bool                  `gorm:"type:boolean;default:false"`
	Creator    *string               `gorm:"type:text"`
	Abilities  []V5DisciplineAbility `gorm:"foreignKey:DisciplineID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Clans      []V5Clan              `gorm:"many2many:v5_clan_disciplines;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5DisciplineAbility struct {
	OldVicarID       uint           `gorm:"type:bigint"`
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DisciplineID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Discipline       V5Discipline   `gorm:"foreignKey:DisciplineID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Level            int            `gorm:"type:int;not null;index"`
	Name             string         `gorm:"type:varchar(255);not null"`
	CombinationRefID *uuid.UUID     `gorm:"type:uuid;index"`
	CombinationLevel *int           `gorm:"type:int"`
	Requirement      *int           `gorm:"type:int"`
	MinBloodPotency  *int           `gorm:"type:int"`
	Summary          string         `gorm:"type:text;not null"`
	Costs            string         `gorm:"type:text;not null"`
	DiceSupplies     *string        `gorm:"type:text"`
	System           string         `gorm:"type:text;not null"`
	Alternatives     pq.StringArray `gorm:"type:text[];not null;default:'{}'"`
	Duration         string         `gorm:"type:text;not null"`
}

type V5PredatorType struct {
	OldVicarID  uint               `gorm:"type:bigint"`
	ID          uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookID      uuid.UUID          `gorm:"type:uuid;not null;index"`
	Book        V5Book             `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string             `gorm:"type:varchar(255);not null;unique"`
	Description string             `gorm:"type:text;not null"`
	Restriction *V5Restriction     `gorm:"embedded;embeddedPrefix:restriction_"`
	Actions     []V5PredatorAction `gorm:"foreignKey:PredatorTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5PredatorAction struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PredatorTypeID uuid.UUID       `gorm:"type:uuid;not null;index"`
	PredatorType   *V5PredatorType `gorm:"foreignKey:PredatorTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description    string          `gorm:"type:text;not null"`
	Type           PTActionType    `gorm:"type:text;not null"`
	Data           datatypes.JSON  `gorm:"type:jsonb;not null"`
	Restriction    *V5Restriction  `gorm:"embedded;embeddedPrefix:restriction_"`
}

type V5TraitPack struct {
	OldVicarID   uint               `gorm:"type:bigint"`
	ID           uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookID       uuid.UUID          `gorm:"type:uuid;not null;index"`
	Book         V5Book             `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type         TraitPackType      `gorm:"type:text;not null"`
	Name         string             `gorm:"type:varchar(255);not null"`
	Description  string             `gorm:"type:text;not null"`
	SpecialRules TraitSpecialRules  `gorm:"type:text;not null;default:'none'"`
	Restriction  *V5Restriction     `gorm:"embedded;embeddedPrefix:restriction_"`
	PackTraits   []V5TraitPackTrait `gorm:"foreignKey:TraitPackID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5Trait struct {
	OldVicarID   uint               `gorm:"type:bigint"`
	ID           uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Level        int                `gorm:"type:int;not null"`
	Name         string             `gorm:"type:varchar(255);not null"`
	Description  string             `gorm:"type:text;not null"`
	IsRepeatable bool               `gorm:"type:boolean;not null;default:false"`
	RepeatAmount *int               `gorm:"type:int"`
	RepeatSize   *int               `gorm:"type:int"`
	Requirement  datatypes.JSON     `gorm:"type:jsonb"`
	Actions      datatypes.JSON     `gorm:"type:jsonb;not null"`
	Restriction  *V5Restriction     `gorm:"embedded;embeddedPrefix:restriction_"`
	PackTraits   []V5TraitPackTrait `gorm:"foreignKey:TraitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5TraitPackTrait struct {
	TraitPackID uuid.UUID     `gorm:"type:uuid;not null;primaryKey"`
	TraitPack   V5TraitPack   `gorm:"foreignKey:TraitPackID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TraitID     uuid.UUID     `gorm:"type:uuid;not null;primaryKey"`
	Trait       V5Trait       `gorm:"foreignKey:TraitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Side        TraitPackSide `gorm:"type:text;not null;primaryKey"`
}

type V5BloodPotencyData struct {
	ID                         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Value                      int       `gorm:"type:int;not null;uniqueIndex"`
	BleedingSpurt              int       `gorm:"type:int;not null"`
	HealedDamage               int       `gorm:"type:int;not null"`
	DisciplineBonus            int       `gorm:"type:int;not null"`
	RouseRepeatDisciplineLevel int       `gorm:"type:int;not null"`
	BaneLevel                  int       `gorm:"type:int;not null"`
	Pray                       string    `gorm:"type:text;not null"`
}
