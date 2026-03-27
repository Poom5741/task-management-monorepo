import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { ReactNode } from 'react'
import * as React from 'react'

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
    back: vi.fn(),
  }),
  useParams: () => ({
    id: 'test-project-id',
  }),
}))

vi.mock('../../../lib/api', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}))

import { apiClient } from '../../../lib/api'
import ProjectDetailPage from './page'

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

describe('ProjectDetailPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render project details', async () => {
    const mockProject = {
      id: 'test-project-id',
      name: 'Test Project',
      description: 'Test Description',
      status: 'active' as const,
      task_count: 10,
      completion_percentage: 60.0,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    
    vi.mocked(apiClient.get).mockResolvedValueOnce(mockProject)

    render(<ProjectDetailPage />, { wrapper: createWrapper() })

    await waitFor(() => {
      expect(screen.getByText('Test Project')).toBeInTheDocument()
    })

    expect(screen.getByText('Test Description')).toBeInTheDocument()
    expect(screen.getByText('active')).toBeInTheDocument()
    expect(screen.getByText('10')).toBeInTheDocument()
    expect(screen.getByText('60%')).toBeInTheDocument()
  })

  it('should render loading skeleton', () => {
    vi.mocked(apiClient.get).mockImplementation(() => new Promise(() => {}))

    render(<ProjectDetailPage />, { wrapper: createWrapper() })

    expect(screen.getByTestId('skeleton')).toBeInTheDocument()
  })

  it('should render error state when project not found', async () => {
    vi.mocked(apiClient.get).mockRejectedValueOnce(new Error('Not found'))

    render(<ProjectDetailPage />, { wrapper: createWrapper() })

    await waitFor(() => {
      expect(screen.getByText('Project not found')).toBeInTheDocument()
    })
  })

  it('should display progress bar for completion percentage', async () => {
    const mockProject = {
      id: 'test-project-id',
      name: 'Test Project',
      description: 'Test Description',
      status: 'active' as const,
      task_count: 10,
      completion_percentage: 75.0,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    
    vi.mocked(apiClient.get).mockResolvedValueOnce(mockProject)

    render(<ProjectDetailPage />, { wrapper: createWrapper() })

    await waitFor(() => {
      expect(screen.getByText('75%')).toBeInTheDocument()
    })

    const progressBar = screen.getByRole('progressbar')
    expect(progressBar).toHaveAttribute('aria-valuenow', '75')
  })
})