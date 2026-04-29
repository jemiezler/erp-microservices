'use client'
import { Users, TrendingUp, Lock, Zap } from 'lucide-react'
import { ServiceCard } from '@/modules/home/components/service-card'
import { FeatureCard } from '@/modules/home/components/feature-card'

export default function Home() {
  return (
    <div className='min-h-screen bg-[#F5F5F5]'>
      {/* Hero Section */}
      <section className='mx-auto max-w-7xl px-6 py-16 sm:py-20'>
        <div className='mb-12'>
          <h1 className='mb-6 text-5xl leading-tight font-bold text-slate-900 sm:text-6xl'>
            Unified Enterprise Platform
          </h1>
          <p className='max-w-2xl text-xl text-slate-600'>
            Streamline operations across HR, Finance, and beyond with our
            integrated microservices architecture designed for modern
            organizations.
          </p>
        </div>

        {/* Services Grid */}
        <div className='grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3'>
          <ServiceCard
            title='HR Management'
            description='Centralized employee management, payroll integration, and performance tracking for your organization.'
            icon={Users}
            href='/hr'
          />

          <ServiceCard
            title='Finance & Accounting'
            description='Advanced accounting, invoice management, expense tracking, and financial reporting tools.'
            icon={TrendingUp}
            comingSoon
            accentColor='amber'
          />

          <ServiceCard
            title='Security & Compliance'
            description='Role-based access control, audit logs, compliance monitoring, and security policies.'
            icon={Lock}
            comingSoon
            accentColor='purple'
          />
        </div>
      </section>

      {/* Features Section */}
      <section className='mx-auto max-w-7xl border-t border-slate-200 px-6 py-16'>
        <div className='grid grid-cols-1 gap-8 md:grid-cols-3'>
          <FeatureCard
            title='Fast & Scalable'
            description='Microservices architecture built for high performance and zero downtime.'
            icon={Zap}
          />
          <FeatureCard
            title='Enterprise Secure'
            description='JWT authentication, encrypted communications, and comprehensive audit trails.'
            icon={Lock}
          />
          <FeatureCard
            title='Team Collaboration'
            description='Unified interface for your entire organization to work seamlessly.'
            icon={Users}
          />
        </div>
      </section>

      {/* Footer */}
      <footer className='border-t border-slate-200 bg-white/50 backdrop-blur-sm'>
        <div className='mx-auto flex max-w-7xl items-center justify-between px-6 py-8 text-sm text-slate-600'>
          <p>ERP Suite v1.0.0</p>
          <div className='flex gap-6'>
            <a href='#' className='transition-colors hover:text-slate-900'>
              Privacy
            </a>
            <a href='#' className='transition-colors hover:text-slate-900'>
              Terms
            </a>
            <a href='#' className='transition-colors hover:text-slate-900'>
              Status
            </a>
          </div>
        </div>
      </footer>
    </div>
  )
}