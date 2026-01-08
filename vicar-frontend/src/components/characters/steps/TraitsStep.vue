<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useV5DataStore } from '@/stores/v5data'
import VCard from '@/components/ui/VCard.vue'
import VModal from '@/components/ui/VModal.vue'
import VButton from '@/components/ui/VButton.vue'
import type {
  V5Character,
  V5TraitPack,
  V5Trait,
  V5CharacterTraitPackUsage,
  V5CharacterTrait,
  V5CharacterFlawTrait,
  TraitPackType
} from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

const isLoading = ref(true)
const expandedMeritPacks = ref<Set<string>>(new Set())
const expandedBackgroundPacks = ref<Set<string>>(new Set())
const expandedFlawPacks = ref<Set<string>>(new Set())

// For suffix/specialization input
const showSuffixModal = ref(false)
const currentSuffixTrait = ref<{ pack: V5TraitPack; trait: V5Trait; level: number } | null>(null)
const suffixInput = ref('')

// Base points
const BASE_MERIT_POINTS = 7
const MIN_FLAW_POINTS = 2

onMounted(async () => {
  await Promise.all([
    v5data.fetchTraitPacks(),
    v5data.fetchPredatorTypes(),
  ])
  isLoading.value = false
})

// Separate merit and background packs
const meritPacks = computed(() => {
  return v5data.traitPacks?.filter(p => p.type === 'merits') ?? []
})

const backgroundPacks = computed(() => {
  return v5data.traitPacks?.filter(p => p.type === 'backgrounds') ?? []
})

// Get predator type traits
const predatorTraits = computed(() => {
  if (!character.value.predatorTypeID) return { merits: [], flaws: [] }

  const predator = v5data.predatorTypes?.find(p => p.id === character.value.predatorTypeID)
  if (!predator) return { merits: [], flaws: [] }

  const merits: { action: any; description: string }[] = []
  const flaws: { action: any; description: string }[] = []

  predator.actions.forEach(action => {
    if (action.type === 'add_merit' || action.type === 'add_background') {
      merits.push({ action, description: action.description })
    } else if (action.type === 'add_flaw') {
      flaws.push({ action, description: action.description })
    }
  })

  return { merits, flaws }
})

// Calculate points
const totalMeritPoints = computed(() => {
  let total = 0
  character.value.traitPackUsages?.forEach(usage => {
    if (usage.kind === 'merits') {
      usage.traits?.forEach(t => {
        if (!t.isLocked) {
          total += t.customLevel ?? 0
        }
      })
    }
  })
  return total
})

const totalBackgroundPoints = computed(() => {
  let total = 0
  character.value.traitPackUsages?.forEach(usage => {
    if (usage.kind === 'backgrounds') {
      usage.traits?.forEach(t => {
        if (!t.isLocked) {
          total += t.customLevel ?? 0
        }
      })
    }
  })
  return total
})

const totalFlawPoints = computed(() => {
  let total = 0
  character.value.traitPackUsages?.forEach(usage => {
    usage.flawTraits?.forEach(ft => {
      const pack = v5data.traitPacks?.find(p => p.id === usage.packID)
      const trait = pack?.packTraits?.find(pt => pt.trait.id === ft.traitID)?.trait
      if (trait) {
        total += trait.level
      }
    })
  })
  return total
})

// Extra merit points from flaws above minimum
const bonusMeritPoints = computed(() => {
  return Math.max(0, totalFlawPoints.value - MIN_FLAW_POINTS)
})

// Total available merit points
const availableMeritPoints = computed(() => {
  return BASE_MERIT_POINTS + bonusMeritPoints.value
})

// Points remaining
const meritPointsRemaining = computed(() => {
  return availableMeritPoints.value - totalMeritPoints.value
})

const flawPointsRemaining = computed(() => {
  return MIN_FLAW_POINTS - totalFlawPoints.value
})

// Distribution valid
const distributionValid = computed(() => {
  return meritPointsRemaining.value >= 0 && flawPointsRemaining.value <= 0
})

// Toggle pack expansion
function toggleMeritPack(packId: string) {
  if (expandedMeritPacks.value.has(packId)) {
    expandedMeritPacks.value.delete(packId)
  } else {
    expandedMeritPacks.value.add(packId)
  }
  expandedMeritPacks.value = new Set(expandedMeritPacks.value)
}

function toggleBackgroundPack(packId: string) {
  if (expandedBackgroundPacks.value.has(packId)) {
    expandedBackgroundPacks.value.delete(packId)
  } else {
    expandedBackgroundPacks.value.add(packId)
  }
  expandedBackgroundPacks.value = new Set(expandedBackgroundPacks.value)
}

function toggleFlawPack(packId: string) {
  if (expandedFlawPacks.value.has(packId)) {
    expandedFlawPacks.value.delete(packId)
  } else {
    expandedFlawPacks.value.add(packId)
  }
  expandedFlawPacks.value = new Set(expandedFlawPacks.value)
}

// Helper functions
function getPackUsage(packId: string): V5CharacterTraitPackUsage | undefined {
  return character.value.traitPackUsages?.find(u => u.packID === packId)
}

function isTraitSelected(packId: string, traitId: string): boolean {
  const usage = getPackUsage(packId)
  return usage?.traits?.some(t => t.traitID === traitId) ?? false
}

function isFlawSelected(packId: string, traitId: string): boolean {
  const usage = getPackUsage(packId)
  return usage?.flawTraits?.some(t => t.traitID === traitId) ?? false
}

function getSelectedTraitLevel(packId: string, traitId: string): number | undefined {
  const usage = getPackUsage(packId)
  const trait = usage?.traits?.find(t => t.traitID === traitId)
  return trait?.customLevel
}

function getSelectedTraitSuffix(packId: string, traitId: string): string | undefined {
  const usage = getPackUsage(packId)
  const trait = usage?.traits?.find(t => t.traitID === traitId)
  return trait?.suffix
}

// Select/toggle trait (merit/background)
function selectTrait(pack: V5TraitPack, trait: V5Trait, level: number) {
  // Check if this is a repeatable trait that needs a suffix
  if (trait.isRepeatable) {
    currentSuffixTrait.value = { pack, trait, level }
    suffixInput.value = ''
    showSuffixModal.value = true
    return
  }

  applyTrait(pack, trait.id, level, undefined)
}

// Edit specialization for an existing trait
function editTraitSuffix(pack: V5TraitPack, trait: V5Trait) {
  const currentLevel = getSelectedTraitLevel(pack.id, trait.id)
  const currentSuffix = getSelectedTraitSuffix(pack.id, trait.id)
  currentSuffixTrait.value = { pack, trait, level: currentLevel ?? trait.level }
  suffixInput.value = currentSuffix || ''
  showSuffixModal.value = true
}

function applyTrait(pack: V5TraitPack, traitId: string, level: number, suffix?: string, isUpdate = false) {
  if (!character.value.traitPackUsages) {
    character.value.traitPackUsages = []
  }

  let usageIndex = character.value.traitPackUsages.findIndex(u => u.packID === pack.id)

  if (usageIndex === -1) {
    character.value.traitPackUsages.push({
      id: `usage-${pack.id}`,
      characterID: '',
      kind: pack.type,
      packID: pack.id,
      traits: [],
      flawTraits: [],
    })
    usageIndex = character.value.traitPackUsages.length - 1
  }

  const usage = character.value.traitPackUsages[usageIndex]

  // For updates (editing suffix), find by traitId only
  // For new repeatable traits, find by traitId AND suffix
  const traitIndex = isUpdate
    ? usage.traits.findIndex(t => t.traitID === traitId)
    : usage.traits.findIndex(t => t.traitID === traitId && t.suffix === suffix)

  if (traitIndex !== -1) {
    if (isUpdate) {
      // Update the suffix and/or level
      usage.traits[traitIndex] = {
        ...usage.traits[traitIndex],
        customLevel: level,
        suffix: suffix,
      }
    } else {
      const existingLevel = usage.traits[traitIndex].customLevel
      if (existingLevel === level) {
        // Deselect
        usage.traits.splice(traitIndex, 1)
      } else {
        // Update level
        usage.traits[traitIndex] = {
          ...usage.traits[traitIndex],
          customLevel: level,
        }
      }
    }
  } else {
    // Add new
    usage.traits.push({
      id: `trait-${traitId}-${Date.now()}`,
      usageID: usage.id,
      traitID: traitId,
      isLocked: false,
      isManual: true,
      customLevel: level,
      suffix: suffix,
    })
  }

  character.value.traitPackUsages = [...character.value.traitPackUsages]
}

function confirmSuffix() {
  if (!currentSuffixTrait.value) return

  const { pack, trait, level } = currentSuffixTrait.value
  // Check if this trait already exists (we're editing)
  const isUpdate = isTraitSelected(pack.id, trait.id)
  applyTrait(pack, trait.id, level, suffixInput.value.trim() || undefined, isUpdate)

  showSuffixModal.value = false
  currentSuffixTrait.value = null
  suffixInput.value = ''
}

// Select/toggle flaw
function toggleFlaw(pack: V5TraitPack, traitId: string) {
  if (!character.value.traitPackUsages) {
    character.value.traitPackUsages = []
  }

  let usageIndex = character.value.traitPackUsages.findIndex(u => u.packID === pack.id)

  if (usageIndex === -1) {
    character.value.traitPackUsages.push({
      id: `usage-${pack.id}`,
      characterID: '',
      kind: pack.type,
      packID: pack.id,
      traits: [],
      flawTraits: [],
    })
    usageIndex = character.value.traitPackUsages.length - 1
  }

  const usage = character.value.traitPackUsages[usageIndex]
  const flawIndex = usage.flawTraits.findIndex(t => t.traitID === traitId)

  if (flawIndex !== -1) {
    // Deselect
    usage.flawTraits.splice(flawIndex, 1)
  } else {
    // Add
    usage.flawTraits.push({
      id: `flaw-${traitId}`,
      usageID: usage.id,
      traitID: traitId,
    })
  }

  character.value.traitPackUsages = [...character.value.traitPackUsages]
}

function getTraitLevels(pack: V5TraitPack, traitId: string): number[] {
  const packTrait = pack.packTraits?.find(pt => pt.trait.id === traitId)
  if (!packTrait) return [1, 2, 3]

  const trait = packTrait.trait
  if (trait.isRepeatable && trait.repeatAmount && trait.repeatSize) {
    const levels: number[] = []
    for (let i = 1; i <= trait.repeatAmount; i++) {
      levels.push(i * trait.repeatSize)
    }
    return levels
  }

  return [trait.level]
}

// Get all flaws from packs
const allFlawTraits = computed(() => {
  const flaws: { pack: V5TraitPack; trait: V5Trait }[] = []

  v5data.traitPacks?.forEach(pack => {
    pack.packTraits?.forEach(pt => {
      // Check if trait is a flaw (level is negative or has flaw indicator)
      // For V5, flaws are typically in certain packs or have specific properties
      // We'll assume flaws have negative levels or specific pack rules
      if (pt.trait.level < 0) {
        flaws.push({ pack, trait: pt.trait })
      }
    })
  })

  return flaws
})

// Get flaws grouped by pack
const flawsByPack = computed(() => {
  const grouped: Map<string, { pack: V5TraitPack; flaws: V5Trait[] }> = new Map()

  v5data.traitPacks?.forEach(pack => {
    const packFlaws = pack.packTraits?.filter(pt => pt.trait.level < 0).map(pt => pt.trait) ?? []
    if (packFlaws.length > 0) {
      grouped.set(pack.id, { pack, flaws: packFlaws })
    }
  })

  return grouped
})
</script>

<template>
  <div class="traits-step">
    <!-- Points Overview -->
    <VCard>
      <h2>Punkte-Übersicht</h2>
      <p class="info-text">
        Du hast {{ BASE_MERIT_POINTS }} Punkte für Vorzüge. Mindestens {{ MIN_FLAW_POINTS }} Punkte musst du in Schwächen investieren.
        Für jeden zusätzlichen Schwächenpunkt erhältst du einen weiteren Vorzugspunkt.
      </p>

      <div class="points-overview" :class="{ 'points-overview--valid': distributionValid }">
        <div class="points-row">
          <span class="points-label">Vorzüge:</span>
          <span class="points-value" :class="{ 'points-value--ok': meritPointsRemaining >= 0, 'points-value--over': meritPointsRemaining < 0 }">
            {{ totalMeritPoints }} / {{ availableMeritPoints }}
          </span>
        </div>
        <div class="points-row">
          <span class="points-label">Schwächen:</span>
          <span class="points-value" :class="{ 'points-value--ok': flawPointsRemaining <= 0, 'points-value--under': flawPointsRemaining > 0 }">
            {{ totalFlawPoints }} / {{ MIN_FLAW_POINTS }}+
          </span>
        </div>
        <div v-if="bonusMeritPoints > 0" class="points-row points-row--bonus">
          <span class="points-label">Bonuspunkte:</span>
          <span class="points-value points-value--bonus">+{{ bonusMeritPoints }}</span>
        </div>
        <div class="points-row">
          <span class="points-label">Hintergründe:</span>
          <span class="points-value">{{ totalBackgroundPoints }}</span>
        </div>
      </div>
    </VCard>

    <!-- Predator Type Bonuses -->
    <VCard v-if="predatorTraits.merits.length > 0 || predatorTraits.flaws.length > 0">
      <h2>Jagdverhalten-Boni</h2>
      <p class="info-text">Diese Vorzüge und Schwächen erhältst du durch dein Jagdverhalten automatisch.</p>

      <div class="predator-traits">
        <div v-if="predatorTraits.merits.length > 0" class="predator-section">
          <h4>Vorzüge:</h4>
          <ul class="predator-list">
            <li v-for="(item, idx) in predatorTraits.merits" :key="`merit-${idx}`">
              {{ item.description }}
            </li>
          </ul>
        </div>
        <div v-if="predatorTraits.flaws.length > 0" class="predator-section">
          <h4>Schwächen:</h4>
          <ul class="predator-list">
            <li v-for="(item, idx) in predatorTraits.flaws" :key="`flaw-${idx}`">
              {{ item.description }}
            </li>
          </ul>
        </div>
      </div>
    </VCard>

    <!-- Merits -->
    <VCard>
      <h2>Vorzüge</h2>

      <div v-if="isLoading" class="loading-text">Lade Vorzüge...</div>
      <div v-else-if="meritPacks.length === 0" class="empty-text">Keine Vorzüge verfügbar</div>
      <div v-else class="packs-list">
        <div
          v-for="pack in meritPacks"
          :key="pack.id"
          class="pack-item"
          :class="{ 'pack-item--expanded': expandedMeritPacks.has(pack.id) }"
        >
          <button
            type="button"
            class="pack-header"
            @click="toggleMeritPack(pack.id)"
          >
            <div class="pack-info">
              <h3>{{ pack.name }}</h3>
              <p v-if="pack.description" class="pack-desc">{{ pack.description }}</p>
            </div>
            <span class="pack-toggle">{{ expandedMeritPacks.has(pack.id) ? '−' : '+' }}</span>
          </button>

          <div v-if="expandedMeritPacks.has(pack.id)" class="pack-traits">
            <div
              v-for="packTrait in pack.packTraits?.filter(pt => pt.trait.level > 0)"
              :key="packTrait.id"
              class="trait-item"
              :class="{ 'trait-item--selected': isTraitSelected(pack.id, packTrait.trait.id) }"
            >
              <div class="trait-info">
                <div class="trait-header">
                  <span class="trait-name">{{ packTrait.trait.name }}</span>
                  <span v-if="getSelectedTraitSuffix(pack.id, packTrait.trait.id)" class="trait-suffix">
                    ({{ getSelectedTraitSuffix(pack.id, packTrait.trait.id) }})
                  </span>
                </div>
                <p class="trait-desc">{{ packTrait.trait.description }}</p>
              </div>
              <div class="trait-actions">
                <button
                  v-if="isTraitSelected(pack.id, packTrait.trait.id)"
                  type="button"
                  class="suffix-edit-btn"
                  title="Spezialisierung bearbeiten"
                  @click.stop="editTraitSuffix(pack, packTrait.trait)"
                >
                  ✎
                </button>
                <div class="trait-levels">
                  <button
                    v-for="level in getTraitLevels(pack, packTrait.trait.id)"
                    :key="level"
                    type="button"
                    class="level-btn"
                    :class="{ 'level-btn--selected': getSelectedTraitLevel(pack.id, packTrait.trait.id) === level }"
                    @click="selectTrait(pack, packTrait.trait, level)"
                  >
                    {{ level }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </VCard>

    <!-- Flaws -->
    <VCard>
      <h2>Schwächen</h2>
      <p class="info-text required-hint">Mindestens {{ MIN_FLAW_POINTS }} Punkte erforderlich</p>

      <div v-if="isLoading" class="loading-text">Lade Schwächen...</div>
      <div v-else-if="flawsByPack.size === 0" class="empty-text">Keine Schwächen verfügbar</div>
      <div v-else class="packs-list">
        <div
          v-for="[packId, packData] in flawsByPack"
          :key="packId"
          class="pack-item"
          :class="{ 'pack-item--expanded': expandedFlawPacks.has(packId) }"
        >
          <button
            type="button"
            class="pack-header"
            @click="toggleFlawPack(packId)"
          >
            <div class="pack-info">
              <h3>{{ packData.pack.name }}</h3>
            </div>
            <span class="pack-toggle">{{ expandedFlawPacks.has(packId) ? '−' : '+' }}</span>
          </button>

          <div v-if="expandedFlawPacks.has(packId)" class="pack-traits">
            <div
              v-for="flaw in packData.flaws"
              :key="flaw.id"
              class="trait-item trait-item--flaw"
              :class="{ 'trait-item--selected': isFlawSelected(packId, flaw.id) }"
            >
              <div class="trait-info">
                <div class="trait-header">
                  <span class="trait-name">{{ flaw.name }}</span>
                  <span class="trait-level-badge trait-level-badge--flaw">{{ Math.abs(flaw.level) }}</span>
                </div>
                <p class="trait-desc">{{ flaw.description }}</p>
              </div>
              <button
                type="button"
                class="flaw-toggle-btn"
                :class="{ 'flaw-toggle-btn--selected': isFlawSelected(packId, flaw.id) }"
                @click="toggleFlaw(packData.pack, flaw.id)"
              >
                {{ isFlawSelected(packId, flaw.id) ? '✓' : '+' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </VCard>

    <!-- Backgrounds -->
    <VCard>
      <h2>Hintergründe</h2>

      <div v-if="isLoading" class="loading-text">Lade Hintergründe...</div>
      <div v-else-if="backgroundPacks.length === 0" class="empty-text">Keine Hintergründe verfügbar</div>
      <div v-else class="packs-list">
        <div
          v-for="pack in backgroundPacks"
          :key="pack.id"
          class="pack-item"
          :class="{ 'pack-item--expanded': expandedBackgroundPacks.has(pack.id) }"
        >
          <button
            type="button"
            class="pack-header"
            @click="toggleBackgroundPack(pack.id)"
          >
            <div class="pack-info">
              <h3>{{ pack.name }}</h3>
              <p v-if="pack.description" class="pack-desc">{{ pack.description }}</p>
            </div>
            <span class="pack-toggle">{{ expandedBackgroundPacks.has(pack.id) ? '−' : '+' }}</span>
          </button>

          <div v-if="expandedBackgroundPacks.has(pack.id)" class="pack-traits">
            <div
              v-for="packTrait in pack.packTraits?.filter(pt => pt.trait.level > 0)"
              :key="packTrait.id"
              class="trait-item"
              :class="{ 'trait-item--selected': isTraitSelected(pack.id, packTrait.trait.id) }"
            >
              <div class="trait-info">
                <div class="trait-header">
                  <span class="trait-name">{{ packTrait.trait.name }}</span>
                  <span v-if="getSelectedTraitSuffix(pack.id, packTrait.trait.id)" class="trait-suffix">
                    ({{ getSelectedTraitSuffix(pack.id, packTrait.trait.id) }})
                  </span>
                </div>
                <p class="trait-desc">{{ packTrait.trait.description }}</p>
              </div>
              <div class="trait-actions">
                <button
                  v-if="isTraitSelected(pack.id, packTrait.trait.id)"
                  type="button"
                  class="suffix-edit-btn"
                  title="Spezialisierung bearbeiten"
                  @click.stop="editTraitSuffix(pack, packTrait.trait)"
                >
                  ✎
                </button>
                <div class="trait-levels">
                  <button
                    v-for="level in getTraitLevels(pack, packTrait.trait.id)"
                    :key="level"
                    type="button"
                    class="level-btn"
                    :class="{ 'level-btn--selected': getSelectedTraitLevel(pack.id, packTrait.trait.id) === level }"
                    @click="selectTrait(pack, packTrait.trait, level)"
                  >
                    {{ level }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </VCard>

    <!-- Suffix Modal -->
    <VModal v-model="showSuffixModal" title="Spezialisierung angeben" size="sm">
      <div v-if="currentSuffixTrait" class="suffix-modal">
        <p>{{ currentSuffixTrait.trait.name }} (Level {{ currentSuffixTrait.level }})</p>
        <p class="suffix-hint">Optional: Gib eine Spezialisierung an (z.B. Kontakt-Name, Ort, etc.)</p>
        <input
          type="text"
          class="suffix-input"
          v-model="suffixInput"
          placeholder="Spezialisierung..."
          @keydown.enter="confirmSuffix"
        />
      </div>
      <template #footer>
        <div class="modal-actions">
          <VButton variant="secondary" @click="showSuffixModal = false">Abbrechen</VButton>
          <VButton variant="primary" @click="confirmSuffix">Bestätigen</VButton>
        </div>
      </template>
    </VModal>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.traits-step {
  display: grid;
  gap: $s-5;
}

h2 {
  margin: 0 0 $s-3;
  font-family: $font-head;
  font-size: $fs-2;
}

h4 {
  margin: 0 0 $s-2;
  font-family: $font-head;
  font-size: 0.95rem;
  color: $text-1;
}

.info-text {
  margin: 0 0 $s-4;
  color: $text-2;
  font-size: 0.9rem;
  line-height: 1.5;

  &.required-hint {
    color: $red-0;
    font-weight: 500;
  }
}

.points-overview {
  display: grid;
  gap: $s-2;
  padding: $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;

  &--valid {
    border-color: rgba(100, 200, 100, 0.4);
    background: rgba(100, 200, 100, 0.05);
  }
}

.points-row {
  display: flex;
  justify-content: space-between;
  align-items: center;

  &--bonus {
    padding-top: $s-2;
    border-top: 1px dashed $border;
  }
}

.points-label {
  color: $text-2;
  font-weight: 500;
}

.points-value {
  font-weight: 600;
  color: $text-1;

  &--ok {
    color: rgba(100, 200, 100, 0.9);
  }

  &--over, &--under {
    color: $red-0;
  }

  &--bonus {
    color: rgba(100, 200, 100, 0.9);
  }
}

.predator-traits {
  display: grid;
  gap: $s-3;
}

.predator-section {
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
  border: 1px solid $border;
}

.predator-list {
  margin: 0;
  padding-left: $s-4;
  color: $text-1;
  font-size: 0.9rem;

  li {
    margin-bottom: $s-1;
    &:last-child { margin-bottom: 0; }
  }
}

.loading-text, .empty-text {
  padding: $s-4;
  text-align: center;
  color: $text-2;
  background: rgba(255, 255, 255, .02);
  border-radius: $r-md;
}

.packs-list {
  display: grid;
  gap: $s-3;
}

.pack-item {
  border-radius: $r-md;
  border: 1px solid $border;
  overflow: hidden;
  transition: all $t-med $ease;

  &--expanded {
    border-color: $border-strong;
  }
}

.pack-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-4;
  width: 100%;
  padding: $s-4;
  background: rgba(255, 255, 255, .02);
  border: none;
  color: $text-0;
  cursor: pointer;
  text-align: left;
  transition: background $t-fast $ease;

  &:hover {
    background: rgba(255, 255, 255, .04);
  }
}

.pack-info {
  flex: 1;
  min-width: 0;

  h3 {
    margin: 0;
    font-family: $font-head;
    font-size: $fs-1;
    color: $text-0;
  }
}

.pack-desc {
  margin: $s-1 0 0;
  font-size: 0.85rem;
  color: $text-2;
  line-height: 1.4;
}

.pack-toggle {
  font-size: 1.5rem;
  font-weight: 300;
  color: $text-2;
  width: 32px;
  text-align: center;
}

.pack-traits {
  border-top: 1px solid $border;
  padding: $s-3;
  display: grid;
  gap: $s-2;
}

.trait-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: $s-3;
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
  border: 1px solid transparent;
  transition: all $t-fast $ease;

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .06);
  }

  &--flaw {
    border-color: rgba(255, 100, 100, 0.2);
  }

  &--flaw.trait-item--selected {
    border-color: rgba(255, 100, 100, 0.5);
    background: rgba(255, 100, 100, .08);
  }
}

.trait-info {
  flex: 1;
  min-width: 0;
  display: grid;
  gap: $s-1;
}

.trait-header {
  display: flex;
  align-items: center;
  gap: $s-2;
  flex-wrap: wrap;
}

.trait-name {
  font-weight: 600;
  color: $text-0;
  font-size: 0.95rem;
}

.trait-suffix {
  color: $text-2;
  font-size: 0.85rem;
  font-style: italic;
}

.trait-level-badge {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 59, 84, .12);
  color: $red-0;
  font-size: 0.75rem;
  font-weight: 600;

  &--flaw {
    background: rgba(255, 100, 100, .12);
    color: rgba(255, 100, 100, 0.9);
  }
}

.trait-desc {
  font-size: 0.85rem;
  color: $text-2;
  line-height: 1.4;
  margin: 0;
}

.trait-actions {
  display: flex;
  align-items: center;
  gap: $s-2;
  flex-shrink: 0;
}

.trait-levels {
  display: flex;
  gap: $s-1;
  flex-shrink: 0;
}

.suffix-edit-btn {
  width: 28px;
  height: 28px;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-2;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all $t-fast $ease;
  display: grid;
  place-items: center;

  &:hover {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
    color: $red-0;
  }
}

.level-btn {
  width: 32px;
  height: 32px;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-1;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all $t-fast $ease;

  &:hover {
    border-color: $border-strong;
    background: rgba(255, 255, 255, .04);
  }

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .12);
    color: $text-0;
  }
}

.flaw-toggle-btn {
  width: 36px;
  height: 36px;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-1;
  cursor: pointer;
  font-weight: 600;
  font-size: 1rem;
  transition: all $t-fast $ease;
  flex-shrink: 0;

  &:hover {
    border-color: rgba(255, 100, 100, 0.5);
    background: rgba(255, 100, 100, .08);
  }

  &--selected {
    border-color: rgba(255, 100, 100, 0.5);
    background: rgba(255, 100, 100, .15);
    color: rgba(255, 100, 100, 0.9);
  }
}

// Suffix Modal
.suffix-modal {
  display: grid;
  gap: $s-3;

  p {
    margin: 0;
    color: $text-1;
  }
}

.suffix-hint {
  color: $text-2 !important;
  font-size: 0.9rem;
}

.suffix-input {
  padding: $s-2 $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 0.95rem;

  &:focus {
    outline: none;
    border-color: $red-0;
  }

  &::placeholder {
    color: $text-2;
  }
}

.modal-actions {
  display: flex;
  gap: $s-3;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .trait-item {
    flex-direction: column;
    align-items: stretch;
  }

  .trait-levels, .flaw-toggle-btn {
    align-self: flex-end;
  }
}

@media (max-width: 420px) {
  .modal-actions {
    flex-direction: column-reverse;
  }
}
</style>
