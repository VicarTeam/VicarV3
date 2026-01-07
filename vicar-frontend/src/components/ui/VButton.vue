<script setup lang="ts">
import {computed} from "vue";

type Variant = "primary" | "secondary" | "ghost" | "danger"
type Size = "sm" | "md" | "lg"
type AsTag = "button" | "a" | "div" | "span"

const props = withDefaults(
  defineProps<{
    as?: AsTag
    variant?: Variant
    size?: Size
    block?: boolean
    loading?: boolean
    disabled?: boolean
    nativeType?: "button" | "submit" | "reset"
  }>(),
  {
    as: "button",
    variant: "secondary",
    size: "md",
    block: false,
    loading: false,
    disabled: false,
    nativeType: "button"
  }
)

const emit = defineEmits<{
  (e: "click", ev: MouseEvent): void
}>()

const isButtonTag = computed(() => props.as === "button")

function onClick(ev: MouseEvent) {
  if (props.disabled || props.loading) {
    ev.preventDefault()
    ev.stopPropagation()
    return
  }
  emit("click", ev)
}
</script>

<template>
  <component
    :is="as"
    class="btn"
    :class="[
      `btn--${variant}`,
      `btn--${size}`,
      block && 'btn--block',
      loading && 'btn--loading'
    ]"
    :type="isButtonTag ? nativeType : undefined"
    :disabled="isButtonTag ? disabled || loading : undefined"
    :aria-disabled="!isButtonTag ? String(disabled || loading) : undefined"
    :tabindex="!isButtonTag && (disabled || loading) ? -1 : undefined"
    @click="onClick"
  >
    <span v-if="loading" class="btn__spinner" aria-hidden="true" />
    <span class="btn__content">
      <slot />
    </span>
  </component>
</template>

<style scoped lang="scss">
.btn {
  position: relative;
  user-select: none;

  &__spinner {
    width: 1em;
    height: 1em;
    border-radius: 999px;
    border: 2px solid rgba(255,255,255,.25);
    border-top-color: rgba(255,255,255,.75);
    animation: spin 900ms linear infinite;
  }

  &__content--hidden {
    opacity: 0;
  }

  &--block {
    width: 100%;
  }

  &--sm {
    padding: 0.55rem 0.85rem;
    font-size: 0.95rem;
  }

  &--md {
    padding: 0.75rem 1rem;
    font-size: 1rem;
  }

  &--lg {
    padding: 0.9rem 1.15rem;
    font-size: 1.05rem;
  }

  &--primary {
    background: var(--grad-accent);
    border-color: rgba(255,59,84,.45);
    box-shadow: 0 12px 40px rgba(209,21,44,.2);
  }

  &--secondary {
    background: rgba(255,255,255,.04);
    border-color: rgba(255,255,255,.12);
  }

  &--ghost {
    background: transparent;
    border-color: rgba(255,255,255,.10);
  }

  &--danger {
    background: linear-gradient(135deg, rgba(255,59,84,.85), rgba(140,12,27,.75));
    border-color: rgba(255,59,84,.38);
  }

  &--loading {
    pointer-events: none;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
