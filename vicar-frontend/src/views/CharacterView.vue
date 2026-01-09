<script setup lang="ts">
import {onMounted, computed, ref, onUnmounted} from "vue"
import { useRoute, useRouter } from "vue-router"
import { useCharactersStore } from "@/stores/characters"
import { useV5DataStore } from "@/stores/v5data"
import VButton from "@/components/ui/VButton.vue"
import VCard from "@/components/ui/VCard.vue"
import VInput from "@/components/ui/VInput.vue"
import VTextarea from "@/components/ui/VTextarea.vue"
import VDotRating from "@/components/characters/VDotRating.vue"
import VTracker from "@/components/characters/VTracker.vue"
import VHumanityTracker from "@/components/characters/VHumanityTracker.vue"
import VHungerTracker from "@/components/characters/VHungerTracker.vue"
import type {
  AttributeKey,
  SkillKey,
  CategoryKey,
  Resonance,
  V5DisciplineAbility,
  V5BloodRitual,
  V5OblivionCeremony
} from "@/@types/v5"
import VCharacterReadyAlert from "@/components/characters/VCharacterReadyAlert.vue";
import VExperienceBox from "@/components/characters/VExperienceBox.vue";
import VLevelHistoryModal from "@/components/characters/VLevelHistoryModal.vue";
import VCharacterViewersModal from "@/components/characters/VCharacterViewersModal.vue";
import VInventoryBox from "@/components/characters/VInventoryBox.vue";
import {resetTitle, useTitle} from "@/composables/useTitle.ts";
import VDisciplineAbilityInfoModal from "@/components/characters/VDisciplineAbilityInfoModal.vue";
import VBloodRitualInfoModal from "@/components/characters/VBloodRitualInfoModal.vue";
import VOblivionCeremonyInfoModal from "@/components/characters/VOblivionCeremonyInfoModal.vue";

const route = useRoute()
const router = useRouter()
const store = useCharactersStore()
const v5data = useV5DataStore()

const isLoading = ref(true)
const isSaving = ref(false)

const showHistory = ref(false)
const showViewers = ref(false)

const abilityModalOpen = ref(false)
const selectedAbility = ref<V5DisciplineAbility | null>(null)

const bloodRitualModalOpen = ref(false)
const selectedBloodRitual = ref<V5BloodRitual | null>(null)

const oblivionCeremonyModalOpen = ref(false)
const selectedOblivionCeremony = ref<V5OblivionCeremony | null>(null)

const characterId = computed(() => route.params.id as string)
const character = computed({
  get: () => store.currentCharacter,
  set: (val) => {
    store.currentCharacter = val
  }
})

const attributeLabels: Record<AttributeKey, string> = {
  str: "Stärke",
  dex: "Geschicklichkeit",
  sta: "Ausdauer",
  cha: "Charisma",
  man: "Manipulation",
  com: "Fassung",
  int: "Intelligenz",
  wit: "Witz",
  res: "Entschlossenheit"
}

const skillLabels: Record<SkillKey, string> = {
  ath: "Athletik",
  bra: "Handgemenge",
  cra: "Handwerk",
  dri: "Fahren",
  fir: "Schusswaffen",
  mel: "Nahkampf",
  lar: "Heimlichkeit",
  ste: "Überleben",
  sur: "Überlebensinstinkt",
  ani: "Tierkunde",
  eti: "Etikette",
  ins: "Einblick",
  int: "Einschüchtern",
  lea: "Anführen",
  per: "Überzeugen",
  prf: "Darbietung",
  sub: "Gassenwissen",
  str: "Täuschen",
  aca: "Akademik",
  awa: "Aufmerksamkeit",
  fin: "Finanzen",
  inv: "Untersuchung",
  med: "Medizin",
  occ: "Okkultismus",
  pol: "Politik",
  sci: "Wissenschaft",
  tec: "Technologie"
}

const categoryLabels: Record<CategoryKey, string> = {
  physical: "Körperlich",
  social: "Sozial",
  mental: "Geistig"
}

const generationEraLabels: Record<string, string> = {
  children: "Kinder",
  newborn: "Neugeborene",
  ancillae: "Ancillae",
  older: "Ältere",
  elder: "Älteste",
  cainesinheritance: "Kains Erbe"
}

const resonanceLabels: Record<Resonance, string> = {
  choleric: "Cholerisch",
  melancholic: "Melancholisch",
  phlegmatic: "Phlegmatisch",
  sanguine: "Sanguinisch",
  animalblood: "Tierblut",
  "": "Leer"
}

const resonanceOptions: Resonance[] = ["", "choleric", "melancholic", "phlegmatic", "sanguine", "animalblood"]

const clan = computed(() => {
  if (!character.value?.clanID) return null
  return v5data.getClanById(character.value.clanID)
})

const predatorType = computed(() => {
  if (!character.value?.predatorTypeID) return null
  return v5data.getPredatorTypeById(character.value.predatorTypeID)
})

const predatorBonusTraitsForDisplay = computed(() => {
  const bonuses = getPredatorBonusTraits()

  console.log(bonuses)

  return bonuses.filter(b => !characterHasTrait(b.packID, b.traitID, b.isFlaw))
})

const traitPackUsagesForDisplay = computed(() => {
  const usages = character.value?.traitPackUsages ? [...character.value.traitPackUsages] : []

  const packsNeeded = new Set(predatorBonusTraitsForDisplay.value.map(b => b.packID))
  for (const packID of packsNeeded) {
    if (!usages.some(u => u.packID === packID)) {
      usages.push({
        id: `usage-display-${packID}`,
        characterID: characterId.value,
        kind: (v5data.traitPacks?.find(p => p.id === packID)?.type as any) ?? "merits",
        packID,
        traits: [],
        flawTraits: []
      } as any)
    }
  }

  return usages
})

const calculatedHealth = computed(() => {
  if (!character.value) return 0
  const stamina = character.value.attributes?.find((a) => a.key === "sta")?.value ?? 0
  let fortitudeBonus = 0

  const fortitudeDisc = character.value.disciplineSelections?.find((d) => {
    const disc = getDisciplineData(d.disciplineID)
    return disc?.name?.toLowerCase().includes("seelenstärke") || disc?.name?.toLowerCase().includes("fortitude")
  })

  if (fortitudeDisc) {
    const hasResilience = fortitudeDisc.abilities?.some((a) => {
      const ability = getAbilityData(fortitudeDisc.disciplineID, a.abilityID)
      return ability?.name?.toLowerCase().includes("resilienz") || ability?.name?.toLowerCase().includes("resilience")
    })
    if (hasResilience) fortitudeBonus = fortitudeDisc.currentLevel
  }

  return stamina + 3 + fortitudeBonus
})

const calculatedWillpower = computed(() => {
  if (!character.value) return 0
  const composure = character.value.attributes?.find((a) => a.key === "com")?.value ?? 0
  const resolve = character.value.attributes?.find((a) => a.key === "res")?.value ?? 0
  return composure + resolve
})

const willpowerRegenAmount = computed(() => {
  if (!character.value) return 0
  const composure = character.value.attributes?.find((a) => a.key === "com")?.value ?? 0
  const resolve = character.value.attributes?.find((a) => a.key === "res")?.value ?? 0
  return Math.max(composure, resolve)
})

const hasBloodSorcery = computed(() => {
  if (!character.value) return false
  return (
    character.value.disciplineSelections?.some((d) => {
      const disc = getDisciplineData(d.disciplineID)
      return (
        disc?.name?.toLowerCase().includes("blutmagie") ||
        disc?.name?.toLowerCase().includes("blood sorcery") ||
        disc?.name?.toLowerCase().includes("blutzauberei")
      )
    }) ?? false
  )
})

const bloodSorceryLevel = computed(() => {
  if (!character.value) return 0
  const selection = character.value.disciplineSelections?.find((d) => {
    const disc = getDisciplineData(d.disciplineID)
    return (
      disc?.name?.toLowerCase().includes("blutmagie") ||
      disc?.name?.toLowerCase().includes("blood sorcery") ||
      disc?.name?.toLowerCase().includes("blutzauberei")
    )
  })
  return selection?.currentLevel ?? 0
})

const hasOblivion = computed(() => {
  if (!character.value) return false
  return (
    character.value.disciplineSelections?.some((d) => {
      const disc = getDisciplineData(d.disciplineID)
      return disc?.name?.toLowerCase().includes("vergessenheit") || disc?.name?.toLowerCase().includes("oblivion")
    }) ?? false
  )
})

const oblivionLevel = computed(() => {
  if (!character.value) return 0
  const selection = character.value.disciplineSelections?.find((d) => {
    const disc = getDisciplineData(d.disciplineID)
    return disc?.name?.toLowerCase().includes("vergessenheit") || disc?.name?.toLowerCase().includes("oblivion")
  })
  return selection?.currentLevel ?? 0
})

const availableBloodRituals = computed(() => {
  if (!v5data.bloodRituals) return []
  return v5data.bloodRituals.filter((r) => r.level <= bloodSorceryLevel.value)
})

const availableOblivionCeremonies = computed(() => {
  if (!v5data.oblivionCeremonies) return []
  return v5data.oblivionCeremonies.filter((c) => c.level <= oblivionLevel.value)
})

onMounted(async () => {
  await Promise.all([
    v5data.fetchClans(),
    v5data.fetchPredatorTypes(),
    v5data.fetchTraitPacks(),
    v5data.fetchDisciplines(),
    v5data.fetchBloodRituals(),
    v5data.fetchOblivionCeremonies()
  ])

  if (characterId.value) {
    await store.fetchCharacter(characterId.value)
  }

  if (character.value) {
    useTitle(character.value.name)
  }
  isLoading.value = false
})

onUnmounted(() => {
  resetTitle()
})

function getDisciplineData(disciplineId: string) {
  if (v5data.disciplines) {
    const disc = v5data.disciplines.find((d) => d.id === disciplineId)
    if (disc) return disc
  }
  if (!character.value?.clanID) return null
  const selectedClan = v5data.getClanById(character.value.clanID)
  return selectedClan?.disciplines?.find((d) => d.id === disciplineId) ?? null
}

function getAbilityData(disciplineId: string, abilityId: string): V5DisciplineAbility | null {
  const disc = getDisciplineData(disciplineId)
  return disc?.abilities?.find((a) => a.id === abilityId) ?? null
}

function getDisciplineName(disciplineId: string): string {
  return getDisciplineData(disciplineId)?.name ?? "Unbekannt"
}

function getSelectedAbilities(disciplineId: string, abilities: { abilityID: string; level: number }[]) {
  const disc = getDisciplineData(disciplineId)
  if (!disc) return []
  return abilities
    .map((a) => {
      const ability = disc.abilities?.find((ab) => ab.id === a.abilityID)
      return ability ? { ...ability, usedLevel: a.level } : null
    })
    .filter(Boolean) as (V5DisciplineAbility & { usedLevel: number })[]
}

let saveTimeout: ReturnType<typeof setTimeout> | null = null

async function saveField(field: string, value: any) {
  if (!characterId.value) return

  if (saveTimeout) clearTimeout(saveTimeout)

  saveTimeout = setTimeout(async () => {
    isSaving.value = true
    try {
      await store.updateCharacter(characterId.value, { [field]: value })
    } finally {
      isSaving.value = false
    }
  }, 450)
}

function goToEditor() {
  router.push(`/characters/${characterId.value}/edit`)
}

function getAttributesByCategory(category: CategoryKey) {
  return character.value?.attributes?.filter((a) => a.category === category) ?? []
}

function getSkillsByCategory(category: CategoryKey) {
  return character.value?.skills?.filter((s) => s.category === category) ?? []
}

function getTraitPackName(packId: string): string {
  const pack = v5data.traitPacks?.find((p) => p.id === packId)
  return pack?.name ?? "Unbekannt"
}

function getTrackPackSummary(packId: string): string {
  const pack = v5data.traitPacks?.find((p) => p.id === packId)
  return pack?.description ?? "";
}

function getTraitInfo(packId: string, traitId: string): { name: string; isFlaw: boolean; description: string } {
  const pack = v5data.traitPacks?.find((p) => p.id === packId)
  const packTrait = pack?.packTraits?.find((pt) => pt.trait.id === traitId)
  return {
    name: packTrait?.trait.name ?? "Unbekannt",
    isFlaw: packTrait?.trait.isFlaw ?? false,
    description: packTrait?.trait.description ?? ""
  }
}

function updateHunger(value: number) {
  if (!character.value) return
  character.value.hunger = value
  saveField("hunger", value)
}

function updateHumanity(value: number) {
  if (!character.value) return
  character.value.humanity = value
  saveField("humanity", value)
}

function updateStains(value: number) {
  if (!character.value) return
  character.value.stains = value
  saveField("stains", value)
}

function updateHealthDamage(value: string[]) {
  if (!character.value) return
  character.value.healthDamage = value
  saveField("healthDamage", value)
}

function updateWillpowerDamage(value: string[]) {
  if (!character.value) return
  character.value.willpowerDamage = value
  saveField("willpowerDamage", value)
}

function updateResonance(value: Resonance) {
  if (!character.value) return
  character.value.resonance = value
  saveField("resonance", value)
}

function updateTextField(field: string, value: string) {
  if (!character.value) return
    ;(character.value as any)[field] = value
  saveField(field, value)
}

function handleAvatarUpload(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return

  const file = input.files[0]
  const reader = new FileReader()
  reader.onload = (e) => {
    const base64 = e.target?.result as string
    if (character.value) {
      character.value.avatar = base64
      saveField("avatar", base64)
    }
  }
  reader.readAsDataURL(file!)
}

function regenerateWillpowerSession() {
  if (!character.value) return
  const heal = Math.max(0, willpowerRegenAmount.value)
  if (!heal) return

  const src = Array.isArray(character.value.willpowerDamage) ? [...character.value.willpowerDamage] : []
  const max = Math.max(0, calculatedWillpower.value)

  const padded = new Array(max).fill("")
  for (let i = 0; i < max; i++) {
    const v = src[i]
    padded[i] = v === "superficial" || v === "aggravated" ? v : ""
  }

  let remaining = heal
  for (let i = 0; i < padded.length && remaining > 0; i++) {
    if (padded[i] === "superficial") {
      padded[i] = ""
      remaining--
    }
  }

  character.value.willpowerDamage = padded
  saveField("willpowerDamage", padded)
}

function saveExperience() {
  if (!character.value) return
  saveField("exp", character.value.exp)
}

function openAbilityInfo(ability: V5DisciplineAbility) {
  selectedAbility.value = ability
  abilityModalOpen.value = true
}

function openBloodRitualInfo(ritual: V5BloodRitual) {
  selectedBloodRitual.value = ritual
  bloodRitualModalOpen.value = true
}

function openOblivionCeremonyInfo(ceremony: V5OblivionCeremony) {
  selectedOblivionCeremony.value = ceremony
  oblivionCeremonyModalOpen.value = true
}

function getPredatorBonusTraits() {
  const predatorId = character.value?.predatorTypeID
  if (!predatorId || !v5data.predatorTypes || !v5data.traitPacks) return []

  const predator = v5data.predatorTypes.find(p => p.id === predatorId)
  if (!predator) return []

  const bonusActions = predator.actions.filter(a =>
    a.type === "add_merit" || a.type === "add_background" || a.type === "add_flaw"
  )

  const bonuses: { packID: string; traitID: string; isFlaw: boolean; level?: number, suffix?: string }[] = []

  for (const action of bonusActions) {
    const data = action.data as any
    let traitOldVicarID: number | null = null
    let flawId: number | null = null
    let flawType: string | null = null
    let suffix: string | undefined = undefined
    if (action.type === "add_background" || action.type === "add_merit") {
      const levelId = data?.levelId as number | null
      const traitId = (action.type === "add_merit" ? data?.meritId : data?.backgroundId) as number | null
      if (!traitId || !levelId) continue
      traitOldVicarID = traitId * 1000 + levelId
    } else {
      const choices = data?.choices as Array<{ id: number; type: 'background'|'merit'; flawId: number; suffix?: string }> | null
      if (choices) {
        if (choices.length === 1) {
          const choice = choices[0]!
          traitOldVicarID = choice.id
          suffix = choice.suffix ?? undefined
          flawId = choice.flawId ?? null
          flawType = choice.type + 's'
        } else {

        }
      }
    }

    if (!traitOldVicarID) continue
    const isFlaw = action.type === "add_flaw"

    for (const pack of v5data.traitPacks) {
      if (isFlaw) {
        if (pack.oldVicarID === traitOldVicarID && pack.type === flawType) {
          const packTrait = pack.packTraits?.find(pt => pt.trait.level === flawId && pt.trait.isFlaw)
          if (!packTrait) continue
          bonuses.push({ packID: pack.id, traitID: packTrait.trait.id, isFlaw, level: flawId!, suffix: suffix })
        }
        continue
      }

      const packTrait = pack.packTraits?.find(pt => pt.trait.oldVicarID === traitOldVicarID)
      if (!packTrait) continue
      const level = data?.level ?? packTrait.trait.level
      bonuses.push({ packID: pack.id, traitID: packTrait.trait.id, isFlaw, level, suffix: suffix })
      break
    }
  }

  return bonuses
}

function characterHasTrait(packID: string, traitID: string, isFlaw: boolean) {
  const usage = character.value?.traitPackUsages?.find(u => u.packID === packID)
  if (!usage) return false
  return isFlaw
    ? usage.flawTraits?.some(t => t.traitID === traitID) ?? false
    : usage.traits?.some(t => t.traitID === traitID) ?? false
}

</script>

<template>
  <div v-if="isLoading" class="character-view character-view--loading">
    <p>Lade Charakter…</p>
  </div>

  <div class="character-view" v-else-if="!!character">
    <VDisciplineAbilityInfoModal v-model="abilityModalOpen" :ability="selectedAbility" />
    <VBloodRitualInfoModal v-model="bloodRitualModalOpen" :ritual="selectedBloodRitual" />
    <VOblivionCeremonyInfoModal v-model="oblivionCeremonyModalOpen" :ceremony="selectedOblivionCeremony" />

    <div class="character-view__header">
      <div class="header-left">
        <VButton variant="ghost" @click="router.push('/')">← Zurück</VButton>
        <h1>{{ character.name }}</h1>
        <p class="subtitle" v-if="clan">{{ clan.slogan }}</p>
      </div>

      <div class="header-right">
        <div style="display: flex; gap: 0.5rem">
          <VButton variant="secondary" @click="showViewers = true">Viewer verwalten</VButton>
          <VButton variant="primary" @click="goToEditor">Bearbeiten</VButton>
        </div>

        <div class="avatar-container">
          <label class="avatar-upload">
            <img v-if="character.avatar" :src="character.avatar" :alt="character.name" class="avatar-image" />
            <div v-else class="avatar-placeholder">
              <span>{{ character.name.charAt(0).toUpperCase() }}</span>
            </div>

            <input type="file" accept="image/*" class="avatar-input" @change="handleAvatarUpload" />
            <div class="avatar-overlay">
              <span>Ändern</span>
            </div>
          </label>
        </div>
      </div>
    </div>

    <VCharacterReadyAlert :character="character" style="margin-bottom: 1rem"/>

    <div class="character-view__content">
      <div class="trackers-section">

        <VCard class="tracker-card">
          <div class="tracker-card__head">
            <h2>Gesundheit</h2>
          </div>
          <div class="tracker-info">
            <span class="tracker-formula">
              Ausdauer ({{ character.attributes?.find((a) => a.key === "sta")?.value ?? 0 }}) + 3
            </span>
          </div>
          <VTracker :max="calculatedHealth" :damage="character.healthDamage ?? []" @update:damage="updateHealthDamage" />
        </VCard>

        <VCard class="tracker-card">
          <div class="tracker-card__head">
            <h2>Willenskraft</h2>
          </div>

          <div class="tracker-info tracker-info--row">
            <span class="tracker-formula">
              Fassung ({{ character.attributes?.find((a) => a.key === "com")?.value ?? 0 }}) + Entschlossenheit
              ({{ character.attributes?.find((a) => a.key === "res")?.value ?? 0 }})
            </span>
          </div>

          <VTracker
            :max="calculatedWillpower"
            :damage="character.willpowerDamage ?? []"
            @update:damage="updateWillpowerDamage"
          />

          <div style="display: flex; justify-content: center; margin-top: 1rem">
            <VButton variant="ghost" @click="regenerateWillpowerSession">
              Sitzungsstart: {{ willpowerRegenAmount }} oberflächlich heilen
            </VButton>
          </div>
        </VCard>

        <VCard class="tracker-card">
          <div class="tracker-card__head">
            <h2>Hunger & Resonanz</h2>
          </div>
          <VHungerTracker :model-value="character.hunger" @update:model-value="updateHunger" />


          <div class="form-field" style="margin-top: 1rem">
            <label class="field-label">Resonanz</label>
            <div class="resonance-select" role="group" aria-label="Resonanz">
              <button
                v-for="res in resonanceOptions"
                :key="res"
                type="button"
                class="resonance-option"
                :class="{ 'resonance-option--active': character.resonance === res }"
                @click="updateResonance(res)"
              >
                {{ resonanceLabels[res] }}
              </button>
            </div>
          </div>
        </VCard>

        <VCard class="tracker-card">
          <div class="tracker-card__head">
            <h2>Menschlichkeit</h2>
          </div>
          <VHumanityTracker
            :humanity="character.humanity"
            :stains="character.stains"
            @update:humanity="updateHumanity"
            @update:stains="updateStains"
          />
        </VCard>
      </div>

      <div class="content-grid">
        <VCard>
          <h2>Grundinformationen</h2>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">Clan:</span>
              <span>{{ clan?.name || "—" }}</span>
            </div>
            <div class="info-item">
              <span class="label">Jagdverhalten:</span>
              <span>{{ predatorType?.name || "—" }}</span>
            </div>
            <div class="info-item">
              <span class="label">Generation:</span>
              <span>{{ character.generation }}. ({{ generationEraLabels[character.generationEra] || character.generationEra }})</span>
            </div>
            <div class="info-item" v-if="character.sire">
              <span class="label">Sire:</span>
              <span>{{ character.sire }}</span>
            </div>
            <div class="info-item">
              <span class="label">Blutmacht:</span>
              <span>{{ character.bloodPotency }}</span>
            </div>
            <div class="info-item">
              <span class="label">Konzept:</span>
              <span>
                <VInput
                  :model-value="character.concept"
                  @update:model-value="updateTextField('concept', $event as string)"
                  placeholder="Das Konzept deines Charakters…"
                  style="width: 100%"
                />
              </span>
            </div>
          </div>
        </VCard>

        <VCard>
          <h2>Leveln</h2>
          <VButton variant="secondary" style="margin-bottom: 1rem" @click="showHistory = true">Level Verlauf</VButton>

          <VExperienceBox v-model="character" @update-xp="saveExperience" style="flex: 1"/>
        </VCard>
      </div>

      <VCard v-if="character.attributes && character.attributes.length > 0">
        <h2>Attribute</h2>
        <div class="attributes-grid">
          <div v-for="category in (['physical', 'social', 'mental'] as CategoryKey[])" :key="category" class="attr-category">
            <h3>{{ categoryLabels[category] }}</h3>
            <div v-for="attr in getAttributesByCategory(category)" :key="attr.id" class="attr-row">
              <span class="attr-name">{{ attributeLabels[attr.key] || attr.key }}</span>
              <VDotRating :model-value="attr.value" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.skills && character.skills.length > 0">
        <h2>Fähigkeiten</h2>
        <div class="skills-grid">
          <div v-for="category in (['physical', 'social', 'mental'] as CategoryKey[])" :key="category" class="skill-category">
            <h3>{{ categoryLabels[category] }}</h3>
            <div v-for="skill in getSkillsByCategory(category)" :key="skill.id" class="skill-row">
              <div class="skill-info">
                <span class="skill-name">{{ skillLabels[skill.key] || skill.key }}</span>
                <span v-if="skill.specialization && skill.specialization.length > 0" class="skill-spec">
                  {{ skill.specialization.join(", ") }}
                </span>
              </div>
              <VDotRating :model-value="skill.value" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.disciplineSelections && character.disciplineSelections.length > 0">
        <h2>Disziplinen</h2>
        <div class="disciplines-list">
          <div v-for="disc in character.disciplineSelections" :key="disc.id" class="disc-item">
            <div class="disc-header">
              <div class="disc-title">
                <span class="disc-name">{{ getDisciplineName(disc.disciplineID) }}</span>
              </div>
              <VDotRating :model-value="disc.currentLevel" :readonly="true" />
            </div>

            <div v-if="disc.abilities && disc.abilities.length > 0" class="disc-abilities">
              <div
                v-for="ability in getSelectedAbilities(disc.disciplineID, disc.abilities)"
                :key="ability.id"
                type="button"
                class="ability-item"
              >
                <div class="ability-header">
                  <span class="ability-level">{{ ability.level }}</span>
                  <span class="ability-name">{{ ability.name }}</span>
                  <span class="ability-pill" @click="openAbilityInfo(ability)">Details</span>
                </div>
                <p class="ability-summary" v-if="ability.summary">{{ ability.summary }}</p>
                <div class="ability-details">
                  <span v-if="ability.costs" class="ability-detail">
                    <strong>Kosten:</strong> {{ ability.costs }}
                  </span>
                  <span v-if="ability.duration" class="ability-detail">
                    <strong>Dauer:</strong> {{ ability.duration }}
                  </span>
                  <span v-if="ability.diceSupplies" class="ability-detail">
                    <strong>Würfelpool:</strong> {{ ability.diceSupplies }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="hasBloodSorcery">
        <h2>Blutrituale</h2>
        <p class="section-info">Verfügbar bis Level {{ bloodSorceryLevel }}</p>

        <div class="rituals-list" v-if="availableBloodRituals.length > 0">
          <div v-for="ritual in availableBloodRituals" :key="ritual.id" class="ritual-item">
            <div class="ritual-header">
              <span class="ritual-name">{{ ritual.name }}</span>
              <div style="display: flex; gap: 0.5rem">
                <span class="ritual-level">Level {{ ritual.level }}</span>
                <span class="ritual-level" style="cursor: pointer" @click="openBloodRitualInfo(ritual)">Details</span>
              </div>
            </div>
            <p class="ritual-desc">{{ ritual.description }}</p>
          </div>
        </div>

        <div v-else class="empty-note">
          <p>Keine Blutrituale verfügbar.</p>
        </div>
      </VCard>

      <VCard v-if="hasOblivion">
        <h2>Vergessenheitszeremonien</h2>
        <p class="section-info">Verfügbar bis Level {{ oblivionLevel }}</p>

        <div class="rituals-list" v-if="availableOblivionCeremonies.length > 0">
          <div v-for="ceremony in availableOblivionCeremonies" :key="ceremony.id" class="ritual-item">
            <div class="ritual-header">
              <span class="ritual-name">{{ ceremony.name }}</span>
              <div style="display: flex; gap: 0.5rem">
                <span class="ritual-level">Level {{ ceremony.level }}</span>
                <span class="ritual-level" style="cursor: pointer" @click="openOblivionCeremonyInfo(ceremony)">Details</span>
              </div>
            </div>
            <p class="ritual-desc">{{ ceremony.summary }}</p>
          </div>
        </div>

        <div v-else class="empty-note">
          <p>Keine Vergessenheitszeremonien verfügbar.</p>
        </div>
      </VCard>

      <VCard
        v-if="(character.traitPackUsages && character.traitPackUsages.length > 0) || predatorBonusTraitsForDisplay.length > 0"
      >
      <h2>Vorzüge & Hintergründe</h2>
        <div class="traits-section">
          <div v-for="usage in traitPackUsagesForDisplay" :key="usage.id" class="trait-pack">
            <h3>{{ getTraitPackName(usage.packID) }}</h3>
            <small class="pack-summary">{{ getTrackPackSummary(usage.packID) }}</small>
            <div class="traits-list">
              <div
                v-for="bonus in predatorBonusTraitsForDisplay.filter(b => b.packID === usage.packID && !b.isFlaw)"
                :key="`pred-${bonus.packID}-${bonus.traitID}`"
                class="trait-item"
              >
                <div style="display: flex; justify-content: space-between; align-items: center; gap: 0.5rem; flex: 1">
                  <div class="trait-main">
                    <span class="trait-name">
                      {{ getTraitInfo(usage.packID, bonus.traitID).name }}
                      <small v-if="bonus.suffix" class="trait-suffix">({{bonus.suffix}})</small>
                    </span>
                    <span class="trait-badge" style="background: rgba(255, 200, 50, 0.12); color: rgba(255, 200, 50, 0.95);">
                      Jagdverh.
                    </span>
                  </div>
                  <span class="trait-level">{{ bonus.level ?? "—" }}</span>
                </div>
                <small class="trait-description">
                  {{ getTraitInfo(usage.packID, bonus.traitID).description }}
                </small>
              </div>

              <div
                v-for="bonus in predatorBonusTraitsForDisplay.filter(b => b.packID === usage.packID && b.isFlaw)"
                :key="`predflaw-${bonus.packID}-${bonus.traitID}`"
                class="trait-item trait-item--flaw"
              >
                <div style="display: flex; justify-content: space-between; align-items: center; gap: 0.5rem; flex: 1">
                  <div class="trait-main">
                  <span class="trait-name">
                    {{ getTraitInfo(usage.packID, bonus.traitID).name }}
                    <small v-if="bonus.suffix" class="trait-suffix">({{bonus.suffix}})</small>
                  </span>
                    <span class="trait-badge trait-badge--flaw">Schwäche</span>
                    <span class="trait-badge" style="background: rgba(255, 200, 50, 0.12); color: rgba(255, 200, 50, 0.95);">
                    Jagdverh.
                  </span>
                  </div>
                  <span class="trait-level trait-level--flaw">{{ bonus.level ?? "—" }}</span>
                </div>
                <small class="trait-description">
                  {{ getTraitInfo(usage.packID, bonus.traitID).description }}
                </small>
              </div>

              <div v-for="trait in usage.traits" :key="trait.id" class="trait-item">
                <div style="display: flex; justify-content: space-between; align-items: center; gap: 0.5rem; flex: 1">
                  <div class="trait-main">
                    <span class="trait-name">
                      {{ getTraitInfo(usage.packID, trait.traitID).name }}
                      <small v-if="trait.suffix" class="trait-suffix">({{trait.suffix}})</small></span>
                    <span v-if="getTraitInfo(usage.packID, trait.traitID).isFlaw" class="trait-badge trait-badge--flaw">
                      Schwäche
                    </span>
                  </div>
                  <span class="trait-level" v-if="trait.customLevel">{{ trait.customLevel }}</span>
                </div>
                <small class="trait-description">
                  {{ getTraitInfo(usage.packID, trait.traitID).description }}
                </small>
              </div>

              <div v-for="flaw in usage.flawTraits" :key="flaw.id" class="trait-item trait-item--flaw">
                <div style="display: flex; justify-content: space-between; align-items: center; gap: 0.5rem; flex: 1">
                  <div class="trait-main">
                    <span class="trait-name">
                      {{ getTraitInfo(usage.packID, flaw.traitID).name }}
                      <small v-if="flaw.suffix" class="trait-suffix">({{flaw.suffix}})</small>
                    </span>
                    <span class="trait-badge trait-badge--flaw">Schwäche</span>
                  </div>
                  <span class="trait-level trait-level--flaw" v-if="flaw.customLevel">{{ flaw.customLevel }}</span>
                </div>
                <small class="trait-description">
                  {{ getTraitInfo(usage.packID, flaw.traitID).description }}
                </small>
              </div>
            </div>
          </div>
        </div>
      </VCard>

      <VInventoryBox
        v-model="character"
        @update-inventory="(inv) => saveField('inventory', inv)"
      />



      <VCard>
        <h2>Chronik & Geschichte</h2>
        <div class="chronicle-grid">
          <VInput
            label="Chronik"
            :model-value="character.chronicle"
            @update:model-value="updateTextField('chronicle', $event as string)"
            placeholder="Name der Chronik…"
          />
          <VInput
            label="Sire"
            :model-value="character.sire"
            @update:model-value="updateTextField('sire', $event as string)"
            placeholder="Name des Sires…"
          />
          <VInput
            label="Ehrgeiz"
            :model-value="character.ambition"
            @update:model-value="updateTextField('ambition', $event as string)"
            placeholder="Langfristiges Ziel…"
          />
          <VInput
            label="Begehren"
            :model-value="character.desire"
            @update:model-value="updateTextField('desire', $event as string)"
            placeholder="Kurzfristiges Verlangen…"
          />
          <VTextarea
            label="Chronikprinzipien"
            :model-value="character.chroniclePrinciples"
            @update:model-value="updateTextField('chroniclePrinciples', $event)"
            placeholder="Die Prinzipien der Chronik…"
            :rows="3"
          />
          <VTextarea
            label="Anker & Überzeugungen"
            :model-value="character.anchorsAndBeliefs"
            @update:model-value="updateTextField('anchorsAndBeliefs', $event)"
            placeholder="Wichtige Überzeugungen und Anker…"
            :rows="3"
          />
          <VTextarea
            class="full-width"
            label="Hintergrundgeschichte"
            :model-value="character.backstory"
            @update:model-value="updateTextField('backstory', $event)"
            placeholder="Die Geschichte deines Charakters…"
            :rows="6"
          />
        </div>
      </VCard>

      <VCard>
        <h2>Notizen</h2>
        <VTextarea
          :model-value="character.notes"
          @update:model-value="updateTextField('notes', $event)"
          placeholder="Persönliche Notizen zum Charakter…"
          :rows="6"
        />
      </VCard>
    </div>

    <div v-if="isSaving" class="save-indicator">Speichert…</div>
  </div>

  <div v-else class="character-view character-view--loading">
    <p>Charakter nicht gefunden.</p>
    <VButton variant="primary" @click="router.push('/')">Zur Übersicht</VButton>
  </div>

  <VLevelHistoryModal v-model="showHistory" :character="character as any" />
  <VCharacterViewersModal v-model="showViewers" :character="character as any" :readonly="!character?.isOwner" />
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.character-view {
  padding: $s-6;
  max-width: 1120px;
  margin: 0 auto;
  position: relative;

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: $s-6;
    gap: $s-4;

    .header-left {
      display: flex;
      flex-direction: column;
      gap: $s-2;

      h1 {
        margin: 0;
        font-family: $font-head;
        font-size: $fs-3;
      }

      .subtitle {
        margin: 0;
        font-style: italic;
        color: $text-2;
      }
    }

    .header-right {
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      gap: $s-3;
    }
  }

  &__content {
    display: grid;
    gap: $s-5;

    h2 {
      margin: 0 0 $s-4;
      font-family: $font-head;
      font-size: $fs-2;
    }

    h3 {
      margin: 0 0 $s-3;
      font-family: $font-head;
      font-size: $fs-1;
      color: $text-1;
    }
  }

  &--loading {
    display: grid;
    place-items: center;
    min-height: 50vh;
    color: $text-1;
    gap: $s-4;
  }
}

.avatar-container {
  width: 80px;
  height: 80px;
}

.avatar-upload {
  display: block;
  width: 100%;
  height: 100%;
  position: relative;
  cursor: pointer;
  border-radius: $r-md;
  overflow: hidden;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: top;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: $bg-2;
  border: 2px solid $border;
  font-size: 2rem;
  font-weight: 700;
  font-family: $font-head;
  color: $text-1;
}

.avatar-input {
  display: none;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity $t-fast $ease;

  span {
    font-size: 0.75rem;
    color: $text-0;
    text-transform: uppercase;
  }
}

.avatar-upload:hover .avatar-overlay {
  opacity: 1;
}

.trackers-section {
  display: grid;
  gap: $s-4;
  grid-template-columns: repeat(2, 1fr);
}

.tracker-card {
  h2 {
    margin: 0;
  }
}

.tracker-card__head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: $s-3;
  margin-bottom: $s-3;
}

.tracker-card__meta {
  font-family: $font-head;
  color: rgba($red-0, 0.9);
  font-weight: 700;
}

.tracker-info {
  margin-bottom: $s-3;
}

.tracker-info--row {
  display: grid;
  gap: $s-2;
}

.tracker-formula {
  font-size: 0.82rem;
  color: $text-2;
}

.content-grid {
  display: grid;
  gap: $s-5;
  grid-template-columns: repeat(2, 1fr);
}

.info-grid {
  display: grid;
  gap: $s-2;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: $s-2 0;
  border-bottom: 1px solid $border;

  &:last-child {
    border-bottom: none;
  }

  .label {
    color: $text-2;
  }
}

.form-grid {
  display: grid;
  gap: $s-4;
}

.form-field {
  display: grid;
  gap: $s-2;
}

.field-label {
  font-size: 0.95rem;
  color: rgba(255, 255, 255, 0.72);
}

.resonance-select {
  display: flex;
  flex-wrap: wrap;
  gap: $s-2;
}

.resonance-option {
  padding: $s-2 $s-3;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid $border;
  border-radius: 999px;
  color: $text-1;
  cursor: pointer;
  font-size: 0.85rem;
  transition: transform $t-fast $ease, border-color $t-fast $ease;

  &:hover {
    transform: translateY(-1px);
    border-color: rgba($red-0, 0.35);
    color: $text-0;
  }

  &--active {
    background: rgba($red-0, 0.12);
    border-color: rgba($red-0, 0.55);
    color: rgba($red-0, 0.95);
  }
}

.chronicle-grid {
  display: grid;
  gap: $s-4;
  grid-template-columns: repeat(2, 1fr);

  .full-width {
    grid-column: 1 / -1;
  }
}

.attributes-grid,
.skills-grid {
  display: grid;
  gap: $s-5;
  grid-template-columns: repeat(3, 1fr);
}

.attr-category,
.skill-category {
  display: grid;
  gap: $s-3;
}

.attr-row,
.skill-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-3;
  padding: $s-1 0;
}

.attr-name,
.skill-name {
  font-weight: 500;
  color: $text-0;
}

.skill-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.skill-spec {
  font-size: 0.8rem;
  color: rgba($red-0, 0.95);
  font-style: italic;
}

.disciplines-list {
  display: grid;
  gap: $s-4;
}

.disc-item {
  padding: $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;
}

.disc-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: $s-3;
  gap: $s-3;
}

.disc-title {
  display: flex;
  flex-direction: column;
  gap: $s-1;
}

.disc-name {
  font-weight: 700;
  color: $text-0;
  font-family: $font-head;
  font-size: 1.1rem;
}

.disc-summary {
  font-size: 0.85rem;
  color: $text-2;
  font-style: italic;
}

.disc-abilities {
  display: grid;
  gap: $s-3;
  margin-top: $s-3;
  padding-top: $s-3;
  border-top: 1px solid $border;
}

.ability-item {
  width: 100%;
  text-align: left;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid rgba(255, 255, 255, 0.10);
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  transition: transform $t-fast $ease, border-color $t-fast $ease;

  &:hover {
    transform: translateY(-1px);
    border-color: rgba($red-0, 0.28);
  }
}

.ability-header {
  display: flex;
  align-items: center;
  gap: $s-2;
  margin-bottom: $s-2;
}

.ability-level {
  width: 26px;
  height: 26px;
  display: grid;
  place-items: center;
  background: rgba($red-0, 0.95);
  color: white;
  font-weight: 800;
  font-size: 0.85rem;
  border-radius: 999px;
  box-shadow: 0 10px 28px rgba($red-0, 0.18);
}

.ability-name {
  font-weight: 700;
  color: $text-0;
  flex: 1;
}

.ability-pill {
  font-size: 0.75rem;
  padding: 0.25rem 0.55rem;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: $text-1;
}

.ability-summary {
  margin: 0 0 $s-2;
  font-size: 0.9rem;
  color: $text-1;
  line-height: 1.5;
}

.ability-details {
  display: flex;
  flex-wrap: wrap;
  gap: $s-2;
}

.ability-detail {
  font-size: 0.8rem;
  color: $text-2;
  padding: $s-1 $s-2;
  background: rgba(255, 255, 255, 0.04);
  border-radius: $r-sm;

  strong {
    color: $text-1;
  }
}

.traits-section {
  display: grid;
  gap: $s-4;
}

.trait-pack {
  h3 {
    margin-bottom: 0;
    padding-bottom: 0;
  }
  .pack-summary {
    display: block;
    margin-bottom: $s-3;
    padding-bottom: $s-2;
    color: $text-2;
    border-bottom: 1px solid $border;
  }
}

.traits-list {
  display: grid;
  gap: $s-2;
}

.trait-item {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, 0.02);

  &--flaw {
    background: rgba($red-0, 0.08);
  }
}

.trait-main {
  display: flex;
  align-items: center;
  gap: $s-2;
}

.trait-name {
  font-weight: 500;
  color: $text-0;
}

.trait-badge {
  font-size: 0.7rem;
  padding: 2px $s-2;
  border-radius: $r-sm;
  text-transform: uppercase;
  font-weight: 700;

  &--flaw {
    background: rgba($red-0, 0.2);
    color: rgba($red-0, 0.95);
  }
}

.trait-level {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba($red-0, 0.12);
  color: rgba($red-0, 0.95);
  font-weight: 700;
  font-size: 0.85rem;

  &--flaw {
    background: rgba($red-0, 0.2);
  }
}

.section-info {
  font-size: 0.85rem;
  color: $text-2;
  margin-bottom: $s-3;
}

.rituals-list {
  display: grid;
  gap: $s-3;
}

.ritual-item {
  padding: $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;
}

.ritual-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-3;
  margin-bottom: $s-2;
}

.ritual-name {
  font-weight: 700;
  color: $text-0;
  font-family: $font-head;
}

.ritual-level {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, 0.06);
  font-size: 0.85rem;
  font-weight: 700;
  color: $text-1;
}

.ritual-desc {
  margin: 0;
  font-size: 0.9rem;
  color: $text-2;
  line-height: 1.5;
}

.empty-note {
  color: $text-2;
  font-size: 0.9rem;
  font-style: italic;
}

.save-indicator {
  position: fixed;
  bottom: $s-4;
  right: $s-4;
  padding: $s-2 $s-4;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid $border;
  border-radius: $r-md;
  font-size: 0.85rem;
  color: $text-1;
  backdrop-filter: blur(10px);
}

@media (max-width: 900px) {
  .trackers-section {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .content-grid {
    grid-template-columns: 1fr;
  }

  .chronicle-grid {
    grid-template-columns: 1fr;
  }

  .attributes-grid,
  .skills-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 600px) {
  .trackers-section {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 420px) {
  .character-view {
    padding: $s-4;

    &__header {
      flex-direction: column;
      align-items: stretch;

      .header-right {
        flex-direction: row;
        justify-content: space-between;
        align-items: center;
      }
    }
  }

  .avatar-container {
    width: 68px;
    height: 68px;
  }
}
</style>
