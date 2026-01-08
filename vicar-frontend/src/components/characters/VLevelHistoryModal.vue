<script setup lang="ts">
import { computed } from "vue"
import VModal from "@/components/ui/VModal.vue"
import VCard from "@/components/ui/VCard.vue"
import type {V5Character, V5LevelChange} from "@/@types/v5"

const open = defineModel<boolean>({ required: true })
const character = defineProps<{ character: Partial<V5Character> }>()

const history = computed<V5LevelChange[]>(() => ((character.character as any).levelHistory ?? []) as V5LevelChange[])

const sorted = computed(() => {
  return [...history.value].sort((a, b) => {
    const da = Date.parse(a.date || "")
    const db = Date.parse(b.date || "")
    return (isNaN(db) ? 0 : db) - (isNaN(da) ? 0 : da)
  })
})

function fmtDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString()
}

function typeLabel(t: V5LevelChange["type"]): string {
  switch (t) {
    case "attribute": return "Attribut"
    case "skill": return "Fähigkeit"
    case "discipline": return "Disziplin"
    case "trait": return "Vorzug/Hintergrund"
    case "flaw": return "Schwäche"
    case "blood_potency": return "Blutpotenz"
    case "specialization": return "Spezialisierung"
    default: return "Unbekannt"
  }
}
</script>

<template>
  <VModal v-model="open" title="Level Verlauf" size="lg">
    <div class="history">
      <div v-if="sorted.length === 0" class="empty">Noch keine Einträge.</div>

      <div v-else class="list">
        <VCard v-for="h in sorted" :key="h.id" class="item">
          <div class="top">
            <span class="type">{{ typeLabel(h.type) }}</span>
            <span class="date">{{ fmtDate(h.date) }}</span>
          </div>

          <div class="text">{{ h.text }}</div>

          <div class="xp">
            <span class="tag">XP vorher: {{ h.expBefore }}</span>
            <span class="tag tag--used">−{{ h.expUsed }}</span>
            <span class="tag tag--after">XP nachher: {{ h.expAfter }}</span>
          </div>
        </VCard>
      </div>
    </div>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.history { display: grid; gap: $s-3; }

.empty {
  padding: $s-4;
  border-radius: $r-md;
  background: rgba(255, 255, 255, .02);
  color: $text-2;
  text-align: center;
}

.list { display: grid; gap: $s-3; }

.item {
  padding: $s-3;
  background: rgba(255, 255, 255, .02);
}

.top {
  display: flex;
  justify-content: space-between;
  gap: $s-3;
  margin-bottom: $s-2;
}

.type { font-weight: 700; color: $text-1; }
.date { font-size: .85rem; color: $text-2; }

.text { color: $text-1; margin-bottom: $s-2; line-height: 1.4; }

.xp { display: flex; flex-wrap: wrap; gap: $s-2; }

.tag {
  font-size: .8rem;
  padding: 2px $s-2;
  border-radius: 999px;
  border: 1px solid $border;
  background: rgba(255, 255, 255, .03);
  color: $text-2;
}

.tag--used {
  border-color: rgba($red-0, .25);
  background: rgba($red-0, .10);
  color: $red-0;
}

.tag--after {
  border-color: rgba(100, 200, 100, .25);
  background: rgba(100, 200, 100, .08);
  color: rgba(100, 200, 100, .95);
}
</style>
