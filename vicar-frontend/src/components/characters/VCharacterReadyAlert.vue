<script setup lang="ts">
import {computed} from "vue"
import {useV5DataStore} from "@/stores/v5data"
import VCard from "@/components/ui/VCard.vue"
import type {
  AttributeKey,
  CategoryKey,
  SkillKey,
  SkillSpreadType,
  V5Character,
  V5CharacterInternalData,
  V5TraitPack,
} from "@/@types/v5"

const props = defineProps<{
  character: Partial<V5Character>
}>()

const v5data = useV5DataStore()

const BASE_BACKGROUND_POINTS = 7
const MIN_FLAW_POINTS = 2

const attributeDefinitions: Record<CategoryKey, { key: AttributeKey; label: string }[]> = {
  physical: [
    { key: "str", label: "Stärke" },
    { key: "dex", label: "Geschicklichkeit" },
    { key: "sta", label: "Ausdauer" },
  ],
  social: [
    { key: "cha", label: "Charisma" },
    { key: "man", label: "Manipulation" },
    { key: "com", label: "Erscheinung" },
  ],
  mental: [
    { key: "int", label: "Intelligenz" },
    { key: "wit", label: "Witz" },
    { key: "res", label: "Entschlossenheit" },
  ],
}

const attributeCounts = computed(() => {
  const counts: Record<number, number> = { 1: 0, 2: 0, 3: 0, 4: 0 }
  for (const [category, defs] of Object.entries(attributeDefinitions)) {
    for (const def of defs) {
      const v = props.character.attributes?.find(a => a.category === category && a.key === def.key)?.value ?? 0
      if (v >= 1 && v <= 4) counts[v] = (counts[v] ?? 0) + 1
    }
  }
  return counts as { 1: number; 2: number; 3: number; 4: number }
})

const attributesValid = computed(() => {
  const c = attributeCounts.value
  return c[4] === 1 && c[3] === 3 && c[2] === 4 && c[1] === 1
})

const skillDefinitions: Record<CategoryKey, { key: SkillKey; label: string }[]> = {
  physical: [
    { key: "ath", label: "Athletik" },
    { key: "bra", label: "Handgemenge" },
    { key: "cra", label: "Handwerk" },
    { key: "dri", label: "Fahren" },
    { key: "fir", label: "Schusswaffen" },
    { key: "lar", label: "Heimlichkeit" },
    { key: "mel", label: "Nahkampf" },
    { key: "ste", label: "Überleben" },
    { key: "sur", label: "Überlebensinstinkt" },
  ],
  social: [
    { key: "ani", label: "Tierkunde" },
    { key: "eti", label: "Etikette" },
    { key: "ins", label: "Einblick" },
    { key: "int", label: "Einschüchtern" },
    { key: "lea", label: "Anführen" },
    { key: "per", label: "Überzeugen" },
    { key: "prf", label: "Darbietung" },
    { key: "sub", label: "Gassenwissen" },
    { key: "str", label: "Täuschen" },
  ],
  mental: [
    { key: "aca", label: "Geisteswissenschaften" },
    { key: "awa", label: "Aufmerksamkeit" },
    { key: "fin", label: "Finanzen" },
    { key: "inv", label: "Untersuchung" },
    { key: "med", label: "Medizin" },
    { key: "occ", label: "Okkultismus" },
    { key: "pol", label: "Politik" },
    { key: "sci", label: "Naturwissenschaften" },
    { key: "tec", label: "Technologie" },
  ],
}

const spreadTargets: Record<SkillSpreadType, Record<number, number>> = {
  jack_of_all_trades: { 3: 1, 2: 8, 1: 10, 0: 8 },
  balanced: { 3: 3, 2: 5, 1: 7, 0: 12 },
  specialist: { 4: 1, 3: 3, 2: 3, 1: 3, 0: 17 },
}

function getSkillValue(key: SkillKey): number {
  return props.character.skills?.find(s => s.key === key)?.value ?? 0
}

const skillCounts = computed(() => {
  const counts: Record<number, number> = { 0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 }
  for (const defs of Object.values(skillDefinitions)) {
    for (const def of defs) {
      const v = getSkillValue(def.key)
      counts[v] = (counts[v] ?? 0) + 1
    }
  }
  return counts
})

const skillsValid = computed(() => {
  const spread = (props.character.skillSpreadType ?? "balanced") as SkillSpreadType
  const target = spreadTargets[spread] ?? spreadTargets.balanced

  for (const [lvl, cnt] of Object.entries(target)) {
    if ((skillCounts.value[Number(lvl)] ?? 0) !== cnt) return false
  }
  return true
})

const internalData = computed<V5CharacterInternalData>(() => props.character.internalData ?? {})
const requiredFreeSpecSkills: SkillKey[] = ["aca", "cra", "prf", "sci"]

const freeSpecsValid = computed(() => {
  const free = (internalData.value.freeSpecializations ?? []).filter(s => (s.specialization ?? "").trim().length > 0)

  const hasRequired = requiredFreeSpecSkills.every(k =>
    free.some(s => s.skillKey === k && (s.specialization ?? "").trim().length > 0)
  )
  if (!hasRequired) return false

  return free.some(s => !requiredFreeSpecSkills.includes(s.skillKey as SkillKey) && (s.specialization ?? "").trim().length > 0)
})

const selectedPredatorType = computed(() => {
  if (!props.character.predatorTypeID) return null
  return v5data.predatorTypes?.find(p => p.id === props.character.predatorTypeID) ?? null
})

function getPacksForSpendBetweenRaw(allowedOldVicarIDs: number[]): V5TraitPack[] {
  return v5data.traitPacks?.filter(p => p.oldVicarID !== undefined && allowedOldVicarIDs.includes(p.oldVicarID)) ?? []
}

const spendPointsBetweenActions = computed(() => {
  if (!selectedPredatorType.value) return []
  return selectedPredatorType.value.actions
    .filter(a => a.type === "spend_background_points_between" || a.type === "spend_flaw_points_between")
    .map(a => {
      const data = a.data as any
      const allowedPackOldVicarIDs = data?.choices ?? []
      const allowedPacks = getPacksForSpendBetweenRaw(allowedPackOldVicarIDs)
      return {
        type: a.type === "spend_background_points_between" ? ("backgrounds" as const) : ("flaws" as const),
        bonusPoints: data?.points ?? 0,
        allowedPackIds: allowedPacks.map(p => p.id),
      }
    })
})

const backgroundCouponPackIds = computed(() => {
  const ids: string[] = []
  spendPointsBetweenActions.value.forEach(spa => { if (spa.type === "backgrounds") ids.push(...spa.allowedPackIds) })
  return ids
})

const flawCouponPackIds = computed(() => {
  const ids: string[] = []
  spendPointsBetweenActions.value.forEach(spa => { if (spa.type === "flaws") ids.push(...spa.allowedPackIds) })
  return ids
})

const pointsSpentOnBackgroundCouponPacks = computed(() => {
  let total = 0
  props.character.traitPackUsages?.forEach(usage => {
    if (backgroundCouponPackIds.value.includes(usage.packID)) {
      usage.traits?.forEach(t => { if (!t.isLocked) total += t.customLevel ?? 0 })
    }
  })
  return total
})

const pointsSpentOnFlawCouponPacks = computed(() => {
  let total = 0
  props.character.traitPackUsages?.forEach(usage => {
    if (flawCouponPackIds.value.includes(usage.packID)) {
      usage.flawTraits?.forEach(ft => {
        if (!ft.isLocked) {
          const pack = v5data.traitPacks?.find(p => p.id === usage.packID)
          const trait = pack?.packTraits?.find(pt => pt.trait.id === ft.traitID)?.trait
          if (trait) total += trait.level
        }
      })
    }
  })
  return total
})

const totalBackgroundCouponPoints = computed(() =>
  spendPointsBetweenActions.value.filter(s => s.type === "backgrounds").reduce((sum, s) => sum + s.bonusPoints, 0)
)
const totalFlawCouponPoints = computed(() =>
  spendPointsBetweenActions.value.filter(s => s.type === "flaws").reduce((sum, s) => sum + s.bonusPoints, 0)
)

const totalBackgroundPointsRaw = computed(() => {
  let total = 0
  props.character.traitPackUsages?.forEach(usage => {
    if (usage.kind === "backgrounds") {
      usage.traits?.forEach(t => { if (!t.isLocked) total += t.customLevel ?? 0 })
    }
  })
  return total
})

const totalBackgroundPoints = computed(() => {
  const couponCoverage = Math.min(pointsSpentOnBackgroundCouponPacks.value, totalBackgroundCouponPoints.value)
  return totalBackgroundPointsRaw.value - couponCoverage
})

const totalMeritPoints = computed(() => {
  let total = 0
  props.character.traitPackUsages?.forEach(usage => {
    if (usage.kind === "merits") {
      usage.traits?.forEach(t => { if (!t.isLocked) total += t.customLevel ?? 0 })
    }
  })
  return total
})

const totalFlawPointsRaw = computed(() => {
  let total = 0
  props.character.traitPackUsages?.forEach(usage => {
    usage.flawTraits?.forEach(ft => {
      if (!ft.isLocked) {
        const pack = v5data.traitPacks?.find(p => p.id === usage.packID)
        const trait = pack?.packTraits?.find(pt => pt.trait.id === ft.traitID)?.trait
        if (trait) total += trait.level
      }
    })
  })
  return total
})

const totalFlawPoints = computed(() => {
  const couponCoverage = Math.min(pointsSpentOnFlawCouponPacks.value, totalFlawCouponPoints.value)
  return totalFlawPointsRaw.value - couponCoverage
})

const bonusMeritPoints = computed(() => Math.max(0, totalFlawPoints.value - MIN_FLAW_POINTS))

const availableBackgroundPoints = computed(() => BASE_BACKGROUND_POINTS)
const availableMeritPoints = computed(() => bonusMeritPoints.value)
const availableFlawPoints = computed(() => MIN_FLAW_POINTS)

const backgroundPointsRemaining = computed(() => availableBackgroundPoints.value - totalBackgroundPoints.value)
const meritPointsRemaining = computed(() => availableMeritPoints.value - totalMeritPoints.value)
const flawPointsRemaining = computed(() => availableFlawPoints.value - totalFlawPoints.value)

const traitsValid = computed(() => {
  return backgroundPointsRemaining.value >= 0 &&
    meritPointsRemaining.value >= 0 &&
    flawPointsRemaining.value <= 0
})

type Check = { key: string; label: string; ok: boolean; hintOk?: string; hintBad?: string }

const checks = computed<Check[]>(() => ([
  {
    key: "attrs",
    label: "Attribute korrekt verteilt",
    ok: attributesValid.value,
    hintBad: `Erwartet: 1×4, 3×3, 4×2, 1×1 (aktuell: ${attributeCounts.value[4]}×4, ${attributeCounts.value[3]}×3, ${attributeCounts.value[2]}×2, ${attributeCounts.value[1]}×1)`,
  },
  {
    key: "skills",
    label: "Fähigkeiten korrekt verteilt",
    ok: skillsValid.value,
    hintBad: `Verteilung passt nicht zu "${props.character.skillSpreadType ?? "balanced"}".`,
  },
  {
    key: "specs",
    label: "Freie Spezialisierungen vollständig (5)",
    ok: freeSpecsValid.value,
    hintBad: `Pflicht: Aca/Cra/Prf/Sci + 1 freie Wahl (mit Wert).`,
  },
  {
    key: "traits",
    label: "Vorzüge/Hintergründe/Schwächen korrekt angerechnet",
    ok: traitsValid.value,
    hintBad: `Budget verletzt: BG ${totalBackgroundPoints.value}/${availableBackgroundPoints.value}, Vorzüge ${totalMeritPoints.value}/${availableMeritPoints.value}, Schwächen ${totalFlawPoints.value}/${availableFlawPoints.value}+`,
  },
]))

const allOk = computed(() => checks.value.every(c => c.ok))

const isDataReady = computed(() => {
  return !!v5data.traitPacks && !!v5data.predatorTypes
})

</script>

<template>
  <VCard v-if="!allOk" class="readyAlert" :class="{ 'readyAlert--ok': allOk, 'readyAlert--warn': !allOk }">
    <div class="readyAlert__header">
      <div class="readyAlert__title">
        <span class="readyAlert__icon" aria-hidden="true">{{ allOk ? "✅" : "⚠️" }}</span>
        <span>{{ allOk ? "Charakter ist bereit" : "Charakter ist noch nicht fertig" }}</span>
      </div>

      <div class="readyAlert__meta" v-if="!isDataReady">
        Hinweis: Daten werden noch geladen (Traits-Check kann unvollständig sein).
      </div>
    </div>

    <ul class="readyAlert__list">
      <li v-for="c in checks" :key="c.key" class="readyAlert__item">
        <span class="readyAlert__dot" :class="{ 'readyAlert__dot--ok': c.ok, 'readyAlert__dot--bad': !c.ok }" />
        <div class="readyAlert__text">
          <div class="readyAlert__label">
            {{ c.label }}
          </div>
          <div v-if="!c.ok && c.hintBad" class="readyAlert__hint">
            {{ c.hintBad }}
          </div>
        </div>
      </li>
    </ul>
  </VCard>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.readyAlert {
  border: 1px solid $border;
  background: rgba(255, 255, 255, 0.02);
}

.readyAlert--ok {
  border-color: rgba(100, 200, 100, 0.35);
  background: rgba(100, 200, 100, 0.05);
}

.readyAlert--warn {
  border-color: rgba($red-0, 0.35);
  background: rgba($red-0, 0.06);
}

.readyAlert__header {
  display: grid;
  gap: $s-2;
  margin-bottom: $s-3;
}

.readyAlert__title {
  display: flex;
  align-items: center;
  gap: $s-2;
  font-family: $font-head;
  font-size: 1.05rem;
  color: $text-0;
}

.readyAlert__icon {
  font-size: 1.1rem;
}

.readyAlert__meta {
  font-size: 0.85rem;
  color: $text-2;
}

.readyAlert__list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: $s-2;
}

.readyAlert__item {
  display: flex;
  gap: $s-2;
  align-items: flex-start;
}

.readyAlert__dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  margin-top: 0.35rem;
  background: rgba(255, 255, 255, 0.18);
}

.readyAlert__dot--ok {
  background: rgba(100, 200, 100, 0.9);
}

.readyAlert__dot--bad {
  background: rgba($red-0, 0.9);
}

.readyAlert__label {
  font-weight: 600;
  color: $text-1;
}

.readyAlert__hint {
  margin-top: 2px;
  font-size: 0.85rem;
  color: $text-2;
}
</style>
