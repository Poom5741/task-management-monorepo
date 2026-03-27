'use client'

import * as React from 'react'
import { Modal } from '../ui/Modal'
import { Button } from '../ui/Button'
import { useDeleteProject } from '../../lib/api/projects'

interface DeleteProjectModalProps {
  isOpen: boolean
  onClose: () => void
  projectId: string
  projectName: string
  onSuccess?: () => void
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

function DeleteProjectModal({ isOpen, onClose, projectId, projectName, onSuccess }: DeleteProjectModalProps) {
  const [error, setError] = React.useState<string>('')
  const deleteProject = useDeleteProject()

  React.useEffect(() => {
    if (isOpen) {
      setError('')
    }
  }, [isOpen])

  const handleClose = () => {
    setError('')
    onClose()
  }

  const handleDelete = async () => {
    setError('')

    try {
      await deleteProject.mutateAsync(projectId)
      handleClose()
      onSuccess?.()
    } catch (err) {
      const apiError = err as ApiError
      
      if (apiError.response?.status === 404) {
        setError('Project not found')
        return
      }
      
      setError(apiError.response?.data?.message || 'Failed to delete project. Please try again.')
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title="Delete Project"
      description="This action cannot be undone"
      size="md"
    >
      <div className="space-y-4">
        <div className="flex items-start gap-3 p-4 bg-red-50 border border-red-200 rounded-lg">
          <svg className="w-6 h-6 text-red-600 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <div>
            <p className="font-medium text-red-900">
              Are you sure you want to delete this project?
            </p>
            <p className="mt-1 text-sm text-red-700">
              This will also delete all tasks associated with "{projectName}". This action cannot be undone.
            </p>
          </div>
        </div>

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex items-start gap-2">
              <svg className="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p className="text-sm text-red-600">{error}</p>
            </div>
          </div>
        )}

        <div className="flex justify-end gap-3 pt-2">
          <Button
            type="button"
            variant="ghost"
            onClick={handleClose}
            disabled={deleteProject.isPending}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="danger"
            onClick={handleDelete}
            loading={deleteProject.isPending}
            disabled={deleteProject.isPending}
          >
            Delete Project
          </Button>
        </div>
      </div>
    </Modal>
  )
}

export { DeleteProjectModal }
