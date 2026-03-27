export interface Project {
  id: string
  name: string
  description: string
  status: 'active' | 'archived'
  task_count: number
  completion_percentage: number
  created_at: string
  updated_at: string
}

export interface CreateProjectInput {
  name: string
  description?: string
}

export interface UpdateProjectInput {
  name?: string
  description?: string
  status?: 'active' | 'archived'
}

export interface ProjectListFilter {
  page?: number
  page_size?: number
  search?: string
}

export interface ProjectStatistics {
  total_tasks: number
  todo_tasks: number
  in_progress_tasks: number
  done_tasks: number
  cancelled_tasks: number
  overdue_tasks: number
  completion_rate: number
}