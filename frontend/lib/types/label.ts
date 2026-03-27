export interface Label {
  id: string
  name: string
  color: string
  task_count: number
}

export interface CreateLabelInput {
  name: string
  color: string
}

export interface UpdateLabelInput {
  name?: string
  color?: string
}