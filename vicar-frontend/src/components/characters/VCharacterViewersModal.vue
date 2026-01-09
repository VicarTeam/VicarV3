<script setup lang="ts">
import { computed, ref, watch } from "vue"
import VModal from "@/components/ui/VModal.vue"
import type { V5Character } from "@/@types/v5"
import type { User } from "@/@types/user"
import { addCharacterViewer, removeCharacterViewer } from "@/rest/api/characters"
import { autocompleteUsers } from "@/rest/api/users"

const open = defineModel<boolean>({ required: true })

const props = defineProps<{
  character: Partial<V5Character>
  readonly?: boolean
}>()

const adding = ref(false)
const removing = ref<Set<string>>(new Set())
const errMsg = ref<string | null>(null)
const okMsg = ref<string | null>(null)

const query = ref("")
const searching = ref(false)
const suggestions = ref<User[]>([])
let searchTimer: ReturnType<typeof setTimeout> | null = null
let lastReq = 0

const viewers = computed<User[]>(() => ((props.character as any).viewers ?? []) as User[])

function label(u: User) {
  return u.username || u.id
}

function isViewer(userId: string) {
  return viewers.value.some((v) => v.id === userId)
}

watch(
  () => open.value,
  (v) => {
    if (!v) {
      query.value = ""
      suggestions.value = []
      errMsg.value = null
      okMsg.value = null
    }
  }
)

watch(
  () => query.value,
  () => {
    errMsg.value = null
    okMsg.value = null

    if (searchTimer) clearTimeout(searchTimer)

    const q = query.value.trim()
    if (q.length < 3) {
      suggestions.value = []
      searching.value = false
      return
    }

    const reqId = ++lastReq
    searching.value = true

    searchTimer = setTimeout(async () => {
      try {
        const res = await autocompleteUsers(q)
        if (reqId !== lastReq) return
        suggestions.value = res.filter((u) => !isViewer(u.id))
      } catch (e: any) {
        if (reqId !== lastReq) return
        errMsg.value = e?.message ?? "Suche fehlgeschlagen."
        suggestions.value = []
      } finally {
        if (reqId === lastReq) searching.value = false
      }
    }, 250)
  }
)

async function addByUser(u: User) {
  errMsg.value = null
  okMsg.value = null

  const charId = (props.character as any).id as string | undefined
  if (props.readonly) return
  if (!charId) return (errMsg.value = "Charakter-ID fehlt.")
  if (isViewer(u.id)) return (errMsg.value = "User ist bereits Viewer.")

  adding.value = true
  try {
    await addCharacterViewer(charId, u.id)
    ;(props.character as any).viewers = [...viewers.value, u]
    okMsg.value = "Viewer hinzugefügt."
    query.value = ""
    suggestions.value = []
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
    ;(props.character as any).viewers = viewers.value.filter((v) => v.id !== userId)
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

        <div v-if="viewers.length === 0" class="empty">Keine Viewer eingetragen.</div>

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

        <input
          v-model="query"
          class="input"
          placeholder="Mindestens 3 Zeichen (Username)…"
          autocomplete="off"
        />

        <div v-if="query.trim().length >= 3" class="results">
          <div v-if="searching" class="resultHint">Suche…</div>
          <div v-else-if="suggestions.length === 0" class="resultHint">Keine Treffer.</div>

          <button
            v-for="u in suggestions"
            :key="u.id"
            type="button"
            class="resultRow"
            :disabled="adding"
            @click="addByUser(u)"
          >
            <div class="resultMain">
              <span class="resultName">{{ u.username }}</span>
              <span class="resultId">{{ u.id }}</span>
            </div>
            <span class="resultAction">Hinzufügen</span>
          </button>
        </div>

        <p class="hint">Tippe einen Usernamen, dann auswählen.</p>
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

.input {
  padding: $s-2 $s-3;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255,255,255,.02);
  color: $text-0;

  &:focus { outline: none; border-color: $red-0; }
}

.results {
  border: 1px solid $border;
  border-radius: $r-md;
  background: rgba(255,255,255,.02);
  overflow: hidden;
}

.resultHint {
  padding: $s-3;
  color: $text-2;
  font-size: .9rem;
}

.resultRow {
  width: 100%;
  display: flex;
  justify-content: space-between;
  gap: $s-3;
  align-items: center;
  padding: $s-3;
  border: none;
  background: transparent;
  color: $text-0;
  cursor: pointer;
  text-align: left;

  &:hover { background: rgba(255,255,255,.03); }
  &:disabled { opacity: .6; cursor: not-allowed; }
}

.resultMain {
  display: grid;
  gap: 2px;
  min-width: 0;
}

.resultName {
  font-weight: 800;
  color: $text-0;
  font-family: $font-head;
}

.resultId {
  font-size: .78rem;
  color: $text-2;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.resultAction {
  font-size: .78rem;
  font-weight: 800;
  padding: 4px 10px;
  border-radius: 999px;
  border: 1px solid rgba($red-0,.25);
  background: rgba($red-0,.10);
  color: rgba($red-0,.95);
  flex-shrink: 0;
}

.hint { margin: 0; color: $text-2; font-size: .85rem; }

@media (max-width: 520px) {
  .resultRow { align-items: flex-start; flex-direction: column; }
  .resultAction { align-self: flex-start; }
}
</style>
