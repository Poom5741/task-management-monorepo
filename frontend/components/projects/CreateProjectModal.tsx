'use client'

import * as React from 'react'
import { Modal } from '../ui/Modal'
import { Button } from '../ui/Button'
import { Input } from '../ui/Input'
import { useCreateProject } from '../../lib/api/projects'

interface CreateProjectModalProps {
  isOpen: boolean
  onClose: () => void
}

function CreateProjectModal({ isOpen, onClose }: CreateProjectModalProps) {
  const [name, setName] = React.useState('')
  const [description, setDescription] = React.useState('')
  const [errors, setErrors] = React.useState<{ name?: string; description?: string }>({})

  const createProject = useCreateProject()

  const resetForm = () => {
    setName('')
    setDescription('')
    setErrors({})
  }

  const handleClose = () => {
    resetForm()
    onClose()
  }

  const validate = (): boolean => {
    const newErrors: { name?: string; description?: string } = {}

    if (!name.trim()) {
      newErrors.name = 'Name is required'
    } else if (name.length > 100) {
      newErrors.name = 'Name must be at most 100 characters'
    }

    if (description.length > 500) {
      newErrors.description = 'Description must be at most 500 characters'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validate()) return

    try {
      await createProject.mutateAsync({
        name: name.trim(),
        description: description.trim(),
      })
      handleClose()
    } catch (error) {
      console.error('Failed to create project:', error)
    }
  }

  return React.createElement(
    Modal,
    {
      isOpen,
      onClose: handleClose,
      title: 'Create New Project',
      size: 'md',
    },
    React.createElement(
      'form',
      { onSubmit: handleSubmit, className: 'space-y-4' },
      React.createElement(Input, {
        label: 'Project Name',
        placeholder: 'Enter project name',
        value: name,
        onChange: (e) => setName(e.target.value),
        error: errors.name,
        disabled: createProject.isPending,
      }),
      React.createElement(
        'div',
        { className: 'w-full' },
        React.createElement(
          'label',
          { className: 'block text-sm font-medium text-gray-700 mb-1' },
          'Description (optional)'
        ),
        React.createElement('textarea', {
          className:
            'w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed',
          rows: 3,
          placeholder: 'Enter project description',
          value: description,
          onChange: (e) => setDescription(e.target.value),
          disabled: createProject.isPending,
        }),
        errors.description &&
          React.createElement(
            'p',
            { className: 'mt-1 text-sm text-red-600' },
            errors.description
          )
      ),
      createProject.isError &&
        React.createElement(
          'p',
          { className: 'text-sm text-red-600' },
          'Failed to create project. Please try again.'
        ),
      React.createElement(
        'div',
        { className: 'flex justify-end gap-3 pt-4' },
        React.createElement(
          Button,
          {
            type: 'button',
            variant: 'secondary',
            onClick: handleClose,
            disabled: createProject.isPending,
          },
          'Cancel'
        ),
        React.createElement(
          Button,
          {
            type: 'submit',
            loading: createProject.isPending,
            disabled: createProject.isPending,
          },
          'Create Project'
        )
      )
    )
  )
}

export { CreateProjectModal }