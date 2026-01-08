import type {User} from "@/@types/user.ts";

export type GenerationEra = 'children' | 'newborn' | 'ancillae' | 'older' | 'elder' | 'cainesinheritance'
export type Resonance = 'choleric' | 'melancholic' | 'phlegmatic' | 'sanguine' | 'animalblood' | ''
export type CategoryKey = 'physical' | 'social' | 'mental'
export type AttributeKey = 'str' | 'dex' | 'sta' | 'cha' | 'man' | 'com' | 'int' | 'wit' | 'res'
export type SkillKey = 'ath' | 'bra' | 'cra' | 'dri' | 'fir' | 'mel' | 'lar' | 'ste' | 'sur' |
  'ani' | 'eti' | 'ins' | 'int' | 'lea' | 'per' | 'prf' | 'sub' | 'str' |
  'aca' | 'awa' | 'fin' | 'inv' | 'med' | 'occ' | 'pol' | 'sci' | 'tec'
export type SkillSpreadType = 'balanced' | 'specialist' | 'jack_of_all_trades'
export type TraitPackType = 'merits' | 'backgrounds'
export type Sex = 'm' | 'f' | 'd'

export interface V5Book {
  id: string
  oldVicarID?: number
  name: string
  isOfficial: boolean
}

export interface V5Discipline {
  id: string
  oldVicarID?: number
  name: string
  summary?: string
  note?: string
  isHomebrew: boolean
  creator?: string
  abilities: V5DisciplineAbility[]
}

export interface V5DisciplineAbility {
  id: string
  disciplineID: string
  level: number
  name: string
  combinationRefID?: string
  combinationLevel?: number
  requirement?: number
  minBloodPotency?: number
  summary: string
  costs: string
  diceSupplies?: string
  system: string
  alternatives?: string[]
  duration: string
}

export interface V5Clan {
  id: string
  oldVicarID?: number
  bookID: string
  name: string
  slogan: string
  description: string
  curse: string
  symbol?: string
  isHomebrew: boolean
  creator?: string
  disciplines: V5Discipline[]
}

export interface V5PredatorAction {
  id: string
  predatorTypeID: string
  description: string
  type: 'additional_specialization' | 'discipline_point' | 'humanity_change' | 'add_merit' |
    'add_background' | 'add_background_points' | 'add_flaw' | 'blood_potency_change' |
    'attribute_change' | 'skill_change' | 'spend_background_points_between' | 'spend_flaw_points_between'
  data: any
  restriction?: any
}

export interface V5PredatorType {
  id: string
  bookID: string
  name: string
  description: string
  restriction?: any
  actions: V5PredatorAction[]
}

export interface V5CharacterAttribute {
  id: string
  category: CategoryKey
  key: AttributeKey
  value: number
}

export interface V5CharacterSkill {
  id: string
  characterID: string
  category: CategoryKey
  key: SkillKey
  value: number
  specialization: string[]
}

export interface V5CharacterDisciplineAbility {
  id: string
  selectionID: string
  abilityID: string
  level: number
  usedLevel: number
}

export interface V5CharacterDisciplineSelection {
  id: string
  characterID: string
  disciplineID: string
  points: number
  currentLevel: number
  abilities: V5CharacterDisciplineAbility[]
}

export interface V5CharacterTrait {
  id: string
  usageID: string
  traitID: string
  isLocked: boolean
  isManual: boolean
  customLevel?: number
  suffix?: string
}

export interface V5CharacterFlawTrait {
  id: string
  usageID: string
  traitID: string
  isLocked?: boolean
  isManual?: boolean
  customLevel?: number
  suffix?: string
}

export interface V5CharacterTraitPackUsage {
  id: string
  characterID: string
  kind: 'merits' | 'backgrounds'
  packID: string
  traits: V5CharacterTrait[]
  flawTraits: V5CharacterFlawTrait[]
}

export interface V5BloodRitual {
  id: string
  bookID: string
  level: number
  name: string
  description: string
  ingredients: string
  execution: string
  system: string
}

export interface V5OblivionCeremony {
  id: string
  bookID: string
  level: number
  name: string
  cost: string
  roll: string
  summary: string
  requires?: number
  cult?: string
  ingredients: string
  execution: string
  system: string
  duration?: string
}

export interface V5LevelChange {
  id: string
  characterID: string
  type: 'attribute' | 'skill' | 'discipline' | 'trait' | 'flaw' | 'blood_potency' | 'specialization' | 'unknown'
  date: string
  text: string
  expUsed: number
  expBefore: number
  expAfter: number
}

export interface V5CharacterCategory {
  id: string
  name: string
}

// Internal data for tracking character creation choices
export interface V5CharacterInternalData {
  // Predator type action choices - keyed by action ID
  predatorActionChoices?: Record<string, {
    skillKey?: string
    specialization?: string
    disciplineID?: string
    // For spend_background/flaw_points_between
    packDistribution?: Record<string, number> // packId -> points spent
  }>
  // Predator type bonuses that were applied (for reverting)
  predatorBonusesApplied?: {
    humanityChange?: number
    bloodPotencyChange?: number
    specializations?: { skillKey: string; specialization: string }[]
    disciplinePoints?: { disciplineID: string; points: number }[]
    traits?: { packID: string; traitID: string; level: number; suffix?: string; isFlaw?: boolean }[]
    spendPointsBetween?: {
      actionId: string
      type: 'backgrounds' | 'flaws'
      requiredPoints: number
      allowedPackOldVicarIDs: number[]
      distribution: Record<string, number> // packId -> points spent from predator bonus
    }[]
  }
  // Discipline distribution selection
  disciplineDistribution?: {
    primaryDisciplineID?: string  // Gets 2 points
    secondaryDisciplineID?: string  // Gets 1 point
    predatorBonusDisciplineID?: string  // Gets predator bonus points
  }
  // Free specializations chosen in skills step
  freeSpecializations?: { skillKey: string; specialization: string }[]
}

export interface V5Character {
  userID: string
  id: string
  isOwner?: boolean
  avatar?: string
  name: string
  notes: string
  sex: Sex
  concept: string
  chronicle: string
  exp: number
  usedExp: number
  directoryID?: string
  chroniclePrinciples: string
  anchorsAndBeliefs: string
  backstory: string
  sire: string
  ambition: string
  desire: string
  clanID?: string
  predatorTypeID?: string
  generationEra: GenerationEra
  generation: number
  hunger: number
  humanity: number
  stains: number
  resonance: Resonance
  bloodPotency: number
  health: number
  healthDamage: string[]
  willpower: number
  willpowerDamage: string[]
  useAdvancedDisciplines: boolean
  allowLearningOfAllPowers: boolean
  fullCustomization: boolean
  version: number
  skillSpreadType: SkillSpreadType
  internalData: V5CharacterInternalData
  inventory: V5Inventory
  books: V5Book[]
  categories: V5CharacterCategory[]
  attributes: V5CharacterAttribute[]
  skills: V5CharacterSkill[]
  disciplineSelections: V5CharacterDisciplineSelection[]
  traitPackUsages: V5CharacterTraitPackUsage[]
  bloodRituals: V5BloodRitual[]
  oblivionCeremonies: V5OblivionCeremony[]
  levelHistory: V5LevelChange[]
  viewers: User[]
}

export type InventorySide = "carriedItems" | "ownedItems"

export type V5InventoryItem = {
  id: string
  name: string
  description: string
  amount: number
  category?: string
}

export type V5Inventory = {
  cash: number
  bank: number
  carriedItems: V5InventoryItem[]
  ownedItems: V5InventoryItem[]
}

export interface V5Trait {
  id: string
  oldVicarID?: number
  level: number
  name: string
  description: string
  isRepeatable: boolean
  repeatAmount?: number
  repeatSize?: number
  requirement?: any
  actions?: any
  restriction?: any
  isFlaw: boolean
}

export interface V5TraitPackTrait {
  id: string
  packID: string
  trait: V5Trait
}

export interface V5TraitPack {
  id: string
  oldVicarID?: number
  bookID: string
  type: TraitPackType
  name: string
  description: string
  specialRules: 'none' | 'allies' | 'haven' | 'mask'
  restriction?: any
  packTraits: V5TraitPackTrait[]
}

export interface V5CharacterCreateRequest {
  name: string
  books: string[]
  generationEra: GenerationEra
  generation: number
}

export interface V5CharacterUpdateRequest {
  name?: string
  avatar?: string
  notes?: string
  sex?: Sex
  concept?: string
  chronicle?: string
  chroniclePrinciples?: string
  anchorsAndBeliefs?: string
  backstory?: string
  sire?: string
  ambition?: string
  desire?: string
  clanID?: string
  predatorTypeID?: string
  generationEra?: GenerationEra
  generation?: number
  hunger?: number
  humanity?: number
  stains?: number
  resonance?: Resonance
  bloodPotency?: number
  health?: number
  healthDamage?: string[]
  willpower?: number
  willpowerDamage?: string[]
  useAdvancedDisciplines?: boolean
  allowLearningOfAllPowers?: boolean
  fullCustomization?: boolean
  skillSpreadType?: SkillSpreadType
}

export interface V5Directory {
  id: string
  userID: string
  name: string
  open: boolean
}
