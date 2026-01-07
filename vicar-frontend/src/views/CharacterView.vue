<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useV5DataStore } from '@/stores/v5data'
import VButton from '@/components/ui/VButton.vue'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'

const route = useRoute()
const router = useRouter()
const store = useCharactersStore()
const v5data = useV5DataStore()

const characterId = computed(() => route.params.id as string)
const character = computed(() => store.currentCharacter)
const isEditMode = computed(() => route.name === 'character-edit')

const clan = computed(() => {
  if (!character.value?.clanID) return null
  return v5data.getClanById(character.value.clanID)
})

const predatorType = computed(() => {
  if (!character.value?.predatorTypeID) return null
  return v5data.getPredatorTypeById(character.value.predatorTypeID)
})

onMounted(async () => {
  await Promise.all([
    v5data.fetchClans(),
    v5data.fetchPredatorTypes(),
  ])

  if (characterId.value) {
    store.fetchCharacter(characterId.value)
  }
})

function toggleEdit() {
  if (isEditMode.value) {
    router.push(`/characters/${characterId.value}`)
  } else {
    router.push(`/characters/${characterId.value}/edit`)
  }
}
</script>

<template>
  <div class="character-view" v-if="character">
    <div class="character-view__header">
      <div class="header-left">
        <VButton variant="ghost" @click="router.push('/')">
          ← Zurück
        </VButton>
        <h1>{{ character.name }}</h1>
      </div>
      <VButton variant="primary" @click="toggleEdit">
        {{ isEditMode ? 'Ansicht' : 'Bearbeiten' }}
      </VButton>
    </div>

    <div class="character-view__content">
      <VCard>
        <h2>Grundinformationen</h2>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">Clan:</span>
            <span>{{ clan?.name || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">Jagdverhalten:</span>
            <span>{{ predatorType?.name || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">Generation:</span>
            <span>{{ character.generation }}</span>
          </div>
          <div class="info-item">
            <span class="label">Hunger:</span>
            <span>{{ character.hunger }}</span>
          </div>
          <div class="info-item">
            <span class="label">Menschlichkeit:</span>
            <span>{{ character.humanity }}</span>
          </div>
          <div class="info-item">
            <span class="label">Blutpotenz:</span>
            <span>{{ character.bloodPotency }}</span>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.attributes.length > 0">
        <h2>Attribute</h2>
        <div class="attributes">
          <div class="attr-category">
            <h3>Körperlich</h3>
            <div v-for="attr in character.attributes.filter(a => a.category === 'physical')" :key="attr.id" class="attr-row">
              <span class="attr-name">{{ attr.key.toUpperCase() }}</span>
              <VDotRating :model-value="attr.value" :readonly="true" />
            </div>
          </div>
          <div class="attr-category">
            <h3>Sozial</h3>
            <div v-for="attr in character.attributes.filter(a => a.category === 'social')" :key="attr.id" class="attr-row">
              <span class="attr-name">{{ attr.key.toUpperCase() }}</span>
              <VDotRating :model-value="attr.value" :readonly="true" />
            </div>
          </div>
          <div class="attr-category">
            <h3>Geistig</h3>
            <div v-for="attr in character.attributes.filter(a => a.category === 'mental')" :key="attr.id" class="attr-row">
              <span class="attr-name">{{ attr.key.toUpperCase() }}</span>
              <VDotRating :model-value="attr.value" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.skills.length > 0">
        <h2>Fähigkeiten</h2>
        <div class="skills">
          <div class="skill-category">
            <h3>Körperlich</h3>
            <div v-for="skill in character.skills.filter(s => s.category === 'physical')" :key="skill.id" class="skill-row">
              <div class="skill-info">
                <span class="skill-name">{{ skill.key.toUpperCase() }}</span>
                <span v-if="skill.specialization.length > 0" class="skill-spec">{{ skill.specialization.join(', ') }}</span>
              </div>
              <VDotRating :model-value="skill.value" :readonly="true" />
            </div>
          </div>
          <div class="skill-category">
            <h3>Sozial</h3>
            <div v-for="skill in character.skills.filter(s => s.category === 'social')" :key="skill.id" class="skill-row">
              <div class="skill-info">
                <span class="skill-name">{{ skill.key.toUpperCase() }}</span>
                <span v-if="skill.specialization.length > 0" class="skill-spec">{{ skill.specialization.join(', ') }}</span>
              </div>
              <VDotRating :model-value="skill.value" :readonly="true" />
            </div>
          </div>
          <div class="skill-category">
            <h3>Geistig</h3>
            <div v-for="skill in character.skills.filter(s => s.category === 'mental')" :key="skill.id" class="skill-row">
              <div class="skill-info">
                <span class="skill-name">{{ skill.key.toUpperCase() }}</span>
                <span v-if="skill.specialization.length > 0" class="skill-spec">{{ skill.specialization.join(', ') }}</span>
              </div>
              <VDotRating :model-value="skill.value" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.disciplineSelections.length > 0">
        <h2>Disziplinen</h2>
        <div class="disciplines">
          <div v-for="disc in character.disciplineSelections" :key="disc.id" class="disc-item">
            <div class="disc-header">
              <span class="disc-name">{{ v5data.getDisciplineById(disc.disciplineID)?.name || 'Unbekannt' }}</span>
              <VDotRating :model-value="disc.currentLevel" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>
    </div>
  </div>

  <div v-else class="character-view character-view--loading">
    <p>Lade Charakter...</p>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.character-view {
  padding: $s-6;
  max-width: 1120px;
  margin: 0 auto;

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
    display: grid;
    gap: $s-5;

    h2 {
      margin: 0 0 $s-4;
      font-family: $font-head;
      font-size: $fs-2;
    }

    h3 {
      margin: 0 0 $s-3;
      font-family: $font-head;
      font-size: $fs-1;
      color: $text-1;
    }
  }

  &--loading {
    display: grid;
    place-items: center;
    min-height: 50vh;
    color: $text-1;
  }
}

.info-grid {
  display: grid;
  gap: $s-3;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: $s-2 0;
  border-bottom: 1px solid $border;

  &:last-child {
    border-bottom: none;
  }

  .label {
    color: $text-2;
  }
}

.attributes, .skills {
  display: grid;
  gap: $s-5;
  grid-template-columns: repeat(3, 1fr);
}

.attr-category, .skill-category {
  display: grid;
  gap: $s-3;
}

.attr-row, .skill-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-3;
}

.attr-name, .skill-name {
  font-weight: 600;
  color: $text-0;
  min-width: 80px;
}

.skill-info {
  display: flex;
  flex-direction: column;
  gap: $s-1;
  flex: 1;
}

.skill-spec {
  font-size: 0.85rem;
  color: $text-2;
}

.disciplines {
  display: grid;
  gap: $s-3;
}

.disc-item {
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
}

.disc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.disc-name {
  font-weight: 600;
  color: $text-0;
}

@media (max-width: 420px) {
  .character-view {
    padding: $s-4;

    &__header {
      flex-direction: column;
      align-items: stretch;
    }
  }

  .attributes, .skills {
    grid-template-columns: 1fr;
  }
}
</style>
