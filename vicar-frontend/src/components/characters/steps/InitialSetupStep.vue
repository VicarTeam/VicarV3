<script setup lang="ts">
import { onMounted } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VInput from '@/components/ui/VInput.vue'
import VCard from '@/components/ui/VCard.vue'
import type { V5Character, GenerationEra } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

const generationEras: { value: GenerationEra; label: string }[] = [
  { value: 'children', label: 'Kinder' },
  { value: 'newborn', label: 'Neugeborene' },
  { value: 'ancillae', label: 'Ancillae' },
  { value: 'older', label: 'Ältere' },
  { value: 'elder', label: 'Älteste' },
  { value: 'cainesinheritance', label: 'Kains Erbe' },
]

onMounted(() => {
  v5data.fetchBooks()
})

function toggleBook(bookId: string) {
  if (!character.value.books) {
    character.value.books = []
  }
  const index = character.value.books.findIndex(b => b.id === bookId)
  if (index === -1) {
    const book = v5data.books?.find(b => b.id === bookId)
    if (book) {
      character.value.books.push(book)
    }
  } else {
    character.value.books.splice(index, 1)
  }
}

function isBookSelected(bookId: string): boolean {
  return character.value.books?.some(b => b.id === bookId) ?? false
}
</script>

<template>
  <VCard>
    <h2>Grundeinstellungen</h2>

    <div class="step-content">
      <VInput
        v-model="character.name"
        label="Name *"
        placeholder="Charaktername"
      />

      <div class="field">
        <label class="field__label">Bücher auswählen *</label>
        <div class="books-grid">
          <button
            v-for="book in v5data.books"
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

      <div class="field">
        <label class="field__label">Generation Ära *</label>
        <select v-model="character.generationEra" class="select">
          <option v-for="era in generationEras" :key="era.value" :value="era.value">
            {{ era.label }}
          </option>
        </select>
      </div>

      <VInput
        v-model="character.generation"
        type="number"
        label="Generation *"
        :min="4"
        :max="16"
      />

      <VInput
        v-model="character.sire"
        label="Sire"
        placeholder="Name des Sires (optional)"
      />
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

.field {
  display: grid;
  gap: 0.4rem;

  &__label {
    font-size: 0.95rem;
    color: rgba(255, 255, 255, .72);
  }
}

.books-grid {
  display: grid;
  gap: $s-2;
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
</style>
