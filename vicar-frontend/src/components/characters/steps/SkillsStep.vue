<script setup lang="ts">
import { ref, computed } from 'vue'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'
import type { V5Character, CategoryKey, SkillKey, V5CharacterSkill, SkillSpreadType } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })

const skillDefinitions: Record<CategoryKey, { key: SkillKey; label: string }[]> = {
  physical: [
    { key: 'ath', label: 'Athletik' },
    { key: 'bra', label: 'Handgemenge' },
    { key: 'cra', label: 'Handwerk' },
    { key: 'dri', label: 'Fahren' },
    { key: 'fir', label: 'Schusswaffen' },
    { key: 'lar', label: 'Heimlichkeit' },
    { key: 'mel', label: 'Nahkampf' },
    { key: 'ste', label: 'Überleben' },
    { key: 'sur', label: 'Überlebensinstinkt' },
  ],
  social: [
    { key: 'ani', label: 'Tierkunde' },
    { key: 'eti', label: 'Etikette' },
    { key: 'ins', label: 'Einblick' },
    { key: 'int', label: 'Einschüchtern' },
    { key: 'lea', label: 'Anführen' },
    { key: 'per', label: 'Überzeugen' },
    { key: 'prf', label: 'Darbietung' },
    { key: 'sub', label: 'Gassenwissen' },
    { key: 'str', label: 'Täuschen' },
  ],
  mental: [
    { key: 'aca', label: 'Geisteswissenschaften' },
    { key: 'awa', label: 'Aufmerksamkeit' },
    { key: 'fin', label: 'Finanzen' },
    { key: 'inv', label: 'Untersuchung' },
    { key: 'med', label: 'Medizin' },
    { key: 'occ', label: 'Okkultismus' },
    { key: 'pol', label: 'Politik' },
    { key: 'sci', label: 'Naturwissenschaften' },
    { key: 'tec', label: 'Technologie' },
  ],
}

const spreadOptions: { label: string; value: SkillSpreadType; distribution: Record<number, number> }[] = [
  { label: 'Tausendsassa', value: 'jack_of_all_trades', distribution: { 3: 1, 2: 8, 1: 10, 0: 8 } },
  { label: 'Ausgewogen', value: 'balanced', distribution: { 3: 3, 2: 5, 1: 7, 0: 12 } },
  { label: 'Spezialist', value: 'specialist', distribution: { 4: 1, 3: 3, 2: 3, 1: 3, 0: 17 } },
]

// Input refs for UI only
const specializationInputs = ref<Record<string, string>>({})

const freeSpecSkills: { key: SkillKey; label: string }[] = [
  { key: 'aca', label: 'Geisteswissenschaften' },
  { key: 'cra', label: 'Handwerk' },
  { key: 'prf', label: 'Darbietung' },
  { key: 'sci', label: 'Naturwissenschaften' },
]

// @ts-ignore
const freeSpecInputs = ref<Record<SkillKey, string>>({
// @ts-ignore
  aca: '',
// @ts-ignore
  cra: '',
// @ts-ignore
  prf: '',
// @ts-ignore
  sci: '',
})
const freeChoiceSkill = ref<SkillKey | ''>('')
const freeChoiceSpec = ref('')

const localSkills = computed<V5CharacterSkill[]>({
  get: () => {
    const skills: V5CharacterSkill[] = []
    Object.entries(skillDefinitions).forEach(([category, defs]) => {
      defs.forEach(def => {
        const existing = character.value?.skills?.find(
          s => s.category === category && s.key === def.key
        )
        skills.push({
          id: existing?.id ?? '',
          characterID: existing?.characterID ?? '',
          category: category as CategoryKey,
          key: def.key,
          value: existing?.value ?? 0,
          specialization: existing?.specialization ?? [],
        })
      })
    })
    return skills
  },
  set: (val: V5CharacterSkill[]) => {
    character.value.skills = val
  }
})

const selectedSpread = computed<SkillSpreadType>({
  get: () => character.value?.skillSpreadType ?? 'balanced',
  set: (val: SkillSpreadType) => {
    character.value.skillSpreadType = val
  }
})

function getSkillValue(category: CategoryKey, key: SkillKey): number {
  const skill = localSkills.value.find(s => s.category === category && s.key === key)
  return skill?.value ?? 0
}

function setSkillValue(category: CategoryKey, key: SkillKey, value: number) {
  const index = localSkills.value.findIndex(s => s.category === category && s.key === key)
  if (index !== -1) {
    localSkills.value[index] = { ...localSkills.value[index], value } as any
    localSkills.value = [...localSkills.value]
  }
}

function getSpecializations(category: CategoryKey, key: SkillKey): string[] {
  const skill = localSkills.value.find(s => s.category === category && s.key === key)
  return skill?.specialization ?? []
}

function addSpecialization(category: CategoryKey, key: SkillKey) {
  const inputKey = `${category}-${key}`
  const inputValue = specializationInputs.value[inputKey]?.trim()
  if (!inputValue) return

  const index = localSkills.value.findIndex(s => s.category === category && s.key === key)
  if (index !== -1) {
    const skill = localSkills.value[index]!
    if (!skill.specialization.includes(inputValue)) {
      localSkills.value[index] = {
        ...skill,
        specialization: [...skill.specialization, inputValue]
      }
      localSkills.value = [...localSkills.value]
    }
  }
  specializationInputs.value[inputKey] = ''
}

function removeSpecialization(category: CategoryKey, key: SkillKey, spec: string) {
  const index = localSkills.value.findIndex(s => s.category === category && s.key === key)
  if (index !== -1) {
    const skill = localSkills.value[index]!
    localSkills.value[index] = {
      ...skill,
      specialization: skill.specialization.filter(s => s !== spec)
    }
    localSkills.value = [...localSkills.value]
  }
}

function addFreeSpec(key: SkillKey) {
  const value = freeSpecInputs.value[key]?.trim()
  if (!value) return

  const skill = localSkills.value.find(s => s.key === key)
  if (skill && !skill.specialization.includes(value)) {
    const index = localSkills.value.findIndex(s => s.key === key)
    localSkills.value[index] = {
      ...localSkills.value[index],
      specialization: [...localSkills.value[index]!.specialization, value]
    } as any
    localSkills.value = [...localSkills.value]
  }
  freeSpecInputs.value[key] = ''
}

function addFreeChoiceSpec() {
  if (!freeChoiceSkill.value || !freeChoiceSpec.value.trim()) return

  const skill = localSkills.value.find(s => s.key === freeChoiceSkill.value)
  if (skill && !skill.specialization.includes(freeChoiceSpec.value.trim())) {
    const index = localSkills.value.findIndex(s => s.key === freeChoiceSkill.value)
    localSkills.value[index] = {
      ...localSkills.value[index],
      specialization: [...localSkills.value[index]!.specialization, freeChoiceSpec.value.trim()]
    } as any
    localSkills.value = [...localSkills.value]
  }
  freeChoiceSpec.value = ''
}

const currentDistribution = computed(() => {
  const spread = spreadOptions.find(s => s.value === selectedSpread.value)
  return spread?.distribution ?? {}
})

const distributionCounts = computed(() => {
  const counts: Record<number, number> = { 0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 }
  localSkills.value.forEach(skill => {
    counts[skill.value] = (counts[skill.value] || 0) + 1
  })
  return counts
})

const distributionValid = computed(() => {
  const target = currentDistribution.value
  for (const [level, count] of Object.entries(target)) {
    if ((distributionCounts.value[Number(level)] || 0) !== count) {
      return false
    }
  }
  return true
})

const allSkillsFlat = computed(() => {
  const result: { key: SkillKey; label: string }[] = []
  Object.values(skillDefinitions).forEach(defs => {
    defs.forEach(def => result.push(def))
  })
  return result
})

const categoryLabels: Record<CategoryKey, string> = {
  physical: 'Körperlich',
  social: 'Sozial',
  mental: 'Geistig',
}

function getMaxForSpread(): number {
  return selectedSpread.value === 'specialist' ? 4 : 3
}
</script>

<template>
  <div class="skills-step">
    <VCard>
      <h2>Verteilung wählen</h2>
      <div class="spread-selection">
        <button
          v-for="spread in spreadOptions"
          :key="spread.value"
          type="button"
          class="spread-btn"
          :class="{ 'spread-btn--selected': selectedSpread === spread.value }"
          @click="selectedSpread = spread.value"
        >
          <span class="spread-name">{{ spread.label }}</span>
          <span class="spread-info">
            <template v-for="(count, level) in spread.distribution" :key="level">
              <span v-if="Number(level) > 0 && count > 0" class="spread-detail">
                {{ count }}×{{ level }}
              </span>
            </template>
          </span>
        </button>
      </div>

      <div class="distribution-status" :class="{ 'distribution-status--valid': distributionValid }">
        <template v-for="(count, level) in currentDistribution" :key="level">
          <div v-if="Number(level) > 0" class="dist-item">
            <span class="dist-label">{{ level }}er:</span>
            <span
              class="dist-value"
              :class="{
                'dist-value--ok': distributionCounts[Number(level)] === count,
                'dist-value--over': distributionCounts[Number(level)]! > count
              }"
            >
              {{ distributionCounts[Number(level)] || 0 }}/{{ count }}
            </span>
          </div>
        </template>
      </div>
    </VCard>

    <VCard>
      <h2>Fähigkeiten verteilen</h2>
      <div class="skills-grid">
        <div v-for="category in (['physical', 'social', 'mental'] as CategoryKey[])" :key="category" class="skill-category">
          <h3>{{ categoryLabels[category] }}</h3>
          <div
            v-for="skill in skillDefinitions[category]"
            :key="skill.key"
            class="skill-item"
          >
            <div class="skill-row">
              <span class="skill-label">{{ skill.label }}</span>
              <VDotRating
                :model-value="getSkillValue(category, skill.key)"
                :max="getMaxForSpread()"
                @update:model-value="(v: number) => setSkillValue(category, skill.key, v)"
              />
            </div>
            <div v-if="getSpecializations(category, skill.key).length > 0" class="specializations">
              <div class="spec-tags">
                <span
                  v-for="spec in getSpecializations(category, skill.key)"
                  :key="spec"
                  class="spec-tag"
                >
                  {{ spec }}
                  <button
                    type="button"
                    class="spec-remove"
                    @click="removeSpecialization(category, skill.key, spec)"
                  >×</button>
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </VCard>

    <VCard>
      <h2>Freie Spezialisierungen (5)</h2>
      <p class="info-text">Du erhältst 5 kostenlose Spezialisierungen: je eine für Geisteswissenschaften, Handwerk, Darbietung, Naturwissenschaften und eine beliebige.</p>

      <div class="free-specs">
        <div v-for="spec in freeSpecSkills" :key="spec.key" class="free-spec-row">
          <span class="free-spec-label">{{ spec.label }}:</span>
          <div class="free-spec-input-row">
            <input
              type="text"
              class="input"
              v-model="freeSpecInputs[spec.key]"
              :placeholder="`Spezialisierung für ${spec.label}`"
              @keydown.enter.prevent="addFreeSpec(spec.key)"
            />
            <button type="button" class="btn" @click="addFreeSpec(spec.key)">+</button>
          </div>
        </div>

        <div class="free-spec-row free-spec-row--choice">
          <span class="free-spec-label">Freie Wahl:</span>
          <div class="free-choice-row">
            <select v-model="freeChoiceSkill" class="select">
              <option value="">Fähigkeit wählen...</option>
              <option v-for="skill in allSkillsFlat" :key="skill.key" :value="skill.key">
                {{ skill.label }}
              </option>
            </select>
            <input
              type="text"
              class="input"
              v-model="freeChoiceSpec"
              placeholder="Spezialisierung"
              @keydown.enter.prevent="addFreeChoiceSpec"
            />
            <button type="button" class="btn" @click="addFreeChoiceSpec">+</button>
          </div>
        </div>
      </div>
    </VCard>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.skills-step {
  display: grid;
  gap: $s-5;
}

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
  padding-bottom: $s-2;
  border-bottom: 1px solid $border;
}

.spread-selection {
  display: grid;
  gap: $s-2;
  margin-bottom: $s-4;
}

.spread-btn {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $s-3 $s-4;
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

.spread-name {
  font-weight: 600;
}

.spread-info {
  display: flex;
  gap: $s-2;
  color: $text-2;
  font-size: 0.85rem;
}

.spread-detail {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 255, 255, 0.04);
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

.skills-grid {
  display: grid;
  gap: $s-6;
  grid-template-columns: repeat(3, 1fr);
}

.skill-category {
  display: grid;
  gap: $s-3;
  min-width: 0;
}

.skill-item {
  display: grid;
  gap: $s-2;
}

.skill-row {
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

.skill-label {
  font-weight: 500;
  color: $text-0;
  font-size: 0.9rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.specializations {
  padding: 0 $s-2 $s-2;
  display: grid;
  gap: $s-2;
}

.spec-tags {
  display: flex;
  flex-wrap: wrap;
  gap: $s-1;
}

.spec-tag {
  display: inline-flex;
  align-items: center;
  gap: $s-1;
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  background: rgba(255, 59, 84, .12);
  color: $text-0;
  font-size: 0.75rem;
}

.spec-remove {
  background: none;
  border: none;
  color: $text-2;
  cursor: pointer;
  padding: 0;
  font-size: 0.9rem;
  line-height: 1;

  &:hover {
    color: $red-0;
  }
}

.spec-input-row {
  display: flex;
  gap: $s-1;
}

.spec-input {
  flex: 1;
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-size: 0.75rem;
  font-family: $font-body;

  &::placeholder {
    color: $text-2;
  }

  &:focus {
    outline: none;
    border-color: $red-0;
  }
}

.spec-add-btn {
  padding: $s-1 $s-2;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-1;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all $t-fast $ease;

  &:hover {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
    color: $text-0;
  }
}

.info-text {
  margin: 0 0 $s-4;
  color: $text-2;
  font-size: 0.9rem;
}

.free-specs {
  display: grid;
  gap: $s-3;
}

.free-spec-row {
  display: grid;
  gap: $s-2;

  &--choice {
    padding-top: $s-3;
    border-top: 1px solid $border;
  }
}

.free-spec-label {
  font-weight: 500;
  color: $text-1;
}

.free-spec-input-row {
  display: flex;
  gap: $s-2;
}

.free-choice-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: $s-2;
}

.input {
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 0.9rem;

  &:focus {
    outline: none;
    border-color: $red-0;
  }

  &::placeholder {
    color: $text-2;
  }
}

.select {
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 0.9rem;

  &:focus {
    outline: none;
    border-color: $red-0;
  }

  option {
    background: $bg-1;
    color: $text-0;
  }
}

.btn {
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-1;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all $t-fast $ease;

  &:hover {
    border-color: $red-0;
    background: rgba(255, 59, 84, .08);
    color: $text-0;
  }
}

@media (max-width: 768px) {
  .skills-grid {
    grid-template-columns: 1fr;
  }

  .free-choice-row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 420px) {
  .distribution-status {
    flex-direction: column;
    gap: $s-2;
  }

  .spread-btn {
    flex-direction: column;
    align-items: flex-start;
    gap: $s-2;
  }
}
</style>
