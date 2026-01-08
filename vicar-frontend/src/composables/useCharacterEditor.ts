import { ref, computed, watch } from 'vue'
import type { V5Character } from '@/@types/v5'

export type EditorStep = 'initial' | 'clan' | 'predator' | 'attributes' | 'skills' | 'disciplines' | 'traits' | 'optional'

function calculateBloodPotency(generation: number): number {
  if (generation >= 14) return 0
  if (generation >= 12) return 1
  if (generation >= 10) return 2
  if (generation === 9) return 3
  if (generation === 8) return 4
  if (generation === 7) return 5
  if (generation === 6) return 6
  if (generation === 5) return 7
  if (generation <= 4) return 8
  return 0
}

export function useCharacterEditor(initialChar?: V5Character) {
  const currentStep = ref<EditorStep>('initial')
  const initialGeneration = initialChar?.generation || 13
  const character = ref<Partial<V5Character>>({
    name: initialChar?.name || '',
    books: initialChar?.books || [],
    generationEra: initialChar?.generationEra || 'newborn',
    generation: initialGeneration,
    sire: initialChar?.sire || '',
    clanID: initialChar?.clanID || undefined,
    predatorTypeID: initialChar?.predatorTypeID || undefined,
    attributes: initialChar?.attributes || [],
    skills: initialChar?.skills || [],
    disciplineSelections: initialChar?.disciplineSelections || [],
    traitPackUsages: initialChar?.traitPackUsages || [],
    bloodRituals: initialChar?.bloodRituals || [],
    oblivionCeremonies: initialChar?.oblivionCeremonies || [],
    hunger: initialChar?.hunger || 1,
    humanity: initialChar?.humanity || 7,
    stains: initialChar?.stains || 0,
    bloodPotency: initialChar?.bloodPotency ?? calculateBloodPotency(initialGeneration),
  })

  const steps: EditorStep[] = ['initial', 'clan', 'predator', 'attributes', 'skills', 'disciplines', 'traits']

  const currentStepIndex = computed(() => steps.indexOf(currentStep.value))

  const canGoNext = computed(() => {
    switch (currentStep.value) {
      case 'initial':
        return !!(character.value.name &&
          character.value.books && character.value.books.length > 0 &&
          character.value.generationEra &&
          character.value.generation)
      case 'clan':
        return !!character.value.clanID
      case 'predator':
        return !!character.value.predatorTypeID
      case 'attributes':
        return character.value.attributes && character.value.attributes.length > 0
      case 'skills':
        return character.value.skills && character.value.skills.length > 0
      case 'disciplines':
        return true
      case 'traits':
        return true
      case 'optional':
        return true
      default:
        return false
    }
  })

  const canGoPrevious = computed(() => currentStepIndex.value > 0)

  function nextStep() {
    if (canGoNext.value && currentStepIndex.value < steps.length - 1) {
      currentStep.value = steps[currentStepIndex.value + 1]!
    }
  }

  function previousStep() {
    if (canGoPrevious.value) {
      currentStep.value = steps[currentStepIndex.value - 1]!
    }
  }

  function goToStep(step: EditorStep) {
    currentStep.value = step
  }

  watch(() => character.value.clanID, (newClan, oldClan) => {
    if (newClan !== oldClan && oldClan) {
      character.value.disciplineSelections = []
    }
  })

  watch(() => character.value.predatorTypeID, (newPred, oldPred) => {
    if (newPred !== oldPred && oldPred) {
      character.value.disciplineSelections = []
      character.value.traitPackUsages = []
    }
  })

  watch(() => character.value.generation, (newGen) => {
    if (newGen) {
      character.value.bloodPotency = calculateBloodPotency(newGen)
    }
  })

  return {
    character,
    currentStep,
    currentStepIndex,
    canGoNext,
    canGoPrevious,
    nextStep,
    previousStep,
    goToStep,
    steps,
  }
}
