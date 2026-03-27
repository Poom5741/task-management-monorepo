import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { ReactNode } from 'react'
import * as React from 'react'

vi.mock('../api', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}))

import { apiClient } from '../api'
import { useProjects, useProject, useCreateProject } from './projects'

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })
  
  return function Wrapper({ children }: { children: ReactNode }) {
    return React.createElement(
      QueryClientProvider,
      { client: queryClient },
      children
    )
  }
}

describe('projectApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('useProjects', () => {
    it('should fetch projects list', async () => {
      const mockData = {
        data: [
          {
            id: '1',
            name: 'Test Project',
            description: 'Test Description',
            status: 'active' as const,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          },
        ],
        total: 1,
      }
      
      vi.mocked(apiClient.get).mockResolvedValueOnce(mockData)

      const { result } = renderHook(() => useProjects(), {
        wrapper: createWrapper(),
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      
      expect(result.current.data).toEqual(mockData)
    })
  })

  describe('useProject', () => {
    it('should fetch single project', async () => {
      const mockProject = {
        id: '1',
        name: 'Test Project',
        description: 'Test Description',
        status: 'active' as const,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(apiClient.get).mockResolvedValueOnce(mockProject)

      const { result } = renderHook(() => useProject('1'), {
        wrapper: createWrapper(),
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      
      expect(result.current.data).toEqual(mockProject)
    })

    it('should fetch project with task statistics', async () => {
      const mockProject = {
        id: '1',
        name: 'Test Project',
        description: 'Test Description',
        status: 'active' as const,
        task_count: 10,
        completion_percentage: 60.0,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(apiClient.get).mockResolvedValueOnce(mockProject)

      const { result } = renderHook(() => useProject('1'), {
        wrapper: createWrapper(),
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      
      expect(result.current.data?.task_count).toBe(10)
      expect(result.current.data?.completion_percentage).toBe(60.0)
    })

    it('should not fetch when id is empty', () => {
      const { result } = renderHook(() => useProject(''), {
        wrapper: createWrapper(),
      })

      expect(result.current.isFetching).toBe(false)
    })
  })

  describe('useCreateProject', () => {
    it('should create project', async () => {
      const mockProject = {
        id: '1',
        name: 'New Project',
        description: 'New Description',
        status: 'active' as const,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(apiClient.post).mockResolvedValueOnce(mockProject)

      const { result } = renderHook(() => useCreateProject(), {
        wrapper: createWrapper(),
      })

      result.current.mutate({
        name: 'New Project',
        description: 'New Description',
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      
      expect(result.current.data).toEqual(mockProject)
    })
  })

  describe('useProjects with filter', () => {
    it('should fetch projects with pagination filter', async () => {
      const mockData = {
        data: [
          {
            id: '1',
            name: 'Test Project',
            description: 'Test Description',
            status: 'active' as const,
            task_count: 0,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          },
        ],
        total: 1,
      }
      
      vi.mocked(apiClient.get).mockResolvedValueOnce(mockData)

      const { result } = renderHook(() => useProjects({ page: 2, page_size: 20 }), {
        wrapper: createWrapper(),
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      
      expect(result.current.data).toEqual(mockData)
      expect(apiClient.get).toHaveBeenCalledWith('/projects?page=2&page_size=20')
    })

    it('should fetch projects with search filter', async () => {
      const mockData = {
        data: [],
        total: 0,
      }
      
      vi.mocked(apiClient.get).mockResolvedValueOnce(mockData)

      const { result } = renderHook(() => useProjects({ search: 'test' }), {
        wrapper: createWrapper(),
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      
      expect(result.current.data).toEqual(mockData)
      expect(apiClient.get).toHaveBeenCalledWith('/projects?search=test')
    })
  })
})