<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'
import type { V5Character, V5Discipline, V5CharacterDisciplineSelection } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

const localSelections = ref<V5CharacterDisciplineSelection[]>([])

// Base points: 3 (2+1 distribution between two disciplines)
const BASE_POINTS = 3

onMounted(async () => {
  await Promise.all([
    v5data.fetchClans(),
    v5data.fetchDisciplines(),
    v5data.fetchPredatorTypes(),
  ])
  initializeSelections()
})

function initializeSelections() {
  if (character.value.disciplineSelections && character.value.disciplineSelections.length > 0) {
    localSelections.value = [...character.value.disciplineSelections]
  }
}

watch(() => character.value.disciplineSelections, () => {
  if (character.value.disciplineSelections) {
    initializeSelections()
  }
}, { deep: true })

watch(localSelections, (newSelections) => {
  character.value = { ...character.value, disciplineSelections: [...newSelections] }
}, { deep: true })

// Get clan disciplines
const clanDisciplines = computed(() => {
  if (!character.value.clanID) return []
  const clan = v5data.getClanById(character.value.clanID)
  return clan?.disciplines ?? []
})

// Get predator type discipline bonus
const predatorDisciplineBonus = computed(() => {
  if (!character.value.predatorTypeID) return null

  const predator = v5data.predatorTypes?.find(p => p.id === character.value.predatorTypeID)
  if (!predator) return null

  const disciplineAction = predator.actions.find(a => a.type === 'discipline_point')
  if (!disciplineAction) return null

  return {
    points: disciplineAction.data?.points ?? 1,
    disciplineID: disciplineAction.data?.disciplineID, // undefined means user chose
    choosenDisciplineID: (character.value as any)?._predatorActionChoices?.[disciplineAction.id]?.disciplineID
  }
})

// Total available points
const totalAvailablePoints = computed(() => {
  return BASE_POINTS + (predatorDisciplineBonus.value?.points ?? 0)
})

// Points spent
const pointsSpent = computed(() => {
  return localSelections.value.reduce((sum, sel) => sum + sel.currentLevel, 0)
})

// Points remaining
const pointsRemaining = computed(() => {
  return totalAvailablePoints.value - pointsSpent.value
})

// Check if distribution is valid (for base points: 2 in one, 1 in another)
const distributionValid = computed(() => {
  const clanSelections = localSelections.value.filter(sel =>
    clanDisciplines.value.some(d => d.id === sel.disciplineID)
  )

  if (clanSelections.length < 2) return false

  const levels = clanSelections.map(s => s.currentLevel).sort((a, b) => b - a)

  // Check base distribution (2+1)
  if (levels.length >= 2) {
    const hasValidBase = levels[0] >= 2 && levels[1] >= 1
    return hasValidBase && pointsRemaining.value === 0
  }

  return false
})

function getDisciplineLevel(disciplineId: string): number {
  const selection = localSelections.value.find(d => d.disciplineID === disciplineId)
  return selection?.currentLevel ?? 0
}

function setDisciplineLevel(disciplineId: string, level: number) {
  const existingIndex = localSelections.value.findIndex(d => d.disciplineID === disciplineId)

  if (level === 0) {
    if (existingIndex !== -1) {
      localSelections.value = localSelections.value.filter(d => d.disciplineID !== disciplineId)
    }
    return
  }

  if (existingIndex !== -1) {
    localSelections.value[existingIndex] = {
      ...localSelections.value[existingIndex],
      currentLevel: level,
      points: level
    }
    localSelections.value = [...localSelections.value]
  } else {
    localSelections.value = [
      ...localSelections.value,
      {
        id: `disc-${disciplineId}`,
        characterID: '',
        disciplineID: disciplineId,
        points: level,
        currentLevel: level,
        abilities: [],
      }
    ]
  }
}

// Check if discipline can be increased
function canIncrease(disciplineId: string): boolean {
  const currentLevel = getDisciplineLevel(disciplineId)
  if (currentLevel >= 5) return false
  if (pointsRemaining.value <= 0) return false

  return true
}

// Get discipline abilities for current level
function getDisciplineAbilities(discipline: V5Discipline, level: number) {
  return discipline.abilities.filter(a => a.level <= level)
}

// Check if ability requirements are met
function abilityRequirementsMet(ability: any): boolean {
  if (!ability.combinationRefID) return true

  const requiredLevel = ability.combinationLevel ?? 1
  const requiredDiscLevel = getDisciplineLevel(ability.combinationRefID)

  return requiredDiscLevel >= requiredLevel
}

function getDisciplineName(disciplineId: string): string {
  const disc = v5data.disciplines?.find(d => d.id === disciplineId)
  return disc?.name ?? disciplineId
}
</script>

<template>
  <div class="disciplines-step">
    <VCard>
      <h2>Disziplinen verteilen</h2>
      <p class="info-text">
        Verteile {{ BASE_POINTS }} Punkte auf 2 deiner Clan-Disziplinen: 2 Punkte in eine, 1 Punkt in die andere.
        <template v-if="predatorDisciplineBonus">
          <br/>Dein Jagdverhalten gibt dir {{ predatorDisciplineBonus.points }} zusätzliche(n) Punkt(e).
        </template>
      </p>

      <div class="points-status" :class="{ 'points-status--valid': distributionValid }">
        <div class="points-item">
          <span class="points-label">Verfügbar:</span>
          <span class="points-value">{{ totalAvailablePoints }}</span>
        </div>
        <div class="points-item">
          <span class="points-label">Verteilt:</span>
          <span class="points-value">{{ pointsSpent }}</span>
        </div>
        <div class="points-item">
          <span class="points-label">Übrig:</span>
          <span class="points-value" :class="{ 'points-value--ok': pointsRemaining === 0, 'points-value--warn': pointsRemaining > 0 }">
            {{ pointsRemaining }}
          </span>
        </div>
      </div>
    </VCard>

    <VCard>
      <h2>Clan-Disziplinen</h2>

      <p v-if="!character.clanID" class="info-message">
        Bitte wähle zuerst einen Clan aus, um Disziplinen auswählen zu können.
      </p>

      <div v-else-if="clanDisciplines.length > 0" class="disciplines-list">
        <div
          v-for="discipline in clanDisciplines"
          :key="discipline.id"
          class="discipline-item"
          :class="{ 'discipline-item--active': getDisciplineLevel(discipline.id) > 0 }"
        >
          <div class="discipline-header">
            <div class="discipline-info">
              <h3>{{ discipline.name }}</h3>
              <p v-if="discipline.summary" class="discipline-summary">{{ discipline.summary }}</p>
            </div>
            <VDotRating
              :model-value="getDisciplineLevel(discipline.id)"
              :max="Math.min(5, getDisciplineLevel(discipline.id) + pointsRemaining)"
              @update:model-value="(v) => setDisciplineLevel(discipline.id, v)"
            />
          </div>

          <!-- Show abilities for selected disciplines -->
          <div v-if="getDisciplineLevel(discipline.id) > 0 && discipline.abilities?.length" class="discipline-abilities">
            <h4>Verfügbare Kräfte (bis Level {{ getDisciplineLevel(discipline.id) }}):</h4>
            <div class="abilities-list">
              <div
                v-for="ability in getDisciplineAbilities(discipline, getDisciplineLevel(discipline.id))"
                :key="ability.id"
                class="ability-item"
                :class="{ 'ability-item--locked': !abilityRequirementsMet(ability) }"
              >
                <div class="ability-header">
                  <span class="ability-level">{{ ability.level }}</span>
                  <span class="ability-name">{{ ability.name }}</span>
                  <span v-if="ability.combinationRefID" class="ability-req">
                    (Benötigt: {{ getDisciplineName(ability.combinationRefID) }} {{ ability.combinationLevel ?? 1 }})
                  </span>
                </div>
                <p class="ability-summary">{{ ability.summary }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <p v-else class="info-message">
        Keine Disziplinen für den gewählten Clan verfügbar.
      </p>
    </VCard>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.disciplines-step {
  display: grid;
  gap: $s-5;
}

h2 {
  margin: 0 0 $s-3;
  font-family: $font-head;
  font-size: $fs-2;
}

h4 {
  margin: 0 0 $s-3;
  font-family: $font-head;
  font-size: 0.95rem;
  color: $text-1;
}

.info-text {
  margin: 0 0 $s-4;
  color: $text-2;
  font-size: 0.95rem;
  line-height: 1.5;
}

.points-status {
  display: flex;
  gap: $s-4;
  flex-wrap: wrap;
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;

  &--valid {
    border-color: rgba(100, 200, 100, 0.4);
    background: rgba(100, 200, 100, 0.05);
  }
}

.points-item {
  display: flex;
  align-items: center;
  gap: $s-2;
}

.points-label {
  color: $text-2;
  font-weight: 500;
}

.points-value {
  font-weight: 600;
  color: $text-1;

  &--ok {
    color: rgba(100, 200, 100, 0.9);
  }

  &--warn {
    color: $red-0;
  }
}

.info-message {
  text-align: center;
  padding: $s-6;
  color: $text-2;
  margin: 0;
}

.disciplines-list {
  display: grid;
  gap: $s-4;
}

.discipline-item {
  padding: $s-4;
  border-radius: $r-lg;
  background: rgba(255, 255, 255, .02);
  border: 1px solid $border;
  transition: all $t-med $ease;

  &:hover {
    background: rgba(255, 255, 255, .04);
  }

  &--active {
    border-color: $red-0;
    background: rgba(255, 59, 84, .04);
  }
}

.discipline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-4;
}

.discipline-info {
  flex: 1;
  min-width: 0;

  h3 {
    margin: 0 0 $s-1;
    font-family: $font-head;
    font-size: $fs-1;
    color: $text-0;
  }
}

.discipline-summary {
  margin: 0;
  font-size: 0.9rem;
  color: $text-2;
  line-height: 1.4;
}

.discipline-abilities {
  margin-top: $s-4;
  padding-top: $s-4;
  border-top: 1px solid $border;
}

.abilities-list {
  display: grid;
  gap: $s-2;
}

.ability-item {
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
  border: 1px solid $border;

  &--locked {
    opacity: 0.5;
    border-color: rgba(255, 100, 100, 0.3);
  }
}

.ability-header {
  display: flex;
  align-items: center;
  gap: $s-2;
  margin-bottom: $s-2;
}

.ability-level {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: $grad-accent;
  color: white;
  font-size: 0.75rem;
  font-weight: 700;
  display: grid;
  place-items: center;
  flex-shrink: 0;
}

.ability-name {
  font-weight: 600;
  color: $text-0;
}

.ability-req {
  font-size: 0.8rem;
  color: $text-2;
  font-style: italic;
}

.ability-summary {
  margin: 0;
  font-size: 0.85rem;
  color: $text-2;
  line-height: 1.4;
}

@media (max-width: 420px) {
  .points-status {
    flex-direction: column;
    gap: $s-2;
  }

  .discipline-header {
    flex-direction: column;
    align-items: flex-start;
    gap: $s-3;
  }

  .ability-header {
    flex-wrap: wrap;
  }
}
</style>
