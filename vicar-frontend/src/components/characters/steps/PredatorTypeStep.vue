<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import VModal from '@/components/ui/VModal.vue'
import VButton from '@/components/ui/VButton.vue'
import type { V5Character, V5PredatorType, V5PredatorAction } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

const showActionsModal = ref(false)
const pendingPredatorType = ref<V5PredatorType | null>(null)
const actionsToProcess = ref<V5PredatorAction[]>([])
const currentActionIndex = ref(0)
const actionChoices = ref<Record<string, any>>({})

// Skill key mapping: JSON keys -> Frontend keys
const skillKeyMap: Record<string, string> = {
  ken: 'ani', // Animal Ken
  stw: 'str', // Streetwise
  pes: 'per', // Persuasion
  per: 'prf', // Performance
}

// Reverse mapping
const skillKeyMapReverse: Record<string, string> = {
  ani: 'ken',
  str: 'stw',
  per: 'pes',
  prf: 'per',
}

// Skill labels
const skillLabels: Record<string, string> = {
  ath: 'Athletik', bra: 'Handgemenge', cra: 'Handwerk', dri: 'Fahren',
  fir: 'Schusswaffen', lar: 'Heimlichkeit', mel: 'Nahkampf', ste: 'Diebstahl',
  sur: 'Überleben', ani: 'Tierkunde', eti: 'Etikette', ins: 'Einblick',
  int: 'Einschüchtern', lea: 'Führung', per: 'Überreden', prf: 'Darbietung',
  sub: 'Täuschung', str: 'Szenekenntnis', aca: 'Bildung', awa: 'Aufmerksamkeit',
  fin: 'Finanzen', inv: 'Nachforschung', med: 'Medizin', occ: 'Okkultismus',
  pol: 'Politik', sci: 'Wissenschaft', tec: 'Technologie',
  // Also map JSON keys for display
  ken: 'Tierkunde', stw: 'Szenekenntnis', pes: 'Überreden',
}

onMounted(async () => {
  await Promise.all([
    v5data.fetchPredatorTypes(),
    v5data.fetchClans(),
    v5data.fetchDisciplines(),
  ])
})

// Get character's clan oldVicarID for restriction checking
const characterClanOldVicarID = computed(() => {
  if (!character.value.clanID) return null
  const clan = v5data.clans?.find(c => c.id === character.value.clanID)
  return clan?.oldVicarID ?? null
})

// Get character's selected books oldVicarIDs
const characterBookOldVicarIDs = computed(() => {
  return character.value.books?.map(b => b.oldVicarID).filter(id => id !== undefined) ?? []
})

// Check if a predator type passes its restrictions
function checkPredatorTypeRestriction(predator: V5PredatorType): { allowed: boolean; reason?: string } {
  if (!predator.restriction) return { allowed: true }

  const restriction = predator.restriction as any
  const type = restriction.type
  const data = restriction.data

  switch (type) {
    case 'exclude_clans':
      if (characterClanOldVicarID.value && (data as number[]).includes(characterClanOldVicarID.value)) {
        return { allowed: false, reason: 'Nicht für deinen Clan verfügbar' }
      }
      break
    case 'only_clans':
      if (!characterClanOldVicarID.value || !(data as number[]).includes(characterClanOldVicarID.value)) {
        return { allowed: false, reason: 'Nur für bestimmte Clans verfügbar' }
      }
      break
    case 'max_blood_potency':
      if ((character.value.bloodPotency ?? 1) > (data as number)) {
        return { allowed: false, reason: `Nur für Blutpotenz ${data} oder niedriger` }
      }
      break
    case 'book_activated':
      const requiredBooks = data as number[]
      if (!requiredBooks.some(b => characterBookOldVicarIDs.value.includes(b))) {
        return { allowed: false, reason: 'Benötigt bestimmte Bücher' }
      }
      break
  }

  return { allowed: true }
}

// Filter available predator types
const availablePredatorTypes = computed(() => {
  return v5data.predatorTypes?.map(p => ({
    ...p,
    ...checkPredatorTypeRestriction(p)
  })) ?? []
})

// Current action being processed
const currentAction = computed(() => actionsToProcess.value[currentActionIndex.value] || null)

// Parse specialization choices from action data
function parseSpecializationChoices(action: V5PredatorAction): { skillKey: string; spec: string; needsInput: boolean; inputLabel?: string }[] {
  const data = action.data as any
  if (!data?.choices) return []

  return (data.choices as string[]).map(choice => {
    const [skillKey, specPart] = choice.split('=')
    const frontendKey = skillKeyMap[skillKey] || skillKey

    if (specPart.startsWith('$input:')) {
      return {
        skillKey: frontendKey,
        spec: '',
        needsInput: true,
        inputLabel: specPart.replace('$input:', '')
      }
    }
    return {
      skillKey: frontendKey,
      spec: specPart,
      needsInput: false
    }
  })
}

// Parse discipline choices from action data
function parseDisciplineChoices(action: V5PredatorAction): { disciplineOldVicarID: number; restriction?: any; allowed: boolean }[] {
  const data = action.data as any
  if (!data?.choices) return []

  return (data.choices as any[]).map(choice => {
    const discOldVicarID = choice.id as number
    let allowed = true

    if (choice.restriction) {
      const rType = choice.restriction.type
      const rData = choice.restriction.data

      switch (rType) {
        case 'only_clans':
          allowed = characterClanOldVicarID.value !== null && (rData as number[]).includes(characterClanOldVicarID.value)
          break
        case 'book_activated':
          allowed = (rData as number[]).some(b => characterBookOldVicarIDs.value.includes(b))
          break
      }
    }

    return {
      disciplineOldVicarID: discOldVicarID,
      restriction: choice.restriction,
      allowed
    }
  })
}

// Get discipline by oldVicarID
function getDisciplineByOldVicarID(oldVicarID: number) {
  return v5data.disciplines?.find(d => d.oldVicarID === oldVicarID)
}

// Available discipline choices for current action
const availableDisciplineChoices = computed(() => {
  if (!currentAction.value || currentAction.value.type !== 'discipline_point') return []

  const choices = parseDisciplineChoices(currentAction.value)
  return choices.map(c => ({
    ...c,
    discipline: getDisciplineByOldVicarID(c.disciplineOldVicarID)
  })).filter(c => c.discipline && c.allowed)
})

// Available specialization choices for current action
const availableSpecChoices = computed(() => {
  if (!currentAction.value || currentAction.value.type !== 'additional_specialization') return []
  return parseSpecializationChoices(currentAction.value)
})

// Check which actions need user input
function actionsNeedingChoice(predator: V5PredatorType): V5PredatorAction[] {
  return predator.actions.filter(a => {
    if (a.type === 'additional_specialization') {
      const choices = parseSpecializationChoices(a)
      return choices.length > 1 || choices.some(c => c.needsInput)
    }
    if (a.type === 'discipline_point') {
      const choices = parseDisciplineChoices(a)
      const availableChoices = choices.filter(c => {
        if (!c.restriction) return true
        // Check restriction
        if (c.restriction.type === 'only_clans') {
          return characterClanOldVicarID.value !== null &&
                 (c.restriction.data as number[]).includes(characterClanOldVicarID.value)
        }
        if (c.restriction.type === 'book_activated') {
          return (c.restriction.data as number[]).some(b => characterBookOldVicarIDs.value.includes(b))
        }
        return true
      })
      return availableChoices.length > 1
    }
    if (a.type === 'add_flaw' && (a.data as any)?.choices?.length > 1) {
      return true
    }
    if (a.type === 'spend_background_points_between' || a.type === 'spend_flaw_points_between') {
      return true
    }
    return false
  })
}

function selectPredatorType(predatorTypeId: string) {
  const predator = v5data.predatorTypes?.find(p => p.id === predatorTypeId)
  if (!predator) return

  // Check restrictions
  const restriction = checkPredatorTypeRestriction(predator)
  if (!restriction.allowed) return

  // Get actions that need user choices
  const choiceActions = actionsNeedingChoice(predator)

  if (choiceActions.length > 0) {
    pendingPredatorType.value = predator
    actionsToProcess.value = choiceActions
    currentActionIndex.value = 0
    actionChoices.value = {}
    showActionsModal.value = true
  } else {
    // No choices needed, apply immediately
    applyPredatorType(predator, {})
  }
}

function applyPredatorType(predator: V5PredatorType, choices: Record<string, any>) {
  // Store predator type and choices
  character.value.predatorTypeID = predator.id

  // Apply all actions
  predator.actions.forEach(action => {
    applyAction(action, choices[action.id])
  })
}

function applyAction(action: V5PredatorAction, choice?: any) {
  const data = action.data as any

  switch (action.type) {
    case 'humanity_change':
      if (data?.amount !== undefined) {
        character.value.humanity = Math.max(0, Math.min(10, (character.value.humanity ?? 7) + data.amount))
      }
      break

    case 'blood_potency_change':
      if (data?.amount !== undefined) {
        character.value.bloodPotency = Math.max(0, Math.min(10, (character.value.bloodPotency ?? 1) + data.amount))
      }
      break

    case 'additional_specialization':
      if (choice?.skillKey && choice?.specialization) {
        const skills = [...(character.value.skills || [])]
        const frontendKey = skillKeyMap[choice.skillKey] || choice.skillKey
        const skillIndex = skills.findIndex(s => s.key === frontendKey)
        if (skillIndex !== -1) {
          const skill = skills[skillIndex]
          if (!skill.specialization.includes(choice.specialization)) {
            skills[skillIndex] = {
              ...skill,
              specialization: [...skill.specialization, choice.specialization]
            }
            character.value.skills = skills
          }
        }
      } else if (!choice && data?.choices?.length === 1) {
        // Auto-apply single choice
        const choices = parseSpecializationChoices(action)
        if (choices.length === 1 && !choices[0].needsInput) {
          const c = choices[0]
          const skills = [...(character.value.skills || [])]
          const skillIndex = skills.findIndex(s => s.key === c.skillKey)
          if (skillIndex !== -1) {
            const skill = skills[skillIndex]
            if (!skill.specialization.includes(c.spec)) {
              skills[skillIndex] = {
                ...skill,
                specialization: [...skill.specialization, c.spec]
              }
              character.value.skills = skills
            }
          }
        }
      }
      break

    case 'discipline_point':
      // Store discipline choice for DisciplinesStep to use
      if (choice?.disciplineID) {
        const predatorDisciplines = character.value._predatorDisciplines || []
        character.value._predatorDisciplines = [...predatorDisciplines, {
          disciplineID: choice.disciplineID,
          points: 1
        }]
      }
      break

    // TODO: Implement other action types as needed
    case 'add_merit':
    case 'add_background':
    case 'add_flaw':
    case 'add_background_points':
    case 'spend_background_points_between':
    case 'spend_flaw_points_between':
      // These will be handled in TraitsStep
      break
  }
}

function setActionChoice(actionId: string, key: string, value: any) {
  actionChoices.value = {
    ...actionChoices.value,
    [actionId]: {
      ...(actionChoices.value[actionId] || {}),
      [key]: value
    }
  }
}

function isCurrentActionComplete(): boolean {
  if (!currentAction.value) return true

  const choice = actionChoices.value[currentAction.value.id]

  switch (currentAction.value.type) {
    case 'additional_specialization':
      return !!(choice?.skillKey && choice?.specialization?.trim())
    case 'discipline_point':
      return !!(choice?.disciplineID)
    default:
      return true
  }
}

function nextAction() {
  if (currentActionIndex.value < actionsToProcess.value.length - 1) {
    currentActionIndex.value++
  } else {
    if (pendingPredatorType.value) {
      applyPredatorType(pendingPredatorType.value, actionChoices.value)
    }
    closeModal()
  }
}

function closeModal() {
  showActionsModal.value = false
  pendingPredatorType.value = null
  actionsToProcess.value = []
  currentActionIndex.value = 0
  actionChoices.value = {}
}

function getActionTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    additional_specialization: 'Spezialisierung',
    discipline_point: 'Disziplin',
    humanity_change: 'Menschlichkeit',
    blood_potency_change: 'Blutpotenz',
    add_merit: 'Vorzug',
    add_background: 'Hintergrund',
    add_background_points: 'Hintergrund',
    add_flaw: 'Schwäche',
    spend_background_points_between: 'Hintergründe verteilen',
    spend_flaw_points_between: 'Schwächen verteilen',
  }
  return labels[type] || type
}
</script>

<template>
  <VCard>
    <h2>Jagdverhalten auswählen</h2>

    <div class="step-content">
      <div class="predator-grid">
        <div
          v-for="predator in availablePredatorTypes"
          :key="predator.id"
          class="predator-card"
          :class="{
            'predator-card--selected': character.predatorTypeID === predator.id,
            'predator-card--disabled': !predator.allowed
          }"
          @click="predator.allowed && selectPredatorType(predator.id)"
        >
          <h3>{{ predator.name }}</h3>
          <p v-if="!predator.allowed" class="restriction-warning">{{ predator.reason }}</p>
          <p class="predator-description">{{ predator.description }}</p>

          <div v-if="predator.actions.length > 0" class="predator-actions">
            <span class="label">Effekte:</span>
            <ul class="actions-list">
              <li v-for="action in predator.actions" :key="action.id">
                <span class="action-type">{{ getActionTypeLabel(action.type) }}:</span>
                {{ action.description }}
              </li>
            </ul>
          </div>
        </div>
      </div>

      <p v-if="!v5data.predatorTypes || v5data.predatorTypes.length === 0" class="empty-state">
        Keine Jagdverhalten verfügbar. Stelle sicher, dass du Bücher ausgewählt hast.
      </p>
    </div>
  </VCard>

  <VModal
    v-model="showActionsModal"
    :title="pendingPredatorType?.name ? `${pendingPredatorType.name} - Auswahl` : 'Auswahl'"
    size="md"
    :close-on-backdrop="false"
    :close-on-esc="false"
  >
    <div v-if="currentAction" class="action-modal-content">
      <p class="action-description">{{ currentAction.description }}</p>

      <!-- Additional Specialization -->
      <div v-if="currentAction.type === 'additional_specialization'" class="action-form">
        <div class="form-group">
          <label>Fähigkeit & Spezialisierung wählen:</label>
          <div class="choice-list">
            <label
              v-for="(choice, idx) in availableSpecChoices"
              :key="idx"
              class="choice-item"
              :class="{ 'choice-item--selected': actionChoices[currentAction.id]?.skillKey === choice.skillKey && (!choice.needsInput || actionChoices[currentAction.id]?.specialization) }"
            >
              <input
                type="radio"
                :name="`spec-choice-${currentAction.id}`"
                :value="choice.skillKey"
                :checked="actionChoices[currentAction.id]?.skillKey === choice.skillKey"
                @change="setActionChoice(currentAction.id, 'skillKey', choice.skillKey); if (!choice.needsInput) setActionChoice(currentAction.id, 'specialization', choice.spec)"
              />
              <span class="choice-label">
                {{ skillLabels[choice.skillKey] || choice.skillKey }}
                <span v-if="!choice.needsInput" class="choice-spec">({{ choice.spec }})</span>
              </span>
            </label>
          </div>

          <!-- Custom input if selected choice needs it -->
          <div
            v-if="availableSpecChoices.find(c => c.skillKey === actionChoices[currentAction.id]?.skillKey)?.needsInput"
            class="custom-spec-input"
          >
            <input
              type="text"
              class="form-input"
              :placeholder="availableSpecChoices.find(c => c.skillKey === actionChoices[currentAction.id]?.skillKey)?.inputLabel || 'Spezialisierung eingeben...'"
              :value="actionChoices[currentAction.id]?.specialization || ''"
              @input="(e) => setActionChoice(currentAction.id, 'specialization', (e.target as HTMLInputElement).value)"
            />
          </div>
        </div>
      </div>

      <!-- Discipline Point -->
      <div v-if="currentAction.type === 'discipline_point'" class="action-form">
        <div class="form-group">
          <label>Disziplin wählen:</label>
          <div class="choice-list">
            <label
              v-for="choice in availableDisciplineChoices"
              :key="choice.disciplineOldVicarID"
              class="choice-item"
              :class="{ 'choice-item--selected': actionChoices[currentAction.id]?.disciplineID === choice.discipline?.id }"
            >
              <input
                type="radio"
                :name="`disc-choice-${currentAction.id}`"
                :value="choice.discipline?.id"
                :checked="actionChoices[currentAction.id]?.disciplineID === choice.discipline?.id"
                @change="setActionChoice(currentAction.id, 'disciplineID', choice.discipline?.id)"
              />
              <span class="choice-label">{{ choice.discipline?.name }}</span>
            </label>
          </div>
          <p v-if="availableDisciplineChoices.length === 0" class="no-choices">
            Keine Disziplinen verfügbar (Clan-Einschränkung oder fehlende Bücher)
          </p>
        </div>
      </div>

      <div class="action-progress">
        Schritt {{ currentActionIndex + 1 }} von {{ actionsToProcess.length }}
      </div>
    </div>

    <template #footer>
      <div class="modal-actions">
        <VButton variant="secondary" @click="closeModal">Abbrechen</VButton>
        <VButton
          variant="primary"
          :disabled="!isCurrentActionComplete()"
          @click="nextAction"
        >
          {{ currentActionIndex < actionsToProcess.length - 1 ? 'Weiter' : 'Bestätigen' }}
        </VButton>
      </div>
    </template>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

h2 {
  margin: 0 0 $s-4;
  font-family: $font-head;
  font-size: $fs-2;
}

.step-content {
  display: grid;
  gap: $s-4;
}

.predator-grid {
  display: grid;
  gap: $s-4;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
}

.predator-card {
  padding: $s-4;
  border-radius: $r-lg;
  border: 2px solid $border;
  background: rgba(255, 255, 255, .02);
  cursor: pointer;
  transition: all $t-med $ease;

  &:hover:not(.predator-card--disabled) {
    border-color: $border-strong;
    background: rgba(255, 255, 255, .04);
    transform: translateY(-2px);
  }

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
    box-shadow: 0 0 0 4px rgba(255, 59, 84, .15);
  }

  &--disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  h3 {
    margin: 0 0 $s-2;
    font-family: $font-head;
    font-size: $fs-2;
  }
}

.restriction-warning {
  margin: 0 0 $s-2;
  padding: $s-2;
  border-radius: $r-sm;
  background: rgba(255, 100, 100, .15);
  color: rgba(255, 100, 100, .9);
  font-size: 0.85rem;
  font-weight: 500;
}

.predator-description {
  margin: 0 0 $s-3;
  color: $text-1;
  font-size: 0.95rem;
  line-height: 1.5;
}

.predator-actions {
  padding-top: $s-3;
  border-top: 1px solid $border;
}

.label {
  font-weight: 600;
  color: $text-1;
  font-size: 0.9rem;
  display: block;
  margin-bottom: $s-2;
}

.actions-list {
  margin: 0;
  padding-left: $s-4;
  color: $text-2;
  font-size: 0.9rem;

  li {
    margin-bottom: $s-1;
    &:last-child { margin-bottom: 0; }
  }
}

.action-type {
  font-weight: 500;
  color: $text-1;
}

.empty-state {
  text-align: center;
  padding: $s-6;
  color: $text-2;
}

// Modal
.action-modal-content {
  display: grid;
  gap: $s-4;
}

.action-description {
  margin: 0;
  color: $text-1;
  font-size: 0.95rem;
  line-height: 1.5;
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
  border: 1px solid $border;
}

.action-form {
  display: grid;
  gap: $s-3;
}

.form-group {
  display: grid;
  gap: $s-2;

  > label:first-child {
    font-weight: 500;
    color: $text-1;
    font-size: 0.9rem;
  }
}

.choice-list {
  display: grid;
  gap: $s-2;
}

.choice-item {
  display: flex;
  align-items: center;
  gap: $s-2;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  cursor: pointer;
  transition: all $t-fast $ease;

  &:hover {
    border-color: $border-strong;
    background: rgba(255, 255, 255, .04);
  }

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
  }

  input[type="radio"] {
    accent-color: $red-0;
  }
}

.choice-label {
  color: $text-0;
  font-weight: 500;
}

.choice-spec {
  color: $text-2;
  font-weight: 400;
  margin-left: $s-1;
}

.custom-spec-input {
  margin-top: $s-2;
}

.form-input {
  width: 100%;
  padding: $s-2 $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 0.95rem;

  &:focus {
    outline: none;
    border-color: $red-0;
  }

  &::placeholder {
    color: $text-2;
  }
}

.no-choices {
  margin: 0;
  padding: $s-3;
  text-align: center;
  color: $text-2;
  font-size: 0.9rem;
  background: rgba(255, 100, 100, .05);
  border-radius: $r-md;
}

.action-progress {
  text-align: center;
  color: $text-2;
  font-size: 0.85rem;
  padding-top: $s-3;
  border-top: 1px solid $border;
}

.modal-actions {
  display: flex;
  gap: $s-3;
  justify-content: flex-end;
}

@media (max-width: 420px) {
  .predator-grid {
    grid-template-columns: 1fr;
  }

  .modal-actions {
    flex-direction: column-reverse;
  }
}
</style>
