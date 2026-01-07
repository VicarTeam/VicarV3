<script setup lang="ts">
import { onMounted } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import type { V5Character } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

onMounted(() => {
  v5data.fetchClans()
})

function selectClan(clanId: string) {
  character.value.clanID = clanId
}
</script>

<template>
  <VCard>
    <h2>Clan auswählen</h2>

    <div class="step-content">
      <div class="clans-grid">
        <div
          v-for="clan in v5data.clans"
          :key="clan.id"
          class="clan-card"
          :class="{ 'clan-card--selected': character.clanID === clan.id }"
          @click="selectClan(clan.id)"
        >
          <div class="clan-card__header">
            <h3>{{ clan.name }}</h3>
            <p class="clan-slogan">{{ clan.slogan }}</p>
          </div>

          <div class="clan-card__disciplines">
            <span class="label">Disziplinen:</span>
            <span class="disciplines-list">
              {{ clan.disciplines.map(d => d.name).join(', ') }}
            </span>
          </div>

          <div class="clan-card__curse">
            <span class="label">Fluch:</span>
            <p>{{ clan.curse }}</p>
          </div>
        </div>
      </div>

      <p v-if="!v5data.clans || v5data.clans.length === 0" class="empty-state">
        Keine Clans verfügbar. Stelle sicher, dass du Bücher ausgewählt hast.
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

.clans-grid {
  display: grid;
  gap: $s-4;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
}

.clan-card {
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

  &__header {
    margin-bottom: $s-3;
    padding-bottom: $s-3;
    border-bottom: 1px solid $border;

    h3 {
      margin: 0 0 $s-1;
      font-family: $font-head;
      font-size: $fs-2;
    }
  }

  &__disciplines {
    margin-bottom: $s-3;
  }

  &__curse {
    p {
      margin: $s-1 0 0;
      font-size: 0.9rem;
      color: $text-2;
    }
  }
}

.clan-slogan {
  margin: 0;
  font-style: italic;
  color: $text-1;
  font-size: 0.95rem;
}

.label {
  font-weight: 600;
  color: $text-1;
  font-size: 0.9rem;
  display: block;
  margin-bottom: $s-1;
}

.disciplines-list {
  color: $text-0;
  font-size: 0.95rem;
}

.empty-state {
  text-align: center;
  padding: $s-6;
  color: $text-2;
}

@media (max-width: 420px) {
  .clans-grid {
    grid-template-columns: 1fr;
  }
}
</style>
