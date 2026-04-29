import { LucideIcon } from 'lucide-react'
import { AppCard } from '@erp/ui'

interface FeatureCardProps {
  title: string
  description: string
  icon: LucideIcon
}

export function FeatureCard({ title, description, icon: Icon }: FeatureCardProps) {
  return (
    <AppCard className='border-none shadow-none bg-transparent p-6 hover:bg-slate-50 transition-colors'>
      <div className='mb-4 flex items-center gap-3'>
        <Icon className='text-blue-600' size={20} />
        <h3 className='font-semibold text-slate-900'>{title}</h3>
      </div>
      <p className='text-sm text-slate-600'>
        {description}
      </p>
    </AppCard>
  )
}