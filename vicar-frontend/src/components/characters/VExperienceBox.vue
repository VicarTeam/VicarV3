<script setup lang="ts">
import {computed, ref} from "vue"
import type {V5Character} from "@/@types/v5"

const character = defineModel<Partial<V5Character>>({ required: true })
const emit = defineEmits<{
  "update-xp": [value: number]
}>()

const raw = ref("")
const focused = ref(false)

const xp = computed<number>({
  get: () => (character.value as any).exp ?? 0,
  set: (v) => {
    ;(character.value as any).exp = Math.max(0, Math.floor(Number(v) || 0))
    emit("update-xp", (character.value as any).exp)
  }
})

function evaluateExpression(input: string, base: number): number {
  let s = String(input ?? "").trim()
  if (!s) return base

  s = s.replace(/,/g, ".")
  s = s.replace(/[^0-9.+\-*/() ]/g, "")
  s = s.replace(/\s+/g, "")

  if (!s) return base

  if (/^[+\-*/]/.test(s)) s = `${base}${s}`

  try {
    const fn = new Function(`return (${s});`)
    const result = Number(fn())
    if (!Number.isFinite(result)) return base
    return Math.max(0, Math.floor(result))
  } catch {
    return base
  }
}

function commit() {
  const base = xp.value
  xp.value = evaluateExpression(raw.value, base)
  raw.value = ""
}

function onFocus() {
  raw.value = String(xp.value ?? 0)
  focused.value = true
}

function onBlur() {
  commit()
  focused.value = false
}

function onKeydown(ev: KeyboardEvent) {
  if (ev.key === "Enter") {
    ev.preventDefault()
    commit()
    ;(ev.target as HTMLInputElement)?.blur?.()
  }
}

function inc() {
  if (!raw.value.trim()) {
    xp.value = xp.value + 1
    return
  }
  raw.value = raw.value.trim().endsWith("+") ? raw.value : `${raw.value}+`
  commit()
}

function dec() {
  if (!raw.value.trim()) {
    xp.value = Math.max(0, xp.value - 1)
    return
  }
  raw.value = raw.value.trim().endsWith("-") ? raw.value : `${raw.value}-`
  commit()
}
</script>

<template>
  <div class="xp" style="margin-top: 1rem">
    <div class="xp__head">
      <h2>Erfahrung</h2>
      <div class="xp__stepper">
        <button type="button" class="xp__btn" @click="dec" aria-label="XP minus">−</button>
        <button type="button" class="xp__btn" @click="inc" aria-label="XP plus">+</button>
      </div>
    </div>

    <div class="xp__row">
      <input
        id="xpInput"
        class="xp__input"
        type="text"
        inputmode="numeric"
        autocomplete="off"
        :value="raw"
        @input="raw = ($event.target as HTMLInputElement).value"
        @focus="onFocus"
        @blur="onBlur"
        @keydown="onKeydown"
      />
      <span v-if="!focused" class="xp__text">{{xp}}</span>
      <div class="xp__hint">Du kannst rechnen: <code>+5</code>, <code>-2</code>, <code>10+3</code>, <code>7*2</code> …</div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.xp__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $s-3;

  h2 {
    margin: 0;
    font-family: $font-head;
    font-size: $fs-2;
  }
}

.xp__stepper {
  display: flex;
  gap: $s-2;
}

.xp__btn {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
  color: $text-0;
  cursor: pointer;
  font-weight: 900;
  font-size: 1.1rem;
  display: grid;
  place-items: center;
  transition: transform $t-fast $ease, border-color $t-fast $ease;

  &:hover {
    transform: translateY(-1px);
    border-color: rgba($red-0, 0.35);
  }
}

.xp__row {
  display: flex;
  flex-direction: column;
  gap: $s-2;
  margin-top: $s-3;
  position: relative;
}

.xp__text {
  position: absolute;
  top: 0.5rem;
  left: 0.75rem;
}

.xp__label {
  color: $text-2;
  font-weight: 600;
}

.xp__input {
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, 0.02);
  color: $text-0;
  font-family: $font-body;
  font-size: 1rem;

  &:focus {
    outline: none;
    border-color: $red-0;
  }
}

.xp__current {
  font-size: 0.9rem;
  color: $text-2;

  strong {
    color: $text-0;
    font-family: $font-head;
  }
}

.xp__hint {
  color: $text-2;
  font-size: 0.85rem;

  code {
    padding: 1px 6px;
    border-radius: $r-sm;
    border: 1px solid rgba(255, 255, 255, 0.10);
    background: rgba(255, 255, 255, 0.03);
    color: $text-1;
  }
}
</style>
