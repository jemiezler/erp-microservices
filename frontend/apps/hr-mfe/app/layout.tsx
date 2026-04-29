import type { Metadata } from "next";
import { Geist_Mono, Inter } from "next/font/google";
import "./globals.css";
import Link from 'next/link';
import { Button, ToastProvider, cn } from '@erp/ui';
import { ArrowLeft } from 'lucide-react';

const inter = Inter({subsets:['latin'],variable:'--font-sans'})

const fontMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
})


export const metadata: Metadata = {
  title: "HR Management | ERP Suite",
  description: "Secure HR Module",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={cn("antialiased", fontMono.variable, "font-sans", inter.variable)}
    >
      <body className="min-h-full flex flex-col bg-background text-foreground">
        {/* Domain-specific Header for HR Zone */}
        <header className='border-b border-border bg-card/80 backdrop-blur-md sticky top-0 z-50'>
          <div className='max-w-7xl mx-auto px-6 py-4'>
            <div className='flex items-center justify-between'>
              <div className='flex items-center gap-4'>
                <Link href="/" prefetch={false}>
                  <Button variant="ghost" size="icon" className="rounded-full text-muted-foreground hover:text-foreground">
                    <ArrowLeft size={20} />
                  </Button>
                </Link>
                <div className='flex items-center gap-2'>
                  <div className='w-8 h-8 bg-primary rounded-lg flex items-center justify-center shadow-sm'>
                    <span className='text-primary-foreground font-bold text-sm'>E</span>
                  </div>
                  <span className='font-semibold text-foreground'>ERP Suite</span>
                </div>
              </div>
              <div className="flex items-center gap-2">
                 <span className='px-2 py-0.5 bg-secondary text-secondary-foreground rounded text-[10px] font-black uppercase'>HR Zone</span>
              </div>
            </div>
          </div>
        </header>

        <main className="flex-1">{children}</main>

        <footer className="bg-card border-t border-border py-6">
            <div className="max-w-7xl mx-auto px-6 text-center text-muted-foreground text-xs">
                &copy; 2026 ERP HR Domain • Isolated Process Environment
            </div>
        </footer>
        <ToastProvider />
      </body>
    </html>
  );
}
