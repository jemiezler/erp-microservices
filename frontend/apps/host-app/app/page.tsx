import { logger } from '@erp/logger';
import Link from 'next/link';

export default function Home() {
  logger.info('host-app', 'Home page rendered');
  return (
    <div className='flex flex-col items-center justify-center min-h-screen p-8 font-sans bg-slate-50'>
      <main className='max-w-4xl w-full text-center'>
        <h1 className='text-4xl font-extrabold text-slate-900 mb-4 tracking-tight'>
          ERP Microservices Platform
        </h1>
        <p className='text-lg text-slate-600 mb-12'>
          Consolidated Host Application managing multiple Micro-Frontends.
        </p>

        <div className='grid grid-cols-1 md:grid-cols-2 gap-8'>
          <div className='p-6 bg-white border border-slate-200 rounded-xl shadow-sm hover:shadow-md transition-shadow'>
            <h2 className='text-xl font-bold text-slate-800 mb-2'>HR Management</h2>
            <p className='text-slate-500 mb-6'>
              Manage employees, roles, and status through the HR Micro-Frontend.
            </p>
            <Link 
              href='/hr' 
              className='inline-block px-6 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors'
            >
              Go to HR Dashboard
            </Link>
          </div>

          <div className='p-6 bg-white border border-slate-200 rounded-xl shadow-sm opacity-50'>
            <h2 className='text-xl font-bold text-slate-800 mb-2'>Finance (Coming Soon)</h2>
            <p className='text-slate-500 mb-6'>
              Accounting, payroll, and billing systems are currently being integrated.
            </p>
            <button disabled className='px-6 py-2 bg-slate-300 text-white font-medium rounded-lg cursor-not-allowed'>
              Finance Dashboard
            </button>
          </div>
        </div>
      </main>

      <footer className='mt-16 text-slate-400 text-sm'>
        [SYSTEM] Host App Interface v1.0.0
      </footer>
    </div>
  );
}
