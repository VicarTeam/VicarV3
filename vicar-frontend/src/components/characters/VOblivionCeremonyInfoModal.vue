<script setup lang="ts">
import { computed } from "vue"
import VModal from "@/components/ui/VModal.vue"
import VCard from "@/components/ui/VCard.vue"
import VButton from "@/components/ui/VButton.vue"
import type { V5OblivionCeremony } from "@/@types/v5"

const open = defineModel<boolean>({ default: false })

const props = defineProps<{
  ceremony: V5OblivionCeremony | null
}>()

const hasHeader = computed(() => Boolean(props.ceremony?.name))
</script>

<template>
  <VModal v-model="open" size="md" :title="undefined">
    <VCard v-if="ceremony" class="ceremony" elevated>
      <div class="ceremony__header" v-if="hasHeader">
        <div class="ceremony__titleWrap">
          <h2 class="ceremony__title">{{ ceremony.name }}</h2>
          <span class="ceremony__pill">Level {{ ceremony.level }}</span>
        </div>
        <VButton variant="ghost" @click="open = false">Schließen</VButton>
      </div>

      <div class="ceremony__body">
        <small v-if="ceremony.summary" class="ceremony__summary"><i>{{ ceremony.summary }}</i></small>

        <div class="ceremony__grid">
          <div class="ceremony__row">
            <b>Level:</b>
            <span>{{ ceremony.level }}</span>
          </div>

          <div class="ceremony__row">
            <b>Buch:</b>
            <span>{{ ceremony.bookID || "—" }}</span>
          </div>

          <div class="ceremony__row" v-if="ceremony.cost">
            <b>Kosten:</b>
            <span>{{ ceremony.cost }}</span>
          </div>

          <div class="ceremony__row" v-if="ceremony.roll">
            <b>Wurf:</b>
            <span>{{ ceremony.roll }}</span>
          </div>

          <div class="ceremony__row" v-if="typeof ceremony.requires === 'number'">
            <b>Voraussetzung:</b>
            <span>{{ ceremony.requires }}</span>
          </div>

          <div class="ceremony__row" v-if="ceremony.duration">
            <b>Dauer:</b>
            <span>{{ ceremony.duration }}</span>
          </div>

          <div class="ceremony__row" v-if="ceremony.cult">
            <b>Kult:</b>
            <span>{{ ceremony.cult }}</span>
          </div>
        </div>

        <div class="ceremony__block" v-if="ceremony.ingredients">
          <b>Zutaten:</b>
          <div class="ceremony__text" v-html="ceremony.ingredients" />
        </div>

        <div class="ceremony__block" v-if="ceremony.execution">
          <b>Ausführung:</b>
          <div class="ceremony__text" v-html="ceremony.execution" />
        </div>

        <div class="ceremony__block" v-if="ceremony.system">
          <b>System:</b>
          <div class="ceremony__text" v-html="ceremony.system" />
        </div>
      </div>
    </VCard>

    <VCard v-else class="ceremony" elevated>
      <div class="ceremony__header">
        <h2 class="ceremony__title">Zeremonie</h2>
        <VButton variant="ghost" @click="open = false">Schließen</VButton>
      </div>
      <div class="ceremony__body">
        <p class="ceremony__empty">Keine Details verfügbar.</p>
      </div>
    </VCard>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.ceremony {
  padding: $s-5;
  max-height: min(75vh, 720px);
  overflow: auto;
}

.ceremony__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: $s-3;
  margin-bottom: $s-3;
}

.ceremony__titleWrap {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: $s-2;
  min-width: 0;
}

.ceremony__title {
  margin: 0;
  font-size: 1.25rem;
  min-width: 0;
}

.ceremony__pill {
  font-size: 0.75rem;
  padding: 2px $s-2;
  border-radius: 999px;
  border: 1px solid $border;
  background: rgba(255, 255, 255, 0.03);
  color: $text-2;
  flex-shrink: 0;
}

.ceremony__body {
  display: grid;
  gap: $s-3;
}

.ceremony__summary {
  color: $text-2;
}

.ceremony__grid {
  display: grid;
  gap: $s-2;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid rgba(255, 255, 255, 0.10);
  background: rgba(255, 255, 255, 0.03);
}

.ceremony__row {
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

.ceremony__block {
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

.ceremony__text {
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

.ceremony__empty {
  margin: 0;
  color: $text-2;
}

@media (max-width: 420px) {
  .ceremony {
    padding: $s-4;
  }

  .ceremony__row {
    grid-template-columns: 1fr;
  }
}
</style>
