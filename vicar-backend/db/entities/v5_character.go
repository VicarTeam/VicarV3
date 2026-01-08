package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type GenerationEra string

const (
	GenerationChildren          GenerationEra = "children"
	GenerationNewborn           GenerationEra = "newborn"
	GenerationAncillae          GenerationEra = "ancillae"
	GenerationOlder             GenerationEra = "older"
	GenerationElder             GenerationEra = "elder"
	GenerationCainesInheritance GenerationEra = "cainesinheritance"
)

type V5Resonance string

const (
	V5ResonanceEmpty       V5Resonance = ""
	V5ResonanceCholeric    V5Resonance = "choleric"
	V5ResonanceMelancholic V5Resonance = "melancholic"
	V5ResonancePhlegmatic  V5Resonance = "phlegmatic"
	V5ResonanceSanguine    V5Resonance = "sanguine"
	V5ResonanceAnimalBlood V5Resonance = "animalblood"
)

type V5CharacterTraitPackKind string

const (
	V5TraitPackKindMerits      V5CharacterTraitPackKind = "merits"
	V5TraitPackKindBackgrounds V5CharacterTraitPackKind = "backgrounds"
)

type LevelChangeType string

const (
	LevelChangeAttribute      LevelChangeType = "attribute"
	LevelChangeSkill          LevelChangeType = "skill"
	LevelChangeDiscipline     LevelChangeType = "discipline"
	LevelChangeTrait          LevelChangeType = "trait"
	LevelChangeFlaw           LevelChangeType = "flaw"
	LevelChangeBloodPotency   LevelChangeType = "blood_potency"
	LevelChangeBloodRitual    LevelChangeType = "blood_ritual"
	LevelChangeSpecialization LevelChangeType = "specialization"
	LevelChangeUnknown        LevelChangeType = "unknown"
)

type V5Character struct {
	BaseSheet   `gorm:"embedded"`
	UpdatedAt   time.Time           `gorm:"type:timestamp;not null;default:current_timestamp"`
	CreatedAt   time.Time           `gorm:"type:timestamp;not null;default:current_timestamp"`
	DirectoryID *uuid.UUID          `gorm:"type:uuid;index"`
	Directory   *CharacterDirectory `gorm:"foreignKey:DirectoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ChroniclePrinciples string `gorm:"type:text;not null;default:''"`
	AnchorsAndBeliefs   string `gorm:"type:text;not null;default:''"`
	Backstory           string `gorm:"type:text;not null;default:''"`

	Sire     string `gorm:"type:text;not null;default:''"`
	Ambition string `gorm:"type:text;not null;default:''"`
	Desire   string `gorm:"type:text;not null;default:''"`

	ClanID         *uuid.UUID      `gorm:"type:uuid;index"`
	Clan           *V5Clan         `gorm:"foreignKey:ClanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	PredatorTypeID *uuid.UUID      `gorm:"type:uuid;index"`
	PredatorType   *V5PredatorType `gorm:"foreignKey:PredatorTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	GenerationEra GenerationEra `gorm:"type:text;not null;default:'children'"`
	Generation    int           `gorm:"type:int;not null;default:0"`

	Hunger       int         `gorm:"type:int;not null;default:0"`
	Humanity     int         `gorm:"type:int;not null;default:7"`
	Stains       int         `gorm:"type:int;not null;default:0"`
	Resonance    V5Resonance `gorm:"type:text;not null;default:''"`
	BloodPotency int         `gorm:"type:int;not null;default:0"`

	Health          int            `gorm:"type:int;not null;default:0"`
	HealthDamage    pq.StringArray `gorm:"type:text[];not null;default:'{}'"`
	Willpower       int            `gorm:"type:int;not null;default:0"`
	WillpowerDamage pq.StringArray `gorm:"type:text[];not null;default:'{}'"`

	UseAdavancedDisciplines  bool `gorm:"type:boolean;not null;default:false"`
	AllowLearningOfAllPowers bool `gorm:"type:boolean;not null;default:false"`
	FullCustomization        bool `gorm:"type:boolean;not null;default:false"`
	Version                  int  `gorm:"type:int;not null;default:2"`

	SkillSpreadType SkillSpreadType `gorm:"type:varchar(32);not null;default:'balanced'"`

	InternalData datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`

	Books []V5Book `gorm:"many2many:v5_character_books;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Categories           []V5CharacterCategory            `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Attributes           []V5CharacterAttribute           `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Skills               []V5CharacterSkill               `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RequiredPointSpreads []V5RequiredPointSpread          `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DisciplineSelections []V5CharacterDisciplineSelection `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TraitPackUsages      []V5CharacterTraitPackUsage      `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BloodRituals         []V5BloodRitual                  `gorm:"many2many:v5_character_blood_rituals;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OblivionCeremonies   []V5OblivionCeremony             `gorm:"many2many:v5_character_oblivion_ceremonies;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LevelHistory         []V5LevelChange                  `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Viewers []User `gorm:"many2many:v5_character_viewers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5CharacterCategory struct {
	ID          uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID uuid.UUID   `gorm:"type:uuid;not null;index"`
	Character   V5Character `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        CategoryKey `gorm:"type:text;not null"`
}

type V5CharacterAttribute struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Character   V5Character  `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category    CategoryKey  `gorm:"type:text;not null"`
	Key         AttributeKey `gorm:"type:text;not null"`
	Value       int          `gorm:"type:int;not null;default:0"`
}

type V5CharacterSkill struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Character      V5Character    `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category       CategoryKey    `gorm:"type:text;not null"`
	Key            SkillKey       `gorm:"type:text;not null"`
	Value          int            `gorm:"type:int;not null;default:0"`
	Specialization pq.StringArray `gorm:"type:text[];not null;default:'{}'"`
}

type V5RequiredPointSpread struct {
	ID          uuid.UUID                `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID uuid.UUID                `gorm:"type:uuid;not null;index"`
	Character   V5Character              `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type        V5CharacterTraitPackKind `gorm:"type:text;not null"`
	IsFlaw      bool                     `gorm:"type:boolean;not null;default:false"`
	Points      int                      `gorm:"type:int;not null;default:0"`
	PackID      *uuid.UUID               `gorm:"type:uuid;index"`
	Pack        *V5TraitPack             `gorm:"foreignKey:PackID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type V5CharacterDisciplineSelection struct {
	ID           uuid.UUID                      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID  uuid.UUID                      `gorm:"type:uuid;not null;index"`
	Character    V5Character                    `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DisciplineID uuid.UUID                      `gorm:"type:uuid;not null;index"`
	Discipline   V5Discipline                   `gorm:"foreignKey:DisciplineID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Points       int                            `gorm:"type:int;not null;default:0"`
	CurrentLevel int                            `gorm:"type:int;not null;default:0"`
	Abilities    []V5CharacterDisciplineAbility `gorm:"foreignKey:SelectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5CharacterDisciplineAbility struct {
	ID          uuid.UUID                      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SelectionID uuid.UUID                      `gorm:"type:uuid;not null;index"`
	Selection   V5CharacterDisciplineSelection `gorm:"foreignKey:SelectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AbilityID   uuid.UUID                      `gorm:"type:uuid;not null;index"`
	Ability     V5DisciplineAbility            `gorm:"foreignKey:AbilityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Level       int                            `gorm:"type:int;not null;default:0"`
	UsedLevel   int                            `gorm:"type:int;not null;default:0"`
}

type V5CharacterTraitPackUsage struct {
	ID          uuid.UUID                `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID uuid.UUID                `gorm:"type:uuid;not null;index"`
	Character   V5Character              `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Kind        V5CharacterTraitPackKind `gorm:"type:text;not null"`
	PackID      uuid.UUID                `gorm:"type:uuid;not null;index"`
	Pack        V5TraitPack              `gorm:"foreignKey:PackID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Traits      []V5CharacterTrait       `gorm:"foreignKey:UsageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FlawTraits  []V5CharacterFlawTrait   `gorm:"foreignKey:UsageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5CharacterTrait struct {
	ID          uuid.UUID                 `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UsageID     uuid.UUID                 `gorm:"type:uuid;not null;index"`
	Usage       V5CharacterTraitPackUsage `gorm:"foreignKey:UsageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TraitID     uuid.UUID                 `gorm:"type:uuid;not null;index"`
	Trait       V5Trait                   `gorm:"foreignKey:TraitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	IsLocked    bool                      `gorm:"type:boolean;not null;default:false"`
	IsManual    bool                      `gorm:"type:boolean;not null;default:false"`
	CustomLevel *int                      `gorm:"type:int"`
	Suffix      *string                   `gorm:"type:text"`
}

type V5CharacterFlawTrait struct {
	ID          uuid.UUID                 `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UsageID     uuid.UUID                 `gorm:"type:uuid;not null;index"`
	Usage       V5CharacterTraitPackUsage `gorm:"foreignKey:UsageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TraitID     uuid.UUID                 `gorm:"type:uuid;not null;index"`
	Trait       V5Trait                   `gorm:"foreignKey:TraitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	IsLocked    bool                      `gorm:"type:boolean;not null;default:false"`
	IsManual    bool                      `gorm:"type:boolean;not null;default:false"`
	CustomLevel *int                      `gorm:"type:int"`
	Suffix      *string                   `gorm:"type:text"`
}

type V5CharacterInventory struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex"`
	Character   V5Character   `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Bank        int           `gorm:"type:int;not null;default:0"`
	Cash        int           `gorm:"type:int;not null;default:0"`
	Stacks      []V5ItemStack `gorm:"foreignKey:InventoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type V5ItemStack struct {
	ID          uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InventoryID uuid.UUID              `gorm:"type:uuid;not null;index"`
	Inventory   V5CharacterInventory   `gorm:"foreignKey:InventoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Location    InventoryStackLocation `gorm:"type:text;not null"`
	ItemID      uuid.UUID              `gorm:"type:uuid;not null;index"`
	Item        V5Item                 `gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Amount      int                    `gorm:"type:int;not null;default:1"`
}

type V5Item struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	IsCustom    bool      `gorm:"type:boolean;not null;default:false"`
	Name        string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text;not null;default:''"`
	Category    string    `gorm:"type:text;not null;default:''"`
}

type V5LevelChange struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CharacterID uuid.UUID       `gorm:"type:uuid;not null;index"`
	Character   V5Character     `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type        LevelChangeType `gorm:"type:text;not null"`
	Date        string          `gorm:"type:text;not null"`
	Text        string          `gorm:"type:text;not null"`
	ExpUsed     int             `gorm:"type:int;not null;default:0"`
	ExpBefore   int             `gorm:"type:int;not null;default:0"`
	ExpAfter    int             `gorm:"type:int;not null;default:0"`
}
