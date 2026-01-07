import { inject, type InjectionKey } from "vue"

export type ConfirmVariant = "primary" | "secondary" | "ghost" | "danger"

export type ConfirmRequest = {
  title?: string
  message: string
  yesText?: string
  noText?: string
  yesVariante?: ConfirmVariant
  noVariante?: ConfirmVariant
}

export type ConfirmApi = {
  confirm: (req: ConfirmRequest) => Promise<boolean>
}

export const ConfirmKey: InjectionKey<ConfirmApi> = Symbol("ConfirmKey") as InjectionKey<ConfirmApi>

export function useConfirm() {
  const api = inject(ConfirmKey, null)
  if (!api) {
    throw new Error("useConfirm() muss innerhalb von <VConfirmHost> verwendet werden.")
  }
  return api.confirm
}
