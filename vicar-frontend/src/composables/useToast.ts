import { inject, type InjectionKey } from "vue"

export type ToastVariant = "success" | "error" | "warning" | "info"

export type ToastOptions = {
  title?: string
  message: string
  variant?: ToastVariant
  durationMs?: number
  dismissible?: boolean
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

export const ToastKey: InjectionKey<ToastApi> = Symbol("ToastKey") as InjectionKey<ToastApi>

export function useToast() {
  const api = inject(ToastKey, null)
  if (!api) throw new Error("useToast() muss innerhalb von <VToastHost> verwendet werden.")
  return api
}
