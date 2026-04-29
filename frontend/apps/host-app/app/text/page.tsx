import { AppCard } from "@erp/ui"
import { ReactNode } from "react"

export type CardType = "kpi" | "hero" | "chart" | "reminder"

export interface BaseCardConfig {
  id: string
  type: CardType
  className?: string
}

export interface KPICardConfig extends BaseCardConfig {
  type: "kpi"
  title: string
  subtitle?: string
  value: string
  delta?: string
  deltaColor?: "green" | "red"
}

export interface HeroCardConfig extends BaseCardConfig {
  type: "hero"
  title: string
  description: string
  image?: ReactNode
}

export interface ChartCardConfig extends BaseCardConfig {
  type: "chart"
  title: string
  subtitle?: string
  data: { label: string; value: number }[]
}

export interface ReminderCardConfig extends BaseCardConfig {
  type: "reminder"
  title: string
  items: {
    label: string
    meta?: string
    checked?: boolean
  }[]
}

export type CardConfig =
  | KPICardConfig
  | HeroCardConfig
  | ChartCardConfig
  | ReminderCardConfig
export const dashboardCards: CardConfig[] = [
  // KPI ROW
  {
    id: "production",
    type: "kpi",
    title: "Production",
    subtitle: "Efficiency",
    value: "12%",
    delta: "-$10K",
    deltaColor: "red",
  },
  {
    id: "demand",
    type: "kpi",
    title: "Demand",
    subtitle: "Trends",
    value: "16%",
    delta: "+$18K",
    deltaColor: "green",
  },
  {
    id: "inventory",
    type: "kpi",
    title: "Inventory",
    subtitle: "Turnover",
    value: "18%",
    delta: "+$14K",
    deltaColor: "green",
  },
  {
    id: "shipping",
    type: "kpi",
    title: "Shipping",
    subtitle: "Efficiency",
    value: "22%",
    delta: "+$26K",
    deltaColor: "green",
  },

  // HERO CARD
  {
    id: "billing",
    type: "hero",
    title: "Automate billing",
    description: "Receive, process and send - all within seconds.",
  },

  // CHART CARD
  {
    id: "sales",
    type: "chart",
    title: "Change of Sales",
    subtitle: "Last three months",
    data: [
      { label: "Jan", value: 12620 },
      { label: "Feb", value: 14850 },
      { label: "Mar", value: 17642 },
    ],
  },

  // REMINDER CARD
  {
    id: "reminders",
    type: "reminder",
    title: "Reminders",
    items: [
      {
        label: "Cotton Fabric: 23 left",
        meta: "Added 3 days ago",
        checked: true,
      },
      {
        label: "Shipment #5098 delayed",
        meta: "Added 4 days ago",
        checked: false,
      },
      { label: "Budget review 3/5", meta: "Added 6 days ago", checked: false },
    ],
  },
]

export default function DashboardGrid() {
  return (
    <div className="grid grid-cols-12 gap-4">
      {dashboardCards.map((card) => (
        <CardRenderer key={card.id} card={card} />
      ))}
    </div>
  )
}

const CardRenderer = ({ card }: { card: CardConfig }) => {
  switch (card.type) {
    case "kpi":
      return (
        <AppCard className="col-span-3">
          <div className="text-muted-foreground text-sm">{card.subtitle}</div>
          <div className="text-3xl font-semibold">{card.value}</div>
          <div
            className={
              card.deltaColor === "green"
                ? "text-xs text-green-500"
                : "text-xs text-red-500"
            }
          >
            {card.delta}
          </div>
        </AppCard>
      )

    case "hero":
      return (
        <AppCard className="col-span-6 flex items-center justify-between">
          <div>
            <h2 className="text-xl font-bold">{card.title}</h2>
            <p className="text-muted-foreground">{card.description}</p>
          </div>
        </AppCard>
      )

    case "chart":
      return (
        <AppCard className="col-span-6">
          <div className="mb-4">
            <h3 className="font-semibold">{card.title}</h3>
            <p className="text-muted-foreground text-sm">{card.subtitle}</p>
          </div>
          <div className="bg-muted h-40 rounded-xl" />
        </AppCard>
      )

    case "reminder":
      return (
        <AppCard className="col-span-4">
          <h3 className="mb-3 font-semibold">{card.title}</h3>
          <div className="space-y-3">
            {card.items.map((item, idx) => (
              <div key={idx} className="flex justify-between text-sm">
                <div>
                  <p>{item.label}</p>
                  <p className="text-muted-foreground text-xs">{item.meta}</p>
                </div>
                <div>{item.checked ? "✓" : ""}</div>
              </div>
            ))}
          </div>
        </AppCard>
      )
  }
}
