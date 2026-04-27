import { logger } from '@erp/logger';
import Link from 'next/link';
import { ArrowRight, Users, TrendingUp, Lock, Zap } from 'lucide-react';

export default function Home() {
  logger.info('host-app', 'Home page rendered');
  return (
    <div className='min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-slate-50'>
      {/* Header */}
      <header className='border-b border-slate-200/50 bg-white/80 backdrop-blur-md sticky top-0 z-50'>
        <div className='max-w-7xl mx-auto px-6 py-4 flex items-center justify-between'>
          <div className='flex items-center gap-2'>
            <div className='w-8 h-8 bg-gradient-to-br from-blue-600 to-blue-700 rounded-lg flex items-center justify-center'>
              <span className='text-white font-bold text-sm'>E</span>
            </div>
            <span className='font-semibold text-slate-900'>ERP Suite</span>
          </div>
          <nav className='hidden md:flex items-center gap-8 text-sm'>
            <a href='#' className='text-slate-600 hover:text-slate-900 transition-colors'>Services</a>
            <a href='#' className='text-slate-600 hover:text-slate-900 transition-colors'>Documentation</a>
            <a href='#' className='text-slate-600 hover:text-slate-900 transition-colors'>Support</a>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <section className='max-w-7xl mx-auto px-6 py-16 sm:py-20'>
        <div className='mb-12'>
          <h1 className='text-5xl sm:text-6xl font-bold text-slate-900 mb-6 leading-tight'>
            Unified Enterprise Platform
          </h1>
          <p className='text-xl text-slate-600 max-w-2xl'>
            Streamline operations across HR, Finance, and beyond with our integrated microservices architecture designed for modern organizations.
          </p>
        </div>

        {/* Services Grid */}
        <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'>
          {/* HR Service - Active */}
          <Link href='/hr'>
            <div className='group relative h-96 bg-white rounded-2xl border border-slate-200 hover:border-blue-300 shadow-sm hover:shadow-lg transition-all duration-300 overflow-hidden cursor-pointer'>
              {/* Background accent */}
              <div className='absolute inset-0 bg-gradient-to-br from-blue-50 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300' />
              
              {/* Content */}
              <div className='relative p-8 h-full flex flex-col justify-between'>
                <div>
                  <div className='mb-4 w-12 h-12 bg-gradient-to-br from-blue-100 to-blue-50 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300'>
                    <Users className='text-blue-600' size={24} />
                  </div>
                  <h2 className='text-2xl font-bold text-slate-900 mb-3'>HR Management</h2>
                  <p className='text-slate-600 leading-relaxed'>
                    Centralized employee management, payroll integration, and performance tracking for your organization.
                  </p>
                </div>
                <div className='flex items-center gap-2 text-blue-600 font-medium group-hover:gap-3 transition-all duration-300'>
                  <span>Access Module</span>
                  <ArrowRight size={18} />
                </div>
              </div>

              {/* Bottom border accent */}
              <div className='absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-500 to-transparent transform scale-x-0 group-hover:scale-x-100 transition-transform duration-300' />
            </div>
          </Link>

          {/* Finance Service - Coming Soon */}
          <div className='relative h-96 bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden opacity-60'>
            {/* Badge */}
            <div className='absolute top-4 right-4 px-3 py-1 bg-amber-100 text-amber-700 text-xs font-semibold rounded-full'>
              Coming Soon
            </div>

            <div className='p-8 h-full flex flex-col justify-between'>
              <div>
                <div className='mb-4 w-12 h-12 bg-gradient-to-br from-amber-100 to-amber-50 rounded-xl flex items-center justify-center'>
                  <TrendingUp className='text-amber-600' size={24} />
                </div>
                <h2 className='text-2xl font-bold text-slate-900 mb-3'>Finance & Accounting</h2>
                <p className='text-slate-600 leading-relaxed'>
                  Advanced accounting, invoice management, expense tracking, and financial reporting tools.
                </p>
              </div>
              <div className='text-slate-500 font-medium'>
                Available Soon
              </div>
            </div>
          </div>

          {/* Security Service - Coming Soon */}
          <div className='relative h-96 bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden opacity-60'>
            {/* Badge */}
            <div className='absolute top-4 right-4 px-3 py-1 bg-purple-100 text-purple-700 text-xs font-semibold rounded-full'>
              Coming Soon
            </div>

            <div className='p-8 h-full flex flex-col justify-between'>
              <div>
                <div className='mb-4 w-12 h-12 bg-gradient-to-br from-purple-100 to-purple-50 rounded-xl flex items-center justify-center'>
                  <Lock className='text-purple-600' size={24} />
                </div>
                <h2 className='text-2xl font-bold text-slate-900 mb-3'>Security & Compliance</h2>
                <p className='text-slate-600 leading-relaxed'>
                  Role-based access control, audit logs, compliance monitoring, and security policies.
                </p>
              </div>
              <div className='text-slate-500 font-medium'>
                Available Soon
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className='max-w-7xl mx-auto px-6 py-16 border-t border-slate-200'>
        <div className='grid grid-cols-1 md:grid-cols-3 gap-8'>
          <div className='p-6'>
            <div className='flex items-center gap-3 mb-4'>
              <Zap className='text-blue-600' size={20} />
              <h3 className='font-semibold text-slate-900'>Fast & Scalable</h3>
            </div>
            <p className='text-slate-600 text-sm'>Microservices architecture built for high performance and zero downtime.</p>
          </div>
          <div className='p-6'>
            <div className='flex items-center gap-3 mb-4'>
              <Lock className='text-blue-600' size={20} />
              <h3 className='font-semibold text-slate-900'>Enterprise Secure</h3>
            </div>
            <p className='text-slate-600 text-sm'>JWT authentication, encrypted communications, and comprehensive audit trails.</p>
          </div>
          <div className='p-6'>
            <div className='flex items-center gap-3 mb-4'>
              <Users className='text-blue-600' size={20} />
              <h3 className='font-semibold text-slate-900'>Team Collaboration</h3>
            </div>
            <p className='text-slate-600 text-sm'>Unified interface for your entire organization to work seamlessly.</p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className='border-t border-slate-200 bg-white/50 backdrop-blur-sm'>
        <div className='max-w-7xl mx-auto px-6 py-8 flex items-center justify-between text-sm text-slate-600'>
          <p>ERP Suite v1.0.0</p>
          <div className='flex gap-6'>
            <a href='#' className='hover:text-slate-900 transition-colors'>Privacy</a>
            <a href='#' className='hover:text-slate-900 transition-colors'>Terms</a>
            <a href='#' className='hover:text-slate-900 transition-colors'>Status</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
