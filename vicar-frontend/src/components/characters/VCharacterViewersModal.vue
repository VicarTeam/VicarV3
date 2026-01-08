<script setup lang="ts">
import { computed, ref } from "vue"
import VModal from "@/components/ui/VModal.vue"
import VButton from "@/components/ui/VButton.vue"
import type { V5Character } from "@/@types/v5"
import type {User} from "@/@types/user.ts";
import {addCharacterViewer, removeCharacterViewer} from "@/rest/api/characters.ts";

const open = defineModel<boolean>({ required: true })

const props = defineProps<{
  character: Partial<V5Character>
  readonly?: boolean
}>()

const addUserId = ref("")
const adding = ref(false)
const removing = ref<Set<string>>(new Set())
const errMsg = ref<string | null>(null)
const okMsg = ref<string | null>(null)

const viewers = computed<User[]>(() => ((props.character as any).viewers ?? []) as User[])

function label(u: User) {
  return u.username || u.id
}

async function add() {
  errMsg.value = null
  okMsg.value = null

  const charId = (props.character as any).id as string | undefined
  const userId = addUserId.value.trim()

  if (props.readonly) return
  if (!charId) return (errMsg.value = "Charakter-ID fehlt.")
  if (!userId) return (errMsg.value = "Bitte User-ID eingeben.")
  if (viewers.value.some(v => v.id === userId)) return (errMsg.value = "User ist bereits Viewer.")

  adding.value = true
  try {
    await addCharacterViewer(charId, userId)
    ;(props.character as any).viewers = [...viewers.value, { id: userId }]
    okMsg.value = "Viewer hinzugefügt."
    addUserId.value = ""
  } catch (e: any) {
    errMsg.value = e?.message ?? "Viewer konnte nicht hinzugefügt werden."
  } finally {
    adding.value = false
  }
}

async function remove(userId: string) {
  errMsg.value = null
  okMsg.value = null

  const charId = (props.character as any).id as string | undefined
  if (props.readonly) return
  if (!charId) return (errMsg.value = "Charakter-ID fehlt.")

  removing.value = new Set(removing.value).add(userId)
  try {
    await removeCharacterViewer(charId, userId)
    ;(props.character as any).viewers = viewers.value.filter(v => v.id !== userId)
    okMsg.value = "Viewer entfernt."
  } catch (e: any) {
    errMsg.value = e?.message ?? "Viewer konnte nicht entfernt werden."
  } finally {
    const s = new Set(removing.value)
    s.delete(userId)
    removing.value = s
  }
}
</script>

<template>
  <VModal v-model="open" title="Zugriff verwalten" size="md">
    <div class="wrap">
      <div v-if="errMsg" class="msg msg--err">{{ errMsg }}</div>
      <div v-else-if="okMsg" class="msg msg--ok">{{ okMsg }}</div>

      <div class="section">
        <h4>Viewer</h4>

        <div v-if="viewers.length === 0" class="empty">
          Keine Viewer eingetragen.
        </div>

        <div v-else class="chips">
          <div v-for="v in viewers" :key="v.id" class="chip">
            <span class="chip__name">{{ label(v) }}</span>
            <button
              v-if="!readonly"
              type="button"
              class="chip__remove"
              :disabled="removing.has(v.id)"
              @click="remove(v.id)"
              title="Entfernen"
            >
              ✕
            </button>
          </div>
        </div>
      </div>

      <div v-if="!readonly" class="section section--add">
        <h4>Hinzufügen</h4>
        <div class="row">
          <input v-model="addUserId" class="input" placeholder="User-ID (UUID) …" />
          <VButton variant="primary" :disabled="adding" @click="add">Hinzufügen</VButton>
        </div>
        <p class="hint">Später können wir das durch User-Suche/Autocomplete ersetzen.</p>
      </div>
    </div>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.wrap { display: grid; gap: $s-3; }

.msg {
  padding: $s-2 $s-3;
  border-radius: $r-md;
  border: 1px solid $border;
  font-size: 0.9rem;

  &--ok {
    border-color: rgba(100, 200, 100, .35);
    background: rgba(100, 200, 100, .08);
    color: rgba(100, 200, 100, .95);
  }
  &--err {
    border-color: rgba($red-0, .35);
    background: rgba($red-0, .08);
    color: $red-0;
  }
}

.section { display: grid; gap: $s-2; }

h4 {
  margin: 0;
  font-family: $font-head;
  color: $text-1;
}

.empty {
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255,255,255,.02);
  color: $text-2;
  text-align: center;
}

.chips { display: flex; flex-wrap: wrap; gap: $s-2; }

.chip {
  display: inline-flex;
  align-items: center;
  gap: $s-2;
  padding: $s-1 $s-2;
  border-radius: 999px;
  border: 1px solid $border;
  background: rgba(255,255,255,.03);
}

.chip__name { font-weight: 600; font-size: .85rem; color: $text-1; }

.chip__remove {
  width: 22px;
  height: 22px;
  border-radius: 999px;
  border: 1px solid rgba($red-0,.25);
  background: rgba($red-0,.10);
  color: $red-0;
  cursor: pointer;
  display: grid;
  place-items: center;

  &:disabled { opacity: .5; cursor: not-allowed; }
}

.section--add { padding-top: $s-3; border-top: 1px solid $border; }

.row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: $s-2;
  align-items: center;
}

.input {
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255,255,255,.02);
  color: $text-0;

  &:focus { outline: none; border-color: $red-0; }
}

.hint { margin: 0; color: $text-2; font-size: .85rem; }

@media (max-width: 520px) {
  .row { grid-template-columns: 1fr; }
}
</style>
