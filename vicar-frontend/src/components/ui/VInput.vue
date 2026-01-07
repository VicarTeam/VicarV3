<script setup lang="ts">
import { computed, useAttrs } from "vue"

type Size = "sm" | "md" | "lg"

const model = defineModel<string | number | null>({
  default: null
})

const props = withDefaults(
  defineProps<{
    label?: string
    help?: string
    error?: string
    size?: Size
  }>(),
  {
    label: undefined,
    help: undefined,
    error: undefined,
    size: "md"
  }
)

const emit = defineEmits<{
  (e: "blur", ev: FocusEvent): void
  (e: "focus", ev: FocusEvent): void
}>()

const attrs = useAttrs()

const invalid = computed(() => Boolean(props.error))

const sizeClass = computed(() => {
  if (props.size === "sm") return "input--sm"
  if (props.size === "lg") return "input--lg"
  return "input--md"
})

const inputAttrs = computed(() => {
  const { class: _class, style: _style, ...rest } = attrs
  return rest
})

function onInput(ev: Event) {
  const el = ev.target as HTMLInputElement
  const type = (attrs.type as string | undefined) ?? "text"
  if (type === "number") {
    model.value = el.value === "" ? null : Number(el.value)
  } else {
    model.value = el.value
  }
}
</script>

<template>
  <label v-if="label" class="field">
    <span class="field__label">{{ label }}</span>

    <input
      v-bind="inputAttrs"
      class="input"
      :class="[
        sizeClass,
        invalid && 'input--invalid'
      ]"
      :value="model"
      @input="onInput"
      @blur="emit('blur', $event)"
      @focus="emit('focus', $event)"
    />

    <span v-if="error" class="field__error">{{ error }}</span>
    <span v-else-if="help" class="field__help">{{ help }}</span>
  </label>

  <div v-else class="field">
    <input
      v-bind="inputAttrs"
      class="input"
      :class="[
        sizeClass,
        invalid && 'input--invalid'
      ]"
      :value="model"
      @input="onInput"
      @blur="emit('blur', $event)"
      @focus="emit('focus', $event)"
    />
    <span v-if="error" class="field__error">{{ error }}</span>
    <span v-else-if="help" class="field__help">{{ help }}</span>
  </div>
</template>

<style scoped lang="scss">
.field {
  display: grid;
  gap: 0.4rem;

  &__label {
    font-size: 0.95rem;
    color: rgba(255,255,255,.72);
  }

  &__help {
    font-size: 0.9rem;
    color: rgba(255,255,255,.54);
  }

  &__error {
    font-size: 0.9rem;
    color: rgba(255,59,84,.92);
  }
}

.input {
  &--sm { padding: 0.55rem 0.75rem; }
  &--md { padding: 0.7rem 0.85rem; }
  &--lg { padding: 0.85rem 1rem; }

  &--invalid {
    border-color: rgba(255,59,84,.45);
  }
}
</style>
