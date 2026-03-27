import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '../api'
import type { Project, CreateProjectInput, UpdateProjectInput, ProjectListFilter } from '../types/project'

const PROJECTS_KEY = 'projects'

export const projectApi = {
  list: async (filter?: ProjectListFilter): Promise<{ data: Project[]; total: number }> => {
    const params = new URLSearchParams()
    if (filter?.page) params.append('page', String(filter.page))
    if (filter?.page_size) params.append('page_size', String(filter.page_size))
    if (filter?.search) params.append('search', filter.search)
    
    const queryString = params.toString()
    const url = queryString ? `/projects?${queryString}` : '/projects'
    return apiClient.get<{ data: Project[]; total: number }>(url)
  },

  get: async (id: string): Promise<Project> => {
    return apiClient.get<Project>(`/projects/${id}`)
  },

  create: async (input: CreateProjectInput): Promise<Project> => {
    return apiClient.post<Project>('/projects', input)
  },

  update: async (id: string, input: UpdateProjectInput): Promise<Project> => {
    return apiClient.patch<Project>(`/projects/${id}`, input)
  },

  delete: async (id: string): Promise<void> => {
    return apiClient.delete(`/projects/${id}`)
  },
}

export function useProjects(filter?: ProjectListFilter) {
  return useQuery({
    queryKey: [PROJECTS_KEY, filter],
    queryFn: () => projectApi.list(filter),
  })
}

export function useProject(id: string) {
  return useQuery({
    queryKey: [PROJECTS_KEY, id],
    queryFn: () => projectApi.get(id),
    enabled: !!id,
  })
}

export function useCreateProject() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (input: CreateProjectInput) => projectApi.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY] })
    },
  })
}

export function useUpdateProject() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateProjectInput }) => 
      projectApi.update(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY] })
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY, id] })
    },
  })
}

export function useDeleteProject() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (id: string) => projectApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY] })
    },
  })
}