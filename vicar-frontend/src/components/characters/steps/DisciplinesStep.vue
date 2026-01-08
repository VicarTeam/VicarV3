<script setup lang="ts">
import {ref, onMounted, computed, watch, nextTick} from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import VButton from '@/components/ui/VButton.vue'
import type { V5Character, V5Discipline, V5DisciplineAbility, V5CharacterDisciplineSelection, V5CharacterInternalData, V5CharacterDisciplineAbility } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

onMounted(async () => {
  await Promise.all([
    v5data.fetchClans(),
    v5data.fetchDisciplines(),
    v5data.fetchPredatorTypes(),
  ])
})

const internalData = computed<V5CharacterInternalData>({
  get: () => character.value?.internalData ?? {},
  set: (val: V5CharacterInternalData) => {
    character.value.internalData = val
  }
})

const disciplineDistribution = computed(() => internalData.value.disciplineDistribution ?? {})

const localSelections = computed<V5CharacterDisciplineSelection[]>({
  get: () => character.value?.disciplineSelections ?? [],
  set: (val: V5CharacterDisciplineSelection[]) => {
    character.value.disciplineSelections = val
  }
})

const clanDisciplines = computed(() => {
  if (!character.value.clanID) return []
  const clan = v5data.getClanById(character.value.clanID)
  return clan?.disciplines ?? []
})

const clanDisciplineIds = computed(() => clanDisciplines.value.map(d => d.id))

const predatorDisciplineBonus = computed(() => {
  const bonuses = internalData.value.predatorBonusesApplied
  if (!bonuses?.disciplinePoints?.length) return null
  const bonus = bonuses.disciplinePoints[0]
  return bonus ? {
    disciplineID: bonus.disciplineID,
    points: bonus.points
  } : null
})

const predatorBonusDisciplineID = computed(() => predatorDisciplineBonus.value?.disciplineID ?? null)
const predatorBonusPoints = computed(() => predatorDisciplineBonus.value?.points ?? 0)

const predatorBonusIsClanDiscipline = computed(() => {
  if (!predatorBonusDisciplineID.value) return false
  return clanDisciplineIds.value.includes(predatorBonusDisciplineID.value)
})

const predatorBonusDiscipline = computed(() => {
  if (!predatorBonusDisciplineID.value) return null
  return v5data.disciplines?.find(d => d.id === predatorBonusDisciplineID.value) ?? null
})

const selectedPrimaryId = computed(() => disciplineDistribution.value.primaryDisciplineID)
const selectedSecondaryId = computed(() => disciplineDistribution.value.secondaryDisciplineID)

const editingAbilitiesFor = ref<string | null>(null)

watch(
  () => [disciplineDistribution.value, predatorBonusDisciplineID.value],
  () => {
    updateDisciplineSelections()
  },
  { deep: true, immediate: true }
)

function getDisciplineLevel(disciplineId: string): number {
  let level = 0

  if (disciplineId === selectedPrimaryId.value) {
    level = 2
  } else if (disciplineId === selectedSecondaryId.value) {
    level = 1
  }

  if (predatorBonusDisciplineID.value === disciplineId) {
    if (predatorBonusIsClanDiscipline.value) {
      level += predatorBonusPoints.value
    } else {
      level = predatorBonusPoints.value
    }
  }

  return level
}

function selectPrimary(disciplineId: string) {
  if (predatorBonusDisciplineID.value === disciplineId && !predatorBonusIsClanDiscipline.value) {
    return
  }

  if (disciplineId === selectedSecondaryId.value) {
    internalData.value = {
      ...internalData.value,
      disciplineDistribution: {
        ...disciplineDistribution.value,
        primaryDisciplineID: disciplineId,
        secondaryDisciplineID: selectedPrimaryId.value
      }
    }
  } else {
    internalData.value = {
      ...internalData.value,
      disciplineDistribution: {
        ...disciplineDistribution.value,
        primaryDisciplineID: disciplineId
      }
    }
  }

  updateDisciplineSelections()
}

function selectSecondary(disciplineId: string) {
  if (predatorBonusDisciplineID.value === disciplineId && !predatorBonusIsClanDiscipline.value) {
    return
  }

  if (disciplineId === selectedPrimaryId.value) {
    internalData.value = {
      ...internalData.value,
      disciplineDistribution: {
        ...disciplineDistribution.value,
        secondaryDisciplineID: disciplineId,
        primaryDisciplineID: selectedSecondaryId.value
      }
    }
  } else {
    internalData.value = {
      ...internalData.value,
      disciplineDistribution: {
        ...disciplineDistribution.value,
        secondaryDisciplineID: disciplineId
      }
    }
  }

  updateDisciplineSelections()
}

function canSelectAsPrimary(disciplineId: string): boolean {
  if (disciplineId === selectedPrimaryId.value) return false
  if (predatorBonusDisciplineID.value === disciplineId && !predatorBonusIsClanDiscipline.value) return false
  return true
}

function canSelectAsSecondary(disciplineId: string): boolean {
  if (disciplineId === selectedSecondaryId.value) return false
  if (predatorBonusDisciplineID.value === disciplineId && !predatorBonusIsClanDiscipline.value) return false
  return true
}

function updateDisciplineSelections() {
  const newSelections: V5CharacterDisciplineSelection[] = []

  if (selectedPrimaryId.value) {
    const level = getDisciplineLevel(selectedPrimaryId.value)
    const existingSelection = localSelections.value.find(s => s.disciplineID === selectedPrimaryId.value)
    newSelections.push({
      id: existingSelection?.id ?? `disc-${selectedPrimaryId.value}`,
      characterID: character.value.id ?? '',
      disciplineID: selectedPrimaryId.value,
      points: level,
      currentLevel: level,
      abilities: existingSelection?.abilities ?? []
    })
  }

  if (selectedSecondaryId.value) {
    const level = getDisciplineLevel(selectedSecondaryId.value)
    const existingSelection = localSelections.value.find(s => s.disciplineID === selectedSecondaryId.value)
    newSelections.push({
      id: existingSelection?.id ?? `disc-${selectedSecondaryId.value}`,
      characterID: character.value.id ?? '',
      disciplineID: selectedSecondaryId.value,
      points: level,
      currentLevel: level,
      abilities: existingSelection?.abilities ?? []
    })
  }

  if (predatorBonusDisciplineID.value && !predatorBonusIsClanDiscipline.value) {
    const existingSelection = localSelections.value.find(s => s.disciplineID === predatorBonusDisciplineID.value)
    newSelections.push({
      id: existingSelection?.id ?? `disc-${predatorBonusDisciplineID.value}`,
      characterID: character.value.id ?? '',
      disciplineID: predatorBonusDisciplineID.value,
      points: predatorBonusPoints.value,
      currentLevel: predatorBonusPoints.value,
      abilities: existingSelection?.abilities ?? []
    })
  }

  for (const selection of newSelections) {
    const existingSelection = localSelections.value.find(s => s.disciplineID === selection.disciplineID)
    if (existingSelection) {
      const requiredAbilityCount = selection.points
      const existingAbilities = existingSelection.abilities ?? []
      const filteredAbilities = existingAbilities.filter(a => a.usedLevel <= requiredAbilityCount)
      if (filteredAbilities.length < existingAbilities.length) {
        selection.abilities = filteredAbilities
      }
    }
  }

  localSelections.value = newSelections
}

function startAbilitySelection(disciplineId: string) {
  editingAbilitiesFor.value = disciplineId
}

function closeAbilitySelection() {
  editingAbilitiesFor.value = null
}

const editingDiscipline = computed(() => {
  if (!editingAbilitiesFor.value) return null
  return v5data.disciplines?.find(d => d.id === editingAbilitiesFor.value) ?? null
})

const editingSelection = computed(() => {
  if (!editingAbilitiesFor.value) return null
  return localSelections.value.find(s => s.disciplineID === editingAbilitiesFor.value) ?? null
})

function getSelectedAbilities(disciplineId: string): V5CharacterDisciplineAbility[] {
  const selection = localSelections.value.find(s => s.disciplineID === disciplineId)
  return selection?.abilities ?? []
}

function getRequiredAbilityCount(disciplineId: string): number {
  return getDisciplineLevel(disciplineId)
}

function getAvailableAbilitiesForSlot(discipline: V5Discipline, slotIndex: number): V5DisciplineAbility[] {
  const maxLevel = Math.min(slotIndex, getDisciplineLevel(discipline.id))
  return discipline.abilities.filter(a => {
    if (a.level > maxLevel) return false
    if (!abilityRequirementsMet(a)) return false
    if (getSelectedAbilities(discipline.id).some(sa => sa.abilityID === a.id)) {
      const currentAbility = getAbilityAtSlot(discipline.id, slotIndex - 1)
      return currentAbility?.abilityID === a.id;
    }
    return true
  })
}

function abilityRequirementsMet(ability: V5DisciplineAbility): boolean {
  if (!ability.combinationRefID) return true
  const requiredLevel = ability.combinationLevel ?? 1
  const requiredDiscLevel = getDisciplineLevel(ability.combinationRefID)
  return requiredDiscLevel >= requiredLevel
}

function selectAbility(disciplineId: string, slotIndex: number, abilityId: string) {
  const selections = [...localSelections.value]
  const selectionIndex = selections.findIndex(s => s.disciplineID === disciplineId)
  if (selectionIndex === -1) return

  const selection = selections[selectionIndex]!
  const abilities = [...(selection.abilities ?? [])]

  const discipline = v5data.disciplines?.find(d => d.id === disciplineId)
  const ability = discipline?.abilities.find(a => a.id === abilityId)
  if (!ability) return

  if (abilities.length > slotIndex) {
    abilities[slotIndex] = {
      id: abilities[slotIndex]?.id ?? `ability-${slotIndex}`,
      selectionID: selection.id,
      abilityID: abilityId,
      level: ability.level,
      usedLevel: slotIndex + 1
    }
  } else {
    while (abilities.length < slotIndex) {
      abilities.push({
        id: `ability-${abilities.length}`,
        selectionID: selection.id,
        abilityID: '',
        level: 0,
        usedLevel: abilities.length + 1
      })
    }
    abilities.push({
      id: `ability-${slotIndex}`,
      selectionID: selection.id,
      abilityID: abilityId,
      level: ability.level,
      usedLevel: slotIndex + 1
    })
  }

  selections[selectionIndex] = { ...selection, abilities }
  localSelections.value = selections
}

function clearAbilitySlot(disciplineId: string, slotIndex: number) {
  const selections = [...localSelections.value]
  const selectionIndex = selections.findIndex(s => s.disciplineID === disciplineId)
  if (selectionIndex === -1) return

  const selection = selections[selectionIndex]!
  const abilities = [...(selection.abilities ?? [])]

  if (slotIndex < abilities.length) {
    abilities[slotIndex] = {
      ...abilities[slotIndex]!,
      abilityID: '',
      level: 0
    }
    selections[selectionIndex] = { ...selection, abilities }
    localSelections.value = selections
  }
}

function getAbilityAtSlot(disciplineId: string, slotIndex: number): V5CharacterDisciplineAbility | null {
  const selection = localSelections.value.find(s => s.disciplineID === disciplineId)
  if (!selection?.abilities?.[slotIndex]) return null
  return selection.abilities[slotIndex]!.abilityID ? selection.abilities[slotIndex]! : null
}

function getAbilityById(disciplineId: string, abilityId: string): V5DisciplineAbility | null {
  const discipline = v5data.disciplines?.find(d => d.id === disciplineId)
  return discipline?.abilities.find(a => a.id === abilityId) ?? null
}

function getDisciplineName(disciplineId: string): string {
  const disc = v5data.disciplines?.find(d => d.id === disciplineId)
  return disc?.name ?? disciplineId
}

function getSelectionStatus(disciplineId: string): 'primary' | 'secondary' | 'predator-only' | 'none' {
  if (disciplineId === selectedPrimaryId.value) return 'primary'
  if (disciplineId === selectedSecondaryId.value) return 'secondary'
  if (disciplineId === predatorBonusDisciplineID.value && !predatorBonusIsClanDiscipline.value) return 'predator-only'
  return 'none'
}

function allAbilitiesSelected(disciplineId: string): boolean {
  const required = getRequiredAbilityCount(disciplineId)
  const selected = getSelectedAbilities(disciplineId).filter(a => a.abilityID).length
  return selected >= required
}
</script>

<template>
  <div class="disciplines-step">
    <VCard>
      <h2>Disziplinen wählen</h2>
      <p class="info-text">
        Wähle eine Disziplin als <strong>Primär</strong> (2 Punkte) und eine als <strong>Sekundär</strong> (1 Punkt).
        <template v-if="predatorBonusDisciplineID">
          <br/>
          <template v-if="predatorBonusIsClanDiscipline">
            Dein Jagdverhalten gibt <strong>{{ getDisciplineName(predatorBonusDisciplineID) }}</strong>
            +{{ predatorBonusPoints }} zusätzliche(n) Punkt(e).
          </template>
          <template v-else>
            Dein Jagdverhalten gibt dir <strong>{{ getDisciplineName(predatorBonusDisciplineID) }}</strong>
            auf Stufe {{ predatorBonusPoints }} (separate Disziplin).
          </template>
        </template>
      </p>

      <div class="selection-summary">
        <div class="summary-item" :class="{ 'summary-item--filled': selectedPrimaryId }">
          <span class="summary-label">Primär (2 Pkt)</span>
          <span class="summary-value">
            <template v-if="selectedPrimaryId">
              {{ getDisciplineName(selectedPrimaryId) }}
              <span v-if="predatorBonusDisciplineID === selectedPrimaryId && predatorBonusIsClanDiscipline" class="predator-badge">
                +{{ predatorBonusPoints }}
              </span>
              <span class="total-level">= {{ getDisciplineLevel(selectedPrimaryId) }} Punkte</span>
            </template>
            <template v-else>Nicht gewählt</template>
          </span>
        </div>
        <div class="summary-item" :class="{ 'summary-item--filled': selectedSecondaryId }">
          <span class="summary-label">Sekundär (1 Pkt)</span>
          <span class="summary-value">
            <template v-if="selectedSecondaryId">
              {{ getDisciplineName(selectedSecondaryId) }}
              <span v-if="predatorBonusDisciplineID === selectedSecondaryId && predatorBonusIsClanDiscipline" class="predator-badge">
                +{{ predatorBonusPoints }}
              </span>
              <span class="total-level">= {{ getDisciplineLevel(selectedSecondaryId) }} Punkte</span>
            </template>
            <template v-else>Nicht gewählt</template>
          </span>
        </div>
        <div v-if="predatorBonusDisciplineID && !predatorBonusIsClanDiscipline" class="summary-item summary-item--predator">
          <span class="summary-label">Jagdverh. Bonus</span>
          <span class="summary-value">
            {{ getDisciplineName(predatorBonusDisciplineID) }}
            <span class="total-level">= {{ predatorBonusPoints }} Punkt(e)</span>
          </span>
        </div>
      </div>
    </VCard>

    <VCard>
      <h2>Clan-Disziplinen</h2>

      <p v-if="!character.clanID" class="info-message">
        Bitte wähle zuerst einen Clan aus.
      </p>

      <div v-else-if="clanDisciplines.length > 0" class="disciplines-list">
        <div
          v-for="discipline in clanDisciplines"
          :key="discipline.id"
          class="discipline-item"
          :class="{
            'discipline-item--primary': getSelectionStatus(discipline.id) === 'primary',
            'discipline-item--secondary': getSelectionStatus(discipline.id) === 'secondary'
          }"
        >
          <div class="discipline-header">
            <div class="discipline-info">
              <h3>{{ discipline.name }}</h3>
              <p v-if="discipline.summary" class="discipline-summary">{{ discipline.summary }}</p>
            </div>

            <div class="discipline-actions">
              <button
                class="select-btn select-btn--primary"
                :class="{
                  'select-btn--active': getSelectionStatus(discipline.id) === 'primary',
                  'select-btn--disabled': !canSelectAsPrimary(discipline.id)
                }"
                :disabled="!canSelectAsPrimary(discipline.id) && getSelectionStatus(discipline.id) !== 'primary'"
                @click="selectPrimary(discipline.id)"
              >
                <span class="btn-points">2</span>
                <span class="btn-label">Primär</span>
              </button>
              <button
                class="select-btn select-btn--secondary"
                :class="{
                  'select-btn--active': getSelectionStatus(discipline.id) === 'secondary',
                  'select-btn--disabled': !canSelectAsSecondary(discipline.id)
                }"
                :disabled="!canSelectAsSecondary(discipline.id) && getSelectionStatus(discipline.id) !== 'secondary'"
                @click="selectSecondary(discipline.id)"
              >
                <span class="btn-points">1</span>
                <span class="btn-label">Sekundär</span>
              </button>
            </div>
          </div>

          <div v-if="getDisciplineLevel(discipline.id) > 0" class="discipline-level-row">
            <div class="level-display">
              <span class="level-label">Stufe:</span>
              <div class="level-dots">
                <span
                  v-for="i in 5"
                  :key="i"
                  class="level-dot"
                  :class="{ 'level-dot--filled': i <= getDisciplineLevel(discipline.id) }"
                />
              </div>
              <span class="level-value">{{ getDisciplineLevel(discipline.id) }}</span>
            </div>
            <VButton
              size="sm"
              :variant="allAbilitiesSelected(discipline.id) ? 'secondary' : 'primary'"
              @click="startAbilitySelection(discipline.id)"
            >
              {{ allAbilitiesSelected(discipline.id) ? 'Kräfte ändern' : 'Kräfte wählen' }}
              ({{ getSelectedAbilities(discipline.id).filter(a => a.abilityID).length }}/{{ getRequiredAbilityCount(discipline.id) }})
            </VButton>
          </div>

          <div v-if="getSelectedAbilities(discipline.id).filter(a => a.abilityID).length > 0" class="selected-abilities-preview">
            <div
              v-for="(charAbility) in getSelectedAbilities(discipline.id).filter(a => a.abilityID)"
              :key="charAbility.id"
              class="ability-preview-item"
            >
              <span class="ability-preview-level">{{ charAbility.level }}</span>
              <span class="ability-preview-name">
                {{ getAbilityById(discipline.id, charAbility.abilityID)?.name }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <p v-else class="info-message">
        Keine Disziplinen für den gewählten Clan verfügbar.
      </p>
    </VCard>

    <VCard v-if="predatorBonusDiscipline && !predatorBonusIsClanDiscipline">
      <h2>Jagdverhalten-Disziplin</h2>
      <p class="info-text">
        Diese Disziplin erhältst du durch dein Jagdverhalten. Du kannst keine Verteilungspunkte darauf verwenden.
      </p>

      <div class="discipline-item discipline-item--predator-only">
        <div class="discipline-header">
          <div class="discipline-info">
            <h3>{{ predatorBonusDiscipline.name }}</h3>
            <p v-if="predatorBonusDiscipline.summary" class="discipline-summary">{{ predatorBonusDiscipline.summary }}</p>
          </div>
        </div>

        <div class="discipline-level-row">
          <div class="level-display">
            <span class="level-label">Stufe:</span>
            <div class="level-dots">
              <span
                v-for="i in 5"
                :key="i"
                class="level-dot"
                :class="{ 'level-dot--filled': i <= predatorBonusPoints }"
              />
            </div>
            <span class="level-value">{{ predatorBonusPoints }}</span>
          </div>
          <VButton
            size="sm"
            :variant="allAbilitiesSelected(predatorBonusDiscipline.id) ? 'secondary' : 'primary'"
            @click="startAbilitySelection(predatorBonusDiscipline.id)"
          >
            {{ allAbilitiesSelected(predatorBonusDiscipline.id) ? 'Kräfte ändern' : 'Kräfte wählen' }}
            ({{ getSelectedAbilities(predatorBonusDiscipline.id).filter(a => a.abilityID).length }}/{{ predatorBonusPoints }})
          </VButton>
        </div>

        <div v-if="getSelectedAbilities(predatorBonusDiscipline.id).filter(a => a.abilityID).length > 0" class="selected-abilities-preview">
          <div
            v-for="charAbility in getSelectedAbilities(predatorBonusDiscipline.id).filter(a => a.abilityID)"
            :key="charAbility.id"
            class="ability-preview-item"
          >
            <span class="ability-preview-level">{{ charAbility.level }}</span>
            <span class="ability-preview-name">
              {{ getAbilityById(predatorBonusDiscipline.id, charAbility.abilityID)?.name }}
            </span>
          </div>
        </div>
      </div>
    </VCard>

    <div v-if="editingAbilitiesFor && editingDiscipline && editingSelection" class="ability-selection-overlay" @click.self="closeAbilitySelection">
      <div class="ability-selection-panel">
        <div class="panel-header">
          <h3>{{ editingDiscipline.name }} - Kräfte wählen</h3>
          <button class="close-btn" @click="closeAbilitySelection">&times;</button>
        </div>

        <div class="panel-content">
          <p class="panel-info">
            Wähle {{ getRequiredAbilityCount(editingDiscipline.id) }} Kräfte für diese Disziplin.
            Jeder Punkt muss sequentiell vergeben werden: Punkt 1 = max. Stufe 1, Punkt 2 = max. Stufe 2, usw.
          </p>

          <div class="ability-slots">
            <div
              v-for="slotIdx in getRequiredAbilityCount(editingDiscipline.id)"
              :key="slotIdx"
              class="ability-slot"
            >
              <div class="slot-header">
                <span class="slot-number">Punkt {{ slotIdx }}</span>
                <span class="slot-max">max. Stufe {{ slotIdx }}</span>
              </div>

              <div class="slot-selection">
                <select
                  class="ability-select"
                  :value="getAbilityAtSlot(editingDiscipline.id, slotIdx - 1)?.abilityID ?? ''"
                  @change="(e) => {
                    const val = (e.target as HTMLSelectElement).value
                    if (val) selectAbility(editingDiscipline!.id, slotIdx - 1, val)
                    else clearAbilitySlot(editingDiscipline!.id, slotIdx - 1)
                  }"
                >
                  <option value="">-- Kraft wählen --</option>
                  <option
                    v-for="ability in getAvailableAbilitiesForSlot(editingDiscipline, slotIdx)"
                    :key="ability.id"
                    :value="ability.id"
                  >
                    [{{ ability.level }}] {{ ability.name }}
                    <template v-if="ability.combinationRefID">
                      (Kombo: {{ getDisciplineName(ability.combinationRefID) }} {{ ability.combinationLevel ?? 1 }})
                    </template>
                  </option>
                </select>
              </div>

              <div v-if="getAbilityAtSlot(editingDiscipline.id, slotIdx - 1)" class="slot-ability-info">
                <p class="ability-summary">
                  {{ getAbilityById(editingDiscipline.id, getAbilityAtSlot(editingDiscipline.id, slotIdx - 1)!.abilityID)?.summary }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div class="panel-footer">
          <VButton variant="primary" @click="closeAbilitySelection">Fertig</VButton>
        </div>
      </div>
    </div>
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

h3 {
  margin: 0;
  font-family: $font-head;
  font-size: $fs-1;
  color: $text-0;
}

.info-text {
  margin: 0 0 $s-4;
  color: $text-2;
  font-size: 0.95rem;
  line-height: 1.5;

  strong {
    color: $text-0;
  }
}

.selection-summary {
  display: grid;
  gap: $s-3;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: $s-3;
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;

  &--filled {
    border-color: rgba(100, 200, 100, 0.4);
    background: rgba(100, 200, 100, 0.05);
  }

  &--predator {
    border-color: rgba(255, 59, 84, 0.4);
    background: rgba(255, 59, 84, 0.05);
  }
}

.summary-label {
  font-weight: 600;
  color: $text-1;
  min-width: 120px;
}

.summary-value {
  color: $text-0;
  display: flex;
  align-items: center;
  gap: $s-2;
  flex-wrap: wrap;
}

.predator-badge {
  font-size: 0.75rem;
  padding: 2px 6px;
  border-radius: $r-sm;
  background: rgba(255, 59, 84, 0.15);
  color: $red-0;
  font-weight: 500;
}

.total-level {
  font-size: 0.85rem;
  color: $text-2;
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

  &--primary {
    border-color: $red-0;
    background: rgba(255, 59, 84, .06);
  }

  &--secondary {
    border-color: rgba(100, 150, 255, 0.6);
    background: rgba(100, 150, 255, .06);
  }

  &--predator-only {
    border-color: rgba(255, 59, 84, 0.6);
    background: rgba(255, 59, 84, .08);
  }
}

.discipline-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: $s-4;
}

.discipline-info {
  flex: 1;
  min-width: 0;
}

.discipline-summary {
  margin: $s-1 0 0;
  font-size: 0.9rem;
  color: $text-2;
  line-height: 1.4;
}

.discipline-actions {
  display: flex;
  gap: $s-2;
  flex-shrink: 0;
}

.select-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: $s-2 $s-3;
  border: 1px solid $border;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  color: $text-1;
  cursor: pointer;
  transition: all $t-fast $ease;
  min-width: 60px;

  &:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.06);
  }

  &--primary.select-btn--active {
    border-color: $red-0;
    background: $grad-accent;
    color: white;
  }

  &--secondary.select-btn--active {
    border-color: rgba(100, 150, 255, 0.8);
    background: linear-gradient(135deg, rgba(100, 150, 255, 0.8), rgba(80, 130, 230, 0.9));
    color: white;
  }

  &--disabled, &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
}

.btn-points {
  font-size: 1.1rem;
  font-weight: 700;
}

.btn-label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.discipline-level-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-4;
  margin-top: $s-3;
  padding-top: $s-3;
  border-top: 1px solid $border;
}

.level-display {
  display: flex;
  align-items: center;
  gap: $s-2;
}

.level-label {
  font-size: 0.85rem;
  color: $text-2;
}

.level-dots {
  display: flex;
  gap: 4px;
}

.level-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid $text-2;
  transition: all $t-fast $ease;

  &--filled {
    background: $grad-accent;
    border-color: $red-0;
  }
}

.level-value {
  font-weight: 600;
  color: $text-0;
}

.selected-abilities-preview {
  display: flex;
  flex-wrap: wrap;
  gap: $s-2;
  margin-top: $s-3;
  padding-top: $s-3;
  border-top: 1px solid $border;
}

.ability-preview-item {
  display: inline-flex;
  align-items: center;
  gap: $s-1;
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, .05);
  font-size: 0.85rem;
}

.ability-preview-level {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: $grad-accent;
  color: white;
  font-size: 0.7rem;
  font-weight: 700;
  display: grid;
  place-items: center;
}

.ability-preview-name {
  color: $text-0;
}

// Ability selection overlay
.ability-selection-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: $s-4;
}

.ability-selection-panel {
  background: $bg-1;
  border-radius: $r-xl;
  border: 1px solid $border;
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $s-4;
  border-bottom: 1px solid $border;

  h3 {
    font-size: $fs-2;
  }
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(255, 255, 255, .05);
  border-radius: $r-sm;
  color: $text-1;
  font-size: 1.5rem;
  cursor: pointer;
  display: grid;
  place-items: center;

  &:hover {
    background: rgba(255, 255, 255, .1);
  }
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: $s-4;
}

.panel-info {
  margin: 0 0 $s-4;
  color: $text-2;
  font-size: 0.9rem;
  line-height: 1.5;
}

.ability-slots {
  display: grid;
  gap: $s-4;
}

.ability-slot {
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
  border: 1px solid $border;
}

.slot-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $s-2;
}

.slot-number {
  font-weight: 600;
  color: $text-0;
}

.slot-max {
  font-size: 0.8rem;
  color: $text-2;
}

.ability-select {
  width: 100%;
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 0.9rem;

  &:focus {
    outline: none;
    border-color: $red-0;
  }

  option {
    background: $bg-1;
    color: $text-0;
  }
}

.slot-ability-info {
  margin-top: $s-2;
  padding: $s-2;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, .03);
}

.ability-summary {
  margin: 0;
  font-size: 0.85rem;
  color: $text-2;
  line-height: 1.4;
}

.panel-footer {
  padding: $s-4;
  border-top: 1px solid $border;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 420px) {
  .summary-item {
    flex-direction: column;
    align-items: flex-start;
    gap: $s-1;
  }

  .discipline-header {
    flex-direction: column;
    gap: $s-3;
  }

  .discipline-actions {
    width: 100%;

    .select-btn {
      flex: 1;
    }
  }

  .discipline-level-row {
    flex-direction: column;
    align-items: flex-start;
    gap: $s-3;
  }
}
</style>
