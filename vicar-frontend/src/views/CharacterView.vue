<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useV5DataStore } from '@/stores/v5data'
import VButton from '@/components/ui/VButton.vue'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'
import type { AttributeKey, SkillKey, CategoryKey } from '@/@types/v5'

const route = useRoute()
const router = useRouter()
const store = useCharactersStore()
const v5data = useV5DataStore()
const isLoading = ref(true)

const characterId = computed(() => route.params.id as string)
const character = computed(() => store.currentCharacter)

const attributeLabels: Record<AttributeKey, string> = {
  str: 'Stärke',
  dex: 'Geschicklichkeit',
  sta: 'Ausdauer',
  cha: 'Charisma',
  man: 'Manipulation',
  com: 'Erscheinung',
  int: 'Intelligenz',
  wit: 'Witz',
  res: 'Entschlossenheit',
}

const skillLabels: Record<SkillKey, string> = {
  ath: 'Athletik',
  bra: 'Handgemenge',
  cra: 'Handwerk',
  dri: 'Fahren',
  fir: 'Schusswaffen',
  mel: 'Nahkampf',
  lar: 'Heimlichkeit',
  ste: 'Überleben',
  sur: 'Überlebensinstinkt',
  ani: 'Tierkunde',
  eti: 'Etikette',
  ins: 'Einblick',
  int: 'Einschüchtern',
  lea: 'Anführen',
  per: 'Überzeugen',
  prf: 'Darbietung',
  sub: 'Gassenwissen',
  str: 'Täuschen',
  aca: 'Akademik',
  awa: 'Aufmerksamkeit',
  fin: 'Finanzen',
  inv: 'Untersuchung',
  med: 'Medizin',
  occ: 'Okkultismus',
  pol: 'Politik',
  sci: 'Wissenschaft',
  tec: 'Technologie',
}

const categoryLabels: Record<CategoryKey, string> = {
  physical: 'Körperlich',
  social: 'Sozial',
  mental: 'Geistig',
}

const generationEraLabels: Record<string, string> = {
  children: 'Kinder',
  newborn: 'Neugeborene',
  ancillae: 'Ancillae',
  older: 'Ältere',
  elder: 'Älteste',
  cainesinheritance: 'Kains Erbe',
}

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
    v5data.fetchTraitPacks(),
  ])

  if (characterId.value) {
    await store.fetchCharacter(characterId.value)
  }
  isLoading.value = false
})

function goToEditor() {
  router.push(`/characters/${characterId.value}/edit`)
}

function getAttributesByCategory(category: CategoryKey) {
  return character.value?.attributes?.filter(a => a.category === category) ?? []
}

function getSkillsByCategory(category: CategoryKey) {
  return character.value?.skills?.filter(s => s.category === category) ?? []
}

function getTraitPackName(packId: string): string {
  const pack = v5data.traitPacks?.find(p => p.id === packId)
  return pack?.name ?? 'Unbekannt'
}

function getTraitName(packId: string, traitId: string): string {
  const pack = v5data.traitPacks?.find(p => p.id === packId)
  const packTrait = pack?.packTraits?.find(pt => pt.trait.id === traitId)
  return packTrait?.trait.name ?? 'Unbekannt'
}

function getDisciplineName(disciplineId: string): string {
  if (!character.value?.clanID) return 'Unbekannt'
  const selectedClan = v5data.getClanById(character.value.clanID)
  const disc = selectedClan?.disciplines?.find(d => d.id === disciplineId)
  return disc?.name ?? 'Unbekannt'
}
</script>

<template>
  <div v-if="isLoading" class="character-view character-view--loading">
    <p>Lade Charakter...</p>
  </div>

  <div class="character-view" v-else-if="character">
    <div class="character-view__header">
      <div class="header-left">
        <VButton variant="ghost" @click="router.push('/')">
          ← Zurück
        </VButton>
        <h1>{{ character.name }}</h1>
        <p class="subtitle" v-if="clan">{{ clan.slogan }}</p>
      </div>
      <VButton variant="primary" @click="goToEditor">
        Bearbeiten
      </VButton>
    </div>

    <div class="character-view__content">
      <div class="content-grid">
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
              <span>{{ character.generation }}. ({{ generationEraLabels[character.generationEra] || character.generationEra }})</span>
            </div>
            <div class="info-item" v-if="character.sire">
              <span class="label">Sire:</span>
              <span>{{ character.sire }}</span>
            </div>
          </div>
        </VCard>

        <VCard>
          <h2>Werte</h2>
          <div class="stats-grid">
            <div class="stat-item">
              <span class="stat-label">Hunger</span>
              <div class="stat-value stat-value--hunger">{{ character.hunger }}</div>
            </div>
            <div class="stat-item">
              <span class="stat-label">Menschlichkeit</span>
              <div class="stat-value">{{ character.humanity }}</div>
            </div>
            <div class="stat-item">
              <span class="stat-label">Blutmacht</span>
              <div class="stat-value stat-value--potency">{{ character.bloodPotency }}</div>
            </div>
            <div class="stat-item" v-if="character.stains > 0">
              <span class="stat-label">Flecken</span>
              <div class="stat-value stat-value--stains">{{ character.stains }}</div>
            </div>
          </div>
        </VCard>
      </div>

      <VCard v-if="character.attributes && character.attributes.length > 0">
        <h2>Attribute</h2>
        <div class="attributes-grid">
          <div v-for="category in (['physical', 'social', 'mental'] as CategoryKey[])" :key="category" class="attr-category">
            <h3>{{ categoryLabels[category] }}</h3>
            <div v-for="attr in getAttributesByCategory(category)" :key="attr.id" class="attr-row">
              <span class="attr-name">{{ attributeLabels[attr.key] || attr.key }}</span>
              <VDotRating :model-value="attr.value" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.skills && character.skills.length > 0">
        <h2>Fähigkeiten</h2>
        <div class="skills-grid">
          <div v-for="category in (['physical', 'social', 'mental'] as CategoryKey[])" :key="category" class="skill-category">
            <h3>{{ categoryLabels[category] }}</h3>
            <div v-for="skill in getSkillsByCategory(category)" :key="skill.id" class="skill-row">
              <div class="skill-info">
                <span class="skill-name">{{ skillLabels[skill.key] || skill.key }}</span>
                <span v-if="skill.specialization && skill.specialization.length > 0" class="skill-spec">
                  {{ skill.specialization.join(', ') }}
                </span>
              </div>
              <VDotRating :model-value="skill.value" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.disciplineSelections && character.disciplineSelections.length > 0">
        <h2>Disziplinen</h2>
        <div class="disciplines-list">
          <div v-for="disc in character.disciplineSelections" :key="disc.id" class="disc-item">
            <div class="disc-header">
              <span class="disc-name">{{ getDisciplineName(disc.disciplineID) }}</span>
              <VDotRating :model-value="disc.currentLevel" :readonly="true" />
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.traitPackUsages && character.traitPackUsages.length > 0">
        <h2>Vorzüge & Hintergründe</h2>
        <div class="traits-section">
          <div v-for="usage in character.traitPackUsages" :key="usage.id" class="trait-pack">
            <h3>{{ getTraitPackName(usage.packID) }}</h3>
            <div class="traits-list">
              <div v-for="trait in usage.traits" :key="trait.id" class="trait-item">
                <span class="trait-name">{{ getTraitName(usage.packID, trait.traitID) }}</span>
                <span class="trait-level" v-if="trait.customLevel">{{ trait.customLevel }}</span>
              </div>
            </div>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.bloodRituals && character.bloodRituals.length > 0">
        <h2>Blutrituale</h2>
        <div class="rituals-list">
          <div v-for="ritual in character.bloodRituals" :key="ritual.id" class="ritual-item">
            <div class="ritual-header">
              <span class="ritual-name">{{ ritual.name }}</span>
              <span class="ritual-level">Level {{ ritual.level }}</span>
            </div>
            <p class="ritual-desc">{{ ritual.description }}</p>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.oblivionCeremonies && character.oblivionCeremonies.length > 0">
        <h2>Oblivion Zeremonien</h2>
        <div class="rituals-list">
          <div v-for="ceremony in character.oblivionCeremonies" :key="ceremony.id" class="ritual-item">
            <div class="ritual-header">
              <span class="ritual-name">{{ ceremony.name }}</span>
              <span class="ritual-level">Level {{ ceremony.level }}</span>
            </div>
            <p class="ritual-desc">{{ ceremony.summary }}</p>
          </div>
        </div>
      </VCard>

      <VCard v-if="character.notes">
        <h2>Notizen</h2>
        <div class="notes-content">
          {{ character.notes }}
        </div>
      </VCard>
    </div>
  </div>

  <div v-else class="character-view character-view--loading">
    <p>Charakter nicht gefunden.</p>
    <VButton variant="primary" @click="router.push('/')">
      Zur Übersicht
    </VButton>
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

      .subtitle {
        margin: 0;
        font-style: italic;
        color: $text-2;
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
    gap: $s-4;
  }
}

.content-grid {
  display: grid;
  gap: $s-5;
  grid-template-columns: repeat(2, 1fr);
}

.info-grid {
  display: grid;
  gap: $s-2;
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: $s-4;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $s-2;
}

.stat-label {
  font-size: 0.85rem;
  color: $text-2;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  font-family: $font-head;
  color: $text-0;
  padding: $s-2 $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.04);

  &--hunger {
    background: rgba(255, 59, 84, 0.12);
    color: $red-0;
  }

  &--potency {
    background: rgba(140, 12, 27, 0.2);
    color: $red-1;
  }

  &--stains {
    background: rgba(140, 12, 27, 0.15);
    color: $text-2;
  }
}

.attributes-grid, .skills-grid {
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
  padding: $s-1 0;
}

.attr-name, .skill-name {
  font-weight: 500;
  color: $text-0;
}

.skill-info {
  display: flex;
  flex-direction: column;
  gap: 0;
  flex: 1;
}

.skill-spec {
  font-size: 0.8rem;
  color: $red-0;
  font-style: italic;
}

.disciplines-list {
  display: grid;
  gap: $s-3;
}

.disc-item {
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;
}

.disc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.disc-name {
  font-weight: 600;
  color: $text-0;
  font-family: $font-head;
}

.traits-section {
  display: grid;
  gap: $s-4;
}

.trait-pack {
  h3 {
    margin-bottom: $s-2;
    padding-bottom: $s-2;
    border-bottom: 1px solid $border;
  }
}

.traits-list {
  display: grid;
  gap: $s-2;
}

.trait-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, 0.02);
}

.trait-name {
  font-weight: 500;
  color: $text-0;
}

.trait-level {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 59, 84, 0.12);
  color: $red-0;
  font-weight: 600;
  font-size: 0.85rem;
}

.rituals-list {
  display: grid;
  gap: $s-3;
}

.ritual-item {
  padding: $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $border;
}

.ritual-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $s-2;
}

.ritual-name {
  font-weight: 600;
  color: $text-0;
  font-family: $font-head;
}

.ritual-level {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, 0.06);
  font-size: 0.85rem;
  font-weight: 600;
  color: $text-1;
}

.ritual-desc {
  margin: 0;
  font-size: 0.9rem;
  color: $text-2;
  line-height: 1.5;
}

.notes-content {
  white-space: pre-wrap;
  color: $text-1;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .content-grid {
    grid-template-columns: 1fr;
  }

  .attributes-grid, .skills-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 420px) {
  .character-view {
    padding: $s-4;

    &__header {
      flex-direction: column;
      align-items: stretch;
    }
  }
}
</style>
