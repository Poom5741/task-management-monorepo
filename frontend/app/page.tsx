import Link from 'next/link'

export default function Home() {
  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-8">Task Management System</h1>
        
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Link
            href="/projects"
            className="p-6 border rounded-lg hover:bg-gray-50 transition"
          >
            <h2 className="text-xl font-semibold mb-2">Projects</h2>
            <p className="text-gray-600">Manage your projects</p>
          </Link>
          
          <Link
            href="/tasks"
            className="p-6 border rounded-lg hover:bg-gray-50 transition"
          >
            <h2 className="text-xl font-semibold mb-2">Tasks</h2>
            <p className="text-gray-600">View and manage tasks</p>
          </Link>
          
          <Link
            href="/labels"
            className="p-6 border rounded-lg hover:bg-gray-50 transition"
          >
            <h2 className="text-xl font-semibold mb-2">Labels</h2>
            <p className="text-gray-600">Organize with labels</p>
          </Link>
        </div>
      </div>
    </main>
  )
}