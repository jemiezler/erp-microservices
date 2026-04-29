import Link from 'next/link'
import { ArrowRight, LucideIcon } from 'lucide-react'
import { AppCard } from '@erp/ui'
import { cn } from '@/lib/utils'

interface ServiceCardProps {
  title: string
  description: string
  icon: LucideIcon
  href?: string
  comingSoon?: boolean
  accentColor?: string
}

export function ServiceCard({
  title,
  description,
  icon: Icon,
  href,
  comingSoon,
  accentColor = 'blue'
}: ServiceCardProps) {
  const CardContent = (
    <div className={cn(
      'relative flex h-full flex-col justify-between p-8',
      comingSoon && 'opacity-60'
    )}>
      {comingSoon && (
        <div className={cn(
          'absolute top-4 right-4 rounded-full px-3 py-1 text-xs font-semibold',
          accentColor === 'amber' ? 'bg-amber-100 text-amber-700' : 
          accentColor === 'purple' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'
        )}>
          Coming Soon
        </div>
      )}
      
      <div>
        <div className={cn(
          'mb-4 flex h-12 w-12 items-center justify-center rounded-xl transition-transform duration-300',
          accentColor === 'blue' ? 'bg-linear-to-br from-blue-100 to-blue-50' :
          accentColor === 'amber' ? 'bg-linear-to-br from-amber-100 to-amber-50' :
          accentColor === 'purple' ? 'bg-linear-to-br from-purple-100 to-purple-50' : 'bg-slate-100',
          !comingSoon && 'group-hover:scale-110'
        )}>
          <Icon className={cn(
            accentColor === 'blue' ? 'text-blue-600' :
            accentColor === 'amber' ? 'text-amber-600' :
            accentColor === 'purple' ? 'text-purple-600' : 'text-slate-600'
          )} size={24} />
        </div>
        <h2 className='mb-3 text-2xl font-bold text-slate-900'>
          {title}
        </h2>
        <p className='leading-relaxed text-slate-600'>
          {description}
        </p>
      </div>

      {!comingSoon ? (
        <div className='flex items-center gap-2 font-medium text-blue-600 transition-all duration-300 group-hover:gap-3'>
          <span>Access Module</span>
          <ArrowRight size={18} />
        </div>
      ) : (
        <div className='font-medium text-slate-500'>Available Soon</div>
      )}
    </div>
  )

  const card = (
    <AppCard className={cn(
      'group relative h-96 overflow-hidden p-0 transition-all duration-300',
      !comingSoon && 'cursor-pointer hover:border-blue-300 hover:shadow-lg'
    )}>
      {!comingSoon && (
        <>
          <div className='absolute inset-0 bg-linear-to-br from-blue-50 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100' />
          <div className='absolute right-0 bottom-0 left-0 h-1 scale-x-0 transform bg-linear-to-r from-blue-500 to-transparent transition-transform duration-300 group-hover:scale-x-100' />
        </>
      )}
      {CardContent}
    </AppCard>
  )

  if (href && !comingSoon) {
    return (
      <Link href={href}>
        {card}
      </Link>
    )
  }

  return card
}