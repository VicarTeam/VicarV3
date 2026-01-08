<script setup lang="ts">
import { computed, useAttrs } from "vue"

const model = defineModel<string>({
  default: ''
})

const props = withDefaults(
  defineProps<{
    label?: string
    help?: string
    error?: string
    rows?: number
  }>(),
  {
    label: undefined,
    help: undefined,
    error: undefined,
    rows: 4
  }
)

const emit = defineEmits<{
  (e: "blur", ev: FocusEvent): void
  (e: "focus", ev: FocusEvent): void
}>()

const attrs = useAttrs()

const invalid = computed(() => Boolean(props.error))

const inputAttrs = computed(() => {
  const { class: _class, style: _style, ...rest } = attrs
  return rest
})

function onInput(ev: Event) {
  const el = ev.target as HTMLTextAreaElement
  model.value = el.value
}
</script>

<template>
  <label v-if="label" class="field">
    <span class="field__label">{{ label }}</span>

    <textarea
      v-bind="inputAttrs"
      class="textarea"
      :class="[invalid && 'textarea--invalid']"
      :rows="rows"
      :value="model"
      @input="onInput"
      @blur="emit('blur', $event)"
      @focus="emit('focus', $event)"
    />

    <span v-if="error" class="field__error">{{ error }}</span>
    <span v-else-if="help" class="field__help">{{ help }}</span>
  </label>

  <div v-else class="field">
    <textarea
      v-bind="inputAttrs"
      class="textarea"
      :class="[invalid && 'textarea--invalid']"
      :rows="rows"
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
@use "@/styles/variables" as *;

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

.textarea {
  width: 100%;
  padding: 0.7rem 0.85rem;
  font-family: inherit;
  font-size: 1rem;
  color: $text-0;
  background: $bg-1;
  border: 1px solid $border;
  border-radius: $r-md;
  resize: vertical;
  min-height: 80px;
  transition: border-color $t-fast $ease;

  &:focus {
    outline: none;
    border-color: $red-1;
  }

  &::placeholder {
    color: $text-2;
  }

  &--invalid {
    border-color: rgba(255,59,84,.45);
  }
}
</style>
