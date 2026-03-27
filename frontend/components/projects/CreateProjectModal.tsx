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

interface ApiError {
  response?: {
    status: number
    data?: {
      error?: string
      message?: string
    }
  }
}

function CreateProjectModal({ isOpen, onClose }: CreateProjectModalProps) {
  const [name, setName] = React.useState('')
  const [description, setDescription] = React.useState('')
  const [errors, setErrors] = React.useState<{ name?: string; description?: string; general?: string }>({})

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

    setErrors(prev => ({ ...prev, ...newErrors }))
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})

    if (!validate()) return

    try {
      await createProject.mutateAsync({
        name: name.trim(),
        description: description.trim(),
      })
      handleClose()
    } catch (error) {
      const apiError = error as ApiError
      
      // Handle 409 Conflict - duplicate name
      if (apiError.response?.status === 409) {
        setErrors(prev => ({
          ...prev,
          name: apiError.response?.data?.message || 'A project with this name already exists'
        }))
        return
      }
      
      // Handle validation errors
      if (apiError.response?.status === 400) {
        setErrors(prev => ({
          ...prev,
          general: apiError.response?.data?.message || 'Please check your input and try again'
        }))
        return
      }
      
      // Handle other errors
      setErrors(prev => ({
        ...prev,
        general: apiError.response?.data?.message || 'Failed to create project. Please try again.'
      }))
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title="Create New Project"
      description="Add a new project to organize your tasks"
      size="md"
    >
      <form onSubmit={handleSubmit} className="space-y-5">
        <Input
          label="Project Name"
          placeholder="Enter project name"
          value={name}
          onChange={(e) => {
            setName(e.target.value)
            if (errors.name) {
              setErrors(prev => ({ ...prev, name: undefined }))
            }
          }}
          error={errors.name}
          disabled={createProject.isPending}
          leftIcon={
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
            </svg>
          }
        />
        
        <div className="w-full">
          <label className="block text-sm font-medium text-gray-700 mb-1.5">
            Description <span className="text-gray-400">(optional)</span>
          </label>
          <textarea
            className="w-full px-4 py-2.5 bg-white border border-gray-300 rounded-lg text-gray-900 placeholder-gray-400 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 resize-none disabled:bg-gray-100 disabled:cursor-not-allowed"
            rows={3}
            placeholder="Enter project description"
            value={description}
            onChange={(e) => {
              setDescription(e.target.value)
              if (errors.description) {
                setErrors(prev => ({ ...prev, description: undefined }))
              }
            }}
            disabled={createProject.isPending}
          />
          {errors.description && (
            <p className="mt-1.5 text-sm text-red-600">{errors.description}</p>
          )}
        </div>

        {errors.general && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex items-start gap-2">
              <svg className="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p className="text-sm text-red-600">{errors.general}</p>
            </div>
          </div>
        )}

        <div className="flex justify-end gap-3 pt-2">
          <Button
            type="button"
            variant="ghost"
            onClick={handleClose}
            disabled={createProject.isPending}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            loading={createProject.isPending}
            disabled={createProject.isPending}
          >
            Create Project
          </Button>
        </div>
      </form>
    </Modal>
  )
}

export { CreateProjectModal }
