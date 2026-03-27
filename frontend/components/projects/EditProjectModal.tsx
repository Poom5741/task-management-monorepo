'use client'

import * as React from 'react'
import { Modal } from '../ui/Modal'
import { Button } from '../ui/Button'
import { Input } from '../ui/Input'
import { useUpdateProject } from '../../lib/api/projects'
import type { Project } from '../../lib/types/project'

interface EditProjectModalProps {
  isOpen: boolean
  onClose: () => void
  project: Project
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

function EditProjectModal({ isOpen, onClose, project }: EditProjectModalProps) {
  const [name, setName] = React.useState(project.name)
  const [description, setDescription] = React.useState(project.description)
  const [status, setStatus] = React.useState<'active' | 'archived'>(project.status)
  const [errors, setErrors] = React.useState<{ name?: string; description?: string; general?: string }>({})

  const updateProject = useUpdateProject()

  React.useEffect(() => {
    if (isOpen) {
      setName(project.name)
      setDescription(project.description)
      setStatus(project.status)
      setErrors({})
    }
  }, [isOpen, project])

  const resetForm = () => {
    setName(project.name)
    setDescription(project.description)
    setStatus(project.status)
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
      const updateInput: { name?: string; description?: string; status?: 'active' | 'archived' } = {}
      
      if (name.trim() !== project.name) {
        updateInput.name = name.trim()
      }
      if (description !== project.description) {
        updateInput.description = description
      }
      if (status !== project.status) {
        updateInput.status = status
      }

      await updateProject.mutateAsync({
        id: project.id,
        input: updateInput,
      })
      handleClose()
    } catch (error) {
      const apiError = error as ApiError
      
      if (apiError.response?.status === 409) {
        setErrors(prev => ({
          ...prev,
          name: apiError.response?.data?.message || 'A project with this name already exists'
        }))
        return
      }
      
      if (apiError.response?.status === 404) {
        setErrors(prev => ({
          ...prev,
          general: 'Project not found'
        }))
        return
      }
      
      if (apiError.response?.status === 400) {
        setErrors(prev => ({
          ...prev,
          general: apiError.response?.data?.message || 'Please check your input and try again'
        }))
        return
      }
      
      setErrors(prev => ({
        ...prev,
        general: apiError.response?.data?.message || 'Failed to update project. Please try again.'
      }))
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title="Edit Project"
      description="Update project information"
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
          disabled={updateProject.isPending}
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
            disabled={updateProject.isPending}
          />
          {errors.description && (
            <p className="mt-1.5 text-sm text-red-600">{errors.description}</p>
          )}
        </div>

        <div className="w-full">
          <label className="block text-sm font-medium text-gray-700 mb-1.5">
            Status
          </label>
          <select
            className="w-full px-4 py-2.5 bg-white border border-gray-300 rounded-lg text-gray-900 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
            value={status}
            onChange={(e) => setStatus(e.target.value as 'active' | 'archived')}
            disabled={updateProject.isPending}
          >
            <option value="active">Active</option>
            <option value="archived">Archived</option>
          </select>
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
            disabled={updateProject.isPending}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            loading={updateProject.isPending}
            disabled={updateProject.isPending}
          >
            Save Changes
          </Button>
        </div>
      </form>
    </Modal>
  )
}

export { EditProjectModal }
