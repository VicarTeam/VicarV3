<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, watch } from "vue"

type ModalSize = "sm" | "md" | "lg"

const open = defineModel<boolean>({ default: false })

const props = withDefaults(
  defineProps<{
    title?: string
    size?: ModalSize
    closeOnEsc?: boolean
    closeOnBackdrop?: boolean
    preventScroll?: boolean
  }>(),
  {
    title: undefined,
    size: "md",
    closeOnEsc: true,
    closeOnBackdrop: true,
    preventScroll: true
  }
)

const emit = defineEmits<{
  (e: "close"): void
}>()

const sizeClass = computed(() => `modal__panel--${props.size}`)

function schliessen() {
  open.value = false
  emit("close")
}

function onBackdropClick() {
  if (!props.closeOnBackdrop) return
  schliessen()
}

function onKeydown(e: KeyboardEvent) {
  if (!props.closeOnEsc) return
  if (e.key === "Escape" && open.value) schliessen()
}

function setBodyScrollLock(locked: boolean) {
  if (!props.preventScroll) return
  document.documentElement.style.overflow = locked ? "hidden" : ""
}

watch(
  () => open.value,
  (v) => setBodyScrollLock(v),
  { immediate: true }
)

onMounted(() => window.addEventListener("keydown", onKeydown))
onBeforeUnmount(() => {
  window.removeEventListener("keydown", onKeydown)
  setBodyScrollLock(false)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="open" class="modal" role="dialog" aria-modal="true">
        <button class="modal__backdrop" type="button" @click="onBackdropClick" aria-label="Schließen" />
        <div class="modal__wrap">
          <div class="modal__panel" :class="sizeClass">
            <div class="modal__header" v-if="title || $slots.header">
              <slot name="header">
                <h2 class="modal__title">{{ title }}</h2>
              </slot>

              <button class="modal__x" type="button" aria-label="Schließen" @click="schliessen">
                <span aria-hidden="true">×</span>
              </button>
            </div>

            <div class="modal__body">
              <slot :close="schliessen" />
            </div>

            <div v-if="$slots.footer" class="modal__footer">
              <slot name="footer" :close="schliessen" />
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.modal {
  position: fixed;
  inset: 0;
  z-index: 60;
}

.modal__backdrop {
  position: absolute;
  inset: 0;
  border: 0;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(10px);
  cursor: pointer;
}

.modal__wrap {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  padding: 1rem;
}

.modal__panel {
  width: min(720px, 100%);
  border-radius: var(--r-lg, 20px);
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(10, 7, 14, 0.72);
  backdrop-filter: blur(14px);
  box-shadow: var(--shadow-2, 0 18px 60px rgba(0, 0, 0, 0.55));
  overflow: hidden;
}

.modal__panel--sm {
  width: min(420px, 100%);
}

.modal__panel--md {
  width: min(720px, 100%);
}

.modal__panel--lg {
  width: min(980px, 100%);
}

.modal__header {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1rem 0.85rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.modal__title {
  margin: 0;
  font-size: clamp(1.1rem, 1vw + 1rem, 1.5rem);
}

.modal__x {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
  cursor: pointer;
  display: grid;
  place-items: center;
  transition: transform 120ms cubic-bezier(.2,.8,.2,1), border-color 220ms cubic-bezier(.2,.8,.2,1);

  &:hover {
    transform: translateY(-1px);
    border-color: rgba(255, 255, 255, 0.18);
  }

  span {
    font-size: 1.35rem;
    line-height: 1;
    transform: translateY(-1px);
  }
}

.modal__body {
  padding: 1rem;
}

.modal__footer {
  padding: 0.9rem 1rem 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 180ms cubic-bezier(.2,.8,.2,1);
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

@media (max-width: 420px) {
  .modal__wrap {
    padding: 0.75rem;
    place-items: end center;
  }

  .modal__panel {
    border-bottom-left-radius: 0;
    border-bottom-right-radius: 0;
  }
}
</style>
