'use client'

import * as React from 'react'
import { useParams, useRouter } from 'next/navigation'
import { Button } from '../../../components/ui/Button'
import { Card } from '../../../components/ui/Card'
import { Badge } from '../../../components/ui/Badge'
import { Skeleton } from '../../../components/ui/Skeleton'
import { EmptyState } from '../../../components/ui/EmptyState'
import { EditProjectModal } from '../../../components/projects/EditProjectModal'
import { DeleteProjectModal } from '../../../components/projects/DeleteProjectModal'
import { useProject, useDeleteProject } from '../../../lib/api/projects'

function ProjectDetailPage() {
  const params = useParams()
  const router = useRouter()
  const projectId = params.id as string
  const [isEditModalOpen, setIsEditModalOpen] = React.useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = React.useState(false)

  const { data: project, isLoading, error, refetch } = useProject(projectId)
  const deleteProject = useDeleteProject()

  React.useEffect(() => {
    if (!isLoading && !error && project) {
    }
  }, [isLoading, error, project])

  React.useEffect(() => {
    if (!deleteProject.isPending && !deleteProject.error) {
      if (deleteProject.isSuccess) {
        router.push('/projects')
      }
    }
  }, [deleteProject.isPending, deleteProject.error, deleteProject.isSuccess, router])

  if (isLoading) {
    return (
      <div className="min-h-screen pt-24 pb-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <div data-testid="skeleton" className="space-y-6">
            <Skeleton className="h-8 w-48" />
            <Skeleton className="h-4 w-32" />
            <div className="card p-6 space-y-4">
              <Skeleton className="h-6 w-32" />
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-3/4" />
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !project) {
    return (
      <div className="min-h-screen pt-24 pb-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <EmptyState
            icon={
              <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            }
            title="Project not found"
            description="The project you're looking for doesn't exist or has been deleted."
            action={{
              label: 'Back to Projects',
              onClick: () => router.push('/projects')
            }}
          />
        </div>
      </div>
    )
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  return (
    <div className="min-h-screen pt-24 pb-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <Button
            variant="ghost"
            onClick={() => router.push('/projects')}
            leftIcon={
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            }
          >
            Back to Projects
          </Button>
        </div>

        <div className="space-y-6">
          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{project.name}</h1>
              <p className="mt-1 text-sm text-gray-500">
                Created {formatDate(project.created_at)}
                {project.updated_at !== project.created_at && (
                  <> · Updated {formatDate(project.updated_at)}</>
                )}
              </p>
            </div>
            <div className="flex items-center gap-3">
              <Button
                variant="secondary"
                onClick={() => setIsEditModalOpen(true)}
                leftIcon={
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                }
              >
                Edit
              </Button>
              <Button
                variant="danger"
                onClick={() => setIsDeleteModalOpen(true)}
                leftIcon={
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                }
              >
                Delete
              </Button>
              <Badge 
                variant={project.status === 'active' ? 'success' : 'neutral'}
                size="md"
              >
                {project.status}
              </Badge>
            </div>
          </div>

          <EditProjectModal
            isOpen={isEditModalOpen}
            onClose={() => setIsEditModalOpen(false)}
            project={project}
          />

          <DeleteProjectModal
            isOpen={isDeleteModalOpen}
            onClose={() => setIsDeleteModalOpen(false)}
            projectId={project.id}
            projectName={project.name}
            onSuccess={() => {
              router.push('/projects')
            }}
          />

          {project.description && (
            <Card className="p-6">
              <h2 className="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">
                Description
              </h2>
              <p className="text-gray-700">{project.description}</p>
            </Card>
          )}

          <Card className="p-6">
            <h2 className="text-sm font-medium text-gray-500 uppercase tracking-wide mb-4">
              Progress
            </h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                  <span className="text-gray-700">Tasks</span>
                </div>
                <span className="font-semibold text-gray-900">{project.task_count}</span>
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-600">Completion</span>
                  <span className="text-sm font-medium text-gray-900">{Math.round(project.completion_percentage)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2.5">
                  <div
                    role="progressbar"
                    aria-valuenow={Math.round(project.completion_percentage)}
                    aria-valuemin={0}
                    aria-valuemax={100}
                    className="bg-primary-600 h-2.5 rounded-full transition-all duration-300"
                    style={{ width: `${Math.round(project.completion_percentage)}%` }}
                  />
                </div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}

export default ProjectDetailPage