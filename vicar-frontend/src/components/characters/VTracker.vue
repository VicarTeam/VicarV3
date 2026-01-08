<script setup lang="ts">
import { computed } from "vue"

type DamageKind = "" | "superficial" | "aggravated"

const props = defineProps<{
  max: number
  damage: string[]
  readonly?: boolean
  label?: string
}>()

const emit = defineEmits<{
  "update:damage": [value: string[]]
}>()

const padded = computed<DamageKind[]>(() => {
  const max = Math.max(0, props.max ?? 0)
  const src = Array.isArray(props.damage) ? props.damage : []
  const result: DamageKind[] = new Array(max).fill("")
  for (let i = 0; i < max; i++) {
    const v = src[i]
    if (v === "superficial" || v === "aggravated") result[i] = v
  }
  return result
})

const boxes = computed(() => {
  return padded.value.map((type) => ({ type }))
})

function setAt(index: number, kind: DamageKind) {
  const next = [...padded.value]
  next[index] = kind
  emit("update:damage", next)
}

function cycleBox(index: number) {
  if (props.readonly) return
  const cur = padded.value[index]
  if (cur === "") setAt(index, "superficial")
  else if (cur === "superficial") setAt(index, "aggravated")
  else setAt(index, "")
}

function reverseCycle(index: number, ev: MouseEvent) {
  if (props.readonly) return
  ev.preventDefault()
  const cur = padded.value[index]
  if (cur === "") setAt(index, "aggravated")
  else if (cur === "aggravated") setAt(index, "superficial")
  else setAt(index, "")
}
</script>

<template>
  <div class="tracker">
    <span v-if="label" class="tracker__label">{{ label }}</span>

    <div class="tracker__boxes" role="group">
      <button
        v-for="(box, idx) in boxes"
        :key="idx"
        type="button"
        class="tracker__box"
        :class="{
          'tracker__box--superficial': box.type === 'superficial',
          'tracker__box--aggravated': box.type === 'aggravated',
          'tracker__box--readonly': readonly
        }"
        :disabled="readonly"
        @click="cycleBox(idx)"
        @contextmenu="reverseCycle(idx, $event)"
        :aria-label="`Feld ${idx + 1}`"
      >
        <svg v-if="box.type === 'superficial'" viewBox="0 0 20 20" class="tracker__mark">
          <line x1="4" y1="16" x2="16" y2="4" />
        </svg>
        <svg v-else-if="box.type === 'aggravated'" viewBox="0 0 20 20" class="tracker__mark tracker__mark--aggravated">
          <line x1="4" y1="16" x2="16" y2="4" />
          <line x1="4" y1="4" x2="16" y2="16" />
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.tracker {
  display: grid;
  gap: $s-2;
}

.tracker__label {
  font-size: 0.85rem;
  color: $text-2;
  text-transform: uppercase;
  letter-spacing: 0.6px;
}

.tracker__boxes {
  display: flex;
  gap: $s-2;
  flex-wrap: wrap;
}

.tracker__box {
  width: 28px;
  height: 28px;
  border-radius: $r-sm;
  border: 2px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  display: grid;
  place-items: center;
  transition: transform $t-fast $ease, border-color $t-fast $ease, background $t-fast $ease;

  &:hover:not(:disabled) {
    transform: translateY(-1px);
    border-color: rgba($red-0, 0.45);
  }
}

.tracker__box--superficial {
  border-color: rgba($text-1, 0.35);
}

.tracker__box--aggravated {
  border-color: rgba($red-0, 0.55);
  background: rgba($red-0, 0.10);
}

.tracker__box--readonly {
  cursor: default;
  &:hover {
    transform: none;
  }
}

.tracker__mark {
  width: 14px;
  height: 14px;

  line {
    stroke: rgba(255, 255, 255, 0.72);
    stroke-width: 2.2;
    stroke-linecap: round;
  }
}

.tracker__mark--aggravated line {
  stroke: rgba($red-0, 0.95);
}

@media (max-width: 420px) {
  .tracker__box {
    width: 30px;
    height: 30px;
  }
}
</style>
