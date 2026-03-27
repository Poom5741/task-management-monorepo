import Link from 'next/link'
import { cn } from '../lib/utils'
import { Card, CardTitle, CardDescription } from '../components/ui/Card'

const navigationItems = [
  {
    title: 'Projects',
    description: 'Manage your projects',
    href: '/projects',
    icon: (
      <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
      </svg>
    ),
    color: 'text-primary-600',
    bgColor: 'bg-primary-50',
  },
  {
    title: 'Tasks',
    description: 'View and manage tasks',
    href: '/tasks',
    icon: (
      <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
      </svg>
    ),
    color: 'text-cta-600',
    bgColor: 'bg-cta-50',
  },
  {
    title: 'Labels',
    description: 'Organize with labels',
    href: '/labels',
    icon: (
      <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
      </svg>
    ),
    color: 'text-emerald-600',
    bgColor: 'bg-emerald-50',
  },
]

export default function Home() {
  return (
    <main className="min-h-screen pt-24 pb-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-12">
          <h1 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4">
            Welcome to{' '}
            <span className="text-primary-600">TaskFlow</span>
          </h1>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            A powerful, intuitive task management system designed for modern teams. 
            Organize your work, track progress, and achieve your goals.
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {navigationItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="block group"
            >
              <Card hover padding="lg">
                <div className="flex flex-col items-center text-center">
                  <div className={cn(
                    'p-4 rounded-2xl mb-4 transition-transform duration-200 group-hover:scale-110',
                    item.bgColor
                  )}>
                    <div className={item.color}>
                      {item.icon}
                    </div>
                  </div>
                  <CardTitle className="mb-2">{item.title}</CardTitle>
                  <CardDescription>{item.description}</CardDescription>
                </div>
              </Card>
            </Link>
          ))}
        </div>

        <div className="mt-16 text-center">
          <p className="text-sm text-gray-500">
            Get started by selecting one of the options above
          </p>
        </div>
      </div>
    </main>
  )
}
