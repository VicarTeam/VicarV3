<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'
import type { V5Character } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

onMounted(async () => {
  await Promise.all([
    v5data.fetchClans(),
    v5data.fetchDisciplines(),
  ])
})

const availableDisciplines = computed(() => {
  if (!character.value.clanID) return []
  const clan = v5data.getClanById(character.value.clanID)
  return clan?.disciplines ?? []
})

function getDisciplineLevel(disciplineId: string): number {
  const selection = character.value.disciplineSelections?.find(d => d.disciplineID === disciplineId)
  return selection?.currentLevel ?? 0
}

function setDisciplineLevel(disciplineId: string, level: number) {
  if (!character.value.disciplineSelections) {
    character.value.disciplineSelections = []
  }

  const existingIndex = character.value.disciplineSelections.findIndex(d => d.disciplineID === disciplineId)

  if (level === 0) {
    if (existingIndex !== -1) {
      character.value.disciplineSelections.splice(existingIndex, 1)
    }
    return
  }

  if (existingIndex !== -1) {
    character.value.disciplineSelections[existingIndex].currentLevel = level
    character.value.disciplineSelections[existingIndex].points = level
  } else {
    character.value.disciplineSelections.push({
      id: `disc-${disciplineId}`,
      characterID: '',
      disciplineID: disciplineId,
      points: level,
      currentLevel: level,
      abilities: [],
    })
  }
}

const totalPoints = computed(() => {
  return character.value.disciplineSelections?.reduce((sum, disc) => sum + disc.currentLevel, 0) ?? 0
})
</script>

<template>
  <VCard>
    <div class="header">
      <h2>Disziplinen auswählen</h2>
      <div class="points-display">
        Gesamt: {{ totalPoints }} Punkte
      </div>
    </div>

    <div class="step-content">
      <p v-if="!character.clanID" class="info-message">
        Bitte wähle zuerst einen Clan aus, um Disziplinen auswählen zu können.
      </p>

      <div v-else-if="availableDisciplines.length > 0" class="disciplines-list">
        <div
          v-for="discipline in availableDisciplines"
          :key="discipline.id"
          class="discipline-item"
        >
          <div class="discipline-header">
            <div class="discipline-info">
              <h3>{{ discipline.name }}</h3>
              <p v-if="discipline.summary" class="discipline-summary">{{ discipline.summary }}</p>
            </div>
            <VDotRating
              :model-value="getDisciplineLevel(discipline.id)"
              @update:model-value="(v) => setDisciplineLevel(discipline.id, v)"
            />
          </div>
        </div>
      </div>

      <p v-else class="info-message">
        Keine Disziplinen für den gewählten Clan verfügbar.
      </p>
    </div>
  </VCard>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $s-4;

  h2 {
    margin: 0;
    font-family: $font-head;
    font-size: $fs-2;
  }
}

.points-display {
  padding: $s-2 $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .04);
  font-weight: 600;
  color: $text-0;
}

.step-content {
  display: grid;
  gap: $s-4;
}

.info-message {
  text-align: center;
  padding: $s-6;
  color: $text-2;
  margin: 0;
}

.disciplines-list {
  display: grid;
  gap: $s-3;
}

.discipline-item {
  padding: $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
  border: 1px solid $border;
  transition: all $t-med $ease;

  &:hover {
    background: rgba(255, 255, 255, .04);
  }
}

.discipline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-4;
}

.discipline-info {
  flex: 1;
  min-width: 0;

  h3 {
    margin: 0 0 $s-1;
    font-family: $font-head;
    font-size: $fs-1;
    color: $text-0;
  }
}

.discipline-summary {
  margin: 0;
  font-size: 0.9rem;
  color: $text-2;
  line-height: 1.4;
}

@media (max-width: 420px) {
  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: $s-3;
  }

  .discipline-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
