import type { Metadata } from 'next'
import './globals.css'
import { Providers } from './providers'
import { Navbar } from '../components/layout/Navbar'

export const metadata: Metadata = {
  title: 'TaskFlow - Modern Task Management',
  description: 'A powerful, intuitive task management system designed for modern teams',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="font-sans bg-primary-50">
        <Providers>
          <Navbar />
          {children}
        </Providers>
      </body>
    </html>
  )
}
