<script setup lang="ts">
import {computed, provide, reactive, ref} from "vue"
import VModal from "@/components/ui/VModal.vue"
import VCard from "@/components/ui/VCard.vue"
import VButton from "@/components/ui/VButton.vue"
import {ConfirmKey} from "@/composables/useConfirm.ts";

export type ConfirmVariant = "primary" | "secondary" | "ghost" | "danger"

export type ConfirmRequest = {
  title?: string
  message: string
  yesText?: string
  noText?: string
  yesVariante?: ConfirmVariant
  noVariante?: ConfirmVariant
}

type ConfirmInternal = ConfirmRequest & {
  id: number
  resolve: (v: boolean) => void
}

export type ConfirmApi = {
  confirm: (req: ConfirmRequest) => Promise<boolean>
}

const open = ref(false)
const queue = reactive<ConfirmInternal[]>([])
const current = computed(() => queue[0])

let seq = 0

function confirm(req: ConfirmRequest): Promise<boolean> {
  return new Promise<boolean>((resolve) => {
    queue.push({ ...req, id: ++seq, resolve })
    open.value = true
  })
}

function finish(v: boolean) {
  const item = queue.shift()
  if (item) item.resolve(v)
  open.value = queue.length > 0
}

provide(ConfirmKey, { confirm })
</script>

<template>
  <slot />

  <VModal v-model="open" :title="undefined" size="sm" :closeOnBackdrop="false" :closeOnEsc="false">
    <VCard class="confirm" elevated>
      <header class="confirm__header">
        <h2 class="confirm__title">{{ current?.title ?? "Bestätigen" }}</h2>
      </header>

      <div class="confirm__body">
        <p class="confirm__message">{{ current?.message ?? "" }}</p>
      </div>

      <div class="confirm__actions">
        <VButton
          :variant="(current?.noVariante ?? 'ghost')"
          size="md"
          block
          @click="finish(false)"
        >
          {{ current?.noText ?? "Nein" }}
        </VButton>

        <VButton
          :variant="(current?.yesVariante ?? 'primary')"
          size="md"
          block
          @click="finish(true)"
        >
          {{ current?.yesText ?? "Ja" }}
        </VButton>
      </div>
    </VCard>
  </VModal>
</template>

<style scoped lang="scss">
.confirm {
  padding: 1rem;
}

.confirm__header {
  margin-bottom: 0.6rem;
}

.confirm__title {
  margin: 0;
  font-size: 1.25rem;
}

.confirm__body {
  margin-bottom: 1rem;
}

.confirm__message {
  margin: 0;
  color: rgba(255, 255, 255, 0.72);
}

.confirm__actions {
  display: grid;
  gap: 0.6rem;
}

@media (min-width: 520px) {
  .confirm__actions {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
