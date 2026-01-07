package entities

import (
	"github.com/google/uuid"
)

type Sex string

const (
	SexMale   Sex = "m"
	SexFemale Sex = "f"
	SexDivers Sex = "d"
)

type DamageType string

const (
	DamageNone        DamageType = "none"
	DamageSuperficial DamageType = "superficial"
	DamageHeavy       DamageType = "heavy"
	DamageFull        DamageType = "full"
)

type CategoryKey string

const (
	CategoryPhysical CategoryKey = "physical"
	CategorySocial   CategoryKey = "social"
	CategoryMental   CategoryKey = "mental"
)

type AttributeKey string

const (
	AttrStrength     AttributeKey = "str"
	AttrDexterity    AttributeKey = "dex"
	AttrStamina      AttributeKey = "sta"
	AttrCharisma     AttributeKey = "cha"
	AttrManipulation AttributeKey = "man"
	AttrComposure    AttributeKey = "com"
	AttrIntelligence AttributeKey = "int"
	AttrWits         AttributeKey = "wit"
	AttrResolve      AttributeKey = "res"
)

type SkillKey string

const (
	SkillAthletics     SkillKey = "ath"
	SkillBrawl         SkillKey = "bra"
	SkillCraft         SkillKey = "cra"
	SkillDrive         SkillKey = "dri"
	SkillFirearms      SkillKey = "fir"
	SkillMelee         SkillKey = "mel"
	SkillLarceny       SkillKey = "lar"
	SkillStealth       SkillKey = "ste"
	SkillSurvival      SkillKey = "sur"
	SkillAnimalKen     SkillKey = "ken"
	SkillEtiquette     SkillKey = "eti"
	SkillInsight       SkillKey = "ins"
	SkillIntimidation  SkillKey = "int"
	SkillLeadership    SkillKey = "lea"
	SkillPerformance   SkillKey = "per"
	SkillPersuasion    SkillKey = "pes"
	SkillStreetwise    SkillKey = "stw"
	SkillSubterfuge    SkillKey = "sub"
	SkillAcademics     SkillKey = "aca"
	SkillAwareness     SkillKey = "awa"
	SkillFinance       SkillKey = "fin"
	SkillInvestigation SkillKey = "inv"
	SkillMedicine      SkillKey = "med"
	SkillOccult        SkillKey = "occ"
	SkillPolitics      SkillKey = "pol"
	SkillScience       SkillKey = "sci"
	SkillTechnology    SkillKey = "tec"
)

type InventoryStackLocation string

const (
	InventoryCarried InventoryStackLocation = "carried"
	InventoryOwned   InventoryStackLocation = "owned"
)

type SkillSpreadType string

const (
	SkillSpreadTypeBalanced        SkillSpreadType = "balanced"
	SkillSpreadTypeSpecialist      SkillSpreadType = "specialist"
	SkillSpreadTypeJackOfAllTrades SkillSpreadType = "jack_of_all_trades"
)

type BaseSheet struct {
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	User       *User     `gorm:"foreignKey:UserID"`
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OldVicarID uint      `gorm:"type:bigint"`
	Avatar     string    `gorm:"type:text;not null;default:''"`
	Name       string    `gorm:"type:text;not null;default:''"`
	Notes      string    `gorm:"type:text;not null;default:''"`
	Sex        Sex       `gorm:"type:text;not null;default:'d'"`
	Concept    string    `gorm:"type:text;not null;default:''"`
	Chronicle  string    `gorm:"type:text;not null;default:''"`
	Exp        int       `gorm:"type:int;not null;default:0"`
	UsedExp    int       `gorm:"type:int;not null;default:0"`
}

type CharacterDirectory struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string        `gorm:"type:text;not null"`
	Open         bool          `gorm:"type:boolean;not null;default:false"`
	V5Characters []V5Character `gorm:"foreignKey:DirectoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
