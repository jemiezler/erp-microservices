"use client"
import { useState } from "react"
import { Button } from "@erp/ui"
import { LayoutDashboard } from "lucide-react"

const navItems = [
  { label: "Dashboard" },
  { label: "Services" },
  { label: "Reports" },
]

export default function AppHeader() {
  const [active, setActive] = useState("Dashboard")

  return (
    <header className="flex items-center justify-between px-10 py-6">
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-black">
          <LayoutDashboard className="text-white" size={20} />
        </div>
        <span className="text-xl font-bold tracking-tight text-black uppercase">
          ERP CORE
        </span>
      </div>

      <nav className="hidden items-stretch gap-2 rounded-full border p-1 md:grid md:grid-cols-3">
        {navItems.map((item) => {
          const isActive = active === item.label

          return (
            <Button
              key={item.label}
              onClick={() => setActive(item.label)}
              className={[
                "rounded-full px-6 font-medium transition-all",
                isActive
                  ? "bg-black text-white shadow-sm"
                  : "bg-transparent text-black hover:bg-black/10",
              ].join(" ")}
            >
              {item.label}
            </Button>
          )
        })}
      </nav>

      <div className="flex items-center gap-4 rounded-full border bg-white">
        <Button
          variant="ghost"
          size="sm"
          className="hidden rounded-full px-4 font-medium text-black/60 sm:inline-flex"
        >
          Support
        </Button>

        <div className="h-10 w-10 overflow-hidden rounded-full border-2 border-white bg-slate-200 shadow-sm">
          <img
            src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix"
            alt="avatar"
          />
        </div>
      </div>
    </header>
  )
}
