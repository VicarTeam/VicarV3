<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: number
  readonly?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const boxes = computed(() => {
  const result: boolean[] = []
  for (let i = 1; i <= 5; i++) {
    result.push(i <= props.modelValue)
  }
  return result
})

function clickBox(index: number) {
  if (props.readonly) return

  const clickedValue = index + 1
  if (clickedValue === props.modelValue) {
    emit('update:modelValue', clickedValue - 1)
  } else {
    emit('update:modelValue', clickedValue)
  }
}
</script>

<template>
  <div class="hunger-tracker">
    <button
      v-for="(filled, idx) in boxes"
      :key="idx"
      type="button"
      class="hunger-tracker__box"
      :class="{
        'hunger-tracker__box--filled': filled,
        'hunger-tracker__box--readonly': readonly
      }"
      :disabled="readonly"
      @click="clickBox(idx)"
    />
    <span class="hunger-tracker__value">{{ modelValue }}</span>
  </div>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.hunger-tracker {
  display: flex;
  gap: $s-2;
  align-items: center;

  &__box {
    width: 28px;
    height: 28px;
    border: 2px solid rgba($red-0, 0.4);
    border-radius: 50%;
    background: transparent;
    cursor: pointer;
    transition: all $t-fast $ease;

    &:hover:not(:disabled) {
      border-color: $red-0;
      transform: scale(1.1);
    }

    &--filled {
      background: $red-0;
      border-color: $red-0;
      box-shadow: 0 0 8px rgba($red-0, 0.4);
    }

    &--readonly {
      cursor: default;

      &:hover {
        transform: none;
      }
    }
  }

  &__value {
    font-size: 1.2rem;
    font-weight: 700;
    font-family: $font-head;
    color: $red-0;
    margin-left: $s-2;
  }
}
</style>
