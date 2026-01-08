<script setup lang="ts">
import { computed } from "vue"
import type { V5Character } from "@/@types/v5"

const character = defineModel<Partial<V5Character>>({ required: true })
const emit = defineEmits<{
  "update-xp": [value: number]
}>()

const xp = computed<number>({
  get: () => (character.value as any).exp ?? 0,
  set: (v) => {
    (character.value as any).exp = Math.max(0, Math.floor(Number(v) || 0))
    emit("update-xp", (character.value as any).exp)
  }
})
</script>

<template>
  <div class="xp__row">
    <label class="xp__label" for="xpInput">XP</label>
    <input
      id="xpInput"
      class="xp__input"
      type="number"
      min="0"
      step="1"
      v-model.number="xp"
    />
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.xp h2 {
  margin: 0 0 $s-3;
  font-family: $font-head;
  font-size: $fs-2;
}

.xp__row {
  display: flex;
  flex-direction: column;
}

.xp__label {
  color: $text-2;
  font-weight: 600;
}

.xp__input {
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .02);
  color: $text-0;
  font-family: $font-body;
  font-size: 1rem;

  &:focus {
    outline: none;
    border-color: $red-0;
  }
}

.xp__hint {
  margin: $s-3 0 0;
  color: $text-2;
  font-size: 0.85rem;
}
</style>
