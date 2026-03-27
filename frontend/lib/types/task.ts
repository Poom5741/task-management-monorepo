export type TaskStatus = 'todo' | 'in-progress' | 'done' | 'cancelled'
export type TaskPriority = 'low' | 'medium' | 'high' | 'urgent'

export interface Task {
  id: string
  project_id: string
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
  due_date?: string
  labels?: string[]
  created_at: string
  updated_at: string
}

export interface CreateTaskInput {
  project_id: string
  title: string
  description?: string
  priority?: TaskPriority
  due_date?: string
  labels?: string[]
}

export interface UpdateTaskInput {
  title?: string
  description?: string
  status?: TaskStatus
  priority?: TaskPriority
  due_date?: string
  labels?: string[]
}

export interface TaskListFilter {
  page?: number
  page_size?: number
  status?: TaskStatus[]
  priority?: TaskPriority[]
  labels?: string[]
  overdue?: boolean
  sort_by?: 'created' | 'due' | 'priority'
  sort_order?: 'asc' | 'desc'
}