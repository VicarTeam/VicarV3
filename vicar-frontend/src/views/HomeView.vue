<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import VCard from '@/components/ui/VCard.vue'
import VButton from '@/components/ui/VButton.vue'
import CharacterCreateDialog from '@/components/characters/CharacterCreateDialog.vue'
import {useToast} from "@/composables/useToast.ts";
import {migrateCharacter} from "@/rest/api/characters.ts";

const router = useRouter()
const store = useCharactersStore()
const toast = useToast()

const showCreateDialog = ref(false)
const importFiles = ref<HTMLInputElement|null>(null)

onMounted(() => {
  store.fetchCharacters()
})

function openCharacter(id: string) {
  router.push(`/characters/${id}`)
}

function onCreate(characterId: string) {
  showCreateDialog.value = false
  router.push(`/characters/${characterId}/edit`)
}

async function importChar(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files) {
    const file = input.files[0]!
    try {
      const text = await file.text()
      const charData = JSON.parse(text)
      const newChar = await migrateCharacter(charData)
      store.characters.push(newChar)
      toast.success('Charakter erfolgreich importiert.')
    } catch (e) {
      console.error(e)
      toast.error('Fehler beim Importieren des Charakters.')
    } finally {
      if (importFiles.value) {
        importFiles.value.value = ''
      }
    }
  }
}
</script>

<template>
  <div class="home">
    <div class="home__header">
      <h1>Meine Charaktere</h1>
      <div style="display: flex; gap: 0.5rem">
<!--        <VButton variant="ghost" @click="importFiles?.click()">
          Charakter v1 importieren
        </VButton>-->

        <VButton variant="primary" @click="showCreateDialog = true">
          Neuen Charakter erstellen
        </VButton>

        <input type="file" accept=".json" style="display: none" ref="importFiles" @change="importChar($event)" />
      </div>
    </div>

    <div v-if="store.loading" class="home__loading">
      Lade Charaktere...
    </div>

    <div v-else-if="store.characters.length === 0" class="home__empty">
      <p>Du hast noch keine Charaktere erstellt.</p>
      <VButton variant="primary" @click="showCreateDialog = true">
        Ersten Charakter erstellen
      </VButton>
    </div>

    <div v-else class="home__grid">
      <VCard
        v-for="char in store.characters"
        :key="char.id"
        interactive
        @click="openCharacter(char.id)"
      >
        <div class="char-card">
          <img v-if="char.avatar" :src="char.avatar" class="char-card__avatar" alt="" />
          <div v-else class="char-card__avatar char-card__avatar--placeholder">
            {{ char.name.charAt(0).toUpperCase() }}
          </div>
          <div class="char-card__info">
            <h3>{{ char.name }}</h3>
            <p class="muted">Generation {{ char.generation }}</p>
          </div>
        </div>
      </VCard>
    </div>

    <CharacterCreateDialog
      v-model="showCreateDialog"
      @created="onCreate"
    />
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.home {
  padding: $s-6;

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: $s-6;

    h1 {
      margin: 0;
      font-family: $font-head;
      font-size: $fs-3;
    }
  }

  &__grid {
    display: grid;
    gap: $s-4;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }

  &__empty, &__loading {
    text-align: center;
    padding: $s-8 $s-4;

    p {
      margin-bottom: $s-4;
      color: $text-1;
    }
  }
}

.char-card {
  display: flex;
  gap: $s-4;
  align-items: center;

  &__avatar {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    object-fit: cover;
    object-position: top;
    flex-shrink: 0;

    &--placeholder {
      background: $grad-accent;
      display: grid;
      place-items: center;
      font-size: 1.5rem;
      font-weight: 700;
      color: white;
    }
  }

  &__info {
    flex: 1;
    min-width: 0;

    h3 {
      margin: 0 0 $s-1;
      font-family: $font-head;
      font-size: $fs-1;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    p {
      margin: 0;
    }
  }
}

.muted {
  color: $text-2;
  font-size: 0.9rem;
}

@media (max-width: 420px) {
  .home {
    padding: $s-4;

    &__header {
      flex-direction: column;
      gap: $s-3;
      align-items: stretch;
    }

    &__grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
