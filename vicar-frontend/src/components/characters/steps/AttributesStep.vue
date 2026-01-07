<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'
import type { V5Character, CategoryKey, AttributeKey, V5CharacterAttribute } from '@/@types/v5'

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

const localAttributes = ref<V5CharacterAttribute[]>([])

function initializeAttributes() {
  const attrs: V5CharacterAttribute[] = []
  Object.entries(attributeDefinitions).forEach(([category, defs]) => {
    defs.forEach(def => {
      const existing = character.value?.attributes?.find(
        a => a.category === category && a.key === def.key
      )
      attrs.push({
        id: existing?.id ?? `${category}-${def.key}`,
        characterID: existing?.characterID ?? '',
        category: category as CategoryKey,
        key: def.key,
        value: existing?.value ?? 1,
      })
    })
  })
  localAttributes.value = attrs
}

onMounted(() => {
  initializeAttributes()
})


function getAttributeValue(category: CategoryKey, key: AttributeKey): number {
  const attr = localAttributes.value.find(a => a.category === category && a.key === key)
  return attr?.value ?? 1
}

function setAttributeValue(category: CategoryKey, key: AttributeKey, value: number) {
  const index = localAttributes.value.findIndex(a => a.category === category && a.key === key)
  if (index !== -1) {
    localAttributes.value[index] = { ...localAttributes.value[index], value }
    localAttributes.value = [...localAttributes.value]
  }
}

const distributionCounts = computed(() => {
  const counts = { 4: 0, 3: 0, 2: 0, 1: 0 }
  localAttributes.value.forEach(attr => {
    if (attr.value >= 1 && attr.value <= 4) {
      counts[attr.value as 1 | 2 | 3 | 4]++
    }
  })
  return counts
})

const distributionValid = computed(() => {
  return distributionCounts.value[4] === 1 &&
         distributionCounts.value[3] === 3 &&
         distributionCounts.value[2] === 4 &&
         distributionCounts.value[1] === 1
})

const categoryLabels: Record<CategoryKey, string> = {
  physical: 'Körperlich',
  social: 'Sozial',
  mental: 'Geistig',
}
</script>

<template>
  <div class="attributes-step">
    <VCard>
      <h2>Attribute verteilen</h2>
      <p class="distribution-info">
        Verteile deine Attribute: 1× auf 4, 3× auf 3, 4× auf 2, 1× auf 1
      </p>

      <div class="distribution-status" :class="{ 'distribution-status--valid': distributionValid }">
        <div class="dist-item">
          <span class="dist-label">4er:</span>
          <span class="dist-value" :class="{ 'dist-value--ok': distributionCounts[4] === 1, 'dist-value--over': distributionCounts[4] > 1 }">
            {{ distributionCounts[4] }}/1
          </span>
        </div>
        <div class="dist-item">
          <span class="dist-label">3er:</span>
          <span class="dist-value" :class="{ 'dist-value--ok': distributionCounts[3] === 3, 'dist-value--over': distributionCounts[3] > 3 }">
            {{ distributionCounts[3] }}/3
          </span>
        </div>
        <div class="dist-item">
          <span class="dist-label">2er:</span>
          <span class="dist-value" :class="{ 'dist-value--ok': distributionCounts[2] === 4, 'dist-value--over': distributionCounts[2] > 4 }">
            {{ distributionCounts[2] }}/4
          </span>
        </div>
        <div class="dist-item">
          <span class="dist-label">1er:</span>
          <span class="dist-value" :class="{ 'dist-value--ok': distributionCounts[1] === 1, 'dist-value--over': distributionCounts[1] > 1 }">
            {{ distributionCounts[1] }}/1
          </span>
        </div>
      </div>
    </VCard>

    <VCard>
      <div class="attributes-grid">
        <div v-for="category in (['physical', 'social', 'mental'] as CategoryKey[])" :key="category" class="attr-category">
          <h3>{{ categoryLabels[category] }}</h3>
          <div
            v-for="attr in attributeDefinitions[category]"
            :key="attr.key"
            class="attr-row"
          >
            <span class="attr-label">{{ attr.label }}</span>
            <VDotRating
              :model-value="getAttributeValue(category, attr.key)"
              :max="4"
              :min="1"
              @update:model-value="(v: number) => setAttributeValue(category, attr.key, v)"
            />
          </div>
        </div>
      </div>
    </VCard>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.attributes-step {
  display: grid;
  gap: $s-5;
}

h2 {
  margin: 0 0 $s-3;
  font-family: $font-head;
  font-size: $fs-2;
}

h3 {
  margin: 0 0 $s-3;
  font-family: $font-head;
  font-size: $fs-1;
  color: $text-1;
  padding-bottom: $s-2;
  border-bottom: 1px solid $border;
}

.distribution-info {
  margin: 0 0 $s-4;
  color: $text-2;
  font-size: 0.95rem;
}

.distribution-status {
  display: flex;
  gap: $s-4;
  flex-wrap: wrap;
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;

  &--valid {
    border-color: rgba(100, 200, 100, 0.4);
    background: rgba(100, 200, 100, 0.05);
  }
}

.dist-item {
  display: flex;
  align-items: center;
  gap: $s-2;
}

.dist-label {
  color: $text-2;
  font-weight: 500;
}

.dist-value {
  font-weight: 600;
  color: $text-1;

  &--ok {
    color: rgba(100, 200, 100, 0.9);
  }

  &--over {
    color: $red-0;
  }
}

.attributes-grid {
  display: grid;
  gap: $s-6;
  grid-template-columns: repeat(3, 1fr);
}

.attr-category {
  display: grid;
  gap: $s-3;
  min-width: 0;
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
  font-size: 0.95rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 768px) {
  .attributes-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 420px) {
  .distribution-status {
    flex-direction: column;
    gap: $s-2;
  }
}
</style>
