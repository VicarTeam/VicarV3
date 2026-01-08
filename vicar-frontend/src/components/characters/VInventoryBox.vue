<script setup lang="ts">
import { computed, ref, watch } from "vue"
import VCard from "@/components/ui/VCard.vue"
import VButton from "@/components/ui/VButton.vue"
import VInput from "@/components/ui/VInput.vue"
import VTextarea from "@/components/ui/VTextarea.vue"
import VModal from "@/components/ui/VModal.vue"
import type {InventorySide, V5Character, V5Inventory, V5InventoryItem} from "@/@types/v5"

const props = defineProps<{
  modelValue: Partial<V5Character>
  readonly?: boolean
}>()

const emit = defineEmits<{
  "update:modelValue": [value: Partial<V5Character>]
  "update-inventory": [inv: V5Inventory]
}>()

const showItemModal = ref(false)
const editingSide = ref<InventorySide>("carriedItems")
const editingItemId = ref<string | null>(null)

const name = ref("")
const description = ref("")
const amount = ref<number>(1)

const showTransferModal = ref(false)
const transferFrom = ref<"cash" | "bank">("bank")
const transferAmount = ref<number>(1)

const showZeroModal = ref(false)
const zeroItemName = ref("")
const zeroConfirm = ref<null | (() => void)>(null)

const inventory = computed<V5Inventory>(() => {
  const inv = (props.modelValue as any).inventory as V5Inventory | undefined
  return {
    cash: inv?.cash ?? 0,
    bank: inv?.bank ?? 0,
    carriedItems: Array.isArray(inv?.carriedItems) ? inv!.carriedItems : [],
    ownedItems: Array.isArray(inv?.ownedItems) ? inv!.ownedItems : []
  }
})

watch(
  () => props.modelValue,
  () => {
    if (!(props.modelValue as any).inventory) {
      const next = { ...(props.modelValue as any) }
      next.inventory = { cash: 0, bank: 0, carriedItems: [], ownedItems: [] } satisfies V5Inventory
      emit("update:modelValue", next)
      emit("update-inventory", next.inventory)
    }
  },
  { immediate: true, deep: true }
)

function setInventory(nextInv: V5Inventory) {
  const next = { ...(props.modelValue as any) }
  next.inventory = nextInv
  emit("update:modelValue", next)
  emit("update-inventory", nextInv)
}

function toNumberSafe(v: any): number {
  const n = typeof v === "number" ? v : parseInt(String(v ?? ""), 10)
  return Number.isFinite(n) ? n : 0
}

function updateMoney(key: "cash" | "bank", raw: any) {
  if (props.readonly) return
  const nextInv = { ...inventory.value, [key]: Math.max(0, toNumberSafe(raw)) }
  setInventory(nextInv)
}

function beginAdd(side: InventorySide) {
  if (props.readonly) return
  editingSide.value = side
  editingItemId.value = null
  name.value = ""
  description.value = ""
  amount.value = 1
  showItemModal.value = true
}

function beginEdit(side: InventorySide, item: V5InventoryItem) {
  if (props.readonly) return
  editingSide.value = side
  editingItemId.value = item.id
  name.value = item.name
  description.value = item.description
  amount.value = Math.max(1, item.amount ?? 1)
  showItemModal.value = true
}

function saveItem() {
  if (props.readonly) return
  const n = name.value.trim()
  const d = description.value.trim()
  const a = Math.max(1, toNumberSafe(amount.value))

  if (!n || !d) return

  const nextInv: V5Inventory = {
    ...inventory.value,
    carriedItems: [...inventory.value.carriedItems],
    ownedItems: [...inventory.value.ownedItems]
  }

  const list = nextInv[editingSide.value]

  if (editingItemId.value) {
    const idx = list.findIndex((x) => x.id === editingItemId.value)
    if (idx !== -1) list[idx] = { ...list[idx], name: n, description: d, amount: a } as any
  } else {
    list.push({
      id: `inv-${Date.now()}-${Math.random().toString(16).slice(2)}`,
      name: n,
      description: d,
      amount: a,
      category: "Custom"
    })
  }

  sortList(list)
  setInventory(nextInv)
  showItemModal.value = false
}

function sortList(list: V5InventoryItem[]) {
  list.sort((a, b) => {
    const ca = (a.category ?? "").toLowerCase()
    const cb = (b.category ?? "").toLowerCase()
    if (ca < cb) return -1
    if (ca > cb) return 1
    const na = (a.name ?? "").toLowerCase()
    const nb = (b.name ?? "").toLowerCase()
    if (na < nb) return -1
    if (na > nb) return 1
    return 0
  })
}

function removeItem(side: InventorySide, itemId: string) {
  if (props.readonly) return
  const nextInv: V5Inventory = {
    ...inventory.value,
    carriedItems: [...inventory.value.carriedItems],
    ownedItems: [...inventory.value.ownedItems]
  }
  nextInv[side] = nextInv[side].filter((x) => x.id !== itemId)
  setInventory(nextInv)
}

function cloneItem(side: InventorySide, item: V5InventoryItem) {
  if (props.readonly) return
  const nextInv: V5Inventory = {
    ...inventory.value,
    carriedItems: [...inventory.value.carriedItems],
    ownedItems: [...inventory.value.ownedItems]
  }
  nextInv[side].push({
    ...item,
    id: `inv-${Date.now()}-${Math.random().toString(16).slice(2)}`,
    name: item.name
  })
  sortList(nextInv[side])
  setInventory(nextInv)
}

function moveItem(side: InventorySide, itemId: string) {
  if (props.readonly) return
  const other: InventorySide = side === "carriedItems" ? "ownedItems" : "carriedItems"

  const nextInv: V5Inventory = {
    ...inventory.value,
    carriedItems: [...inventory.value.carriedItems],
    ownedItems: [...inventory.value.ownedItems]
  }

  const idx = nextInv[side].findIndex((x) => x.id === itemId)
  if (idx === -1) return

  const [item] = nextInv[side].splice(idx, 1)
  nextInv[other].push(item!)

  sortList(nextInv[side])
  sortList(nextInv[other])
  setInventory(nextInv)
}

function setAmount(side: InventorySide, itemId: string, v: number) {
  if (props.readonly) return
  const nextInv: V5Inventory = {
    ...inventory.value,
    carriedItems: [...inventory.value.carriedItems],
    ownedItems: [...inventory.value.ownedItems]
  }
  const list = nextInv[side]
  const idx = list.findIndex((x) => x.id === itemId)
  if (idx === -1) return

  const nextAmount = toNumberSafe(v)

  if (nextAmount >= 1) {
    list[idx] = { ...list[idx], amount: nextAmount } as any
    setInventory(nextInv)
    return
  }

  zeroItemName.value = list[idx]!.name
  zeroConfirm.value = () => {
    removeItem(side, itemId)
    showZeroModal.value = false
    zeroConfirm.value = null
  }
  showZeroModal.value = true
}

function openTransfer(from: "bank" | "cash") {
  if (props.readonly) return
  transferFrom.value = from
  transferAmount.value = 1
  showTransferModal.value = true
}

const canTransfer = computed(() => {
  const amt = Math.max(1, toNumberSafe(transferAmount.value))
  return inventory.value[transferFrom.value] >= amt
})

function confirmTransfer() {
  if (props.readonly) return
  if (!canTransfer.value) return

  const amt = Math.max(1, toNumberSafe(transferAmount.value))
  const from = transferFrom.value
  const to = from === "bank" ? "cash" : "bank"

  const nextInv: V5Inventory = {
    ...inventory.value,
    [from]: inventory.value[from] - amt,
    [to]: inventory.value[to] + amt
  }
  setInventory(nextInv)
  showTransferModal.value = false
}

const predefinedGroups = computed(() => {
  return [] as any[]
})
</script>

<template>
  <VCard>
    <div class="head">
      <div>
        <h2>Inventar</h2>
        <p class="sub">Bargeld, Konto und Gegenstände.</p>
      </div>
    </div>

    <div class="money">
      <div class="moneyBox">
        <div class="moneyTop">
          <span class="moneyLabel">Bargeld</span>
          <div class="moneyActions">
            <VButton v-if="!readonly" variant="ghost" @click="openTransfer('cash')">→ Bank</VButton>
          </div>
        </div>
        <VInput
          label=""
          type="number"
          :min="0"
          :model-value="inventory.cash"
          @update:model-value="updateMoney('cash', $event)"
        />
      </div>

      <div class="moneyBox">
        <div class="moneyTop">
          <span class="moneyLabel">Bank</span>
          <div class="moneyActions">
            <VButton v-if="!readonly" variant="ghost" @click="openTransfer('bank')">→ Bargeld</VButton>
          </div>
        </div>
        <VInput
          label=""
          type="number"
          :min="0"
          :model-value="inventory.bank"
          @update:model-value="updateMoney('bank', $event)"
        />
      </div>
    </div>

    <div class="twoCols">
      <div class="col">
        <div class="colHead">
          <h3>Getragen</h3>
          <VButton v-if="!readonly" variant="secondary" @click="beginAdd('carriedItems')">+ Item</VButton>
        </div>

        <div v-if="inventory.carriedItems.length === 0" class="empty">
          Keine Items.
        </div>

        <div v-else class="list">
          <div v-for="it in inventory.carriedItems" :key="it.id" class="row">
            <div class="main" @click="beginEdit('carriedItems', it)">
              <div class="title">
                <span class="name">{{ it.name }}</span>
                <span class="pill">{{ it.amount }}×</span>
              </div>
              <div class="desc">{{ it.description }}</div>
            </div>

            <div class="actions">
              <input
                v-if="!readonly"
                class="amount"
                type="number"
                min="0"
                step="1"
                :value="it.amount"
                @input="setAmount('carriedItems', it.id, Number(($event.target as HTMLInputElement).value))"
              />
              <VButton v-if="!readonly" variant="ghost" @click.stop="cloneItem('carriedItems', it)">Klonen</VButton>
              <VButton v-if="!readonly" variant="ghost" @click.stop="moveItem('carriedItems', it.id)">→</VButton>
              <VButton v-if="!readonly" variant="ghost" @click.stop="removeItem('carriedItems', it.id)">✕</VButton>
            </div>
          </div>
        </div>
      </div>

      <div class="col">
        <div class="colHead">
          <h3>Besitz</h3>
          <VButton v-if="!readonly" variant="secondary" @click="beginAdd('ownedItems')">+ Item</VButton>
        </div>

        <div v-if="inventory.ownedItems.length === 0" class="empty">
          Keine Items.
        </div>

        <div v-else class="list">
          <div v-for="it in inventory.ownedItems" :key="it.id" class="row">
            <div class="main" @click="beginEdit('ownedItems', it)">
              <div class="title">
                <span class="name">{{ it.name }}</span>
                <span class="pill">{{ it.amount }}×</span>
              </div>
              <div class="desc">{{ it.description }}</div>
            </div>

            <div class="actions">
              <input
                v-if="!readonly"
                class="amount"
                type="number"
                min="0"
                step="1"
                :value="it.amount"
                @input="setAmount('ownedItems', it.id, Number(($event.target as HTMLInputElement).value))"
              />
              <VButton v-if="!readonly" variant="ghost" @click.stop="cloneItem('ownedItems', it)">Klonen</VButton>
              <VButton v-if="!readonly" variant="ghost" @click.stop="moveItem('ownedItems', it.id)">←</VButton>
              <VButton v-if="!readonly" variant="ghost" @click.stop="removeItem('ownedItems', it.id)">✕</VButton>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div style="display:none">
      {{ predefinedGroups }}
    </div>
  </VCard>

  <VModal v-model="showItemModal" title="Item" size="md">
    <div class="modalBody">
      <VInput label="Name" :model-value="name" @update:model-value="name = $event as any" placeholder="z.B. Pistole, Laptop…" />
      <VTextarea
        label="Beschreibung"
        :model-value="description"
        @update:model-value="description = $event as any"
        placeholder="Kurzbeschreibung…"
        :rows="4"
      />
      <VInput
        label="Anzahl"
        type="number"
        :min="1"
        :model-value="amount"
        @update:model-value="amount = Number($event)"
      />

      <div class="modalActions">
        <VButton variant="secondary" @click="showItemModal = false">Schließen</VButton>
        <VButton variant="primary" :disabled="!name.trim() || !description.trim()" @click="saveItem">
          {{ editingItemId ? "Speichern" : "Hinzufügen" }}
        </VButton>
      </div>
    </div>
  </VModal>

  <VModal v-model="showTransferModal" title="Überweisen" size="sm">
    <div class="modalBody">
      <p class="hint">
        {{ transferFrom === "bank" ? "Von Bank → Cash" : "Von Cash → Bank" }}
      </p>

      <VInput
        label="Betrag"
        type="number"
        :min="1"
        :model-value="transferAmount"
        @update:model-value="transferAmount = Number($event)"
      />

      <p v-if="!canTransfer" class="error">Nicht genug Guthaben.</p>

      <div class="modalActions">
        <VButton variant="secondary" @click="showTransferModal = false">Abbrechen</VButton>
        <VButton variant="primary" :disabled="!canTransfer" @click="confirmTransfer">Bestätigen</VButton>
      </div>
    </div>
  </VModal>

  <VModal v-model="showZeroModal" title="Anzahl ist 0" size="sm">
    <div class="modalBody">
      <p class="hint">„{{ zeroItemName }}“ hat jetzt 0. Löschen?</p>

      <div class="modalActions">
        <VButton variant="secondary" @click="showZeroModal = false">Behalten</VButton>
        <VButton variant="primary" @click="zeroConfirm && zeroConfirm()">Löschen</VButton>
      </div>
    </div>
  </VModal>
</template>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.head {
  display: flex;
  justify-content: space-between;
  gap: $s-4;
  align-items: flex-start;

  h2 { margin: 0 0 $s-2; font-family: $font-head; font-size: $fs-2; }
}

.sub { margin: 0; color: $text-2; font-size: .9rem; }

.money {
  margin-top: $s-4;
  display: grid;
  gap: $s-3;
  grid-template-columns: repeat(2, 1fr);
}

.moneyBox {
  border: 1px solid $border;
  background: rgba(255,255,255,.02);
  border-radius: $r-md;
  padding: $s-3;
  display: grid;
  gap: $s-2;
}

.moneyTop {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-2;
}

.moneyLabel {
  font-family: $font-head;
  font-weight: 700;
  color: $text-1;
}

.moneyActions { display: flex; gap: $s-2; }

.twoCols {
  margin-top: $s-4;
  display: grid;
  gap: $s-4;
  grid-template-columns: repeat(2, 1fr);
}

.col {
  border: 1px solid $border;
  background: rgba(255,255,255,.02);
  border-radius: $r-md;
  padding: $s-3;
  display: grid;
  gap: $s-3;
}

.colHead {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: $s-3;

  h3 {
    margin: 0;
    font-family: $font-head;
    font-size: $fs-1;
    color: $text-1;
  }
}

.empty {
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255,255,255,.02);
  color: $text-2;
  text-align: center;
}

.list { display: grid; gap: $s-2; }

.row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: $s-3;
  padding: $s-3;
  border-radius: $r-md;
  background: rgba(255,255,255,.02);
  border: 1px solid rgba(255,255,255,.08);
}

.main {
  cursor: pointer;
  display: grid;
  gap: $s-1;
  min-width: 0;
}

.title {
  display: flex;
  align-items: center;
  gap: $s-2;
  min-width: 0;
}

.name {
  font-weight: 700;
  color: $text-0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pill {
  font-size: .75rem;
  padding: 2px $s-2;
  border-radius: 999px;
  border: 1px solid $border;
  background: rgba(255,255,255,.03);
  color: $text-2;
  flex-shrink: 0;
}

.desc {
  color: $text-2;
  font-size: .85rem;
  line-height: 1.35;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.actions {
  display: flex;
  align-items: center;
  gap: $s-2;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.amount {
  width: 78px;
  text-align: center;
  padding: .45rem .5rem;
  border-radius: $r-sm;
  border: 1px solid $border;
  background: rgba(255,255,255,.02);
  color: $text-0;

  &:focus { outline: none; border-color: $red-0; }
}

.modalBody { display: grid; gap: $s-3; }

.modalActions {
  display: flex;
  gap: $s-2;
  justify-content: flex-end;
  flex-wrap: wrap;
  padding-top: $s-2;
}

.hint { margin: 0; color: $text-2; font-size: .9rem; }
.error { margin: 0; color: $red-0; font-weight: 600; }

@media (max-width: 900px) {
  .money { grid-template-columns: 1fr; }
  .twoCols { grid-template-columns: 1fr; }
  .row { grid-template-columns: 1fr; }
  .actions { justify-content: flex-start; }
}
</style>
