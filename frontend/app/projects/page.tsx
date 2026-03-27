'use client'

import * as React from 'react'
import { Button } from '../../components/ui/Button'
import { CreateProjectModal } from '../../components/projects/CreateProjectModal'
import { useProjects } from '../../lib/api/projects'

function ProjectsPage() {
  const [isModalOpen, setIsModalOpen] = React.useState(false)
  const { data, isLoading, error } = useProjects()

  return React.createElement(
    'div',
    { className: 'container mx-auto px-4 py-8' },
    React.createElement(
      'div',
      { className: 'flex justify-between items-center mb-8' },
      React.createElement('h1', { className: 'text-2xl font-bold text-gray-900' }, 'Projects'),
      React.createElement(
        Button,
        { onClick: () => setIsModalOpen(true) },
        '+ Create Project'
      )
    ),
    isLoading &&
      React.createElement(
        'div',
        { className: 'text-center py-12' },
        React.createElement('p', { className: 'text-gray-500' }, 'Loading projects...')
      ),
    error &&
      React.createElement(
        'div',
        { className: 'text-center py-12' },
        React.createElement('p', { className: 'text-red-500' }, 'Failed to load projects')
      ),
    data?.data &&
      data.data.length === 0 &&
      React.createElement(
        'div',
        { className: 'text-center py-12' },
        React.createElement(
          'p',
          { className: 'text-gray-500 mb-4' },
          'No projects yet. Create your first project!'
        )
      ),
    data?.data &&
      data.data.length > 0 &&
      React.createElement(
        'div',
        { className: 'grid gap-4 md:grid-cols-2 lg:grid-cols-3' },
        data.data.map((project) =>
          React.createElement(
            'div',
            {
              key: project.id,
              className:
                'bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow',
            },
            React.createElement(
              'div',
              { className: 'flex items-start justify-between' },
              React.createElement('h3', { className: 'font-semibold text-gray-900' }, project.name),
              React.createElement(
                'span',
                {
                  className:
                    project.status === 'active'
                      ? 'text-xs px-2 py-1 bg-green-100 text-green-700 rounded-full'
                      : 'text-xs px-2 py-1 bg-gray-100 text-gray-700 rounded-full',
                },
                project.status
              )
            ),
            project.description &&
              React.createElement(
                'p',
                { className: 'mt-2 text-sm text-gray-600 line-clamp-2' },
                project.description
              )
          )
        )
      ),
    React.createElement(CreateProjectModal, {
      isOpen: isModalOpen,
      onClose: () => setIsModalOpen(false),
    })
  )
}

export default ProjectsPage