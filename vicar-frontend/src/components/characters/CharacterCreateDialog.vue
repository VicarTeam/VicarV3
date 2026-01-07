<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import { useCharactersStore } from '@/stores/characters'
import { useToast } from '@/composables/useToast'
import VModal from '@/components/ui/VModal.vue'
import VInput from '@/components/ui/VInput.vue'
import VButton from '@/components/ui/VButton.vue'
import type { GenerationEra } from '@/@types/v5'

const open = defineModel<boolean>({ default: false })
const emit = defineEmits<{
  (e: 'created', id: string): void
}>()

const v5data = useV5DataStore()
const store = useCharactersStore()
const toast = useToast()

const name = ref('')
const selectedBooks = ref<string[]>([])
const generationEra = ref<GenerationEra>('newborn')
const generation = ref(13)
const loading = ref(false)

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
  const index = selectedBooks.value.indexOf(bookId)
  if (index === -1) {
    selectedBooks.value.push(bookId)
  } else {
    selectedBooks.value.splice(index, 1)
  }
}

async function create() {
  if (!name.value || selectedBooks.value.length === 0) {
    toast.error('Bitte fülle alle Pflichtfelder aus')
    return
  }

  loading.value = true
  try {
    const char = await store.createCharacter({
      name: name.value,
      books: selectedBooks.value,
      generationEra: generationEra.value,
      generation: generation.value,
    })

    if (char) {
      toast.success('Charakter erfolgreich erstellt')
      emit('created', char.id)
      name.value = ''
      selectedBooks.value = []
      generationEra.value = 'newborn'
      generation.value = 13
    } else {
      toast.error('Fehler beim Erstellen des Charakters')
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <VModal v-model="open" title="Neuen Charakter erstellen" size="md">
    <div class="create-dialog">
      <VInput
        v-model="name"
        label="Name"
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
            :class="{ 'book-btn--selected': selectedBooks.includes(book.id) }"
            @click="toggleBook(book.id)"
          >
            {{ book.name }}
          </button>
        </div>
      </div>

      <div class="field">
        <label class="field__label">Generation Ära</label>
        <select v-model="generationEra" class="select">
          <option v-for="era in generationEras" :key="era.value" :value="era.value">
            {{ era.label }}
          </option>
        </select>
      </div>

      <VInput
        v-model="generation"
        type="number"
        label="Generation"
        :min="4"
        :max="16"
      />
    </div>

    <template #footer="{ close }">
      <VButton variant="secondary" @click="close">
        Abbrechen
      </VButton>
      <VButton
        variant="primary"
        :loading="loading"
        :disabled="!name || selectedBooks.length === 0"
        @click="create"
      >
        Erstellen
      </VButton>
    </template>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.create-dialog {
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
