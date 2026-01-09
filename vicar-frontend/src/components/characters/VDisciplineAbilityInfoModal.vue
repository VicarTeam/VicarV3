<script setup lang="ts">
import { computed } from "vue"
import VModal from "@/components/ui/VModal.vue"
import VCard from "@/components/ui/VCard.vue"
import VButton from "@/components/ui/VButton.vue"
import type { V5DisciplineAbility } from "@/@types/v5"

const open = defineModel<boolean>({ default: false })

const props = defineProps<{
  ability: V5DisciplineAbility | null
}>()

const hasHeader = computed(() => Boolean(props.ability?.name))
</script>

<template>
  <VModal v-model="open" size="md" :title="undefined">
    <VCard v-if="ability" class="ability" elevated>
      <div class="ability__header" v-if="hasHeader">
        <h2 class="ability__title">{{ ability.name }}</h2>
        <VButton variant="ghost" @click="open = false">Schließen</VButton>
      </div>

      <div class="ability__body">
        <small v-if="ability.summary" class="ability__summary"><i>{{ ability.summary }}</i></small>

        <div class="ability__grid">
          <div class="ability__row">
            <b>Kosten:</b>
            <span>{{ ability.costs || "—" }}</span>
          </div>

          <div class="ability__row" v-if="ability.diceSupplies">
            <b>Würfelpool:</b>
            <span>{{ ability.diceSupplies }}</span>
          </div>

          <div class="ability__row">
            <b>Dauer:</b>
            <span>{{ ability.duration || "—" }}</span>
          </div>
        </div>

        <div class="ability__system" v-if="ability.system">
          <b>System:</b>
          <div class="ability__systemText" v-html="ability.system" />
        </div>

        <div class="ability__alts" v-if="ability.alternatives && ability.alternatives.length > 0">
          <b>Alternativen:</b>
          <span>{{ ability.alternatives.join(", ") }}</span>
        </div>
      </div>
    </VCard>

    <VCard v-else class="ability" elevated>
      <div class="ability__header">
        <h2 class="ability__title">Disziplin</h2>
        <VButton variant="ghost" @click="open = false">Schließen</VButton>
      </div>
      <div class="ability__body">
        <p class="ability__empty">Keine Details verfügbar.</p>
      </div>
    </VCard>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.ability {
  padding: $s-5;
  max-height: min(75vh, 720px);
  overflow: auto;
}

.ability__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: $s-3;
  margin-bottom: $s-3;
}

.ability__title {
  margin: 0;
  font-size: 1.25rem;
}

.ability__body {
  display: grid;
  gap: $s-3;
}

.ability__summary {
  color: $text-2;
}

.ability__rule {
  display: grid;
  gap: $s-1;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid rgba(255, 255, 255, 0.10);
  background: rgba(255, 255, 255, 0.03);

  b {
    color: $text-1;
  }

  span {
    color: $text-0;
  }
}

.ability__grid {
  display: grid;
  gap: $s-2;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid rgba(255, 255, 255, 0.10);
  background: rgba(255, 255, 255, 0.03);
}

.ability__row {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: $s-2;
  align-items: baseline;

  b {
    color: $text-1;
  }

  span {
    color: $text-0;
  }
}

.ability__system {
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

.ability__systemText {
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

.ability__alts {
  display: grid;
  gap: $s-1;
  color: $text-1;

  b {
    color: $text-1;
  }
}

.ability__empty {
  margin: 0;
  color: $text-2;
}

@media (max-width: 420px) {
  .ability {
    padding: $s-4;
  }

  .ability__row {
    grid-template-columns: 1fr;
  }
}
</style>
