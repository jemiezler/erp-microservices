"use client"

import * as React from "react"
import { X, CircleAlert, CircleCheckBig, Info, TriangleAlert } from "lucide-react"

import { cn } from "../../lib/utils"

type ToastVariant = "default" | "success" | "error" | "warning"

type ToastInput = {
  title: string
  description?: string
  variant?: ToastVariant
  duration?: number
}

type ToastRecord = ToastInput & {
  id: number
}

type ToastListener = (toast: ToastRecord | null) => void

const toastListeners = new Set<ToastListener>()
let nextToastId = 0
let currentToast: ToastRecord | null = null

function publishToast(toast: ToastRecord | null) {
  currentToast = toast
  toastListeners.forEach((listener) => listener(toast))
}

function showToast(input: ToastInput) {
  const toast = {
    ...input,
    id: ++nextToastId,
  }

  publishToast(toast)
  return toast.id
}

function dismissToast() {
  publishToast(null)
}

export const toast = {
  show: showToast,
  dismiss: dismissToast,
  success: (title: string, description?: string, duration?: number) =>
    showToast({ title, description, variant: "success", duration }),
  error: (title: string, description?: string, duration?: number) =>
    showToast({ title, description, variant: "error", duration }),
  info: (title: string, description?: string, duration?: number) =>
    showToast({ title, description, variant: "default", duration }),
  warning: (title: string, description?: string, duration?: number) =>
    showToast({ title, description, variant: "warning", duration }),
}

function getVariantStyles(variant: ToastVariant) {
  switch (variant) {
    case "success":
      return {
        shell: "border-emerald-500/20 bg-emerald-50 text-emerald-950 shadow-emerald-950/5",
        icon: CircleCheckBig,
      }
    case "warning":
      return {
        shell: "border-amber-500/20 bg-amber-50 text-amber-950 shadow-amber-950/5",
        icon: TriangleAlert,
      }
    case "error":
      return {
        shell: "border-rose-500/20 bg-rose-50 text-rose-950 shadow-rose-950/5",
        icon: CircleAlert,
      }
    default:
      return {
        shell: "border-slate-200 bg-white text-slate-950 shadow-slate-950/10",
        icon: Info,
      }
  }
}

export function ToastProvider({ children }: { children?: React.ReactNode }) {
  const [activeToast, setActiveToast] = React.useState<ToastRecord | null>(null)

  React.useEffect(() => {
    const listener: ToastListener = (toast) => {
      setActiveToast(toast)
    }

    toastListeners.add(listener)
    listener(currentToast)

    return () => {
      toastListeners.delete(listener)
    }
  }, [])

  React.useEffect(() => {
    if (!activeToast) {
      return
    }

    const timeout = window.setTimeout(() => {
      setActiveToast(null)
    }, activeToast.duration ?? 4500)

    return () => {
      window.clearTimeout(timeout)
    }
  }, [activeToast])

  const variant = getVariantStyles(activeToast?.variant ?? "default")
  const Icon = variant.icon

  return (
    <>
      {children}

      <div
        aria-live="polite"
        aria-atomic="true"
        className="pointer-events-none fixed right-4 top-4 z-[100] flex w-[min(92vw,24rem)] flex-col items-end"
      >
        <div
          className={cn(
            "pointer-events-auto w-full rounded-2xl border px-4 py-3 shadow-lg transition-all duration-200 ease-out",
            activeToast ? "translate-y-0 opacity-100" : "translate-y-2 opacity-0",
            variant.shell
          )}
          data-state={activeToast ? "open" : "closed"}
          hidden={!activeToast}
        >
          {activeToast ? (
            <div className="flex items-start gap-3">
              <div className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-black/5">
                <Icon className="h-4 w-4" aria-hidden="true" />
              </div>

              <div className="min-w-0 flex-1">
                <p className="text-sm font-semibold leading-5">{activeToast.title}</p>
                {activeToast.description ? (
                  <p className="mt-1 text-sm leading-5 opacity-80">
                    {activeToast.description}
                  </p>
                ) : null}
              </div>

              <button
                type="button"
                onClick={() => setActiveToast(null)}
                className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full transition-colors hover:bg-black/5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-black/20"
                aria-label="Dismiss notification"
              >
                <X className="h-4 w-4" aria-hidden="true" />
              </button>
            </div>
          ) : null}
        </div>
      </div>
    </>
  )
}