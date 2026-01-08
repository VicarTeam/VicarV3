<script setup lang="ts">
type AsTag = keyof HTMLElementTagNameMap

withDefaults(
  defineProps<{
    as?: AsTag
    tight?: boolean
    elevated?: boolean
    interactive?: boolean
  }>(),
  {
    as: "div",
    tight: false,
    elevated: false,
    interactive: false
  }
)
</script>

<template>
  <component
    :is="as"
    class="card"
    :class="[
      tight && 'card--tight',
      elevated && 'card--elevated',
      interactive && 'card--interactive'
    ]"
  >
    <slot />
  </component>
</template>

<style scoped lang="scss">
.card--elevated {
  box-shadow: var(--shadow-2, 0 18px 60px rgba(0,0,0,.55));
}

.card--interactive {
  transition: transform var(--t-med, 220ms) var(--ease, cubic-bezier(.2,.8,.2,1)),
  border-color var(--t-med, 220ms) var(--ease, cubic-bezier(.2,.8,.2,1));
  cursor: pointer;

  &:hover {
    transform: translateY(-2px);
    border-color: rgba(255,255,255,.18);
  }

  &:active {
    transform: translateY(0);
  }
}
</style>