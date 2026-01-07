<script setup lang="ts">
import {computed, provide, reactive} from "vue"
import {ToastKey} from "@/composables/useToast.ts";

export type ToastVariant = "success" | "error" | "warning" | "info"

export type ToastOptions = {
  title?: string
  message: string
  variant?: ToastVariant
  durationMs?: number
  dismissible?: boolean
}

export type ToastItem = Required<Pick<ToastOptions, "message">> &
  Omit<ToastOptions, "message"> & {
  id: number
  createdAt: number
}

export type ToastApi = {
  push: (options: ToastOptions) => number
  dismiss: (id: number) => void
  clear: () => void
  success: (message: string, options?: Omit<ToastOptions, "message" | "variant">) => number
  error: (message: string, options?: Omit<ToastOptions, "message" | "variant">) => number
  warning: (message: string, options?: Omit<ToastOptions, "message" | "variant">) => number
  info: (message: string, options?: Omit<ToastOptions, "message" | "variant">) => number
}

const props = withDefaults(
  defineProps<{
    maxToasts?: number
    defaultDurationMs?: number
  }>(),
  {
    maxToasts: 4,
    defaultDurationMs: 4200
  }
)

const toasts = reactive<ToastItem[]>([])
let seq = 0

function normalize(options: ToastOptions): ToastItem {
  return {
    id: ++seq,
    createdAt: Date.now(),
    title: options.title,
    message: options.message,
    variant: options.variant ?? "info",
    durationMs: options.durationMs ?? props.defaultDurationMs,
    dismissible: options.dismissible ?? true
  }
}

function dismiss(id: number) {
  const idx = toasts.findIndex((t) => t.id === id)
  if (idx >= 0) toasts.splice(idx, 1)
}

function clear() {
  toasts.splice(0, toasts.length)
}

function push(options: ToastOptions): number {
  const item = normalize(options)

  if (toasts.length >= props.maxToasts) {
    toasts.shift()
  }

  toasts.push(item)

  if ((item.durationMs ?? 0) > 0) {
    window.setTimeout(() => dismiss(item.id), item.durationMs)
  }

  return item.id
}

function variantShortcut(variant: ToastVariant) {
  return (message: string, options?: Omit<ToastOptions, "message" | "variant">) =>
    push({ message, variant, ...(options ?? {}) })
}

const api: ToastApi = {
  push,
  dismiss,
  clear,
  success: variantShortcut("success"),
  error: variantShortcut("error"),
  warning: variantShortcut("warning"),
  info: variantShortcut("info")
}

provide(ToastKey, api)

const sortedToasts = computed(() => [...toasts].sort((a, b) => a.createdAt - b.createdAt))
</script>

<template>
  <slot />

  <Teleport to="body">
    <div class="toastHost" aria-live="polite" aria-relevant="additions removals">
      <TransitionGroup name="toast" tag="div" class="toastHost__stack">
        <div
          v-for="t in sortedToasts"
          :key="t.id"
          class="toast"
          :class="`toast--${t.variant}`"
          role="status"
        >
          <div class="toast__content">
            <div v-if="t.title" class="toast__title">{{ t.title }}</div>
            <div class="toast__message">{{ t.message }}</div>
          </div>

          <button
            v-if="t.dismissible"
            class="toast__close"
            type="button"
            aria-label="Schließen"
            @click="dismiss(t.id)"
          >
            <span aria-hidden="true">×</span>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped lang="scss">
.toastHost {
  position: fixed;
  z-index: 70;
  inset: 0;
  pointer-events: none;
  display: grid;
  align-items: start;
  justify-items: center;
  padding: 1rem;
}

.toastHost__stack {
  width: min(520px, 100%);
  display: grid;
  gap: 0.6rem;
}

.toast {
  pointer-events: auto;
  width: 100%;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 0.75rem;
  align-items: start;

  border-radius: var(--r-md, 14px);
  border: 1px solid rgba(255,255,255,.12);
  background: rgba(10, 7, 14, 0.72);
  backdrop-filter: blur(14px);
  box-shadow: var(--shadow-2, 0 18px 60px rgba(0,0,0,.55));
  padding: 0.85rem 0.9rem;
}

.toast__content {
  display: grid;
  gap: 0.2rem;
}

.toast__title {
  font-family: var(--font-head, inherit);
  letter-spacing: 0.4px;
  font-size: 1.02rem;
  color: rgba(255,255,255,.92);
}

.toast__message {
  color: rgba(255,255,255,.72);
  line-height: 1.35;
}

.toast__close {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 999px;
  border: 1px solid rgba(255,255,255,.12);
  background: rgba(255,255,255,.04);
  cursor: pointer;
  display: grid;
  place-items: center;
  transition: transform 120ms cubic-bezier(.2,.8,.2,1), border-color 220ms cubic-bezier(.2,.8,.2,1);

  &:hover {
    transform: translateY(-1px);
    border-color: rgba(255,255,255,.18);
  }

  span {
    font-size: 1.35rem;
    line-height: 1;
    transform: translateY(-1px);
    color: rgba(255,255,255,.85);
  }
}

.toast--success {
  border-color: rgba(80, 200, 120, 0.28);
  box-shadow: 0 18px 60px rgba(0,0,0,.55), 0 0 0 1px rgba(80,200,120,.10) inset;
}

.toast--error {
  border-color: rgba(255, 59, 84, 0.38);
  box-shadow: 0 18px 60px rgba(0,0,0,.55), 0 0 0 1px rgba(255,59,84,.12) inset;
}

.toast--warning {
  border-color: rgba(255, 200, 80, 0.30);
  box-shadow: 0 18px 60px rgba(0,0,0,.55), 0 0 0 1px rgba(255,200,80,.10) inset;
}

.toast--info {
  border-color: rgba(120, 170, 255, 0.28);
  box-shadow: 0 18px 60px rgba(0,0,0,.55), 0 0 0 1px rgba(120,170,255,.10) inset;
}

.toast-enter-active,
.toast-leave-active {
  transition: transform 180ms cubic-bezier(.2,.8,.2,1), opacity 180ms cubic-bezier(.2,.8,.2,1);
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@media (max-width: 420px) {
  .toastHost {
    align-items: end;
    padding: 0.8rem;
  }

  .toastHost__stack {
    width: 100%;
  }

  .toast {
    border-bottom-left-radius: 0;
    border-bottom-right-radius: 0;
  }
}
</style>
