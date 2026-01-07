<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import type { V5Character, GenerationEra, Sex } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()
const isLoadingBooks = ref(true)

const generationEras: { value: GenerationEra; label: string }[] = [
  { value: 'children', label: 'Kinder' },
  { value: 'newborn', label: 'Neugeborene' },
  { value: 'ancillae', label: 'Ancillae' },
  { value: 'older', label: 'Ältere' },
  { value: 'elder', label: 'Methusa/Antidiluvian' },
]

const sexOptions: { value: Sex; label: string }[] = [
  { value: 'm', label: 'Männlich' },
  { value: 'f', label: 'Weiblich' },
  { value: 'd', label: 'Divers' },
]

const availableBooks = computed(() => v5data.books ?? [])

onMounted(async () => {
  await v5data.fetchBooks()
  isLoadingBooks.value = false
})

function toggleBook(bookId: string) {
  if (!character.value.books) {
    character.value.books = []
  }

  const index = character.value.books.findIndex(b => b.id === bookId)
  if (index === -1) {
    const book = availableBooks.value.find(b => b.id === bookId)
    if (book) {
      character.value.books = [...character.value.books, book]
    }
  } else {
    character.value.books = character.value.books.filter(b => b.id !== bookId)
  }
}

function isBookSelected(bookId: string): boolean {
  if (!character.value.books) return false
  return character.value.books.some(b => b.id === bookId)
}
</script>

<template>
  <VCard>
    <h2>Grundeinstellungen</h2>

    <div class="step-content">
      <div class="field">
        <label class="field__label">Name *</label>
        <input
          type="text"
          class="input"
          v-model="character.name"
          placeholder="Charaktername"
        />
      </div>

      <div class="field">
        <label class="field__label">Geschlecht</label>
        <div class="radio-group">
          <button
            v-for="opt in sexOptions"
            :key="opt.value"
            type="button"
            class="radio-btn"
            :class="{ 'radio-btn--selected': character.sex === opt.value }"
            @click="character.sex = opt.value"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>

      <div class="field">
        <label class="field__label">Bücher auswählen *</label>
        <div v-if="isLoadingBooks" class="loading-text">Lade Bücher...</div>
        <div v-else-if="availableBooks.length === 0" class="empty-text">Keine Bücher verfügbar</div>
        <div v-else class="books-grid">
          <button
            v-for="book in availableBooks"
            :key="book.id"
            type="button"
            class="book-btn"
            :class="{ 'book-btn--selected': isBookSelected(book.id) }"
            @click="toggleBook(book.id)"
          >
            {{ book.name }}
          </button>
        </div>
      </div>

      <div class="row-2">
        <div class="field">
          <label class="field__label">Generation *</label>
          <input
            type="number"
            class="input"
            v-model.number="character.generation"
            :min="4"
            :max="16"
          />
        </div>

        <div class="field">
          <label class="field__label">Generation Ära *</label>
          <select v-model="character.generationEra" class="select">
            <option v-for="era in generationEras" :key="era.value" :value="era.value">
              {{ era.label }}
            </option>
          </select>
        </div>
      </div>

      <div class="field">
        <label class="field__label">Sire</label>
        <input
          type="text"
          class="input"
          v-model="character.sire"
          placeholder="Name des Sires"
        />
      </div>
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

.row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $s-4;
}

.field {
  display: grid;
  gap: 0.4rem;

  &__label {
    font-size: 0.95rem;
    color: rgba(255, 255, 255, .72);
  }
}

.input {
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 0.95rem;
  transition: all $t-med $ease;

  &:focus {
    outline: none;
    border-color: $red-0;
    box-shadow: 0 0 0 2px $focus-ring;
  }

  &:hover {
    border-color: $border-strong;
  }

  &::placeholder {
    color: $text-2;
  }
}

.radio-group {
  display: flex;
  gap: $s-2;
}

.radio-btn {
  flex: 1;
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-1;
  cursor: pointer;
  transition: all $t-med $ease;
  text-align: center;

  &:hover {
    border-color: $border-strong;
    background: rgba(255, 255, 255, .04);
  }

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
    color: $text-0;
  }
}

.books-grid {
  display: grid;
  gap: $s-2;
}

.loading-text, .empty-text {
  padding: $s-4;
  text-align: center;
  color: $text-2;
  background: rgba(255, 255, 255, .02);
  border-radius: $r-md;
}

.book-btn {
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-1;
  cursor: pointer;
  transition: all $t-med $ease;
  text-align: left;

  &:hover {
    border-color: $border-strong;
    background: rgba(255, 255, 255, .04);
  }

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
    color: $text-0;
  }
}

.select {
  padding: $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 0.95rem;
  transition: all $t-med $ease;

  &:focus {
    outline: none;
    border-color: $red-0;
    box-shadow: 0 0 0 2px $focus-ring;
  }

  &:hover {
    border-color: $border-strong;
  }

  option {
    background: $bg-1;
    color: $text-0;
  }
}

@media (max-width: 420px) {
  .row-2 {
    grid-template-columns: 1fr;
  }

  .radio-group {
    flex-direction: column;
  }
}
</style>
