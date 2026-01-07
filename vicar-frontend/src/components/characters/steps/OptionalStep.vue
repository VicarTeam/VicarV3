<script setup lang="ts">
import { onMounted } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import type { V5Character } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

onMounted(async () => {
  await Promise.all([
    v5data.fetchBloodRituals(),
    v5data.fetchOblivionCeremonies(),
  ])
})

function toggleBloodRitual(ritualId: string) {
  if (!character.value.bloodRituals) {
    character.value.bloodRituals = []
  }

  const index = character.value.bloodRituals.findIndex(r => r.id === ritualId)
  if (index === -1) {
    const ritual = v5data.bloodRituals?.find(r => r.id === ritualId)
    if (ritual) {
      character.value.bloodRituals.push(ritual)
    }
  } else {
    character.value.bloodRituals.splice(index, 1)
  }
}

function toggleOblivionCeremony(ceremonyId: string) {
  if (!character.value.oblivionCeremonies) {
    character.value.oblivionCeremonies = []
  }

  const index = character.value.oblivionCeremonies.findIndex(c => c.id === ceremonyId)
  if (index === -1) {
    const ceremony = v5data.oblivionCeremonies?.find(c => c.id === ceremonyId)
    if (ceremony) {
      character.value.oblivionCeremonies.push(ceremony)
    }
  } else {
    character.value.oblivionCeremonies.splice(index, 1)
  }
}

function isBloodRitualSelected(ritualId: string): boolean {
  return character.value.bloodRituals?.some(r => r.id === ritualId) ?? false
}

function isOblivionCeremonySelected(ceremonyId: string): boolean {
  return character.value.oblivionCeremonies?.some(c => c.id === ceremonyId) ?? false
}
</script>

<template>
  <div class="optional-step">
    <VCard>
      <h2>Blutrituale (Optional)</h2>

      <div class="step-content">
        <p class="info-text">
          Blutrituale sind optional und können später hinzugefügt werden.
        </p>

        <div v-if="v5data.bloodRituals && v5data.bloodRituals.length > 0" class="items-grid">
          <div
            v-for="ritual in v5data.bloodRituals"
            :key="ritual.id"
            class="item-card"
            :class="{ 'item-card--selected': isBloodRitualSelected(ritual.id) }"
            @click="toggleBloodRitual(ritual.id)"
          >
            <div class="item-header">
              <h3>{{ ritual.name }}</h3>
              <span class="item-level">Level {{ ritual.level }}</span>
            </div>
            <p class="item-description">{{ ritual.description }}</p>
          </div>
        </div>

        <p v-else class="empty-state">
          Keine Blutrituale verfügbar.
        </p>
      </div>
    </VCard>

    <VCard>
      <h2>Oblivion Zeremonien (Optional)</h2>

      <div class="step-content">
        <p class="info-text">
          Oblivion Zeremonien sind optional und können später hinzugefügt werden.
        </p>

        <div v-if="v5data.oblivionCeremonies && v5data.oblivionCeremonies.length > 0" class="items-grid">
          <div
            v-for="ceremony in v5data.oblivionCeremonies"
            :key="ceremony.id"
            class="item-card"
            :class="{ 'item-card--selected': isOblivionCeremonySelected(ceremony.id) }"
            @click="toggleOblivionCeremony(ceremony.id)"
          >
            <div class="item-header">
              <h3>{{ ceremony.name }}</h3>
              <span class="item-level">Level {{ ceremony.level }}</span>
            </div>
            <p class="item-description">{{ ceremony.summary }}</p>
          </div>
        </div>

        <p v-else class="empty-state">
          Keine Oblivion Zeremonien verfügbar.
        </p>
      </div>
    </VCard>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.optional-step {
  display: grid;
  gap: $s-5;
}

h2 {
  margin: 0 0 $s-4;
  font-family: $font-head;
  font-size: $fs-2;
}

.step-content {
  display: grid;
  gap: $s-4;
}

.info-text {
  margin: 0;
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .04);
  color: $text-1;
  font-size: 0.95rem;
}

.items-grid {
  display: grid;
  gap: $s-3;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

.item-card {
  padding: $s-4;
  border-radius: $r-md;
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
    box-shadow: 0 0 0 2px rgba(255, 59, 84, .15);
  }
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: $s-3;
  margin-bottom: $s-2;
  padding-bottom: $s-2;
  border-bottom: 1px solid $border;

  h3 {
    margin: 0;
    font-family: $font-head;
    font-size: $fs-1;
    flex: 1;
  }
}

.item-level {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, .06);
  font-size: 0.85rem;
  font-weight: 600;
  color: $text-1;
  white-space: nowrap;
}

.item-description {
  margin: 0;
  font-size: 0.9rem;
  color: $text-2;
  line-height: 1.5;
}

.empty-state {
  text-align: center;
  padding: $s-6;
  color: $text-2;
  margin: 0;
}

@media (max-width: 420px) {
  .items-grid {
    grid-template-columns: 1fr;
  }
}
</style>
