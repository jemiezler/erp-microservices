import React from 'react'
import { cn } from '../lib/utils'

export interface BentoGridProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode
}

export const BentoGrid = ({ children, className, ...props }: BentoGridProps) => {
  return (
    <div
      className={cn(
        'grid grid-cols-1 md:grid-cols-12 gap-6 p-10',
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
}

export interface BentoCardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode
  span?: number
}

export const BentoCard = ({ children, className, span = 4, ...props }: BentoCardProps) => {
  return (
    <div
      className={cn(
        'bg-white rounded-(--radius-bento) p-6 shadow-(--shadow-bento) border border-black/5 overflow-hidden',
        span === 12 ? 'md:col-span-12' :
        span === 8 ? 'md:col-span-8' :
        span === 6 ? 'md:col-span-6' :
        span === 4 ? 'md:col-span-4' :
        span === 3 ? 'md:col-span-3' : 'md:col-span-4',
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
}
