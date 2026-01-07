<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: number
    max?: number
    readonly?: boolean
  }>(),
  {
    max: 5,
    readonly: false
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
}>()

const dots = computed(() => {
  return Array.from({ length: props.max }, (_, i) => i + 1)
})

function setRating(value: number) {
  if (!props.readonly) {
    emit('update:modelValue', value)
  }
}
</script>

<template>
  <div class="dot-rating" :class="{ 'dot-rating--readonly': readonly }">
    <button
      v-for="dot in dots"
      :key="dot"
      type="button"
      class="dot"
      :class="{ 'dot--filled': dot <= modelValue }"
      @click="setRating(dot)"
      :disabled="readonly"
    />
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.dot-rating {
  display: flex;
  gap: $s-2;

  &--readonly {
    pointer-events: none;
  }
}

.dot {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, .3);
  background: transparent;
  cursor: pointer;
  transition: all $t-fast $ease;

  &:hover:not(:disabled) {
    border-color: rgba(255, 255, 255, .5);
    transform: scale(1.1);
  }

  &--filled {
    background: $grad-accent;
    border-color: $red-0;
  }

  &:disabled {
    cursor: default;
  }
}
</style>
