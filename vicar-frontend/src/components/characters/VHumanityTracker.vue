<script setup lang="ts">
import { computed } from "vue"

const props = defineProps<{
  humanity: number
  stains: number
  readonly?: boolean
}>()

const emit = defineEmits<{
  "update:humanity": [value: number]
  "update:stains": [value: number]
}>()

type BoxType = "filled" | "stain" | "empty"

const boxes = computed(() => {
  const result: { value: number; type: BoxType }[] = []

  const clampedHumanity = Math.max(0, Math.min(10, props.humanity ?? 0))
  const maxStains = Math.max(0, Math.min(10 - clampedHumanity, props.stains ?? 0))

  for (let value = 1; value <= 10; value++) {
    if (value <= clampedHumanity) {
      result.push({ value, type: "filled" })
    } else if (value > 10 - maxStains) {
      result.push({ value, type: "stain" })
    } else {
      result.push({ value, type: "empty" })
    }
  }

  return result
})

function setHumanity(value: number) {
  const v = Math.max(0, Math.min(10, value))
  emit("update:humanity", v)

  const maxStains = Math.max(0, 10 - v)
  if ((props.stains ?? 0) > maxStains) emit("update:stains", maxStains)
}

function setStains(value: number) {
  const maxStains = Math.max(0, 10 - (props.humanity ?? 0))
  const v = Math.max(0, Math.min(maxStains, value))
  emit("update:stains", v)
}

function toggleBox(boxValue: number, ev?: MouseEvent) {
  if (props.readonly) return

  const isShift = !!ev?.shiftKey
  if (isShift) {
    if (boxValue > (props.humanity ?? 0)) {
      const currentMaxStains = Math.max(0, 10 - (props.humanity ?? 0))
      const stainIndexFromRight = boxValue - (10 - currentMaxStains)
      const next = (props.stains ?? 0) === stainIndexFromRight ? 0 : stainIndexFromRight
      setStains(Math.min(next, currentMaxStains))
    }
    return
  }

  const current = props.humanity ?? 0
  if (boxValue === current) setHumanity(current - 1)
  else setHumanity(boxValue)
}

function onContext(boxValue: number, ev: MouseEvent) {
  if (props.readonly) return
  ev.preventDefault()

  if (boxValue > (props.humanity ?? 0)) {
    const currentMaxStains = Math.max(0, 10 - (props.humanity ?? 0))
    const stainIndexFromRight = boxValue - (10 - currentMaxStains)
    setStains(Math.min(stainIndexFromRight, currentMaxStains))
  }
}

function incHumanity() {
  if (props.readonly) return
  setHumanity((props.humanity ?? 0) + 1)
}

function decHumanity() {
  if (props.readonly) return
  setHumanity((props.humanity ?? 0) - 1)
}

function incStains() {
  if (props.readonly) return
  setStains((props.stains ?? 0) + 1)
}

function decStains() {
  if (props.readonly) return
  setStains((props.stains ?? 0) - 1)
}
</script>

<template>
  <div class="humanity">
    <div class="humanity__scale" role="group" aria-label="Menschlichkeit">
      <button
        v-for="(box, idx) in boxes"
        :key="idx"
        type="button"
        class="humanity__box"
        :class="{
          'humanity__box--filled': box.type === 'filled',
          'humanity__box--stain': box.type === 'stain',
          'humanity__box--readonly': readonly
        }"
        :disabled="readonly"
        @click="toggleBox(box.value, $event)"
        @contextmenu="onContext(box.value, $event)"
        :aria-label="`Wert ${box.value}`"
      >
        <span v-if="box.type === 'stain'" class="humanity__stainMark" aria-hidden="true" />
      </button>
    </div>

    <div class="humanity__meta">
      <div class="humanity__labels">
        <span class="humanity__value">{{ humanity }}</span>
        <span v-if="stains > 0" class="humanity__stains">{{ stains }} Makel</span>
      </div>

      <div v-if="!readonly" class="humanity__controls" aria-label="Steuerung">
        <button type="button" class="humanity__ctl" @click="decHumanity">Menschlichkeit −</button>
        <button type="button" class="humanity__ctl" @click="incHumanity">Menschlichkeit +</button>
        <button type="button" class="humanity__ctl humanity__ctl--warn" @click="decStains">Makel −</button>
        <button type="button" class="humanity__ctl humanity__ctl--warn" @click="incStains">Makel +</button>
      </div>

      <div v-if="!readonly" class="humanity__hint">
        Tipp: Shift-Klick oder Rechtsklick setzt Makel.
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.humanity {
  display: grid;
  gap: $s-3;
}

.humanity__scale {
  display: flex;
  flex-wrap: wrap;
  gap: $s-2;
  align-items: center;
}

.humanity__box {
  width: 30px;
  height: 30px;
  border-radius: 999px;
  border: 2px solid rgba($red-0, 0.16);
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  display: grid;
  place-items: center;
  transition: transform $t-fast $ease, border-color $t-fast $ease, background $t-fast $ease;

  &:hover:not(:disabled) {
    transform: translateY(-1px);
    border-color: rgba($grey-0, 0.55);
  }
}

.humanity__box--filled {
  border-color: rgba($red-0, 0.85);
  background: rgba($red-0, 0.9);
  box-shadow: 0 10px 26px rgba(0, 0, 0, 0.35);
}

.humanity__box--stain {
  border-color: rgba($grey-0, 0.7);
  background: rgba($grey-0, 0.14);
}

.humanity__stainMark {
  width: 14px;
  height: 14px;
  border-radius: 999px;
  background: rgba($grey-0, 0.95);
  box-shadow: 0 0 12px rgba($grey-0, 0.35);
}

.humanity__box--readonly {
  cursor: default;
  &:hover {
    transform: none;
  }
}

.humanity__meta {
  display: grid;
  gap: $s-2;
}

.humanity__labels {
  display: flex;
  gap: $s-3;
  align-items: center;
}

.humanity__value {
  font-size: 1.1rem;
  font-weight: 700;
  font-family: $font-head;
  color: $text-0;
}

.humanity__stains {
  font-size: 0.85rem;
  color: $red-0;
  padding: $s-1 $s-2;
  background: rgba($red-0, 0.1);
  border-radius: $r-sm;
  border: 1px solid rgba($red-0, 0.18);
}

.humanity__controls {
  display: grid;
  gap: $s-2;
  grid-template-columns: 1fr 1fr;
}

.humanity__ctl {
  appearance: none;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
  color: $text-1;
  padding: 0.55rem 0.7rem;
  border-radius: 999px;
  cursor: pointer;
  transition: transform $t-fast $ease, border-color $t-fast $ease;

  &:hover {
    transform: translateY(-1px);
    border-color: rgba(255, 255, 255, 0.18);
    color: $red-0;
  }
}

.humanity__ctl--warn {
  border-color: rgba($grey-0, 0.25);

  &:hover {
    border-color: rgba($grey-0, 0.45);
  }
}

.humanity__hint {
  font-size: 0.85rem;
  color: $text-2;
}

@media (max-width: 420px) {
  .humanity__controls {
    grid-template-columns: 1fr;
  }
}
</style>
