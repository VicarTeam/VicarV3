import { ref, computed, watch } from 'vue'
import type { V5Character } from '@/@types/v5'

export type EditorStep = 'initial' | 'clan' | 'predator' | 'attributes' | 'skills' | 'disciplines' | 'optional'

export function useCharacterEditor(initialChar?: V5Character) {
  const currentStep = ref<EditorStep>('initial')
  const character = ref<Partial<V5Character>>({
    name: initialChar?.name || '',
    books: initialChar?.books || [],
    generationEra: initialChar?.generationEra || 'newborn',
    generation: initialChar?.generation || 13,
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
    bloodPotency: initialChar?.bloodPotency || 0,
  })

  const steps: EditorStep[] = ['initial', 'clan', 'predator', 'attributes', 'skills', 'disciplines', 'optional']

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
