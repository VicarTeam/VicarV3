<script setup lang="ts">
import {onMounted, computed, shallowRef, watch, ref, nextTick} from 'vue'
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
import TraitsStep from '@/components/characters/steps/TraitsStep.vue'
import OptionalStep from '@/components/characters/steps/OptionalStep.vue'
import type {V5Character} from "@/@types/v5.ts";

const route = useRoute()
const router = useRouter()
const store = useCharactersStore()

const characterId = computed(() => route.params.id as string)
const isLoading = ref(true)
const isSaving = ref(false)
const lastSaved = ref<Date | null>(null)
const lastSavedChar = ref<Partial<V5Character>>()

const editor = useCharacterEditor()

const stepComponents = shallowRef({
  initial: InitialSetupStep,
  clan: ClanSelectionStep,
  predator: PredatorTypeStep,
  attributes: AttributesStep,
  skills: SkillsStep,
  disciplines: DisciplinesStep,
  traits: TraitsStep,
  optional: OptionalStep,
})

onMounted(async () => {
  if (characterId.value) {
    await store.fetchCharacter(characterId.value)

    if (store.currentCharacter) {
      editor.character.value = store.currentCharacter
      lastSavedChar.value = { ...store.currentCharacter }
    }
  }
  await nextTick()
  isLoading.value = false
})

let saveTimeout: ReturnType<typeof setTimeout> | null = null

watch(editor.character, async () => {
  if (!characterId.value || isLoading.value) return

  if (saveTimeout) {
    clearTimeout(saveTimeout)
  }

  const changes: any = {};
  for (const key in editor.character.value) {
    if (JSON.stringify(editor.character.value[key as keyof Partial<V5Character>]) !== JSON.stringify(lastSavedChar.value?.[key as keyof Partial<V5Character>])) {
      changes[key] = editor.character.value[key as keyof Partial<V5Character>];
    }
  }

  saveTimeout = setTimeout(async () => {
    isSaving.value = true
    try {
      await store.updateCharacter(characterId.value, changes)
      lastSaved.value = new Date()
      lastSavedChar.value = { ...editor.character.value }
      setTimeout(() => {
        lastSaved.value = null
      }, 5000)
    } finally {
      isSaving.value = false
    }
  }, 500)
}, { deep: true })

function nextStep() {
  editor.nextStep()
}

function previousStep() {
  editor.previousStep()
}

function finish() {
  router.push(`/characters/${characterId.value}`)
}

function formatTime(date: Date): string {
  return date.toLocaleTimeString('de-DE', { hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <div v-if="isLoading" class="editor editor--loading">
    <p>Lade Charakter...</p>
  </div>

  <div v-else class="editor">
    <div class="editor__header">
      <div class="header-left">
        <VButton variant="ghost" @click="router.push(`/characters/${characterId}`)">
          ← Zurück zur Ansicht
        </VButton>
        <h1>Character Editor</h1>
      </div>
      <div class="save-status">
        <span v-if="isSaving" class="saving">Speichert...</span>
        <span v-else-if="lastSaved" class="saved">Gespeichert um {{ formatTime(lastSaved) }}</span>
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
        v-model="editor.character.value"
      />
    </div>

    <div class="editor__actions">
      <VButton
        v-if="editor.canGoPrevious.value"
        variant="secondary"
        @click="previousStep"
      >
        Zurück
      </VButton>

      <VButton
        v-if="editor.currentStepIndex.value < editor.steps.length - 1"
        variant="primary"
        :disabled="!editor.canGoNext.value"
        @click="nextStep"
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

  &--loading {
    display: grid;
    place-items: center;
    min-height: 50vh;
    color: $text-1;
  }

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

.save-status {
  font-size: 0.85rem;
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, 0.04);

  .saving {
    color: $text-2;
  }

  .saved {
    color: rgba(100, 200, 100, 0.8);
  }
}

@media (max-width: 420px) {
  .editor {
    padding: $s-4;

    &__header {
      flex-direction: column;
    }

    &__actions {
      flex-direction: column-reverse;
    }
  }
}
</style>
