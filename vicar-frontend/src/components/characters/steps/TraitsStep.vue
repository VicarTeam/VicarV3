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
  V5CharacterInternalData,
  V5PredatorAction,
} from '@/@types/v5'

const character = defineModel<Partial<V5Character>>({ required: true })
const v5data = useV5DataStore()

const isLoading = ref(true)
const expandedMeritPacks = ref<Set<string>>(new Set())
const expandedBackgroundPacks = ref<Set<string>>(new Set())

const showSuffixModal = ref(false)
const currentSuffixTrait = ref<{ pack: V5TraitPack; trait: V5Trait; level: number } | null>(null)
const suffixInput = ref('')

const BASE_BACKGROUND_POINTS = 7
const MIN_FLAW_POINTS = 2

onMounted(async () => {
  await Promise.all([
    v5data.fetchTraitPacks(),
    v5data.fetchPredatorTypes(),
  ])
  isLoading.value = false
  applyPredatorTraitBonuses()
})

watch(() => character.value.predatorTypeID, (newId, oldId) => {
  if (newId !== oldId) {
    revertPredatorTraitBonuses()
    applyPredatorTraitBonuses()
  }
})

const internalData = computed<V5CharacterInternalData>({
  get: () => character.value.internalData ?? {},
  set: (val) => { character.value.internalData = val }
})

const meritPacks = computed(() => {
  return v5data.traitPacks?.filter(p => p.type === 'merits') ?? []
})

const backgroundPacks = computed(() => {
  return v5data.traitPacks?.filter(p => p.type === 'backgrounds') ?? []
})

const selectedPredatorType = computed(() => {
  if (!character.value.predatorTypeID) return null
  return v5data.predatorTypes?.find(p => p.id === character.value.predatorTypeID) ?? null
})

const spendPointsBetweenActions = computed(() => {
  if (!selectedPredatorType.value) return []

  const actions = selectedPredatorType.value.actions

  const betweenCoupons = actions
    .filter(a => a.type === "spend_background_points_between" || a.type === "spend_flaw_points_between")
    .map(a => {
      const data = a.data as any

      const couponType = a.type === "spend_background_points_between"
        ? ("backgrounds" as const)
        : ("flaws" as const)

      if (couponType === "backgrounds") {
        const allowedPackOldVicarIDs = normalizeOldVicarChoiceIds(data?.choices)
        const allowedPacks = getBackgroundCouponPacks(allowedPackOldVicarIDs)

        return {
          action: a,
          type: couponType,
          bonusPoints: data?.points ?? 0,
          allowedPackOldVicarIDs,
          allowedPackIds: allowedPacks.map(p => p.id),
          allowedPacks
        }
      } else {
        const allowedTypedChoices = normalizeOldVicarTypedChoices(data?.choices)
        const allowedPacks = getFlawCouponPacksFromTypedChoices(allowedTypedChoices)

        return {
          action: a,
          type: couponType,
          bonusPoints: data?.points ?? 0,
          allowedPackOldVicarIDs: allowedTypedChoices.map(c => c.id),
          allowedPackIds: allowedPacks.map(p => p.id),
          allowedPacks
        }
      }
    })


  const addBackgroundPointsCoupons = actions
    .filter(a => a.type === "add_background_points")
    .map(a => {
      const data = a.data as any
      const points = data?.amount ?? 0
      const packOldVicarID = data?.backgroundId

      const allowedPackOldVicarIDs = typeof packOldVicarID === "number" ? [packOldVicarID] : []
      const allowedPacks = getBackgroundCouponPacks(allowedPackOldVicarIDs)

      return {
        action: a,
        type: "backgrounds" as const,
        bonusPoints: points,
        allowedPackOldVicarIDs,
        allowedPackIds: allowedPacks.map(p => p.id),
        allowedPacks
      }
    })
    .filter(c => c.bonusPoints > 0 && c.allowedPackIds.length > 0)

  return [...betweenCoupons, ...addBackgroundPointsCoupons]
})

function getBackgroundCouponPacks(allowedOldVicarIDs: number[]): V5TraitPack[] {
  return v5data.traitPacks?.filter(p =>
    p.type === "backgrounds" &&
    p.oldVicarID !== undefined &&
    allowedOldVicarIDs.includes(p.oldVicarID)
  ) ?? []
}

function getFlawCouponPacksFromTypedChoices(choices: OldVicarTypedChoice[]): V5TraitPack[] {
  if (!choices.length) return []

  return v5data.traitPacks?.filter(p => {
    if (p.oldVicarID === undefined) return false
    if (!packHasFlaws(p)) return false

    return choices.some(c => c.id === p.oldVicarID && c.type === p.type)
  }) ?? []
}


const backgroundCouponPackIds = computed(() => {
  const ids: string[] = []
  spendPointsBetweenActions.value.forEach(spa => {
    if (spa.type === 'backgrounds') {
      ids.push(...spa.allowedPackIds)
    }
  })
  return ids
})

const flawCouponPackIds = computed(() => {
  const ids: string[] = []
  spendPointsBetweenActions.value.forEach(spa => {
    if (spa.type === 'flaws') {
      ids.push(...spa.allowedPackIds)
    }
  })
  return ids
})

const pointsSpentOnBackgroundCouponPacks = computed(() => {
  let total = 0
  character.value.traitPackUsages?.forEach(usage => {
    if (backgroundCouponPackIds.value.includes(usage.packID)) {
      usage.traits?.forEach(t => {
        if (!t.isLocked) {
          total += t.customLevel ?? 0
        }
      })
    }
  })
  return total
})

const pointsSpentOnFlawCouponPacks = computed(() => {
  let total = 0
  character.value.traitPackUsages?.forEach(usage => {
    if (flawCouponPackIds.value.includes(usage.packID)) {
      usage.flawTraits?.forEach(ft => {
        if (!ft.isLocked) {
          const pack = v5data.traitPacks?.find(p => p.id === usage.packID)
          const trait = pack?.packTraits?.find(pt => pt.trait.id === ft.traitID)?.trait
          if (trait) {
            total += trait.level
          }
        }
      })
    }
  })
  return total
})

const totalBackgroundCouponPoints = computed(() => {
  return spendPointsBetweenActions.value
    .filter(spa => spa.type === 'backgrounds')
    .reduce((sum, spa) => sum + spa.bonusPoints, 0)
})

const totalFlawCouponPoints = computed(() => {
  return spendPointsBetweenActions.value
    .filter(spa => spa.type === 'flaws')
    .reduce((sum, spa) => sum + spa.bonusPoints, 0)
})

const backgroundCouponRemaining = computed(() => {
  return Math.max(0, totalBackgroundCouponPoints.value - pointsSpentOnBackgroundCouponPacks.value)
})

const flawCouponRemaining = computed(() => {
  return Math.max(0, totalFlawCouponPoints.value - pointsSpentOnFlawCouponPacks.value)
})

const backgroundCouponExcess = computed(() => {
  return Math.max(0, pointsSpentOnBackgroundCouponPacks.value - totalBackgroundCouponPoints.value)
})

const flawCouponExcess = computed(() => {
  return Math.max(0, pointsSpentOnFlawCouponPacks.value - totalFlawCouponPoints.value)
})

const predatorTraitActions = computed(() => {
  if (!selectedPredatorType.value) return { merits: [], backgrounds: [], flaws: [] }

  const merits: { action: V5PredatorAction; description: string }[] = []
  const backgrounds: { action: V5PredatorAction; description: string }[] = []
  const flaws: { action: V5PredatorAction; description: string }[] = []

  selectedPredatorType.value.actions.forEach(action => {
    if (action.type === 'add_merit') {
      merits.push({ action, description: action.description })
    } else if (action.type === 'add_background' || action.type === 'add_background_points') {
      backgrounds.push({ action, description: action.description })
    } else if (action.type === 'add_flaw') {
      flaws.push({ action, description: action.description })
    }
  })

  return { merits, backgrounds, flaws }
})

function applyPredatorTraitBonuses() {
  if (!selectedPredatorType.value || !v5data.traitPacks) return

  const bonusTraits: NonNullable<V5CharacterInternalData['predatorBonusesApplied']>['traits'] = []

  selectedPredatorType.value.actions.forEach(action => {
    const data = action.data as any

    if (action.type === 'add_merit' || action.type === 'add_background' || action.type === 'add_flaw') {
      const traitOldVicarID = data?.id
      const level = data?.level ?? 1

      if (traitOldVicarID) {
        for (const pack of v5data.traitPacks!) {
          const packTrait = pack.packTraits?.find(pt => pt.trait.oldVicarID === traitOldVicarID)
          if (packTrait) {
            const isFlaw = packTrait.trait.isFlaw || action.type === 'add_flaw'

            addLockedTrait(pack, packTrait.trait.id, level, undefined, isFlaw)

            bonusTraits.push({
              packID: pack.id,
              traitID: packTrait.trait.id,
              level,
              isFlaw
            })
            break
          }
        }
      }
    }
  })

  if (bonusTraits.length > 0) {
    internalData.value = {
      ...internalData.value,
      predatorBonusesApplied: {
        ...internalData.value.predatorBonusesApplied,
        traits: bonusTraits
      }
    }
  }
}

function revertPredatorTraitBonuses() {
  const appliedTraits = internalData.value.predatorBonusesApplied?.traits
  if (!appliedTraits?.length) return

  appliedTraits.forEach(bonus => {
    removeLockedTrait(bonus.packID, bonus.traitID, bonus.isFlaw)
  })

  internalData.value = {
    ...internalData.value,
    predatorBonusesApplied: {
      ...internalData.value.predatorBonusesApplied,
      traits: []
    }
  }
}

function addLockedTrait(pack: V5TraitPack, traitId: string, level: number, suffix?: string, isFlaw?: boolean) {
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

  const usage = character.value.traitPackUsages[usageIndex]!

  if (isFlaw) {
    if (!usage.flawTraits.some(t => t.traitID === traitId && t.isLocked)) {
      usage.flawTraits.push({
        id: `flaw-${traitId}-locked`,
        usageID: usage.id,
        traitID: traitId,
        isLocked: true,
        isManual: false,
        customLevel: level,
        suffix,
      })
    }
  } else {
    if (!usage.traits.some(t => t.traitID === traitId && t.isLocked)) {
      usage.traits.push({
        id: `trait-${traitId}-locked`,
        usageID: usage.id,
        traitID: traitId,
        isLocked: true,
        isManual: false,
        customLevel: level,
        suffix,
      })
    }
  }

  character.value.traitPackUsages = [...character.value.traitPackUsages]
}

function removeLockedTrait(packId: string, traitId: string, isFlaw?: boolean) {
  const usage = character.value.traitPackUsages?.find(u => u.packID === packId)
  if (!usage) return

  if (isFlaw) {
    const idx = usage.flawTraits.findIndex(t => t.traitID === traitId && t.isLocked)
    if (idx !== -1) usage.flawTraits.splice(idx, 1)
  } else {
    const idx = usage.traits.findIndex(t => t.traitID === traitId && t.isLocked)
    if (idx !== -1) usage.traits.splice(idx, 1)
  }

  character.value.traitPackUsages = [...(character.value.traitPackUsages ?? [])]
}

const totalBackgroundPointsRaw = computed(() => {
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

const totalBackgroundPoints = computed(() => {
  const couponCoverage = Math.min(pointsSpentOnBackgroundCouponPacks.value, totalBackgroundCouponPoints.value)
  return totalBackgroundPointsRaw.value - couponCoverage
})

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

const totalFlawPointsRaw = computed(() => {
  let total = 0
  character.value.traitPackUsages?.forEach(usage => {
    usage.flawTraits?.forEach(ft => {
      if (!ft.isLocked) {
        const pack = v5data.traitPacks?.find(p => p.id === usage.packID)
        const trait = pack?.packTraits?.find(pt => pt.trait.id === ft.traitID)?.trait
        if (trait) {
          total += trait.level
        }
      }
    })
  })
  return total
})

const totalFlawPoints = computed(() => {
  const couponCoverage = Math.min(pointsSpentOnFlawCouponPacks.value, totalFlawCouponPoints.value)
  return totalFlawPointsRaw.value - couponCoverage
})

const bonusMeritPoints = computed(() => {
  return Math.max(0, totalFlawPoints.value - MIN_FLAW_POINTS)
})

const availableBackgroundPoints = computed(() => {
  return BASE_BACKGROUND_POINTS
})

const availableMeritPoints = computed(() => {
  return bonusMeritPoints.value
})

const availableFlawPoints = computed(() => {
  return MIN_FLAW_POINTS
})

const backgroundPointsRemaining = computed(() => {
  return availableBackgroundPoints.value - totalBackgroundPoints.value
})

const meritPointsRemaining = computed(() => {
  return availableMeritPoints.value - totalMeritPoints.value
})

const flawPointsRemaining = computed(() => {
  return availableFlawPoints.value - totalFlawPoints.value
})

const distributionValid = computed(() => {
  return backgroundPointsRemaining.value >= 0 &&
         meritPointsRemaining.value >= 0 &&
         flawPointsRemaining.value <= 0
})

function isBackgroundCouponPack(packId: string): boolean {
  return backgroundCouponPackIds.value.includes(packId)
}

function isFlawCouponPack(packId: string): boolean {
  return flawCouponPackIds.value.includes(packId)
}

const backgroundCouponPackNames = computed(() => {
  const names = spendPointsBetweenActions.value
    .filter(spa => spa.type === 'backgrounds')
    .flatMap(spa => spa.allowedPacks.map(p => p.name))
  return [...new Set(names)]
})

const flawCouponPackNames = computed(() => {
  const names = spendPointsBetweenActions.value
    .filter(spa => spa.type === 'flaws')
    .flatMap(spa => spa.allowedPacks.map(p => p.name))
  return [...new Set(names)]
})


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

function getPackUsage(packId: string): V5CharacterTraitPackUsage | undefined {
  return character.value.traitPackUsages?.find(u => u.packID === packId)
}

function isTraitSelected(packId: string, traitId: string): boolean {
  const usage = getPackUsage(packId)
  return usage?.traits?.some(t => t.traitID === traitId) ?? false
}

function isTraitLocked(packId: string, traitId: string): boolean {
  const usage = getPackUsage(packId)
  return usage?.traits?.some(t => t.traitID === traitId && t.isLocked) ?? false
}

function isFlawSelected(packId: string, traitId: string): boolean {
  const usage = getPackUsage(packId)
  return usage?.flawTraits?.some(t => t.traitID === traitId) ?? false
}

function isFlawLocked(packId: string, traitId: string): boolean {
  const usage = getPackUsage(packId)
  return usage?.flawTraits?.some(t => t.traitID === traitId && t.isLocked) ?? false
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

function selectTrait(pack: V5TraitPack, trait: V5Trait, level: number) {
  if (trait.isFlaw) {
    toggleFlaw(pack, trait.id)
    return
  }

  if (isTraitLocked(pack.id, trait.id)) return

  if (trait.isRepeatable) {
    currentSuffixTrait.value = { pack, trait, level }
    suffixInput.value = ''
    showSuffixModal.value = true
    return
  }

  applyTrait(pack, trait.id, level, undefined)
}

function editTraitSuffix(pack: V5TraitPack, trait: V5Trait) {
  if (isTraitLocked(pack.id, trait.id)) return

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

  const usage = character.value.traitPackUsages[usageIndex]!

  const traitIndex = isUpdate
    ? usage.traits.findIndex(t => t.traitID === traitId && !t.isLocked)
    : usage.traits.findIndex(t => t.traitID === traitId && t.suffix === suffix && !t.isLocked)

  if (traitIndex !== -1) {
    if (isUpdate) {
      usage.traits[traitIndex] = {
        ...usage.traits[traitIndex],
        customLevel: level,
        suffix: suffix,
      } as any
    } else {
      const existingLevel = usage.traits[traitIndex]!.customLevel
      if (existingLevel === level) {
        usage.traits.splice(traitIndex, 1)
      } else {
        usage.traits[traitIndex] = {
          ...usage.traits[traitIndex],
          customLevel: level,
        } as any
      }
    }
  } else {
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
  const isUpdate = isTraitSelected(pack.id, trait.id)
  applyTrait(pack, trait.id, level, suffixInput.value.trim() || undefined, isUpdate)

  showSuffixModal.value = false
  currentSuffixTrait.value = null
  suffixInput.value = ''
}

function toggleFlaw(pack: V5TraitPack, traitId: string) {
  if (isFlawLocked(pack.id, traitId)) return

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

  const usage = character.value.traitPackUsages[usageIndex]!
  const flawIndex = usage.flawTraits.findIndex(t => t.traitID === traitId && !t.isLocked)

  if (flawIndex !== -1) {
    usage.flawTraits.splice(flawIndex, 1)
  } else {
    usage.flawTraits.push({
      id: `flaw-${traitId}`,
      usageID: usage.id,
      traitID: traitId,
      isLocked: false,
      isManual: true,
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

function normalizeOldVicarChoiceIds(raw: any): number[] {
  if (!raw) return []
  if (Array.isArray(raw)) {
    if (raw.every(v => typeof v === "number")) return raw as number[]
    return raw
      .map(v => (typeof v === "object" && v && typeof v.id === "number" ? v.id : null))
      .filter((v): v is number => typeof v === "number")
  }
  return []
}

type OldVicarTypedChoice = { type: 'backgrounds' | 'merits'; id: number }

function normalizeOldVicarTypedChoices(raw: any): OldVicarTypedChoice[] {
  if (!raw || !Array.isArray(raw)) return []

  if (raw.every(v => typeof v === 'number')) {
    return (raw as number[]).map(id => ({ type: 'backgrounds', id }))
  }

  return raw
    .map(v => {
      if (!v || typeof v !== 'object') return null
      const id = typeof v.id === 'number' ? v.id : null

      const type = (v.type === 'backgrounds' || v.type === 'merits')
        ? (v.type as OldVicarTypedChoice['type'])
        : null

      if (!id || !type) return null
      return { type, id }
    })
    .filter((v): v is OldVicarTypedChoice => !!v)
}


function packHasFlaws(pack: V5TraitPack): boolean {
  return pack.packTraits?.some(pt => pt.trait.isFlaw) ?? false
}
</script>

<template>
  <div class="traits-step">
    <VCard>
      <h2>Punkte-Übersicht</h2>
      <p class="info-text">
        Du hast {{ BASE_BACKGROUND_POINTS }} Punkte für Hintergründe. Mindestens {{ MIN_FLAW_POINTS }} Punkte musst du in Schwächen investieren.
        Für jeden zusätzlichen Schwächenpunkt erhältst du einen weiteren Vorzugspunkt.
      </p>

      <div class="points-overview" :class="{ 'points-overview--valid': distributionValid }">
        <div class="points-row">
          <span class="points-label">Hintergründe:</span>
          <span class="points-value" :class="{ 'points-value--ok': backgroundPointsRemaining >= 0, 'points-value--over': backgroundPointsRemaining < 0 }">
            {{ totalBackgroundPoints }} / {{ availableBackgroundPoints }}
          </span>
        </div>
        <div class="points-row">
          <span class="points-label">Schwächen:</span>
          <span class="points-value" :class="{ 'points-value--ok': flawPointsRemaining <= 0, 'points-value--under': flawPointsRemaining > 0 }">
            {{ totalFlawPoints }} / {{ availableFlawPoints }}+
          </span>
        </div>
        <div v-if="bonusMeritPoints > 0" class="points-row points-row--bonus">
          <span class="points-label">Vorzüge (Bonus):</span>
          <span class="points-value" :class="{ 'points-value--ok': meritPointsRemaining >= 0, 'points-value--over': meritPointsRemaining < 0 }">
            {{ totalMeritPoints }} / {{ availableMeritPoints }}
          </span>
        </div>
      </div>

      <div v-if="totalBackgroundCouponPoints > 0 || totalFlawCouponPoints > 0" class="coupon-section">
        <h3>Jagdverhalten-Gutscheine</h3>
        <p class="info-text">
          Diese Punkte kannst du kostenfrei in bestimmten Hintergründen/Schwächen ausgeben.
          Überschreitest du den Gutscheinwert, zählen die zusätzlichen Punkte gegen dein normales Budget.
        </p>

        <div v-if="totalBackgroundCouponPoints > 0" class="coupon-item">
          <div class="coupon-header">
            <span class="coupon-icon">🎫</span>
            <span class="coupon-title">Hintergrund-Gutschein</span>
          </div>
          <div class="coupon-details">
            <span class="coupon-packs">Gültig für: {{ backgroundCouponPackNames.join(', ') }}</span>
            <span class="coupon-value" :class="{ 'coupon-value--used': backgroundCouponRemaining === 0 }">
              {{ pointsSpentOnBackgroundCouponPacks }} / {{ totalBackgroundCouponPoints }} verwendet
            </span>
          </div>
          <div v-if="backgroundCouponRemaining > 0" class="coupon-remaining">
            {{ backgroundCouponRemaining }} Gratispunkt(e) übrig
          </div>
          <div v-if="backgroundCouponExcess > 0" class="coupon-excess">
            {{ backgroundCouponExcess }} Punkt(e) über Gutschein (zählt gegen Budget)
          </div>
        </div>

        <div v-if="totalFlawCouponPoints > 0" class="coupon-item coupon-item--flaw">
          <div class="coupon-header">
            <span class="coupon-icon">🎫</span>
            <span class="coupon-title">Schwächen-Gutschein</span>
          </div>
          <div class="coupon-details">
            <span class="coupon-packs">Gültig für: {{ flawCouponPackNames.join(', ') }}</span>
            <span class="coupon-value" :class="{ 'coupon-value--used': flawCouponRemaining === 0 }">
              {{ pointsSpentOnFlawCouponPacks }} / {{ totalFlawCouponPoints }} verwendet
            </span>
          </div>
          <div v-if="flawCouponRemaining > 0" class="coupon-remaining">
            {{ flawCouponRemaining }} Gratispunkt(e) übrig
          </div>
          <div v-if="flawCouponExcess > 0" class="coupon-excess">
            {{ flawCouponExcess }} Punkt(e) über Gutschein (zählt gegen Budget)
          </div>
        </div>
      </div>
    </VCard>

    <VCard v-if="predatorTraitActions.merits.length > 0 || predatorTraitActions.backgrounds.length > 0 || predatorTraitActions.flaws.length > 0">
      <h2>Jagdverhalten-Boni</h2>
      <p class="info-text">Diese Vorzüge und Schwächen erhältst du durch dein Jagdverhalten automatisch.</p>

      <div class="predator-traits">
        <div v-if="predatorTraitActions.merits.length > 0" class="predator-section">
          <h4>Vorzüge:</h4>
          <ul class="predator-list">
            <li v-for="(item, idx) in predatorTraitActions.merits" :key="`merit-${idx}`">
              {{ item.description }}
            </li>
          </ul>
        </div>
        <div v-if="predatorTraitActions.backgrounds.length > 0" class="predator-section">
          <h4>Hintergründe:</h4>
          <ul class="predator-list">
            <li v-for="(item, idx) in predatorTraitActions.backgrounds" :key="`bg-${idx}`">
              {{ item.description }}
            </li>
          </ul>
        </div>
        <div v-if="predatorTraitActions.flaws.length > 0" class="predator-section predator-section--flaw">
          <h4>Schwächen:</h4>
          <ul class="predator-list">
            <li v-for="(item, idx) in predatorTraitActions.flaws" :key="`flaw-${idx}`">
              {{ item.description }}
            </li>
          </ul>
        </div>
      </div>
    </VCard>

    <VCard v-if="bonusMeritPoints > 0">
      <h2>Vorzüge</h2>
      <p class="info-text">Du hast {{ bonusMeritPoints }} Bonuspunkte durch zusätzliche Schwächen.</p>

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
              v-for="packTrait in pack.packTraits"
              :key="packTrait.id"
              class="trait-item"
              :class="{
                'trait-item--selected': packTrait.trait.isFlaw ? isFlawSelected(pack.id, packTrait.trait.id) : isTraitSelected(pack.id, packTrait.trait.id),
                'trait-item--locked': packTrait.trait.isFlaw ? isFlawLocked(pack.id, packTrait.trait.id) : isTraitLocked(pack.id, packTrait.trait.id),
                'trait-item--flaw': packTrait.trait.isFlaw
              }"
            >
              <div class="trait-info">
                <div class="trait-header">
                  <span class="trait-name">{{ packTrait.trait.name }}</span>
                  <span v-if="packTrait.trait.isFlaw" class="type-badge badge--flaw">Schwäche</span>
                  <span v-else class="type-badge badge--merit">Vorzug</span>
                  <span v-if="packTrait.trait.isFlaw" class="trait-level-badge trait-level-badge--flaw">{{ packTrait.trait.level }}</span>
                  <span v-if="getSelectedTraitSuffix(pack.id, packTrait.trait.id)" class="trait-suffix">
                    ({{ getSelectedTraitSuffix(pack.id, packTrait.trait.id) }})
                  </span>
                  <span v-if="packTrait.trait.isFlaw ? isFlawLocked(pack.id, packTrait.trait.id) : isTraitLocked(pack.id, packTrait.trait.id)" class="locked-badge">🔒 Jagdverh.</span>
                </div>
                <p class="trait-desc">{{ packTrait.trait.description }}</p>
              </div>
              <div class="trait-actions">
                <template v-if="packTrait.trait.isFlaw">
                  <button
                    type="button"
                    class="flaw-toggle-btn"
                    :class="{
                      'flaw-toggle-btn--selected': isFlawSelected(pack.id, packTrait.trait.id),
                      'flaw-toggle-btn--locked': isFlawLocked(pack.id, packTrait.trait.id)
                    }"
                    :disabled="isFlawLocked(pack.id, packTrait.trait.id)"
                    @click="toggleFlaw(pack, packTrait.trait.id)"
                  >
                    {{ isFlawSelected(pack.id, packTrait.trait.id) ? '✓' : '+' }}
                  </button>
                </template>
                <template v-else>
                  <button
                    v-if="isTraitSelected(pack.id, packTrait.trait.id) && !isTraitLocked(pack.id, packTrait.trait.id)"
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
                      :disabled="isTraitLocked(pack.id, packTrait.trait.id)"
                      @click="selectTrait(pack, packTrait.trait, level)"
                    >
                      {{ level }}
                    </button>
                  </div>
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>
    </VCard>

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
            :class="{ 'pack-header--coupon': isBackgroundCouponPack(pack.id) || isFlawCouponPack(pack.id) }"
            @click="toggleBackgroundPack(pack.id)"
          >
            <div class="pack-info">
              <h3>
                {{ pack.name }}
                <span v-if="isBackgroundCouponPack(pack.id) || isFlawCouponPack(pack.id)" class="coupon-badge">🎫 Gutschein</span>
              </h3>
              <p v-if="pack.description" class="pack-desc">{{ pack.description }}</p>
            </div>
            <span class="pack-toggle">{{ expandedBackgroundPacks.has(pack.id) ? '−' : '+' }}</span>
          </button>

          <div v-if="expandedBackgroundPacks.has(pack.id)" class="pack-traits">
            <div
              v-for="packTrait in pack.packTraits"
              :key="packTrait.id"
              class="trait-item"
              :class="{
                'trait-item--selected': packTrait.trait.isFlaw ? isFlawSelected(pack.id, packTrait.trait.id) : isTraitSelected(pack.id, packTrait.trait.id),
                'trait-item--locked': packTrait.trait.isFlaw ? isFlawLocked(pack.id, packTrait.trait.id) : isTraitLocked(pack.id, packTrait.trait.id),
                'trait-item--coupon': isBackgroundCouponPack(pack.id) || isFlawCouponPack(pack.id),
                'trait-item--flaw': packTrait.trait.isFlaw
              }"
            >
              <div class="trait-info">
                <div class="trait-header">
                  <span class="trait-name">{{ packTrait.trait.name }}</span>
                  <span v-if="packTrait.trait.isFlaw" class="type-badge badge--flaw">Schwäche</span>
                  <span v-else class="type-badge badge--background">Hintergrund</span>
                  <span v-if="packTrait.trait.isFlaw" class="trait-level-badge trait-level-badge--flaw">{{ packTrait.trait.level }}</span>
                  <span v-if="getSelectedTraitSuffix(pack.id, packTrait.trait.id)" class="trait-suffix">
                    ({{ getSelectedTraitSuffix(pack.id, packTrait.trait.id) }})
                  </span>
                  <span v-if="packTrait.trait.isFlaw ? isFlawLocked(pack.id, packTrait.trait.id) : isTraitLocked(pack.id, packTrait.trait.id)" class="locked-badge">🔒 Jagdverh.</span>
                </div>
                <p class="trait-desc">{{ packTrait.trait.description }}</p>
              </div>
              <div class="trait-actions">
                <template v-if="packTrait.trait.isFlaw">
                  <button
                    type="button"
                    class="flaw-toggle-btn"
                    :class="{
                      'flaw-toggle-btn--selected': isFlawSelected(pack.id, packTrait.trait.id),
                      'flaw-toggle-btn--locked': isFlawLocked(pack.id, packTrait.trait.id)
                    }"
                    :disabled="isFlawLocked(pack.id, packTrait.trait.id)"
                    @click="toggleFlaw(pack, packTrait.trait.id)"
                  >
                    {{ isFlawSelected(pack.id, packTrait.trait.id) ? '✓' : '+' }}
                  </button>
                </template>
                <template v-else>
                  <button
                    v-if="isTraitSelected(pack.id, packTrait.trait.id) && !isTraitLocked(pack.id, packTrait.trait.id)"
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
                      :disabled="isTraitLocked(pack.id, packTrait.trait.id)"
                      @click="selectTrait(pack, packTrait.trait, level)"
                    >
                      {{ level }}
                    </button>
                  </div>
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>
    </VCard>

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

  strong {
    color: $text-0;
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

.coupon-section {
  margin-top: $s-4;
  padding-top: $s-4;
  border-top: 1px dashed $border;

  h3 {
    margin: 0 0 $s-2;
    font-family: $font-head;
    font-size: $fs-1;
    color: $text-0;
  }

  .info-text {
    margin-bottom: $s-3;
  }
}

.coupon-item {
  padding: $s-3;
  border-radius: $r-md;
  background: linear-gradient(135deg, rgba(255, 200, 50, 0.08), rgba(255, 180, 50, 0.04));
  border: 1px dashed rgba(255, 200, 50, 0.4);
  margin-bottom: $s-3;

  &:last-child {
    margin-bottom: 0;
  }

  &--flaw {
    background: linear-gradient(135deg, rgba(255, 100, 100, 0.08), rgba(255, 80, 80, 0.04));
    border-color: rgba(255, 100, 100, 0.4);
  }
}

.coupon-header {
  display: flex;
  align-items: center;
  gap: $s-2;
  margin-bottom: $s-2;
}

.coupon-icon {
  font-size: 1.2rem;
}

.coupon-title {
  font-weight: 600;
  color: $text-0;
}

.coupon-details {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-3;
  flex-wrap: wrap;
}

.coupon-packs {
  font-size: 0.85rem;
  color: $text-2;
}

.coupon-value {
  font-weight: 600;
  color: rgba(255, 200, 50, 0.9);
  font-size: 0.9rem;

  &--used {
    color: rgba(100, 200, 100, 0.9);
  }
}

.coupon-remaining {
  margin-top: $s-2;
  font-size: 0.85rem;
  color: rgba(255, 200, 50, 0.9);
  font-weight: 500;
}

.coupon-excess {
  margin-top: $s-1;
  font-size: 0.85rem;
  color: $red-0;
  font-weight: 500;
}

.coupon-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px $s-2;
  margin-left: $s-2;
  border-radius: $r-sm;
  background: rgba(255, 200, 50, 0.15);
  color: rgba(255, 200, 50, 0.9);
  font-size: 0.7rem;
  font-weight: 600;
  vertical-align: middle;
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

  &--flaw {
    border-color: rgba(255, 100, 100, 0.2);
    background: rgba(255, 100, 100, .04);
  }
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

  &--coupon {
    background: rgba(255, 200, 50, 0.04);

    &:hover {
      background: rgba(255, 200, 50, 0.08);
    }
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

  &--locked {
    opacity: 0.8;
    background: rgba(255, 59, 84, .04);
    border-color: rgba(255, 59, 84, 0.3);
  }

  &--coupon:not(&--selected) {
    border-color: rgba(255, 200, 50, 0.2);
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

.type-badge {
  padding: 2px $s-2;
  border-radius: $r-sm;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.3px;

  &.badge--merit {
    background: rgba(100, 200, 100, .15);
    color: rgba(100, 200, 100, 0.9);
  }

  &.badge--background {
    background: rgba(100, 150, 255, .15);
    color: rgba(100, 150, 255, 0.9);
  }

  &.badge--flaw {
    background: rgba(255, 100, 100, .15);
    color: rgba(255, 100, 100, 0.9);
  }
}

.locked-badge {
  padding: 2px $s-2;
  border-radius: $r-sm;
  font-size: 0.7rem;
  font-weight: 500;
  background: rgba(255, 59, 84, .12);
  color: $red-0;
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

  &:hover:not(:disabled) {
    border-color: $border-strong;
    background: rgba(255, 255, 255, .04);
  }

  &--selected {
    border-color: $red-0;
    background: rgba(255, 59, 84, .12);
    color: $text-0;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
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

  &:hover:not(:disabled) {
    border-color: rgba(255, 100, 100, 0.5);
    background: rgba(255, 100, 100, .08);
  }

  &--selected {
    border-color: rgba(255, 100, 100, 0.5);
    background: rgba(255, 100, 100, .15);
    color: rgba(255, 100, 100, 0.9);
  }

  &--locked, &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

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
