<script setup lang="ts">
import { computed } from 'vue'
import VCard from '@/components/ui/VCard.vue'
import VDotRating from '@/components/characters/VDotRating.vue'
import type { V5Character, CategoryKey, SkillKey } from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })

const skillDefinitions: Record<CategoryKey, { key: SkillKey; label: string }[]> = {
  physical: [
    { key: 'ath', label: 'Athletik' },
    { key: 'bra', label: 'Handgemenge' },
    { key: 'cra', label: 'Handwerk' },
    { key: 'dri', label: 'Fahren' },
    { key: 'fir', label: 'Schusswaffen' },
    { key: 'mel', label: 'Nahkampf' },
    { key: 'lar', label: 'Heimlichkeit' },
    { key: 'ste', label: 'Überleben' },
    { key: 'sur', label: 'Überlebenskampf' },
  ],
  social: [
    { key: 'ani', label: 'Mit Tieren' },
    { key: 'eti', label: 'Etikette' },
    { key: 'ins', label: 'Einblick' },
    { key: 'int', label: 'Einschüchtern' },
    { key: 'lea', label: 'Anführen' },
    { key: 'per', label: 'Überzeugen' },
    { key: 'sub', label: 'Gassenwissen' },
    { key: 'str', label: 'Täuschen' },
  ],
  mental: [
    { key: 'aca', label: 'Wissenschaft' },
    { key: 'awa', label: 'Aufmerksamkeit' },
    { key: 'fin', label: 'Finanzen' },
    { key: 'inv', label: 'Untersuchung' },
    { key: 'med', label: 'Medizin' },
    { key: 'occ', label: 'Okkultismus' },
    { key: 'pol', label: 'Politik' },
    { key: 'sci', label: 'Wissenschaft (alt)' },
    { key: 'tec', label: 'Technologie' },
  ],
}

function initializeSkills() {
  if (!character.value.skills || character.value.skills.length === 0) {
    character.value.skills = []
    Object.entries(skillDefinitions).forEach(([category, skills]) => {
      skills.forEach(skill => {
        character.value.skills!.push({
          id: `${category}-${skill.key}`,
          characterID: '',
          category: category as CategoryKey,
          key: skill.key,
          value: 0,
          specialization: [],
        })
      })
    })
  }
}

initializeSkills()

function getSkillValue(category: CategoryKey, key: SkillKey): number {
  const skill = character.value.skills?.find(s => s.category === category && s.key === key)
  return skill?.value ?? 0
}

function setSkillValue(category: CategoryKey, key: SkillKey, value: number) {
  if (!character.value.skills) {
    character.value.skills = []
  }
  const skill = character.value.skills.find(s => s.category === category && s.key === key)
  if (skill) {
    skill.value = value
  } else {
    character.value.skills.push({
      id: `${category}-${key}`,
      characterID: '',
      category,
      key,
      value,
      specialization: [],
    })
  }
}

const totalPoints = computed(() => {
  return character.value.skills?.reduce((sum, skill) => sum + skill.value, 0) ?? 0
})
</script>

<template>
  <VCard>
    <div class="header">
      <h2>Fähigkeiten verteilen</h2>
      <div class="points-display">
        Gesamt: {{ totalPoints }} Punkte
      </div>
    </div>

    <div class="step-content">
      <div class="skills-grid">
        <div class="skill-category">
          <h3>Körperlich</h3>
          <div
            v-for="skill in skillDefinitions.physical"
            :key="skill.key"
            class="skill-row"
          >
            <span class="skill-label">{{ skill.label }}</span>
            <VDotRating
              :model-value="getSkillValue('physical', skill.key)"
              @update:model-value="(v) => setSkillValue('physical', skill.key, v)"
            />
          </div>
        </div>

        <div class="skill-category">
          <h3>Sozial</h3>
          <div
            v-for="skill in skillDefinitions.social"
            :key="skill.key"
            class="skill-row"
          >
            <span class="skill-label">{{ skill.label }}</span>
            <VDotRating
              :model-value="getSkillValue('social', skill.key)"
              @update:model-value="(v) => setSkillValue('social', skill.key, v)"
            />
          </div>
        </div>

        <div class="skill-category">
          <h3>Geistig</h3>
          <div
            v-for="skill in skillDefinitions.mental"
            :key="skill.key"
            class="skill-row"
          >
            <span class="skill-label">{{ skill.label }}</span>
            <VDotRating
              :model-value="getSkillValue('mental', skill.key)"
              @update:model-value="(v) => setSkillValue('mental', skill.key, v)"
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

.skills-grid {
  display: grid;
  gap: $s-5;
  grid-template-columns: repeat(3, 1fr);
}

.skill-category {
  display: grid;
  gap: $s-2;

  h3 {
    margin: 0 0 $s-2;
    padding-bottom: $s-2;
    font-family: $font-head;
    font-size: $fs-1;
    border-bottom: 1px solid $border;
    color: $text-1;
  }
}

.skill-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-2;
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
  min-width: 120px;
  font-size: 0.9rem;
}

@media (max-width: 420px) {
  .skills-grid {
    grid-template-columns: 1fr;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: $s-3;
  }
}
</style>
