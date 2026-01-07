<script setup lang="ts">
import { computed } from 'vue'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'
import type { V5Character, CategoryKey, AttributeKey } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })

const attributeDefinitions: Record<CategoryKey, { key: AttributeKey; label: string }[]> = {
  physical: [
    { key: 'str', label: 'Stärke' },
    { key: 'dex', label: 'Geschicklichkeit' },
    { key: 'sta', label: 'Ausdauer' },
  ],
  social: [
    { key: 'cha', label: 'Charisma' },
    { key: 'man', label: 'Manipulation' },
    { key: 'com', label: 'Erscheinung' },
  ],
  mental: [
    { key: 'int', label: 'Intelligenz' },
    { key: 'wit', label: 'Witz' },
    { key: 'res', label: 'Entschlossenheit' },
  ],
}

function initializeAttributes() {
  if (!character.value.attributes || character.value.attributes.length === 0) {
    character.value.attributes = []
    Object.entries(attributeDefinitions).forEach(([category, attrs]) => {
      attrs.forEach(attr => {
        character.value.attributes!.push({
          id: `${category}-${attr.key}`,
          characterID: '',
          category: category as CategoryKey,
          key: attr.key,
          value: 1,
        })
      })
    })
  }
}

initializeAttributes()

function getAttributeValue(category: CategoryKey, key: AttributeKey): number {
  const attr = character.value.attributes?.find(a => a.category === category && a.key === key)
  return attr?.value ?? 1
}

function setAttributeValue(category: CategoryKey, key: AttributeKey, value: number) {
  if (!character.value.attributes) {
    character.value.attributes = []
  }
  const attr = character.value.attributes.find(a => a.category === category && a.key === key)
  if (attr) {
    attr.value = value
  } else {
    character.value.attributes.push({
      id: `${category}-${key}`,
      characterID: '',
      category,
      key,
      value,
    })
  }
}

const totalPoints = computed(() => {
  return character.value.attributes?.reduce((sum, attr) => sum + attr.value, 0) ?? 0
})
</script>

<template>
  <VCard>
    <div class="header">
      <h2>Attribute verteilen</h2>
      <div class="points-display">
        Gesamt: {{ totalPoints }} Punkte
      </div>
    </div>

    <div class="step-content">
      <div class="attributes-grid">
        <div class="attr-category">
          <h3>Körperlich</h3>
          <div
            v-for="attr in attributeDefinitions.physical"
            :key="attr.key"
            class="attr-row"
          >
            <span class="attr-label">{{ attr.label }}</span>
            <VDotRating
              :model-value="getAttributeValue('physical', attr.key)"
              @update:model-value="(v) => setAttributeValue('physical', attr.key, v)"
            />
          </div>
        </div>

        <div class="attr-category">
          <h3>Sozial</h3>
          <div
            v-for="attr in attributeDefinitions.social"
            :key="attr.key"
            class="attr-row"
          >
            <span class="attr-label">{{ attr.label }}</span>
            <VDotRating
              :model-value="getAttributeValue('social', attr.key)"
              @update:model-value="(v) => setAttributeValue('social', attr.key, v)"
            />
          </div>
        </div>

        <div class="attr-category">
          <h3>Geistig</h3>
          <div
            v-for="attr in attributeDefinitions.mental"
            :key="attr.key"
            class="attr-row"
          >
            <span class="attr-label">{{ attr.label }}</span>
            <VDotRating
              :model-value="getAttributeValue('mental', attr.key)"
              @update:model-value="(v) => setAttributeValue('mental', attr.key, v)"
            />
          </div>
        </div>
      </div>
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

.attributes-grid {
  display: grid;
  gap: $s-5;
  grid-template-columns: repeat(3, 1fr);
}

.attr-category {
  display: grid;
  gap: $s-3;

  h3 {
    margin: 0 0 $s-2;
    padding-bottom: $s-2;
    font-family: $font-head;
    font-size: $fs-1;
    border-bottom: 1px solid $border;
    color: $text-1;
  }
}

.attr-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-3;
  padding: $s-2;
  border-radius: $r-sm;
  transition: background $t-fast $ease;

  &:hover {
    background: rgba(255, 255, 255, .02);
  }
}

.attr-label {
  font-weight: 500;
  color: $text-0;
  min-width: 120px;
}

@media (max-width: 420px) {
  .attributes-grid {
    grid-template-columns: 1fr;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: $s-3;
  }
}
</style>
