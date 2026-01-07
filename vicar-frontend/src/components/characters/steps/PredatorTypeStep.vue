<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import VModal from '@/components/ui/VModal.vue'
import VButton from '@/components/ui/VButton.vue'
import type { V5Character, V5PredatorType, V5PredatorAction, SkillKey, CategoryKey } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

const showActionsModal = ref(false)
const pendingPredatorType = ref<V5PredatorType | null>(null)
const pendingActions = ref<V5PredatorAction[]>([])
const currentActionIndex = ref(0)

// User choices for actions
const actionChoices = ref<Record<string, any>>({})

const skillDefinitions: Record<CategoryKey, { key: SkillKey; label: string }[]> = {
  physical: [
    { key: 'ath', label: 'Athletik' },
    { key: 'bra', label: 'Handgemenge' },
    { key: 'cra', label: 'Handwerk' },
    { key: 'dri', label: 'Fahren' },
    { key: 'fir', label: 'Schusswaffen' },
    { key: 'lar', label: 'Heimlichkeit' },
    { key: 'mel', label: 'Nahkampf' },
    { key: 'ste', label: 'Überleben' },
    { key: 'sur', label: 'Überlebensinstinkt' },
  ],
  social: [
    { key: 'ani', label: 'Tierkunde' },
    { key: 'eti', label: 'Etikette' },
    { key: 'ins', label: 'Einblick' },
    { key: 'int', label: 'Einschüchtern' },
    { key: 'lea', label: 'Anführen' },
    { key: 'per', label: 'Überzeugen' },
    { key: 'prf', label: 'Darbietung' },
    { key: 'sub', label: 'Gassenwissen' },
    { key: 'str', label: 'Täuschen' },
  ],
  mental: [
    { key: 'aca', label: 'Geisteswissenschaften' },
    { key: 'awa', label: 'Aufmerksamkeit' },
    { key: 'fin', label: 'Finanzen' },
    { key: 'inv', label: 'Untersuchung' },
    { key: 'med', label: 'Medizin' },
    { key: 'occ', label: 'Okkultismus' },
    { key: 'pol', label: 'Politik' },
    { key: 'sci', label: 'Naturwissenschaften' },
    { key: 'tec', label: 'Technologie' },
  ],
}

const allSkillsFlat = computed(() => {
  const result: { key: SkillKey; label: string }[] = []
  Object.values(skillDefinitions).forEach(defs => {
    defs.forEach(def => result.push(def))
  })
  return result
})

onMounted(() => {
  v5data.fetchPredatorTypes()
})

const currentAction = computed(() => {
  return pendingActions.value[currentActionIndex.value] || null
})

const needsUserChoice = computed(() => {
  if (!currentAction.value) return false
  const type = currentAction.value.type
  return type === 'additional_specialization' || type === 'discipline_point'
})

function selectPredatorType(predatorTypeId: string) {
  const predator = v5data.predatorTypes?.find(p => p.id === predatorTypeId)
  if (!predator) return

  // Check if there are actions that need user choices
  const choiceActions = predator.actions.filter(a =>
    a.type === 'additional_specialization' || a.type === 'discipline_point'
  )

  if (choiceActions.length > 0) {
    pendingPredatorType.value = predator
    pendingActions.value = choiceActions
    currentActionIndex.value = 0
    actionChoices.value = {}
    showActionsModal.value = true
  } else {
    // No choices needed, apply immediately
    applyPredatorType(predator, {})
  }
}

function applyPredatorType(predator: V5PredatorType, choices: Record<string, any>) {
  character.value = {
    ...character.value,
    predatorTypeID: predator.id,
    // Store the user's choices for predator actions
    _predatorActionChoices: choices
  }

  // Apply automatic actions
  predator.actions.forEach(action => {
    applyAction(action, choices[action.id])
  })
}

function applyAction(action: V5PredatorAction, choice?: any) {
  switch (action.type) {
    case 'humanity_change':
      if (action.data?.amount !== undefined) {
        const currentHumanity = character.value.humanity ?? 7
        character.value = {
          ...character.value,
          humanity: Math.max(0, Math.min(10, currentHumanity + action.data.amount))
        }
      }
      break

    case 'blood_potency_change':
      if (action.data?.amount !== undefined) {
        const currentBP = character.value.bloodPotency ?? 1
        character.value = {
          ...character.value,
          bloodPotency: Math.max(0, Math.min(10, currentBP + action.data.amount))
        }
      }
      break

    case 'additional_specialization':
      if (choice?.skillKey && choice?.specialization) {
        const skills = [...(character.value.skills || [])]
        const skillIndex = skills.findIndex(s => s.key === choice.skillKey)
        if (skillIndex !== -1) {
          const skill = skills[skillIndex]
          if (!skill.specialization.includes(choice.specialization)) {
            skills[skillIndex] = {
              ...skill,
              specialization: [...skill.specialization, choice.specialization]
            }
            character.value = { ...character.value, skills }
          }
        }
      }
      break

    case 'discipline_point':
      // This will be handled in DisciplinesStep where we show the predator discipline points
      break

    case 'skill_change':
      if (action.data?.skillKey && action.data?.amount !== undefined) {
        const skills = [...(character.value.skills || [])]
        const skillIndex = skills.findIndex(s => s.key === action.data.skillKey)
        if (skillIndex !== -1) {
          skills[skillIndex] = {
            ...skills[skillIndex],
            value: Math.max(0, Math.min(5, skills[skillIndex].value + action.data.amount))
          }
          character.value = { ...character.value, skills }
        }
      }
      break

    case 'attribute_change':
      if (action.data?.attributeKey && action.data?.amount !== undefined) {
        const attrs = [...(character.value.attributes || [])]
        const attrIndex = attrs.findIndex(a => a.key === action.data.attributeKey)
        if (attrIndex !== -1) {
          attrs[attrIndex] = {
            ...attrs[attrIndex],
            value: Math.max(1, Math.min(5, attrs[attrIndex].value + action.data.amount))
          }
          character.value = { ...character.value, attributes: attrs }
        }
      }
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
  if (currentActionIndex.value < pendingActions.value.length - 1) {
    currentActionIndex.value++
  } else {
    // All actions completed, apply predator type
    if (pendingPredatorType.value) {
      applyPredatorType(pendingPredatorType.value, actionChoices.value)
    }
    closeModal()
  }
}

function closeModal() {
  showActionsModal.value = false
  pendingPredatorType.value = null
  pendingActions.value = []
  currentActionIndex.value = 0
  actionChoices.value = {}
}

function getSkillLabel(key: SkillKey): string {
  return allSkillsFlat.value.find(s => s.key === key)?.label || key
}

// Get available disciplines from character's clan
const availableDisciplines = computed(() => {
  const clanID = character.value.clanID
  if (!clanID) return []
  const clan = v5data.clans?.find(c => c.id === clanID)
  return clan?.disciplines || []
})

function getActionTypeLabel(type: string): string {
  switch (type) {
    case 'additional_specialization': return 'Zusätzliche Spezialisierung'
    case 'discipline_point': return 'Disziplin-Punkt'
    case 'humanity_change': return 'Menschlichkeit'
    case 'blood_potency_change': return 'Blutpotenz'
    case 'add_merit': return 'Vorzug'
    case 'add_background': return 'Hintergrund'
    case 'add_flaw': return 'Schwäche'
    case 'skill_change': return 'Fähigkeit'
    case 'attribute_change': return 'Attribut'
    default: return type
  }
}
</script>

<template>
  <VCard>
    <h2>Jagdverhalten auswählen</h2>

    <div class="step-content">
      <div class="predator-grid">
        <div
          v-for="predator in v5data.predatorTypes"
          :key="predator.id"
          class="predator-card"
          :class="{ 'predator-card--selected': character.predatorTypeID === predator.id }"
          @click="selectPredatorType(predator.id)"
        >
          <h3>{{ predator.name }}</h3>
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

      <div v-if="currentAction.type === 'additional_specialization'" class="action-form">
        <div class="form-group">
          <label>Fähigkeit wählen:</label>
          <select
            class="form-select"
            :value="actionChoices[currentAction.id]?.skillKey || ''"
            @change="(e) => setActionChoice(currentAction!.id, 'skillKey', (e.target as HTMLSelectElement).value)"
          >
            <option value="">-- Fähigkeit auswählen --</option>
            <optgroup v-for="(skills, category) in skillDefinitions" :key="category" :label="category === 'physical' ? 'Körperlich' : category === 'social' ? 'Sozial' : 'Geistig'">
              <option v-for="skill in skills" :key="skill.key" :value="skill.key">
                {{ skill.label }}
              </option>
            </optgroup>
          </select>
        </div>

        <div class="form-group">
          <label>Spezialisierung:</label>
          <input
            type="text"
            class="form-input"
            placeholder="Spezialisierung eingeben..."
            :value="actionChoices[currentAction.id]?.specialization || ''"
            @input="(e) => setActionChoice(currentAction!.id, 'specialization', (e.target as HTMLInputElement).value)"
          />
        </div>
      </div>

      <!-- Discipline Point -->
      <div v-if="currentAction.type === 'discipline_point'" class="action-form">
        <div class="form-group">
          <label>Disziplin wählen:</label>
          <select
            class="form-select"
            :value="actionChoices[currentAction.id]?.disciplineID || ''"
            @change="(e) => setActionChoice(currentAction!.id, 'disciplineID', (e.target as HTMLSelectElement).value)"
          >
            <option value="">-- Disziplin auswählen --</option>
            <option v-for="disc in availableDisciplines" :key="disc.id" :value="disc.id">
              {{ disc.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="action-progress">
        Schritt {{ currentActionIndex + 1 }} von {{ pendingActions.length }}
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
          {{ currentActionIndex < pendingActions.length - 1 ? 'Weiter' : 'Bestätigen' }}
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

  &:hover {
    border-color: $border-strong;
    background: rgba(255, 255, 255, .04);
    transform: translateY(-2px);
  }

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
    box-shadow: 0 0 0 4px rgba(255, 59, 84, .15);
  }

  h3 {
    margin: 0 0 $s-3;
    font-family: $font-head;
    font-size: $fs-2;
  }
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

    &:last-child {
      margin-bottom: 0;
    }
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

// Modal styles
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

  label {
    font-weight: 500;
    color: $text-1;
    font-size: 0.9rem;
  }
}

.form-select,
.form-input {
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

.form-select {
  option {
    background: $bg-1;
    color: $text-0;
  }

  optgroup {
    font-weight: 600;
    color: $text-1;
  }
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
