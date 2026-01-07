<script setup lang="ts">
import { computed } from 'vue'
import type { EditorStep } from '@/composables/useCharacterEditor'

const props = defineProps<{
  steps: EditorStep[]
  currentStep: EditorStep
}>()

const emit = defineEmits<{
  (e: 'step-click', step: EditorStep): void
}>()

const stepLabels: Record<EditorStep, string> = {
  initial: 'Grundlagen',
  clan: 'Clan',
  predator: 'Jagdverhalten',
  attributes: 'Attribute',
  skills: 'Fähigkeiten',
  disciplines: 'Disziplinen',
  optional: 'Optional',
}

const currentStepIndex = computed(() => props.steps.indexOf(props.currentStep))

function getStepClass(step: EditorStep, index: number) {
  if (step === props.currentStep) return 'stepper__step--active'
  if (index < currentStepIndex.value) return 'stepper__step--completed'
  return ''
}
</script>

<template>
  <div class="stepper">
    <div
      v-for="(step, index) in steps"
      :key="step"
      class="stepper__step"
      :class="getStepClass(step, index)"
      @click="emit('step-click', step)"
    >
      <div class="stepper__step-number">{{ index + 1 }}</div>
      <div class="stepper__step-label">{{ stepLabels[step] }}</div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.stepper {
  display: flex;
  gap: $s-2;
  overflow-x: auto;
  padding: $s-3 0;

  &__step {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: $s-2;
    padding: $s-3;
    border-radius: $r-md;
    border: 1px solid $border;
    background: rgba(255, 255, 255, .02);
    cursor: pointer;
    transition: all $t-med $ease;
    min-width: 100px;

    &:hover {
      border-color: $border-strong;
      transform: translateY(-2px);
    }

    &--active {
      border-color: $red-0;
      background: rgba(255, 59, 84, .08);

      .stepper__step-number {
        background: $grad-accent;
        color: white;
      }
    }

    &--completed {
      .stepper__step-number {
        background: rgba(255, 255, 255, .12);
      }
    }
  }

  &__step-number {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: rgba(255, 255, 255, .06);
    display: grid;
    place-items: center;
    font-weight: 600;
    transition: all $t-med $ease;
  }

  &__step-label {
    font-size: 0.9rem;
    color: $text-1;
    text-align: center;
  }
}

@media (max-width: 420px) {
  .stepper {
    gap: $s-1;

    &__step {
      min-width: 80px;
      padding: $s-2;
    }

    &__step-label {
      font-size: 0.8rem;
    }
  }
}
</style>
