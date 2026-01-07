<script setup lang="ts">
import { onMounted, computed, shallowRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useCharacterEditor } from '@/composables/useCharacterEditor'
import VStepper from '@/components/characters/VStepper.vue'
import VButton from '@/components/ui/VButton.vue'
import InitialSetupStep from '@/components/characters/steps/InitialSetupStep.vue'
import ClanSelectionStep from '@/components/characters/steps/ClanSelectionStep.vue'
import PredatorTypeStep from '@/components/characters/steps/PredatorTypeStep.vue'
import AttributesStep from '@/components/characters/steps/AttributesStep.vue'
import SkillsStep from '@/components/characters/steps/SkillsStep.vue'
import DisciplinesStep from '@/components/characters/steps/DisciplinesStep.vue'
import OptionalStep from '@/components/characters/steps/OptionalStep.vue'

const route = useRoute()
const router = useRouter()
const store = useCharactersStore()

const characterId = computed(() => route.params.id as string)

const editor = useCharacterEditor(store.currentCharacter ?? undefined)

const stepComponents = shallowRef({
  initial: InitialSetupStep,
  clan: ClanSelectionStep,
  predator: PredatorTypeStep,
  attributes: AttributesStep,
  skills: SkillsStep,
  disciplines: DisciplinesStep,
  optional: OptionalStep,
})

onMounted(async () => {
  if (characterId.value) {
    await store.fetchCharacter(characterId.value)
    if (store.currentCharacter) {
      editor.character.value = {
        ...editor.character.value,
        ...store.currentCharacter
      }
    }
  }
})

async function saveAndContinue() {
  if (characterId.value) {
    await store.updateCharacter(characterId.value, editor.character.value)
  }
  editor.nextStep()
}

async function saveAndGoBack() {
  if (characterId.value) {
    await store.updateCharacter(characterId.value, editor.character.value)
  }
  editor.previousStep()
}

async function finish() {
  if (characterId.value) {
    await store.updateCharacter(characterId.value, editor.character.value)
  }
  router.push(`/characters/${characterId.value}`)
}
</script>

<template>
  <div class="editor">
    <div class="editor__header">
      <div class="header-left">
        <VButton variant="ghost" @click="router.push(`/characters/${characterId}`)">
          ← Zurück zur Ansicht
        </VButton>
        <h1>Character Editor</h1>
      </div>
    </div>

    <VStepper
      :steps="editor.steps"
      :current-step="editor.currentStep.value"
      @step-click="editor.goToStep"
    />

    <div class="editor__content">
      <component
        :is="stepComponents[editor.currentStep.value]"
        v-model="editor.character"
      />
    </div>

    <div class="editor__actions">
      <VButton
        v-if="editor.canGoPrevious"
        variant="secondary"
        @click="saveAndGoBack"
      >
        Zurück
      </VButton>

      <VButton
        v-if="editor.canGoNext && editor.currentStepIndex.value < editor.steps.length - 1"
        variant="primary"
        :disabled="!editor.canGoNext"
        @click="saveAndContinue"
      >
        Weiter
      </VButton>

      <VButton
        v-else
        variant="primary"
        @click="finish"
      >
        Fertigstellen
      </VButton>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.editor {
  padding: $s-6;
  max-width: 920px;
  margin: 0 auto;

  &__header {
    margin-bottom: $s-6;

    .header-left {
      display: flex;
      flex-direction: column;
      gap: $s-2;

      h1 {
        margin: 0;
        font-family: $font-head;
        font-size: $fs-3;
      }
    }
  }

  &__content {
    margin: $s-6 0;
  }

  &__actions {
    display: flex;
    gap: $s-3;
    justify-content: flex-end;
    padding-top: $s-5;
    border-top: 1px solid $border;
  }
}

@media (max-width: 420px) {
  .editor {
    padding: $s-4;

    &__actions {
      flex-direction: column-reverse;
    }
  }
}
</style>
