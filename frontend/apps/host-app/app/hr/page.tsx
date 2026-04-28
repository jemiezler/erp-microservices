'use client';

import { logger } from '@erp/logger';
import Link from 'next/link';
import { ChevronRight } from 'lucide-react';
import { useEffect, useState } from 'react';

export default function HRPage() {
  logger.info('host-app', 'HR page rendered');
  const [iframeReady, setIframeReady] = useState(false);

  useEffect(() => {
    setIframeReady(true);
  }, []);

  return (
    <div className='min-h-screen bg-slate-50'>
      {/* Header with Navigation */}
      <header className='border-b border-slate-200/50 bg-white/80 backdrop-blur-md sticky top-0 z-50'>
        <div className='max-w-7xl mx-auto px-6 py-4'>
          <div className='flex items-center justify-between mb-4'>
            <div className='flex items-center gap-2'>
              <div className='w-8 h-8 bg-gradient-to-br from-blue-600 to-blue-700 rounded-lg flex items-center justify-center'>
                <span className='text-white font-bold text-sm'>E</span>
              </div>
              <span className='font-semibold text-slate-900'>ERP Suite</span>
            </div>
            <nav className='hidden md:flex items-center gap-8 text-sm'>
              <Link href='/' className='text-slate-600 hover:text-slate-900 transition-colors'>Home</Link>
              <a href='#' className='text-slate-600 hover:text-slate-900 transition-colors'>Documentation</a>
              <a href='#' className='text-slate-600 hover:text-slate-900 transition-colors'>Support</a>
            </nav>
          </div>
          
          {/* Breadcrumb */}
          <div className='flex items-center gap-2 text-sm text-slate-600'>
            <Link href='/' className='hover:text-slate-900 transition-colors'>Services</Link>
            <ChevronRight size={16} />
            <span className='text-slate-900 font-medium'>HR Management</span>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className='w-full'>
        {iframeReady && (
          <iframe
            src='http://localhost:3001'
            title='HR Dashboard'
            className='w-full border-none'
            style={{
              height: 'calc(100vh - 120px)',
              display: 'block',
            }}
            allow='*'
            sandbox='allow-same-origin allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-presentation'
          />
        )}
      </main>
    </div>
  );
}
