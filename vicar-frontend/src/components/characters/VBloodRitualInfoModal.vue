<script setup lang="ts">
import { computed } from "vue"
import VModal from "@/components/ui/VModal.vue"
import VCard from "@/components/ui/VCard.vue"
import VButton from "@/components/ui/VButton.vue"
import type { V5BloodRitual } from "@/@types/v5"

const open = defineModel<boolean>({ default: false })

const props = defineProps<{
  ritual: V5BloodRitual | null
}>()

const hasHeader = computed(() => Boolean(props.ritual?.name))
</script>

<template>
  <VModal v-model="open" size="md" :title="undefined">
    <VCard v-if="ritual" class="ritual" elevated>
      <div class="ritual__header" v-if="hasHeader">
        <div class="ritual__titleWrap">
          <h2 class="ritual__title">{{ ritual.name }}</h2>
          <span class="ritual__pill">Level {{ ritual.level }}</span>
        </div>
        <VButton variant="ghost" @click="open = false">Schließen</VButton>
      </div>

      <div class="ritual__body">
        <div class="ritual__grid">
          <div class="ritual__row">
            <b>Level:</b>
            <span>{{ ritual.level }}</span>
          </div>

          <div class="ritual__row">
            <b>Buch:</b>
            <span>{{ ritual.bookID || "—" }}</span>
          </div>
        </div>

        <div class="ritual__block" v-if="ritual.description">
          <b>Beschreibung:</b>
          <div class="ritual__text" v-html="ritual.description" />
        </div>

        <div class="ritual__block" v-if="ritual.ingredients">
          <b>Zutaten:</b>
          <div class="ritual__text" v-html="ritual.ingredients" />
        </div>

        <div class="ritual__block" v-if="ritual.execution">
          <b>Ausführung:</b>
          <div class="ritual__text" v-html="ritual.execution" />
        </div>

        <div class="ritual__block" v-if="ritual.system">
          <b>System:</b>
          <div class="ritual__text" v-html="ritual.system" />
        </div>
      </div>
    </VCard>

    <VCard v-else class="ritual" elevated>
      <div class="ritual__header">
        <h2 class="ritual__title">Blutritual</h2>
        <VButton variant="ghost" @click="open = false">Schließen</VButton>
      </div>
      <div class="ritual__body">
        <p class="ritual__empty">Keine Details verfügbar.</p>
      </div>
    </VCard>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.ritual {
  padding: $s-5;
  max-height: min(75vh, 720px);
  overflow: auto;
}

.ritual__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: $s-3;
  margin-bottom: $s-3;
}

.ritual__titleWrap {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: $s-2;
  min-width: 0;
}

.ritual__title {
  margin: 0;
  font-size: 1.25rem;
  min-width: 0;
}

.ritual__pill {
  font-size: 0.75rem;
  padding: 2px $s-2;
  border-radius: 999px;
  border: 1px solid $border;
  background: rgba(255, 255, 255, 0.03);
  color: $text-2;
  flex-shrink: 0;
}

.ritual__body {
  display: grid;
  gap: $s-3;
}

.ritual__grid {
  display: grid;
  gap: $s-2;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid rgba(255, 255, 255, 0.10);
  background: rgba(255, 255, 255, 0.03);
}

.ritual__row {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: $s-2;
  align-items: baseline;

  b {
    color: $text-1;
  }

  span {
    color: $text-0;
    word-break: break-word;
  }
}

.ritual__block {
  display: grid;
  gap: $s-2;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid rgba(255, 255, 255, 0.10);
  background: rgba(255, 255, 255, 0.03);

  b {
    color: $text-1;
  }
}

.ritual__text {
  color: $text-1;

  :deep(p) {
    margin: 0 0 $s-2;
  }

  :deep(ul),
  :deep(ol) {
    margin: 0;
    padding-left: 1.2rem;
  }

  :deep(li) {
    margin: 0.15rem 0;
  }
}

.ritual__empty {
  margin: 0;
  color: $text-2;
}

@media (max-width: 420px) {
  .ritual {
    padding: $s-4;
  }

  .ritual__row {
    grid-template-columns: 1fr;
  }
}
</style>
