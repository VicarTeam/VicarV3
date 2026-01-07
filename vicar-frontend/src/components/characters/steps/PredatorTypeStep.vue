<script setup lang="ts">
import { onMounted } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import type { V5Character } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

onMounted(() => {
  v5data.fetchPredatorTypes()
})

function selectPredatorType(predatorTypeId: string) {
  character.value.predatorTypeID = predatorTypeId
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
            <span class="label">Aktionen:</span>
            <ul class="actions-list">
              <li v-for="action in predator.actions" :key="action.id">
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

.empty-state {
  text-align: center;
  padding: $s-6;
  color: $text-2;
}

@media (max-width: 420px) {
  .predator-grid {
    grid-template-columns: 1fr;
  }
}
</style>
